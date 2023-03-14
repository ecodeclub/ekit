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
)

func TestArenaPool(t *testing.T) {

	type Test struct {
		A int
	}
	pool := NewArenaPool[Test]()
	t.Run("no box in pool", func(t *testing.T) {
		testObject, err := pool.Get()
		assert.NoError(t, err)
		assert.Equal(t, 0, testObject.Object().A)

	})

	t.Run("box already in pool", func(t *testing.T) {
		testObject, err := pool.Get()
		assert.NoError(t, err)
		err = pool.Put(testObject)
		assert.NoError(t, err)
		testObject1, err := pool.Get()
		assert.Equal(t, testObject, testObject1)
	})

}
