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
	"io"
	"net/http"
	"testing"

	"github.com/ecodeclub/ekit/iox"
	"github.com/stretchr/testify/assert"
)

func TestResponse_JSONScan(t *testing.T) {
	testCases := []struct {
		name    string
		resp    *Response
		wantVal User
		wantErr error
	}{
		{
			name: "scan成功",
			resp: &Response{
				Response: &http.Response{
					Body: io.NopCloser(iox.NewJSONReader(User{Name: "Tom"})),
				},
			},
			wantVal: User{
				Name: "Tom",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var u User
			err := tc.resp.JSONScan(&u)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantVal, u)
		})
	}
}
