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

package mod_compress

import (
	"io/ioutil"
	"os"
	"testing"
)

import (
	"github.com/baidu/go-lib/web-monitor/web_monitor"
)

import (
	"github.com/baidu/bfe/bfe_basic"
	"github.com/baidu/bfe/bfe_http"
	"github.com/baidu/bfe/bfe_module"
)

func prepareModule() *ModuleCompress {
	m := NewModuleCompress()
	m.Init(bfe_module.NewBfeCallbacks(), web_monitor.NewWebHandlers(), "./testdata")
	return m
}

func prepareRequest() *bfe_basic.Request {
	req := new(bfe_basic.Request)
	req.HttpRequest = new(bfe_http.Request)
	req.HttpRequest.Header = make(bfe_http.Header)
	req.HttpRequest.Header.Set("Accept-Encoding", EncodeGzip)
	req.Route.Product = "unittest"
	req.Session = new(bfe_basic.Session)
	req.Context = make(map[interface{}]interface{})
	return req
}

func prepareResponse(filename string) *bfe_http.Response {
	res := new(bfe_http.Response)
	res.StatusCode = 200
	res.Header = make(bfe_http.Header)
	res.Body, _ = os.Open(filename)
	return res
}

func TestModuleCompress(t *testing.T) {
	filename := "testdata/mod_compress/data.txt"
	rawData, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	m := prepareModule()
	req := prepareRequest()
	res := prepareResponse(filename)

	m.checkHandler(req)
	val := req.GetContext(ReqCtxEncodeInfo)
	if val == nil {
		t.Errorf("unexpected encode info")
		return
	}

	encodeInfo := val.(*EncodeInfo)
	if encodeInfo.Quality != 9 {
		t.Errorf("unexpected compress quality: %d", encodeInfo.Quality)
		return
	}
	if encodeInfo.FlushSize != 512 {
		t.Errorf("unexpected compress flushSize: %d", encodeInfo.FlushSize)
		return
	}

	m.compressHandler(req, res)
	contentEncoding := res.Header.GetDirect("Content-Encoding")
	if !bfe_http.HasToken(contentEncoding, EncodeGzip) {
		t.Errorf("unexpected content encoding: %s", contentEncoding)
		return
	}

	compressedData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	res.Body.Close()
	if len(compressedData) > len(rawData) {
		t.Errorf("unexpected compressed data")
	}
}
