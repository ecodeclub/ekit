// Copyright 2021 gotomicro
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

package copier

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReflectCopier_Copy(t *testing.T) {
	testCases := []struct {
		name     string
		copyFunc func() (any, error)
		wantDst  any
		wantErr  error
	}{
		// {
		// 	name: "simple struct",
		// 	copyFunc: func() (any, error) {
		// 		copier := NewReflectCopier[SimpleSrc, SimpleDst]()
		// 		return copier.Copy(&SimpleSrc{
		// 			Name:    "大明",
		// 			Age:     ekit.ToPtr[int](18),
		// 			Friends: []string{"Tom", "Jerry"},
		// 		})
		// 	},
		// 	wantDst: SimpleDst{
		// 		Name:    "大明",
		// 		Age:     ekit.ToPtr[int](18),
		// 		Friends: []string{"Tom", "Jerry"},
		// 	},
		// },
		// 你还需要测试
		// 1. Src 或者 Dst 类型非法，例如基本类型，内置类型或者接口
		// 2. 测试组合（结构体组合，指针组合，接口组合——接口组合可以直接不支持），深层组合，多重组合
		// 3. 复杂类型字段，如字段是结构体，字段是结构体指针，以及多级指针（不需要支持）
		// 4. 类型不匹配
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.copyFunc()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantDst, res)
		})
	}
}

type SimpleSrc struct {
	Name    string
	Age     *int
	Friends []string
}

type SimpleDst struct {
	Name    string
	Age     *int
	Friends []string
}
