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
	"encoding/json"
	"net/http"
)

type Response struct {
	*http.Response
	err error
}

// JSONScan 将 Body 按照 JSON 反序列化为结构体
func (r *Response) JSONScan(val any) error {
	if r.err != nil {
		return r.err
	}
	err := json.NewDecoder(r.Body).Decode(val)
	return err
}
