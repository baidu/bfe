// Copyright (c) 2019 Baidu, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     bfe_http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mod_static

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

import (
	"github.com/baidu/go-lib/log"
	"github.com/baidu/go-lib/web-monitor/metrics"
	"github.com/baidu/go-lib/web-monitor/web_monitor"
)

import (
	"github.com/baidu/bfe/bfe_basic"
	"github.com/baidu/bfe/bfe_http"
	"github.com/baidu/bfe/bfe_module"
)

const (
	EncodeGzip = "gzip"

	ModStatic = "mod_static"
)

var (
	openDebug = false

	unixEpochTime = time.Unix(0, 0)
)

type ModuleStaticState struct {
	FileBrowseSize             *metrics.Counter
	FileBrowseCount            *metrics.Counter
	FileBrowseNotExist         *metrics.Counter
	FileBrowseContentTypeError *metrics.Counter
	FileCurrentOpened          *metrics.Gauge
}

type ModuleStatic struct {
	name             string
	state            ModuleStaticState
	metrics          metrics.Metrics
	configPath       string
	mimeTypePath     string
	contentDetection bool
	enableCompress   bool
	ruleTable        *StaticRuleTable
	mimeTypeTable    *MimeTypeTable
}

type staticFile struct {
	http.File
	os.FileInfo
	m *ModuleStatic
}

func newStaticFile(root string, filename string, m *ModuleStatic) (*staticFile, error) {
	var err error
	s := new(staticFile)
	s.m = m
	s.File, err = http.Dir(root).Open(filename)
	if err != nil {
		return nil, err
	}

	s.FileInfo, err = s.File.Stat()
	if err != nil {
		s.File.Close()
		return nil, err
	}

	m.state.FileCurrentOpened.Inc(1)
	return s, nil
}

func (s *staticFile) Close() error {
	err := s.File.Close()
	if err != nil {
		return err
	}

	state := s.m.state
	state.FileCurrentOpened.Dec(1)
	return nil
}

func NewModuleStatic() *ModuleStatic {
	m := new(ModuleStatic)
	m.name = ModStatic
	m.metrics.Init(&m.state, ModStatic, 0)
	m.ruleTable = NewStaticRuleTable()
	m.mimeTypeTable = NewMimeTypeTable()
	return m
}

func (m *ModuleStatic) Name() string {
	return m.name
}

func (m *ModuleStatic) loadConfData(query url.Values) error {
	path := query.Get("path")
	if path == "" {
		path = m.configPath
	}

	conf, err := StaticConfLoad(path)
	if err != nil {
		return fmt.Errorf("error in StaticConfLoad(%s): %v", path, err)
	}

	m.ruleTable.Update(conf)
	return nil
}

func (m *ModuleStatic) loadMimeType(query url.Values) error {
	var err error
	path := query.Get("path")
	if path == "" {
		path = m.mimeTypePath
	}

	conf, err := MimeTypeConfLoad(path)
	if err != nil {
		return fmt.Errorf("error in MimeTypeConfLoad(%s): %v", path, err)
	}
	m.mimeTypeTable.Update(conf)

	return nil
}

func (m *ModuleStatic) getState(params map[string][]string) ([]byte, error) {
	s := m.metrics.GetAll()
	return s.Format(params)
}

func (m *ModuleStatic) getStateDiff(params map[string][]string) ([]byte, error) {
	s := m.metrics.GetDiff()
	return s.Format(params)
}

func (m *ModuleStatic) monitorHandlers() map[string]interface{} {
	handlers := map[string]interface{}{
		m.name:           m.getState,
		m.name + ".diff": m.getStateDiff,
	}
	return handlers
}

func (m *ModuleStatic) reloadHandlers() map[string]interface{} {
	handlers := map[string]interface{}{
		m.name:                              m.loadConfData,
		fmt.Sprintf("%s.mime_type", m.name): m.loadMimeType,
	}
	return handlers
}

func errorStatusCode(err error) int {
	if os.IsNotExist(err) {
		return bfe_http.StatusNotFound
	}
	if os.IsPermission(err) {
		return bfe_http.StatusForbidden
	}

	return bfe_http.StatusInternalServerError
}

func (m *ModuleStatic) tryDefaultFile(root string, defaultFile string) (*staticFile, error) {
	if len(defaultFile) != 0 {
		file, err := newStaticFile(root, defaultFile, m)
		if err != nil {
			return nil, err
		}

		if file.IsDir() {
			if openDebug {
				log.Logger.Debug("%s: %s is directory", m.Name(), file.Name())
			}
			file.Close()
			return nil, fmt.Errorf("directory is not supported")
		}

		return file, nil
	}
	m.state.FileBrowseNotExist.Inc(1)
	return nil, os.ErrNotExist
}

func checkSupportCompress(req *bfe_http.Request) bool {
	header := req.Header
	acceptEncoding := header.GetDirect("Accept-Encoding")
	return bfe_http.HasToken(acceptEncoding, EncodeGzip)
}

func (m *ModuleStatic) tryCompressedStaticFile(req *bfe_http.Request,
	root string) (*staticFile, error) {
	filename := req.URL.Path

	if checkSupportCompress(req) {
		_, err := os.Stat(root + filename + ".gz")
		if m.enableCompress && err == nil {
			return newStaticFile(root, filename+".gz", m)
		}
	}

	return newStaticFile(root, filename, m)
}

