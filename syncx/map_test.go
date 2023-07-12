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
	"errors"
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
	mu.Store("found", testCases[0].wantVal)
	mu.Store("found but empty", testCases[1].wantVal)
	mu.Store("found but nil", nil)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, ok := mu.Load(tc.key)
			assert.Equal(t, tc.wantOk, ok)
			assert.Same(t, tc.wantVal, val)
		})
	}
}

func TestMap_LoadOrStore(t *testing.T) {

	t.Run("store non-nil value", func(t *testing.T) {
		m, user := Map[string, *User]{}, &User{Name: "Tom"}
		val, loaded := m.LoadOrStore(user.Name, user)
		assert.False(t, loaded)
		assert.Same(t, user, val)
	})

	t.Run("load non-nil value", func(t *testing.T) {
		m, user := Map[string, *User]{}, &User{Name: "Tom"}
		val, loaded := m.LoadOrStore(user.Name, user)
		assert.False(t, loaded)
		assert.Same(t, user, val)

		val, loaded = m.LoadOrStore("Tom", &User{Name: "Tom-copy"})

		assert.True(t, loaded)
		assert.Same(t, user, val)
	})

	t.Run("store nil value", func(t *testing.T) {
		m, user := Map[string, *User]{}, &User{Name: "Jerry"}
		val, loaded := m.LoadOrStore(user.Name, nil)
		assert.False(t, loaded)
		assert.Nil(t, val)
	})

	t.Run("load nil value", func(t *testing.T) {
		m, user := Map[string, *User]{}, &User{Name: "Jerry"}
		val, loaded := m.LoadOrStore(user.Name, nil)
		assert.False(t, loaded)
		assert.Nil(t, val)

		val, loaded = m.LoadOrStore(user.Name, user)

		assert.True(t, loaded)
		assert.Nil(t, val)
	})
}

func TestMap_LoadOrStoreFunc(t *testing.T) {

	t.Run("store non-nil value returned by func", func(t *testing.T) {
		m, user := Map[string, *User]{}, &User{Name: "Tom"}

		val, loaded, err := m.LoadOrStoreFunc(user.Name, func() (*User, error) {
			return user, nil
		})

		assert.NoError(t, err)
		assert.False(t, loaded)
		assert.Same(t, user, val)
	})

	t.Run("load non-nil value returned by func", func(t *testing.T) {
		m, user := Map[string, *User]{}, &User{Name: "Tom"}
		val, loaded, err := m.LoadOrStoreFunc(user.Name, func() (*User, error) {
			return user, nil
		})
		assert.NoError(t, err)
		assert.False(t, loaded)
		assert.Same(t, user, val)

		val, loaded, err = m.LoadOrStoreFunc(user.Name, func() (*User, error) {
			return &User{Name: "Tom"}, nil
		})

		assert.NoError(t, err)
		assert.True(t, loaded)
		assert.Same(t, user, val)
	})

	t.Run("store nil value returned by func", func(t *testing.T) {
		m, user := Map[string, *User]{}, &User{Name: "Tom"}

		val, loaded, err := m.LoadOrStoreFunc(user.Name, func() (*User, error) {
			return nil, nil
		})

		assert.NoError(t, err)
		assert.False(t, loaded)
		assert.Nil(t, val)
	})

	t.Run("load nil value returned by func", func(t *testing.T) {
		m, user := Map[string, *User]{}, &User{Name: "Tom"}
		val, loaded, err := m.LoadOrStoreFunc(user.Name, func() (*User, error) {
			return nil, nil
		})
		assert.NoError(t, err)
		assert.False(t, loaded)
		assert.Nil(t, val)

		val, loaded, err = m.LoadOrStoreFunc(user.Name, func() (*User, error) {
			return nil, nil
		})

		assert.NoError(t, err)
		assert.True(t, loaded)
		assert.Nil(t, val)
	})

	t.Run("got error returned by func", func(t *testing.T) {
		m := Map[string, *User]{}
		val, loaded, err := m.LoadOrStoreFunc("Jerry", func() (*User, error) {
			return nil, errors.New("初始话失败")
		})
		assert.Equal(t, err, errors.New("初始话失败"))
		assert.False(t, loaded)
		assert.Equal(t, (*User)(nil), val)
	})
}

