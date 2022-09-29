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

package atomicx

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	val := NewValue[int]()
	assert.Equal(t, 0, val.Load())
	val = NewValueOf[int](123)
	assert.Equal(t, 123, val.Load())
	val.Store(456)
	assert.Equal(t, 456, val.Load())
}

func BenchmarkValue_Load(b *testing.B) {
	b.Run("Value", func(b *testing.B) {
		val := NewValueOf[int](123)
		for i := 0; i < b.N; i++ {
			_ = val.Load()
		}
	})

	b.Run("atomic Value", func(b *testing.B) {
		val := &atomic.Value{}
		val.Store(123)
		for i := 0; i < b.N; i++ {
			_ = val.Load()
		}
	})
}
