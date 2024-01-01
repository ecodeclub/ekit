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
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogRoundTrip(t *testing.T) {
	client := http.DefaultClient
	var acceptLog Log
	var acceptError error
	client.Transport = NewLogRoundTrip(&doNothingRoundTrip{}, func(l Log, err error) {
		acceptLog = l
		acceptError = err
	})
	NewRequest(context.Background(),
		http.MethodGet, "http://localhost/test").
		JSONBody(User{Name: "Tom"}).
		Client(client).
		Do()
	assert.Equal(t, nil, acceptError)
	assert.Equal(t, Log{
		URL:        "http://localhost/test",
		ReqBody:    `{"Name":"Tom"}`,
		RespBody:   "resp body",
		RespStatus: "200 OK",
	}, acceptLog)
}

type doNothingRoundTrip struct {
}

func (d *doNothingRoundTrip) RoundTrip(request *http.Request) (*http.Response, error) {
	return &http.Response{
		Status: "200 OK",
		Body:   io.NopCloser(bytes.NewBuffer([]byte("resp body"))),
	}, nil
}