func (m *ModuleStatic) openStaticFile(req *bfe_http.Request, root string,
	defaultFile string) (*staticFile, error) {
	file, err := m.tryCompressedStaticFile(req, root)
	if os.IsNotExist(err) {
		file, err = m.tryDefaultFile(root, defaultFile)
	}
	if err != nil {
		return nil, err
	}

	if file.IsDir() {
		if openDebug {
			log.Logger.Debug("%s: %s is directory", m.Name(), file.Name())
		}

		file.Close()
		file, err = m.tryDefaultFile(root, defaultFile)
		if err != nil {
			return nil, err
		}
	}

	return file, nil
}

func (m *ModuleStatic) detectContentType(file *staticFile) (string, error) {
	ext := filepath.Ext(file.Name())

	if ctype, ok := m.mimeTypeTable.Search(strings.ToLower(ext)); ok {
		return ctype, nil
	}

	ctype := mime.TypeByExtension(ext)
	if ctype != "" {
		return ctype, nil
	}

	if !m.contentDetection {
		if openDebug {
			log.Logger.Debug("%s: unknown file extension: %s", m.Name(), ext)
		}

		return "", fmt.Errorf("unknown file extension: %s", ext)
	}

	var buf [512]byte
	n, err := io.ReadFull(file, buf[:])
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return "", err
	}

	ctype = http.DetectContentType(buf[:n])
	_, err = file.Seek(0, io.SeekStart)
	return ctype, err
}

func isZeroTime(t time.Time) bool {
	return t.IsZero() || t.Equal(unixEpochTime)
}

func setLastModified(resp *bfe_http.Response, modtime time.Time) {
	if !isZeroTime(modtime) {
		resp.Header.Set("Last-Modified", modtime.UTC().Format(bfe_http.TimeFormat))
	}
}

func (m *ModuleStatic) createRespFromStaticFile(req *bfe_basic.Request,
	rule *StaticRule) *bfe_http.Response {
	resp := bfe_basic.CreateInternalResp(req, bfe_http.StatusOK)
	root := rule.Action.Params[0]
	defaultFile := rule.Action.Params[1]

	httpRequest := req.HttpRequest
	if httpRequest.Method != "GET" && httpRequest.Method != "HEAD" {
		resp.StatusCode = bfe_http.StatusMethodNotAllowed
		return resp
	}

	file, err := m.openStaticFile(httpRequest, root, defaultFile)
	if err != nil {
		resp.StatusCode = errorStatusCode(err)
		return resp
	}

	ctype, err := m.detectContentType(file)
	if err != nil {
		m.state.FileBrowseContentTypeError.Inc(1)
		resp.StatusCode = errorStatusCode(err)
		return resp
	}

	resp.Header.Set("Content-Type", ctype)
	setLastModified(resp, file.ModTime())
	resp.Body = file
	m.state.FileBrowseSize.Inc(uint(file.Size()))
	return resp
}

func (m *ModuleStatic) staticFileHandler(req *bfe_basic.Request) (int, *bfe_http.Response) {
	rules, ok := m.ruleTable.Search(req.Route.Product)
	if !ok {
		return bfe_module.BfeHandlerGoOn, nil
	}

	for _, rule := range *rules {
		if rule.Cond.Match(req) {
			switch rule.Action.Cmd {
			case ActionBrowse:
				m.state.FileBrowseCount.Inc(1)
				return bfe_module.BfeHandlerResponse, m.createRespFromStaticFile(req, &rule)
			default:
				continue
			}
		}
	}

	return bfe_module.BfeHandlerGoOn, nil
}

func (m *ModuleStatic) Init(cbs *bfe_module.BfeCallbacks, whs *web_monitor.WebHandlers,
	cr string) error {
	var err error
	var cfg *ConfModStatic

	confPath := bfe_module.ModConfPath(cr, m.name)
	if cfg, err = ConfLoad(confPath, cr); err != nil {
		return fmt.Errorf("%s: conf load err: %v", m.name, err)
	}

	openDebug = cfg.Log.OpenDebug
	m.configPath = cfg.Basic.DataPath
	m.mimeTypePath = cfg.Basic.MimeTypePath
	m.contentDetection = cfg.Basic.ContentDetection
	m.enableCompress = cfg.Basic.EnableCompress

	if err = m.loadConfData(nil); err != nil {
		return fmt.Errorf("err in loadConfData(): %v", err)
	}

	if err = m.loadMimeType(nil); err != nil {
		return fmt.Errorf("err in loadMimeType(): %v", err)
	}

	err = cbs.AddFilter(bfe_module.HandleFoundProduct, m.staticFileHandler)
	if err != nil {
		return fmt.Errorf("%s.Init(): AddFilter(m.staticFileHandler): %v", m.name, err)
	}

	err = web_monitor.RegisterHandlers(whs, web_monitor.WebHandleMonitor, m.monitorHandlers())
	if err != nil {
		return fmt.Errorf("%s.Init():RegisterHandlers(m.monitorHandlers): %v", m.name, err)
	}

	err = web_monitor.RegisterHandlers(whs, web_monitor.WebHandleReload, m.reloadHandlers())
	if err != nil {
		return fmt.Errorf("%s.Init(): RegisterHandlers(): %v", m.name, err)
	}

	return nil
}
