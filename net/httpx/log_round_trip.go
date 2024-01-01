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
	"bytes"
	"io"
	"net/http"
)

type LogRoundTrip struct {
	delegate http.RoundTripper
	// l 绝对不会为 nil
	log func(l Log, err error)
}

func NewLogRoundTrip(rp http.RoundTripper, log func(l Log, err error)) *LogRoundTrip {
	return &LogRoundTrip{
		delegate: rp,
		log:      log,
	}
}

func (l *LogRoundTrip) RoundTrip(request *http.Request) (resp *http.Response, err error) {
	log := Log{
		URL: request.URL.String(),
	}
	defer func() {
		if resp != nil {
			log.RespStatus = resp.Status
			if resp.Body != nil {
				// 出现 error 了这里也不知道怎么处理，暂时忽略
				body, _ := io.ReadAll(resp.Body)
				resp.Body = io.NopCloser(bytes.NewReader(body))
				log.RespBody = string(body)
			}
		}
		l.log(log, err)
	}()
	if request.Body != nil {
		// 出现 error 了这里也不知道怎么处理，暂时忽略
		body, _ := io.ReadAll(request.Body)
		request.Body = io.NopCloser(bytes.NewReader(body))
		log.ReqBody = string(body)
	}
	resp, err = l.delegate.RoundTrip(request)
	return
}

type Log struct {
	URL        string
	ReqBody    string
	RespBody   string
	RespStatus string
}