func TestMap_LoadAndDelete(t *testing.T) {

	t.Run("non-nil value", func(t *testing.T) {
		m, user := Map[string, *User]{}, &User{Name: "Jerry"}
		m.Store("Jerry", user)

		val, loaded := m.LoadAndDelete(user.Name)
		assert.True(t, loaded)
		assert.Same(t, user, val)

		val, loaded = m.LoadAndDelete(user.Name)
		assert.False(t, loaded)
		assert.Nil(t, val)
	})

	t.Run("nil value", func(t *testing.T) {
		m, user := Map[string, *User]{}, &User{Name: "Tom"}
		m.Store(user.Name, nil)

		val, loaded := m.LoadAndDelete(user.Name)
		assert.True(t, loaded)
		assert.Nil(t, val)

		val, loaded = m.LoadAndDelete(user.Name)
		assert.False(t, loaded)
		assert.Nil(t, val)
	})
}

func TestMap_Delete(t *testing.T) {
	m, user := Map[string, *User]{}, &User{Name: "Tom"}
	m.Store(user.Name, user)
	val, ok := m.Load(user.Name)
	assert.True(t, ok)
	assert.Same(t, user, val)

	m.Delete(user.Name)

	val, ok = m.Load(user.Name)
	assert.False(t, ok)
	assert.Nil(t, val)
}

func TestMap_Range(t *testing.T) {

	t.Run("non-pointer type key", func(t *testing.T) {
		m, tom, jerry := Map[string, *User]{}, &User{Name: "Tom"}, &User{Name: "Jerry"}
		var zero *User
		m.Store(tom.Name, tom)
		m.Store(jerry.Name, jerry)
		m.Store("zero", zero)
		m.Store("nil", nil)

		shadow := make(map[string]*User, 4)
		m.Range(func(key string, val *User) bool {
			shadow[key] = val
			return true
		})

		assert.Same(t, tom, shadow[tom.Name])
		assert.Same(t, jerry, shadow[jerry.Name])
		assert.Same(t, zero, shadow["zero"])
		assert.Same(t, (*User)(nil), shadow["nil"])
	})

	t.Run("pointer type key", func(t *testing.T) {
		m, tom := Map[*User, string]{}, &User{Name: "Tom"}
		var zero *User
		m.Store(tom, "Tom")
		m.Store(zero, "nil")

		shadow := make(map[*User]string, 2)
		m.Range(func(key *User, val string) bool {
			shadow[key] = val
			return true
		})

		assert.Equal(t, shadow[tom], tom.Name)
		assert.Equal(t, shadow[zero], "nil")
		assert.Equal(t, shadow[nil], "nil")
	})
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

func ExampleMap_LoadOrStoreFunc() {
	var m = Map[string, *User]{}
	_, loaded, _ := m.LoadOrStoreFunc("Tom", func() (*User, error) {
		return &User{Name: "Tom"}, nil
	})
	// 执行存储
	if !loaded {
		fmt.Println("设置了新值 Tom")
	}

	_, loaded, _ = m.LoadOrStoreFunc("Tom", func() (*User, error) {
		return &User{Name: "Tom-copy"}, nil
	})
	// Tom 这个 key 已经存在，执行加载
	if loaded {
		fmt.Println("加载旧值 Tom")
	}

	_, loaded, _ = m.LoadOrStoreFunc("Jerry", func() (*User, error) {
		return nil, nil
	})
	// 执行存储，注意值是 nil
	if !loaded {
		fmt.Println("设置了新值 nil")
	}
	val, loaded, _ := m.LoadOrStoreFunc("Jerry", func() (*User, error) {
		return &User{Name: "Jerry"}, nil
	})
	// Jerry 这个 key 已经存在，执行加载，于是把原本的 nil 加载出来
	if loaded {
		fmt.Printf("加载旧值 %v\n", val)
	}

	_, _, err := m.LoadOrStoreFunc("Kitty", func() (*User, error) {
		return nil, errors.New("初始化失败")
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	// Output:
	// 设置了新值 Tom
	// 加载旧值 Tom
	// 设置了新值 nil
	// 加载旧值 <nil>
	// 初始化失败
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
