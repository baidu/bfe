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

package mod_doh

import (
	"testing"
)

import (
	"github.com/baidu/go-lib/web-monitor/web_monitor"
	"github.com/miekg/dns"
)

import (
	"github.com/baidu/bfe/bfe_basic"
	"github.com/baidu/bfe/bfe_http"
	"github.com/baidu/bfe/bfe_module"
)

type TestDnsFetcher struct{}

func (f *TestDnsFetcher) Fetch(req *bfe_http.Request, network string, address string) (*dns.Msg, error) {
	return buildDnsMsg(), nil
}

func TestDohHandler(t *testing.T) {
	m := NewModuleDoh()
	cb := bfe_module.NewBfeCallbacks()
	wh := web_monitor.NewWebHandlers()
	err := m.Init(cb, wh, "./testdata")
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}
	m.dnsFetcher = new(TestDnsFetcher)

	httpRequest := buildDohRequest("GET", t)
	req := new(bfe_basic.Request)
	req.HttpRequest = httpRequest
	req.Session = new(bfe_basic.Session)
	req.Route.Product = "unittest"

	ret, resp := m.dohHandler(req)
	if ret != bfe_module.BfeHandlerResponse {
		t.Errorf("ret should be %d, not %d", bfe_module.BfeHandlerResponse, ret)
		return
	}
	if resp.StatusCode != bfe_http.StatusOK {
		t.Errorf("status code should be %d, not %d", bfe_http.StatusOK, resp.StatusCode)
	}
}
