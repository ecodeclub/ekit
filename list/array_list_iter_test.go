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

package list

import (
	"github.com/ecodeclub/ekit/internal/iterator"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayList_Iter_range(t *testing.T) {
	testCases := []struct {
		name      string
		list      *ArrayList[int]
		index     int
		wantSlice []int
		wantErr   error
	}{
		{
			name:      "遍历数组",
			list:      NewArrayListOf[int]([]int{1, 2, 3}),
			wantSlice: []int{1, 2, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//err := tc.list.Add(tc.index, tc.newVal)
			iter := tc.list.Iterator()
			// 这里处理error 很恶心
			ints := make([]int, 0, tc.list.Len())
			var err error
			for iter.HasNext() {
				next, err := iter.Next()
				assert.NoError(t, err)
				ints = append(ints, next)
			}
			assert.Equal(t, tc.wantErr, err)
			// 因为返回了 error，所以我们不用继续往下比较了
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantSlice, ints)
		})
	}
}

func TestArrayList_Iter_delete(t *testing.T) {
	testCases := []struct {
		name      string
		list      *ArrayList[int]
		index     int
		fn        func(ite iterator.Iterator[int]) error
		wantSlice []int
		wantErr   error
	}{
		{
			name:  "删除某个元素",
			list:  NewArrayListOf[int]([]int{1, 2, 3, 4, 5}),
			index: 3,
			fn: func(ite iterator.Iterator[int]) error {
				err := ite.Delete()
				if err != nil {
					return err
				}
				if ite.HasNext() {
					_, err := ite.Next()
					if err != nil {
						return err
					}
				}
				return nil
			},
			wantSlice: []int{1, 2, 4, 5},
		},
		{
			name:  "连续删除元素",
			list:  NewArrayListOf[int]([]int{1, 2, 3, 4, 5}),
			index: 3,
			fn: func(ite iterator.Iterator[int]) error {
				err := ite.Delete()
				if err != nil {
					return err
				}
				err = ite.Delete()
				if err != nil {
					return err
				}
				return nil
			},
			wantErr: iterator.ErrNoSuchData,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//err := tc.list.Add(tc.index, tc.newVal)
			iter := tc.list.Iterator()
			// 这里处理error 很恶心
			ints := make([]int, 0, tc.list.Len())
			var i = 1
			var err error
			for iter.HasNext() {
				next, err := iter.Next()
				i++
				if i == tc.index {
					err = tc.fn(iter)
					if err != nil {
						return
					}
				}

				ints = append(ints, next)

			}
			assert.Equal(t, tc.wantErr, err)
			// 因为返回了 error，所以我们不用继续往下比较了
			if err != nil {
				return
			}
			assert.NoError(t, iter.Err())
			assert.Equal(t, tc.wantSlice, ints)
		})
	}
}

func TestArrayList_Iter_get(t *testing.T) {
	testCases := []struct {
		name      string
		list      *ArrayList[int]
		wantSlice []int
		wantErr   error
	}{
		{
			name:      "获取当前节点",
			list:      NewArrayListOf[int]([]int{1, 2, 3, 4, 5}),
			wantSlice: []int{1, 2, 3, 4, 5},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//err := tc.list.Add(tc.index, tc.newVal)
			iter := tc.list.Iterator()
			// 这里处理error 很恶心
			ints := make([]int, 0, tc.list.Len())
			var err error
			for iter.HasNext() {
				get, err := iter.Get()

				ints = append(ints, get)
				_, err = iter.Next()
				assert.NoError(t, err)

			}
			assert.Equal(t, tc.wantErr, err)
			// 因为返回了 error，所以我们不用继续往下比较了
			if err != nil {
				return
			}
			assert.NoError(t, iter.Err())
			assert.Equal(t, tc.wantSlice, ints)
		})
	}
}

func TestArrayList_Iter_mod(t *testing.T) {
	testCases := []struct {
		name    string
		list    *ArrayList[int]
		index   int
		fn      func(l *ArrayList[int], idx int) error
		wantErr error
	}{
		{
			name:  "在迭代的时候进行追加元素",
			list:  NewArrayListOf[int]([]int{1, 2, 3, 4, 5}),
			index: 3,
			fn: func(l *ArrayList[int], idx int) error {
				err := l.Append(0)
				return err
			},
			wantErr: iterator.ErrStructHasChange,
		},
		{
			name:  "在迭代的时候进行新增数组",
			list:  NewArrayListOf[int]([]int{1, 2, 3, 4, 5}),
			index: 3,
			fn: func(l *ArrayList[int], idx int) error {
				err := l.Add(idx, 0)
				return err
			},
			wantErr: iterator.ErrStructHasChange,
		},
		{
			name:  "在迭代的时候进行设置数组",
			list:  NewArrayListOf[int]([]int{1, 2, 3, 4, 5}),
			index: 3,
			fn: func(l *ArrayList[int], idx int) error {
				err := l.Set(idx, 0)
				return err
			},
			wantErr: iterator.ErrStructHasChange,
		},
		{
			name:  "在迭代的时候进行删除数组",
			list:  NewArrayListOf[int]([]int{1, 2, 3, 4, 5}),
			index: 3,
			fn: func(l *ArrayList[int], idx int) error {
				_, err := l.Delete(idx)
				return err
			},
			wantErr: iterator.ErrStructHasChange,
		},

		{
			name:  "在迭代的时候进行缩容",
			list:  NewArrayListOf[int]([]int{1, 2, 3, 4, 5}),
			index: 3,
			fn: func(l *ArrayList[int], idx int) error {
				l.shrink()
				return nil
			},
			wantErr: iterator.ErrStructHasChange,
		},

		{
			name: "在迭代的时候不进行修改",
			list: NewArrayListOf[int]([]int{1, 2, 3, 4, 5}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//err := tc.list.Add(tc.index, tc.newVal)
			iter := tc.list.Iterator()
			var i = 1
			// 这里处理error 很恶心
			var err error
			for iter.HasNext() {
				_, err = iter.Next()
				assert.NoError(t, err)
				i++
				if i == tc.index {
					err := tc.fn(tc.list, tc.index)
					assert.NoError(t, err)
				}
			}
			assert.Equal(t, tc.wantErr, iter.Err())
			// 因为返回了 error，所以我们不用继续往下比较了
			if err != nil {
				return
			}
		})
	}
}
