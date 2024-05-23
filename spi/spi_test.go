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

package spi

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func Test_LoadService(t *testing.T) {
	testcases := []struct {
		name    string
		dir     string
		svcName string
		want    []string
		wantErr error
	}{
		{
			name:    "有一个插件",
			dir:     "./user_service",
			svcName: "UserSvc",
			want:    []string{"Get"},
		},
		{
			name:    "有两个插件",
			dir:     "./user_service2",
			svcName: "UserSvc",
			want:    []string{"A", "B"},
		},
		{
			name:    "目录不存在",
			dir:     "./notfound",
			wantErr: DirNotFound,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			list, err := LoadService[UserService](tc.dir, tc.svcName)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			ans := make([]string, 0, len(list))
			for _, svc := range list {
				ans = append(ans, svc.Get())
			}
			log.Println(ans)
			assert.Equal(t, tc.want, ans)
		})
	}
}

type UserService interface {
	Get() string
}
