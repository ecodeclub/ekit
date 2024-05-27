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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
)

type LoadServiceSuite struct {
	suite.Suite
}

func (l *LoadServiceSuite) SetupTest() {
	t := l.T()
	wd, err := os.Getwd()
	require.NoError(t, err)
	cmd := exec.Command("go", "generate", "./...")
	cmd.Dir = filepath.Join(wd, "testdata")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, fmt.Sprintf("执行 go generate 失败: %v\n%s", err, output))
}

func (l *LoadServiceSuite) Test_LoadService() {
	t := l.T()
	testcases := []struct {
		name     string
		dir      string
		svcName  string
		want     []string
		checkErr func(err error, t *testing.T)
	}{
		{
			name:    "有一个插件",
			dir:     "./testdata/user_service",
			svcName: "UserSvc",
			want:    []string{"Get"},
			checkErr: func(err error, t *testing.T) {

			},
		},
		{
			name:    "有两个插件",
			dir:     "./testdata/user_service2",
			svcName: "UserSvc",
			want:    []string{"A", "B"},
			checkErr: func(err error, t *testing.T) {

			},
		},
		{
			name: "目录不存在",
			dir:  "./notfound",
			checkErr: func(err error, t *testing.T) {
				assert.Equal(t, DirNotFound, err)
			},
		},
		{
			name:    "svcName为空",
			dir:     "./testdata/user_service2",
			svcName: "",
			checkErr: func(err error, t *testing.T) {
				assert.Equal(t, SymEmptyErr, err)
			},
		},
		{
			name:    "svcName没找到",
			dir:     "./testdata/user_service2",
			svcName: "notfound",
			checkErr: func(err error, t *testing.T) {
				assert.NotNil(t, err)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			list, err := LoadService[UserService](tc.dir, tc.svcName)
			tc.checkErr(err, t)
			if err != nil {
				return
			}
			ans := make([]string, 0, len(list))
			for _, svc := range list {
				ans = append(ans, svc.Get())
			}
			assert.Equal(t, tc.want, ans)
		})
	}
}

func TestLoadServiceSuite(t *testing.T) {
	suite.Run(t, new(LoadServiceSuite))
}

type UserService interface {
	Get() string
}

func ExampleLoadService() {
	getters, err := LoadService[UserService]("./testdata/user_service", "UserSvc")
	fmt.Println(err)
	for _, getter := range getters {
		fmt.Println(getter.Get())
	}
	// Output:
	// <nil>
	// Get
}
