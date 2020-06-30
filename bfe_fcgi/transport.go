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

package bfe_fcgi

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/textproto"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bfenetworks/bfe/bfe_config/bfe_cluster_conf/cluster_conf"
	"github.com/bfenetworks/bfe/bfe_http"
)

type Transport struct{}

func (t *Transport) RoundTrip(req *bfe_http.Request) (*bfe_http.Response, error) {
	metaData := map[string]string{}
	for k, vs := range req.Header {
		metaData[strings.ToUpper(k)] = strings.Join(vs, ",")
	}

	// TODO timeout
	// TODO: reuse conn
	// there are two way to resue conn:
	// 1. implement it in FCGIClient, the same logic like bfe_http.persistConn
	// 2. modify bfe_http.ReadResponse() and bfe_http.Request.write to deal with
	// fastcgi protocol, so we can reuse bfe_http.persistConn directly
	client, err := Dial("tcp", req.URL.Host)
	if err != nil {
		return nil, ConnectError{
			Addr: req.URL.Host,
			Err:  err,
		}
	}

	reader, err := client.Do(metaData, req.Body)
	if err != nil {
		return nil, WriteRequestError{
			Err: err,
		}
	}

	rsp, err := ReadResponse(reader, req)
	if err != nil {
		return nil, ReadRespHeaderError{
			Err: err,
		}
	}
	return rsp, nil
}

// https://tools.ietf.org/html/rfc3875#section-4
func BuildMetaValsAndMethod(r *bfe_http.Request, fc *cluster_conf.ClusterFCGIConf) {
	ip, port := r.RemoteAddr, ""
	if idx := strings.LastIndex(r.RemoteAddr, ":"); idx > -1 {
		ip = r.RemoteAddr[:idx]
		port = r.RemoteAddr[idx+1:]
	}
	ip = strings.Replace(ip, "[", "", 1)
	ip = strings.Replace(ip, "]", "", 1)

	fpath := r.URL.Path

	root := fc.Root
	docURI, pathInfo, scriptName := fpath, "", fpath
	scriptName = strings.TrimSuffix(scriptName, pathInfo)
	scriptFilename := filepath.Join(root, scriptName)

	reqHost, reqPort, err := net.SplitHostPort(r.Host)
	if err != nil {
		reqHost = r.Host
	}

	metaHeader := bfe_http.Header{}
	metaHeader.Add("GATEWAY_INTERFACE", "CGI/1.1")
	metaHeader.Add("SERVER_SOFTWARE", "BFE")

	metaHeader.Add("AUTH_TYPE", "")
	metaHeader.Add("CONTENT_LENGTH", r.Header.Get("Content-Length"))
	metaHeader.Add("CONTENT_TYPE", r.Header.Get("Content-Type"))
	metaHeader.Add("PATH_INFO", pathInfo)
	metaHeader.Add("QUERY_STRING", r.URL.RawQuery)
	metaHeader.Add("REMOTE_ADDR", ip)
	metaHeader.Add("REMOTE_HOST", ip)
	metaHeader.Add("REMOTE_PORT", port)
	metaHeader.Add("REMOTE_IDENT", "")
	metaHeader.Add("REMOTE_USER", "")
	metaHeader.Add("REQUEST_METHOD", r.Method)
	metaHeader.Add("REQUEST_SCHEME", r.URL.Scheme)
	metaHeader.Add("SERVER_NAME", reqHost)
	metaHeader.Add("SERVER_PORT", reqPort)
	metaHeader.Add("SERVER_PROTOCOL", r.Proto)

	metaHeader.Add("DOCUMENT_ROOT", root)
	metaHeader.Add("DOCUMENT_URI", docURI)
	metaHeader.Add("HTTP_HOST", r.Host)
	metaHeader.Add("REQUEST_URI", r.URL.RequestURI())
	metaHeader.Add("SCRIPT_FILENAME", scriptFilename)
	metaHeader.Add("SCRIPT_NAME", scriptName)

	if metaHeader.Get("PATH_INFO") == "" {
		metaHeader.Add("PATH_INFO", filepath.Join(root, pathInfo))
	}

	// add config
	for key, value := range fc.EnvVars {
		metaHeader.Set(key, value)
	}

	// https://tools.ietf.org/html/rfc3875#section-4.1.18
	for key, val := range r.Header {
		header := strings.Replace(strings.ToUpper(key), "-", "_", -1)
		metaHeader.Add("HTTP_"+header, strings.Join(val, ", "))
	}

	// build Method
	metaHeader.Set("REQUEST_METHOD", r.Method)
	metaHeader.Set("CONTENT_LENGTH", fmt.Sprintf("%d", r.ContentLength))
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/x-www-form-urlencoded"
	}
	metaHeader.Set("CONTENT_TYPE", contentType)

	r.Header = metaHeader
}

func ReadResponse(reader io.Reader, req *bfe_http.Request) (*bfe_http.Response, error) {
	rb := bufio.NewReader(reader)
    //	rb.WriteTo(os.Stderr)
	tp := textproto.NewReader(rb)
	// tp := textproto.NewReader(reader)
	resp := &bfe_http.Response{
		Request: req,
	}

	// Parse the response headers.
	mimeHeader, err := tp.ReadMIMEHeader()
	if err != nil && err != io.EOF {
		return nil, err
	}
	resp.Header = bfe_http.Header(mimeHeader)

	if resp.Header.Get("Status") != "" {
		statusParts := strings.SplitN(resp.Header.Get("Status"), " ", 2)
		resp.StatusCode, err = strconv.Atoi(statusParts[0])
		if err != nil {
			return nil, err
		}
		if len(statusParts) > 1 {
			resp.Status = statusParts[1]
		}
	} else {
		resp.StatusCode = http.StatusOK
	}

	// TODO: fixTransferEncoding ?
	resp.TransferEncoding = resp.Header["Transfer-Encoding"]
	resp.ContentLength, _ = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)

	if chunked(resp.TransferEncoding) {
		resp.Body = ioutil.NopCloser(httputil.NewChunkedReader(rb))
	} else {
		resp.Body = ioutil.NopCloser(rb)
	}
	return resp, nil
}
