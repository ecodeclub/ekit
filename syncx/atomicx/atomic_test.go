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

package atomicx

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValueOf(t *testing.T) {
	testCases := []struct {
		name  string
		input *User
	}{
		{
			name: "nil",
		},
		{
			name: "user",
			input: &User{
				Name: "Tom",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val := NewValueOf[*User](tc.input)
			assert.Equal(t, tc.input, val.Load())
		})
	}
}

func TestValue_CompareAndSwap(t *testing.T) {
	testCases := []struct {
		name string
		old  *User
		new  *User
	}{
		{
			name: "both nil",
		},
		{
			name: "old nil",
			new:  &User{Name: "Tom"},
		},
		{
			name: "new nil",
			old:  &User{Name: "Tom"},
		},
		{
			name: "not nil",
			new:  &User{Name: "Tom"},
			old:  &User{Name: "Jerry"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val := NewValueOf[*User](tc.old)
			swapped := val.CompareAndSwap(tc.old, tc.new)
			assert.True(t, swapped)
		})
	}
}

func TestValue_Swap(t *testing.T) {
	testCases := []struct {
		name string
		old  *User
		new  *User
	}{
		{
			name: "both nil",
		},
		{
			name: "old nil",
			new:  &User{Name: "Tom"},
		},
		{
			name: "new nil",
			old:  &User{Name: "Tom"},
		},
		{
			name: "not nil",
			new:  &User{Name: "Tom"},
			old:  &User{Name: "Jerry"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val := NewValueOf[*User](tc.old)
			oldVal := val.Swap(tc.new)
			newVal := val.Load()
			assert.Equal(t, tc.old, oldVal)
			assert.Equal(t, tc.new, newVal)
		})
	}
}

func TestValue_Store_Load(t *testing.T) {
	testCases := []struct {
		name    string
		input   *User
		wantVal *User
	}{
		{
			name: "nil",
		},
		{
			name: "user",
			input: &User{
				Name: "Tom",
			},
			wantVal: &User{
				Name: "Tom",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val := NewValue[*User]()
			val.Store(tc.input)
			v := val.Load()
			assert.Equal(t, tc.wantVal, v)
		})
	}
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

func BenchmarkValue_Store(b *testing.B) {
	b.Run("Value", func(b *testing.B) {
		val := NewValue[int]()
		for i := 0; i < b.N; i++ {
			val.Store(123)
		}
	})

	b.Run("atomic Value", func(b *testing.B) {
		val := &atomic.Value{}

		for i := 0; i < b.N; i++ {
			val.Store(123)
		}
	})
}

func BenchmarkValue_Swap(b *testing.B) {
	b.Run("Value", func(b *testing.B) {
		val := NewValueOf[int](123)
		for i := 0; i < b.N; i++ {
			_ = val.Swap(456)
		}
	})

	b.Run("atomic Value", func(b *testing.B) {
		val := &atomic.Value{}
		val.Store(123)
		for i := 0; i < b.N; i++ {
			_ = val.Swap(456)
		}
	})
}

func BenchmarkValue_CompareAndSwap(b *testing.B) {
	b.Run("Value", func(b *testing.B) {
		b.Run("fail", func(b *testing.B) {
			val := NewValueOf[int](123)
			for i := 0; i < b.N; i++ {
				_ = val.CompareAndSwap(-1, 100)
			}
		})
		b.Run("success", func(b *testing.B) {
			val := NewValueOf[int](0)
			for i := 0; i < b.N; i++ {
				_ = val.CompareAndSwap(i, i+1)
			}
		})
	})

	b.Run("atomic Value", func(b *testing.B) {
		b.Run("fail", func(b *testing.B) {
			val := &atomic.Value{}
			val.Store(123)
			for i := 0; i < b.N; i++ {
				_ = val.CompareAndSwap(-1, 100)
			}
		})
		b.Run("success", func(b *testing.B) {
			val := &atomic.Value{}
			val.Store(0)
			for i := 0; i < b.N; i++ {
				_ = val.CompareAndSwap(i, i+1)
			}
		})
	})
}

type User struct {
	Name string
}
