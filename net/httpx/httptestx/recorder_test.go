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

package httptestx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONResponseRecorder_MustScan(t *testing.T) {
	// 成功案例
	recorder := NewJSONResponseRecorder[User]()
	_, err := recorder.WriteString(`{"name": "Tom"}`)
	require.NoError(t, err)
	u := recorder.MustScan()
	assert.Equal(t, User{Name: "Tom"}, u)

	// panic 案例
	recorder = NewJSONResponseRecorder[User]()
	assert.Panics(t, func() {
		recorder.MustScan()
	})
}

type User struct {
	Name string `json:"name"`
}
