// Copyright 2021 ecodeclub
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httpx

import (
	"context"
	"io"
	"net/http"

	"github.com/ecodeclub/ekit/iox"
)

type Request struct {
	req    *http.Request
	err    error
	client *http.Client
}

func NewRequest(ctx context.Context, method, url string) *Request {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	return &Request{
		req:    req,
		err:    err,
		client: http.DefaultClient,
	}
}

// JSONBody 使用 JSON body
func (req *Request) JSONBody(val any) *Request {
	req.req.Body = io.NopCloser(iox.NewJSONReader(val))
	req.req.Header.Set("Content-Type", "application/json")
	return req
}

func (req *Request) Client(cli *http.Client) *Request {
	req.client = cli
	return req
}

func (req *Request) AddHeader(key string, value string) *Request {
	req.req.Header.Add(key, value)
	return req
}

// AddParam 添加查询参数
// 这个方法性能不好，但是好用
func (req *Request) AddParam(key string, value string) *Request {
	q := req.req.URL.Query()
	q.Add(key, value)
	req.req.URL.RawQuery = q.Encode()
	return req
}

func (req *Request) Do() *Response {
	if req.err != nil {
		return &Response{
			err: req.err,
		}
	}
	resp, err := req.client.Do(req.req)
	return &Response{
		Response: resp,
		err:      err,
	}
}
