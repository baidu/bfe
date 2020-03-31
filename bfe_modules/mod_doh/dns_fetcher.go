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

package mod_doh

import (
	"fmt"
	"time"
)

import (
	"github.com/baidu/go-lib/log"
	"github.com/miekg/dns"
)

import (
	"github.com/baidu/bfe/bfe_basic"
	"github.com/baidu/bfe/bfe_http"
)

const (
	ErrBadRequest = "ErrBadRequest"
)

type DnsFetchErr struct {
	Code   string
	Reason string
}

func (e DnsFetchErr) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Reason)
}

type DnsFetcher interface {
	Fetch(req *bfe_basic.Request) (*bfe_http.Response, error)
}

type DnsClient struct {
	address  string
	retryMax int
	timeout  int
}

func NewDnsClient(dnsConf *DnsConf) *DnsClient {
	dnsClient := new(DnsClient)
	dnsClient.address = dnsConf.Address
	dnsClient.retryMax = dnsConf.RetryMax
	dnsClient.timeout = dnsConf.Timeout
	return dnsClient
}

func (c *DnsClient) exchangeWithRetry(msg *dns.Msg) (*dns.Msg, error) {
	var reply *dns.Msg
	var err error

	for retry := 0; retry < c.retryMax; retry++ {
		client := dns.Client{
			Net:     "udp",
			Timeout: time.Duration(c.timeout) * time.Second,
			UDPSize: dns.MaxMsgSize,
		}

		reply, _, err = client.Exchange(msg, c.address)
		if err == nil {
			return reply, nil
		}

		if openDebug {
			log.Logger.Debug("dns client: Exchange error: %v, retry: %d", err, retry)
		}
	}

	return nil, err
}

func (c *DnsClient) Fetch(req *bfe_basic.Request) (*bfe_http.Response, error) {
	msg, err := RequestToDnsMsg(req)
	if err != nil {
		if openDebug {
			log.Logger.Debug("dns client: RequestToDnsMsg error: %v", err)
		}

		return nil, DnsFetchErr{Code: ErrBadRequest, Reason: err.Error()}
	}

	reply, err := c.exchangeWithRetry(msg)
	if err != nil {
		return nil, err
	}

	resp, err := DnsMsgToResponse(req, reply)
	if err != nil {
		if openDebug {
			log.Logger.Debug("dns client: DnsMsgToResponse error: %v", err)
		}

		return nil, err
	}

	return resp, nil
}
