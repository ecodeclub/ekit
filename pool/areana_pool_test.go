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

//go:build goexperiment.arenas

package pool

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var count int

type Test struct {
	A int
}

func NewTest() *Test {
	count++
	return &Test{
		A: 1,
	}
}

func TestArenaPool(t *testing.T) {
	pool := NewArenaPool[Test](NewTest)
	t.Run("no box in pool", func(t *testing.T) {
		count = 0
		testObject, err := pool.Get()
		assert.NoError(t, err)
		assert.Equal(t, 1, testObject.Object().A)
		assert.Equal(t, 1, count)
	})

	t.Run("box already in pool", func(t *testing.T) {
		count = 0
		testObject, err := pool.Get()
		assert.NoError(t, err)
		pool.Put(testObject)
		testObject1, err := pool.Get()
		assert.Equal(t, testObject, testObject1)
		assert.Equal(t, 1, count)
	})

}
