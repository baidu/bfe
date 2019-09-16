// Copyright (c) 2019 Baidu, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mod_access

import (
	"net"
	"net/url"
	"testing"
)

import (
	"github.com/baidu/bfe/bfe_basic"
	"github.com/baidu/bfe/bfe_http"
	"github.com/baidu/bfe/bfe_module"
)

// test tokenTypeGet
func TestTokenTypeGet(t *testing.T) {
	template := "123$status_code$res_header"

	// $status_code
	logType, end, err := tokenTypeGet(&template, 4)
	if err != nil {
		t.Errorf("tokenTypeGet() error: %v", err)
	}
	if logType != fmtTable["status_code"] {
		t.Errorf("logType error, logType: %d", logType)
	}
	if end != 14 {
		t.Errorf("end error, end: %d", end)
	}

	// $res_header
	logType, end, err = tokenTypeGet(&template, 16)
	if err != nil {
		t.Errorf("tokenTypeGet() error: %v", err)
	}
	if logType != fmtTable["res_header"] {
		t.Errorf("logType error, logType: %d", logType)
	}
	if end != 25 {
		t.Errorf("end error, end: %d", end)
	}
}

// test parseBracketToken
func TestParseBracketToken(t *testing.T) {
	template := "{CLIENTIP}res_cookie, log"

	item, end, err := parseBracketToken(&template, 0)
	if err != nil {
		t.Errorf("parseBracketToken() error: %v", err)
	}

	if end != 19 {
		t.Errorf("end error, end: %d", end)
	}

	if item.Key != "CLIENTIP" || item.Type != fmtTable["res_cookie"] {
		t.Errorf("item error, item: %v", item)
	}
}

func createModule(cfg *ConfModAccess, cbs *bfe_module.BfeCallbacks,
	t *testing.T) (*ModuleAccess, error) {
	var err error

	m := NewModuleAccess()
	err = m.init(cfg, cbs, nil)
	if err != nil {
		t.Errorf("ModuleAccess init failed: %v %v", err, cfg)
		return nil, err
	}

	return m, nil
}

// test requestFinish
func TestRequestFinish(t *testing.T) {
	cfg, err := ConfLoad("testdata/mod_access/mod_access.conf")
	if err != nil {
		t.Errorf("ConfLoad() error: %v", err)
		return
	}

	cbs := bfe_module.NewBfeCallbacks()

	m, err := createModule(cfg, cbs, t)
	if err != nil {
		t.Errorf("createModule failed: %v", err)
		return
	}

	req := &bfe_basic.Request{}
	req.Session = &bfe_basic.Session{}
	res := &bfe_http.Response{}
	m.requestFinish(req, res)
}

// test requestFinish
func TestRequestFinish2(t *testing.T) {
	cfg, err := ConfLoad("testdata/mod_access/mod_access.conf")
	if err != nil {
		t.Errorf("ConfLoad() error: %v", err)
		return
	}

	cbs := bfe_module.NewBfeCallbacks()

	m, err := createModule(cfg, cbs, t)
	if err != nil {
		t.Errorf("createModule failed: %v", err)
		return
	}

	req := &bfe_basic.Request{}
	req.Session = &bfe_basic.Session{}
	req.Session.Vip = net.ParseIP("74.125.19.99")
	req.Stat = &bfe_basic.RequestStat{}
	req.HttpRequest = &bfe_http.Request{
		URL: &url.URL{},
	}
	res := &bfe_http.Response{}

	m.requestFinish(req, res)
}

// test sessionFinish
func TestSessionFinish(t *testing.T) {
	cfg, err := ConfLoad("testdata/mod_access/mod_access.conf")
	if err != nil {
		t.Errorf("ConfLoad() error: %v", err)
		return
	}

	cbs := bfe_module.NewBfeCallbacks()

	m, err := createModule(cfg, cbs, t)
	if err != nil {
		t.Errorf("createModule failed: %v", err)
		return
	}

	session := &bfe_basic.Session{}

	m.sessionFinish(session)
}
