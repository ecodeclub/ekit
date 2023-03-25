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

//go:build goexperiment.arenas

package pool

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArenaPool(t *testing.T) {
	testCases := []struct {
		name    string
		p       func() *ArenaPool[TestStruct]
		wantObj *TestStruct
		wantErr error
	}{
		{
			name: "no obj",
			p: func() *ArenaPool[TestStruct] {
				return NewArenaPool[TestStruct]()
			},
			wantObj: &TestStruct{},
		},
		{
			name: "reuse",
			p: func() *ArenaPool[TestStruct] {
				p := NewArenaPool[TestStruct]()
				obj, err := p.Get()
				require.NoError(t, err)
				obj.Obj().Age = 123
				err = p.Put(obj)
				require.NoError(t, err)
				return p
			},
			wantObj: &TestStruct{
				Age: 123,
			},
		},
		{
			name: "multiple",
			p: func() *ArenaPool[TestStruct] {
				p := NewArenaPool[TestStruct]()
				obj1, err := p.Get()
				require.NoError(t, err)

				obj2, err := p.Get()
				require.NoError(t, err)

				obj3, err := p.Get()
				require.NoError(t, err)

				err = p.Put(obj1)
				require.NoError(t, err)
				err = p.Put(obj2)
				require.NoError(t, err)
				err = p.Put(obj3)
				require.NoError(t, err)

				newObj3, err := p.Get()
				require.NoError(t, err)
				assert.Equal(t, obj3, newObj3)
				return p
			},
			wantObj: &TestStruct{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			obj, err := tc.p().Get()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantObj, obj.Obj())
		})
	}
}

type TestStruct struct {
	Age    int
	AgePtr *int
}
