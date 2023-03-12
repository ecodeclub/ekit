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

package syncx

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string
}

func TestMap_Load(t *testing.T) {
	testCases := []struct {
		name    string
		key     string
		wantOk  bool
		wantVal *User
	}{
		{
			name:    "found",
			key:     "found",
			wantOk:  true,
			wantVal: &User{Name: "found"},
		},
		{
			name:    "found but empty",
			key:     "found but empty",
			wantOk:  true,
			wantVal: &User{},
		},
		{
			name:   "found but nil",
			key:    "found but nil",
			wantOk: true,
		},
		{
			name: "not found",
			key:  "not found",
		},
	}
	var mu Map[string, *User]
	mu.Store("found", &User{Name: "found"})
	mu.Store("found but empty", &User{})
	mu.Store("found but nil", nil)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, ok := mu.Load(tc.key)
			assert.Equal(t, tc.wantOk, ok)
			assert.Equal(t, tc.wantVal, val)
		})
	}
}

func TestMap_LoadOrStore(t *testing.T) {
	var m = Map[string, *User]{}
	val, loaded := m.LoadOrStore("Tom", &User{Name: "Tom"})
	assert.False(t, loaded)
	assert.Equal(t, &User{Name: "Tom"}, val)

	val, loaded = m.LoadOrStore("Tom", &User{Name: "Tom-copy"})
	assert.True(t, loaded)
	assert.Equal(t, &User{Name: "Tom"}, val)

	val, loaded = m.LoadOrStore("Jerry", nil)
	assert.False(t, loaded)
	assert.Nil(t, val)

	val, loaded = m.LoadOrStore("Jerry", &User{Name: "Jerry"})
	assert.True(t, loaded)
	assert.Nil(t, val)
}

func TestMap_LoadAndDelete(t *testing.T) {
	var m = Map[string, *User]{}
	m.Store("Tom", nil)
	val, loaded := m.LoadAndDelete("Tom")
	assert.True(t, loaded)
	assert.Nil(t, val)

	val, loaded = m.LoadAndDelete("Tom")
	assert.False(t, loaded)
	assert.Nil(t, val)

	m.Store("Jerry", &User{Name: "Jerry"})
	val, loaded = m.LoadAndDelete("Jerry")
	assert.True(t, loaded)
	assert.Equal(t, &User{Name: "Jerry"}, val)

	val, loaded = m.LoadAndDelete("Jerry")
	assert.False(t, loaded)
	assert.Nil(t, val)
}

func TestMap_Delete(t *testing.T) {
	var m = Map[string, *User]{}
	m.Store("Tom", &User{Name: "Tom"})
	val, ok := m.Load("Tom")
	assert.True(t, ok)
	assert.Equal(t, &User{Name: "Tom"}, val)
	m.Delete("Tom")
	val, ok = m.Load("Tom")
	assert.False(t, ok)
	assert.Nil(t, val)
}

func TestMap_Range(t *testing.T) {
	var m = Map[string, *User]{}
	m.Store("Tom", &User{Name: "Tom"})
	m.Store("Jerry", &User{Name: "Jerry"})
	m.Store("nil", nil)

	shadow := make(map[string]*User, 3)
	m.Range(func(key string, val *User) bool {
		shadow[key] = val
		return true
	})
	assert.Equal(t, map[string]*User{
		"Tom":   {Name: "Tom"},
		"Jerry": {Name: "Jerry"},
		"nil":   nil,
	}, shadow)

	var ptrKeyMap Map[*User, string]
	key1 := &User{Name: "Tom"}
	var key2 *User
	ptrKeyMap.Store(key1, "Tom")
	ptrKeyMap.Store(key2, "nil")
	ptrShadow := make(map[*User]string, 2)
	ptrKeyMap.Range(func(key *User, val string) bool {
		ptrShadow[key] = val
		return true
	})
	assert.Equal(t, map[*User]string{
		key1: "Tom",
		nil:  "nil",
	}, ptrShadow)
}

func ExampleMap_LoadAndDelete() {
	var m = Map[string, *User]{}
	_, loaded := m.LoadAndDelete("Tom")
	fmt.Println(loaded)
	m.Store("Tom", nil)
	val, loaded := m.LoadAndDelete("Tom")
	if loaded {
		fmt.Printf("key=Tom, val=%v 被删除\n", val)
	}
	m.Store("Tom", &User{Name: "Tom"})
	val, loaded = m.LoadAndDelete("Tom")
	if loaded {
		fmt.Printf("key=Tom, val=%v 被删除\n", val)
	}
	// Output:
	// false
	// key=Tom, val=<nil> 被删除
	// key=Tom, val=&{Tom} 被删除
}

func ExampleMap_LoadOrStore() {
	var m = Map[string, *User]{}
	_, loaded := m.LoadOrStore("Tom", &User{Name: "Tom"})
	// 执行存储
	if !loaded {
		fmt.Println("设置了新值 Tom")
	}

	_, loaded = m.LoadOrStore("Tom", &User{Name: "Tom-copy"})
	// Tom 这个 key 已经存在，执行加载
	if loaded {
		fmt.Println("加载旧值 Tom")
	}

	_, loaded = m.LoadOrStore("Jerry", nil)
	// 执行存储，注意值是 nil
	if !loaded {
		fmt.Println("设置了新值 nil")
	}
	val, loaded := m.LoadOrStore("Jerry", &User{Name: "Jerry"})
	// Jerry 这个 key 已经存在，执行加载，于是把原本的 nil 加载出来
	if loaded {
		fmt.Printf("加载旧值 %v", val)
	}
	// Output:
	// 设置了新值 Tom
	// 加载旧值 Tom
	// 设置了新值 nil
	// 加载旧值 <nil>
}

func ExampleMap_Range() {
	var m Map[string, int]
	m.Store("Tom", 18)
	m.Store("Jerry", 35)
	var sum int
	m.Range(func(key string, val int) bool {
		sum += val
		return true
	})
	fmt.Println(sum)
	// Output:
	// 53
}

func ExampleMap_Store() {
	var m Map[string, int]
	m.Store("key1", 123)
	val, ok := m.Load("key1")
	if ok {
		fmt.Printf("key1 = %d\n", val)
	}
	// Output:
	// key1 = 123
}

func ExampleMap_Delete() {
	var m Map[string, int]
	m.Store("key1", 123)
	val, ok := m.Load("key1")
	if ok {
		fmt.Printf("key1 = %d\n", val)
	}
	m.Delete("key1")
	_, ok = m.Load("key1")
	if !ok {
		fmt.Println("key1 已被删")
	}
	// Output:
	// key1 = 123
	// key1 已被删
}

func ExampleMap_Load() {
	var m Map[string, int]
	m.Store("key1", 123)
	val, ok := m.Load("key1")
	if ok {
		fmt.Println(val)
	}
	// Output:
	// 123
}

func BenchmarkLoad(b *testing.B) {
	var m Map[string, int]
	m.Store("key1", 123)

	var sm sync.Map
	sm.Store("key1", 123)

	b.Run("Map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = m.Load("key1")
		}
	})

	b.Run("sync.Map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sm.Load("key1")
		}
	})
}
