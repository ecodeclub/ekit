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

package tree

import (
	"errors"
	"testing"

	"github.com/gotomicro/ekit"
	"github.com/stretchr/testify/assert"
)

func TestNewRBTree(t *testing.T) {
	tests := []struct {
		name    string
		compare ekit.Comparator[int]
		wantV   bool
	}{
		{
			name:    "int",
			compare: compare(),
			wantV:   true,
		},
		{
			name:    "nil",
			compare: nil,
			wantV:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redBlackTree := NewRBTree[int, string](compare())
			assert.Equal(t, tt.wantV, IsRedBlackTree[int, string](redBlackTree.root))
		})
	}
}

func compare() ekit.Comparator[int] {
	return ekit.ComparatorRealNumber[int]
}

func TestRBTree_Add(t *testing.T) {
	IsRedBlackTreeCase := []struct {
		name string
		node *RBNode[int, string]
		want bool
	}{
		{
			name: "nil",
			node: nil,
			want: true,
		},
		{
			name: "node-nil",
			node: nil,
			want: true,
		},
		{
			name: "root",
			node: &RBNode[int, string]{left: nil, right: nil, color: Black},
			want: true,
		},
		//			 root(b)
		//			/
		//		   a(b)
		{
			name: "root-oneChild",
			node: &RBNode[int, string]{
				left: &RBNode[int, string]{
					right: nil,
					left:  nil,
					color: Red,
				},
				right: nil,
				color: Black,
			},
			want: true,
		},
		//			 root(b)
		//			/	    \
		//		   a(r)	     b(b)
		{
			name: "root-twoChild",
			node: &RBNode[int, string]{
				left: &RBNode[int, string]{
					right: nil,
					left:  nil,
					color: Red,
				},
				right: &RBNode[int, string]{
					right: nil,
					left:  nil,
					color: Black,
				},
				color: Black,
			},
			want: false,
		},
		//			 root(b)
		//			/	    \
		//		   a(b)	     b(b)
		//		 /  \        /    \
		//      nil  c(r)    d(r)   nil
		//           / \     / \
		//          nil nil nil nil
		{
			name: "blackNodeNotSame",
			node: &RBNode[int, string]{
				left: &RBNode[int, string]{
					right: &RBNode[int, string]{
						right: nil,
						left:  nil,
						color: Red,
					},
					left:  nil,
					color: Black,
				},
				right: &RBNode[int, string]{
					right: nil,
					left: &RBNode[int, string]{
						right: nil,
						left:  nil,
						color: Red,
					},
					color: Black,
				},
				color: Black,
			},
			want: true,
		},
		{
			name: "root-grandson",
			node: &RBNode[int, string]{
				parent: nil,
				Key:    7,
				left: &RBNode[int, string]{
					Key:   5,
					color: Black,
					left: &RBNode[int, string]{
						Key:   4,
						color: Red,
					},
					right: &RBNode[int, string]{
						Key:   6,
						color: Red,
					},
				},
				right: &RBNode[int, string]{
					Key:   10,
					color: Red,
					left: &RBNode[int, string]{
						Key:   9,
						color: Black,
						left: &RBNode[int, string]{
							Key:   8,
							color: Red,
						},
					},
					right: &RBNode[int, string]{
						Key:   12,
						color: Black,
						left: &RBNode[int, string]{
							Key:   11,
							color: Red,
						},
					},
				},
				color: Black,
			},
			want: true,
		},
	}
	for _, tt := range IsRedBlackTreeCase {
		t.Run(tt.name, func(t *testing.T) {
			res := IsRedBlackTree[int](tt.node)
			assert.Equal(t, tt.want, res)

		})
	}
	tests := []struct {
		name    string
		k       []int
		want    bool
		wantErr error
		size    int
		wantKey int
	}{
		{
			name: "nil",
			k:    nil,
			want: true,
			size: 0,
		},
		{
			name: "one",
			k:    []int{1},
			want: true,
			size: 1,
		},
		{
			name:    "one",
			k:       []int{1, 2},
			want:    true,
			size:    2,
			wantKey: 1,
		},
		{
			name:    "normal",
			k:       []int{1, 2, 3, 4},
			want:    true,
			size:    4,
			wantKey: 3,
		},
		{
			name:    "same",
			k:       []int{0, 0, 1, 2, 2, 3},
			want:    true,
			size:    0,
			wantErr: errors.New("ekit: RBTree不能添加重复节点Key"),
		},
		{
			name:    "disorder",
			k:       []int{1, 2, 0, 3, 5, 4},
			want:    true,
			wantErr: nil,
			size:    6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redBlackTree := NewRBTree[int, string](compare())
			for i := 0; i < len(tt.k); i++ {
				err := redBlackTree.Add(&RBNode[int, string]{
					Key: tt.k[i],
				})
				if err != nil {
					assert.Equal(t, tt.wantErr, err)
					return
				}
			}
			res := IsRedBlackTree[int, string](redBlackTree.root)
			assert.Equal(t, tt.want, res)
			assert.Equal(t, tt.size, redBlackTree.Size())
		})
	}
}

func TestRBTree_Delete(t *testing.T) {
	tcase := []struct {
		name   string
		delKey int
		Key    []int
		want   bool
		size   int
	}{
		{
			name:   "nil",
			delKey: 0,
			Key:    nil,
			want:   true,
			size:   0,
		},
		{
			name:   "node-empty",
			delKey: 0,
			Key:    []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			want:   true,
			size:   9,
		},
		{
			name:   "左右非空子节点,删除节点为黑色",
			delKey: 11,
			Key:    []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			want:   true,
			size:   8,
		},
		{
			name:   "左右只有一个非空子节点,删除节点为黑色",
			delKey: 11,
			Key:    []int{4, 5, 6, 7, 8, 9, 11, 12},
			want:   true,
			size:   7,
		},
		{
			name:   "左右均为空节点,删除节点为黑色",
			delKey: 12,
			Key:    []int{4, 5, 6, 7, 8, 9, 12},
			want:   true,
			size:   6,
		}, {
			name:   "左右非空子节点,删除节点为红色",
			delKey: 5,
			Key:    []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			want:   true,
			size:   8,
		},
		// 此状态无法构造出正确的红黑树
		// {
		//	name:   "左右只有一个非空子节点,删除节点为红色",
		//	delKey: 5,
		//	Key:    []int{4, 5, 6, 7, 8, 9, 11, 12},
		//	want:   true,
		// },
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			rbTree := NewRBTree[int, string](compare())
			for i := 0; i < len(tt.Key); i++ {
				node := &RBNode[int, string]{
					Key: tt.Key[i],
				}
				err := rbTree.Add(node)
				if err != nil {
					panic(err)
				}
			}
			assert.Equal(t, tt.want, IsRedBlackTree[int](rbTree.root))
			rbTree.Delete(tt.delKey)
			assert.Equal(t, tt.want, IsRedBlackTree[int](rbTree.root))
			assert.Equal(t, tt.size, rbTree.Size())
		})
	}
}

func TestRBTree_Find(t *testing.T) {
	tcase := []struct {
		name      string
		findKey   int
		Key       []int
		wantKey   int
		wantError error
	}{
		{
			name:      "nil",
			findKey:   0,
			Key:       nil,
			wantError: errors.New("未找到0节点"),
		},
		{
			name:      "node-empty",
			findKey:   0,
			Key:       []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			wantError: errors.New("未找到0节点"),
		},
		{
			name:    "find",
			findKey: 11,
			Key:     []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			wantKey: 11,
		}, {
			name:    "find",
			findKey: 12,
			Key:     []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			wantKey: 12,
		}, {
			name:    "find",
			findKey: 7,
			Key:     []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			wantKey: 7,
		},
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			rbTree := NewRBTree[int, string](compare())
			for i := 0; i < len(tt.Key); i++ {
				node := &RBNode[int, string]{
					Key: tt.Key[i],
				}
				err := rbTree.Add(node)
				if err != nil {
					panic(err)
				}
			}
			assert.Equal(t, true, IsRedBlackTree[int](rbTree.root))
			findNode := rbTree.Find(tt.findKey)
			if findNode == nil {
				assert.Equal(t, tt.wantError, errors.New("未找到0节点"))
			} else {
				assert.Equal(t, tt.wantKey, findNode.Key)
			}
		})
	}
}

func TestRBTree_addNode(t *testing.T) {
	tests := []struct {
		name    string
		k       []int
		want    bool
		wantErr error
	}{
		{
			name: "nil",
			k:    nil,
			want: true,
		},
		{
			name: "case1",
			k:    []int{1, 2, 3, 4},
			want: true,
		},
		{
			name:    "same",
			k:       []int{0, 0, 1, 2, 2, 3},
			want:    true,
			wantErr: errors.New("ekit: RBTree不能添加重复节点Key"),
		},
		{
			name:    "disorder",
			k:       []int{1, 2, 0, 3, 5, 4},
			want:    true,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redBlackTree := NewRBTree[int, string](compare())
			for i := 0; i < len(tt.k); i++ {
				err := redBlackTree.addNode(&RBNode[int, string]{
					Key: tt.k[i],
				})
				if err != nil {
					assert.Equal(t, tt.wantErr, err)
				}
			}
			res := IsRedBlackTree[int](redBlackTree.root)
			assert.Equal(t, tt.want, res)

		})
	}
}

func TestRBTree_deleteNode(t *testing.T) {
	tcase := []struct {
		name      string
		delKey    int
		Key       []int
		want      bool
		wantError error
	}{
		{
			name:      "nil",
			delKey:    0,
			Key:       nil,
			wantError: errors.New("未找到节点0"),
		},
		{
			name:      "node-empty",
			delKey:    0,
			Key:       []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			wantError: errors.New("未找到节点0"),
		},
		{
			name:   "左右非空子节点,删除节点为黑色",
			delKey: 11,
			Key:    []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			want:   true,
		},
		{
			name:   "左右只有一个非空子节点,删除节点为黑色",
			delKey: 11,
			Key:    []int{4, 5, 6, 7, 8, 9, 11, 12},
			want:   true,
		},
		{
			name:   "左右均为空节点,删除节点为黑色",
			delKey: 12,
			Key:    []int{4, 5, 6, 7, 8, 9, 12},
			want:   true,
		}, {
			name:   "左右非空子节点,删除节点为红色",
			delKey: 5,
			Key:    []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			want:   true,
		},
		// 此状态无法构造出正确的红黑树
		// {
		//	name:   "左右只有一个非空子节点,删除节点为红色",
		//	delKey: 5,
		//	Key:    []int{4, 5, 6, 7, 8, 9, 11, 12},
		//	want:   true,
		// },
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			rbTree := NewRBTree[int, string](compare())
			for i := 0; i < len(tt.Key); i++ {
				node := &RBNode[int, string]{
					Key: tt.Key[i],
				}
				err := rbTree.Add(node)
				if err != nil {
					panic(err)
				}

			}
			delNode := rbTree.getRBNode(tt.delKey)
			if delNode == nil {
				assert.Equal(t, tt.wantError, errors.New("未找到节点0"))
			} else {
				rbTree.deleteNode(delNode)
				assert.Equal(t, tt.want, IsRedBlackTree[int](rbTree.root))
			}
		})
	}
}

func TestRBTree_getTreeNode(t *testing.T) {
	tcase := []struct {
		name      string
		findKey   int
		Key       []int
		wantKey   int
		wantError error
	}{
		{
			name:      "nil",
			findKey:   0,
			Key:       nil,
			wantError: errors.New("未找到0节点"),
		},
		{
			name:      "node-empty",
			findKey:   0,
			Key:       []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			wantError: errors.New("未找到0节点"),
		},
		{
			name:    "find",
			findKey: 11,
			Key:     []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			wantKey: 11,
		}, {
			name:    "find",
			findKey: 12,
			Key:     []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			wantKey: 12,
		}, {
			name:    "find",
			findKey: 7,
			Key:     []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			wantKey: 7,
		},
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			rbTree := NewRBTree[int, string](compare())
			for i := 0; i < len(tt.Key); i++ {
				node := &RBNode[int, string]{
					Key: tt.Key[i],
				}
				err := rbTree.Add(node)
				if err != nil {
					panic(err)
				}
			}
			assert.Equal(t, true, IsRedBlackTree[int](rbTree.root))
			findNode := rbTree.getRBNode(tt.findKey)
			if findNode == nil {
				assert.Equal(t, tt.wantError, errors.New("未找到0节点"))
			} else {
				assert.Equal(t, tt.wantKey, findNode.Key)
			}
		})
	}
}

func TestRBTree_rotateLeft(t *testing.T) {
	tcase := []struct {
		name        string
		key         []int
		wantKey     int
		wantLeftKey int
		isRBTree    bool
	}{
		{
			name:     "nod-nil",
			key:      nil,
			isRBTree: true,
		},
		{
			name:     "only-root",
			key:      []int{1},
			wantKey:  1,
			isRBTree: true,
		},
		{
			name:        "right node have two child nods",
			key:         []int{1, 2, 3, 4, 5},
			wantKey:     4,
			wantLeftKey: 2,
			isRBTree:    false,
		},
		{
			name:        "right node have a child nod",
			key:         []int{1, 2, 3, 4},
			wantKey:     3,
			wantLeftKey: 2,
			isRBTree:    false,
		},
		{
			name:        "right node have nil child nod",
			key:         []int{1, 2, 3},
			wantKey:     3,
			wantLeftKey: 2,
			isRBTree:    false,
		},
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			rbTree := NewRBTree[int, string](compare())
			for i := 0; i < len(tt.key); i++ {
				node := &RBNode[int, string]{
					Key: tt.key[i],
				}
				err := rbTree.Add(node)
				if err != nil {
					panic(err)
				}
			}
			assert.Equal(t, true, IsRedBlackTree[int](rbTree.root))
			rbTree.rotateLeft(rbTree.root)
			assert.Equal(t, tt.isRBTree, IsRedBlackTree[int](rbTree.root))
			if rbTree.root != nil {
				assert.Equal(t, tt.wantKey, rbTree.Root().Key)
				if rbTree.Root().getLeft() != nil {
					assert.Equal(t, tt.wantLeftKey, rbTree.Root().getLeft().Key)
				}
			}
		})
	}
}

func TestRBTree_rotateRight(t *testing.T) {
	tcase := []struct {
		name         string
		key          []int
		wantKey      int
		wantRightKey int
		isRBTree     bool
	}{
		{
			name:     "nod-nil",
			key:      nil,
			isRBTree: true,
		},
		{
			name:     "only-root",
			key:      []int{1},
			wantKey:  1,
			isRBTree: true,
		},
		{
			name:         "right node have two child nods",
			key:          []int{4, 5, 3, 2, 1},
			wantKey:      2,
			wantRightKey: 4,
			isRBTree:     false,
		},
		{
			name:         "right node have a child nod",
			key:          []int{4, 5, 3, 2},
			wantKey:      3,
			wantRightKey: 4,
			isRBTree:     false,
		},
		{
			name:         "right node have nil child nod",
			key:          []int{4, 5, 3},
			wantKey:      3,
			wantRightKey: 4,
			isRBTree:     false,
		},
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			rbTree := NewRBTree[int, string](compare())
			for i := 0; i < len(tt.key); i++ {
				node := &RBNode[int, string]{
					Key: tt.key[i],
				}
				err := rbTree.Add(node)
				if err != nil {
					panic(err)
				}
			}
			assert.Equal(t, true, IsRedBlackTree[int](rbTree.root))
			rbTree.rotateRight(rbTree.root)
			assert.Equal(t, tt.isRBTree, IsRedBlackTree[int](rbTree.root))
			if rbTree.root != nil {
				assert.Equal(t, tt.wantKey, rbTree.Root().Key)
				if rbTree.Root().getRight() != nil {
					assert.Equal(t, tt.wantRightKey, rbTree.Root().getRight().Key)
				}
			}
		})
	}
}

func TestRBNode_getColor(t *testing.T) {
	tcase := []struct {
		name      string
		node      *RBNode[int, int]
		wantColor bool
	}{
		{
			name:      "nod-nil",
			node:      nil,
			wantColor: true,
		},
		{
			name:      "new node",
			node:      NewRBNode[int, int](1, 1),
			wantColor: false,
		},
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantColor, tt.node.getColor())
		})
	}
}

func TestRBNode_getLeft(t *testing.T) {
	tcase := []struct {
		name     string
		node     *RBNode[int, int]
		wantNode *RBNode[int, int]
	}{
		{
			name:     "nod-nil",
			node:     nil,
			wantNode: nil,
		},
		{
			name:     "new node",
			node:     NewRBNode[int, int](1, 1),
			wantNode: nil,
		},
		{
			name: "new node have left-child",
			node: &RBNode[int, int]{
				Key: 2,
				left: &RBNode[int, int]{
					Key: 1,
				},
			},
			wantNode: &RBNode[int, int]{
				Key: 1,
			},
		},
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantNode, tt.node.getLeft())
		})
	}
}

func TestRBNode_getRight(t *testing.T) {
	tcase := []struct {
		name     string
		node     *RBNode[int, int]
		wantNode *RBNode[int, int]
	}{
		{
			name:     "nod-nil",
			node:     nil,
			wantNode: nil,
		},
		{
			name:     "new node",
			node:     NewRBNode[int, int](1, 1),
			wantNode: nil,
		},
		{
			name: "new node have right-child",
			node: &RBNode[int, int]{
				Key: 1,
				right: &RBNode[int, int]{
					Key: 2,
				},
			},
			wantNode: &RBNode[int, int]{
				Key: 2,
			},
		},
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantNode, tt.node.getRight())
		})
	}
}

func TestRBNode_getParent(t *testing.T) {
	tcase := []struct {
		name     string
		node     *RBNode[int, int]
		wantNode *RBNode[int, int]
	}{
		{
			name:     "nod-nil",
			node:     nil,
			wantNode: nil,
		},
		{
			name:     "new node",
			node:     NewRBNode[int, int](1, 1),
			wantNode: nil,
		},
		{
			name: "new node have parent",
			node: &RBNode[int, int]{
				Key: 2,
				parent: &RBNode[int, int]{
					Key: 3,
				},
			},
			wantNode: &RBNode[int, int]{
				Key: 3,
			},
		},
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantNode, tt.node.getParent())
		})
	}
}

func TestRBNode_setColor(t *testing.T) {
	tcase := []struct {
		name      string
		node      *RBNode[int, int]
		color     bool
		wantColor bool
	}{
		{
			name:      "nod-nil",
			node:      nil,
			color:     false,
			wantColor: Black,
		},
		{
			name:      "new node",
			node:      NewRBNode[int, int](1, 1),
			color:     true,
			wantColor: Black,
		},
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			tt.node.setColor(tt.color)
			assert.Equal(t, tt.wantColor, tt.node.getColor())
		})
	}
}

func TestNewRBNode(t *testing.T) {
	tcase := []struct {
		name     string
		key      int
		value    int
		wantNode *RBNode[int, int]
	}{
		{
			name:  "new node",
			key:   1,
			value: 1,
			wantNode: &RBNode[int, int]{
				Key:    1,
				Value:  1,
				left:   nil,
				right:  nil,
				parent: nil,
				color:  Red,
			},
		},
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			node := NewRBNode[int, int](tt.key, tt.value)
			assert.Equal(t, tt.wantNode, node)
		})
	}
}

func TestRBNode_Left(t *testing.T) {
	tcase := []struct {
		name     string
		node     *RBNode[int, int]
		wantNode *RBNode[int, int]
	}{
		{
			name:     "nod-nil",
			node:     nil,
			wantNode: nil,
		},
		{
			name:     "new node",
			node:     NewRBNode[int, int](1, 1),
			wantNode: nil,
		},
		{
			name: "new node have left-child",
			node: &RBNode[int, int]{
				Key: 2,
				left: &RBNode[int, int]{
					Key: 1,
				},
			},
			wantNode: &RBNode[int, int]{
				Key: 1,
			},
		},
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantNode, tt.node.Left())
		})
	}
}

func TestRBNode_Right(t *testing.T) {
	tcase := []struct {
		name     string
		node     *RBNode[int, int]
		wantNode *RBNode[int, int]
	}{
		{
			name:     "nod-nil",
			node:     nil,
			wantNode: nil,
		},
		{
			name:     "new node",
			node:     NewRBNode[int, int](1, 1),
			wantNode: nil,
		},
		{
			name: "new node have right-child",
			node: &RBNode[int, int]{
				Key: 1,
				right: &RBNode[int, int]{
					Key: 2,
				},
			},
			wantNode: &RBNode[int, int]{
				Key: 2,
			},
		},
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantNode, tt.node.Right())
		})
	}
}

func TestRBNode_Parent(t *testing.T) {
	tcase := []struct {
		name     string
		node     *RBNode[int, int]
		wantNode *RBNode[int, int]
	}{
		{
			name:     "nod-nil",
			node:     nil,
			wantNode: nil,
		},
		{
			name:     "new node",
			node:     NewRBNode[int, int](1, 1),
			wantNode: nil,
		},
		{
			name: "new node have parent",
			node: &RBNode[int, int]{
				Key: 2,
				parent: &RBNode[int, int]{
					Key: 3,
				},
			},
			wantNode: &RBNode[int, int]{
				Key: 3,
			},
		},
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantNode, tt.node.Parent())
		})
	}
}

func TestRBNode_getBrother(t *testing.T) {
	tests := []struct {
		name    string
		k       []int
		nodeKye int
		want    int
	}{
		{
			name: "nil",
			k:    nil,
		},
		{
			name:    "no-brother",
			nodeKye: 1,
			k:       []int{1},
		},
		{
			name:    "no-brother",
			nodeKye: 1,
			k:       []int{1, 2},
		},
		{
			name:    "have brother",
			k:       []int{1, 2, 3, 4},
			nodeKye: 1,
			want:    3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redBlackTree := NewRBTree[int, string](compare())
			for i := 0; i < len(tt.k); i++ {
				err := redBlackTree.Add(&RBNode[int, string]{
					Key: tt.k[i],
				})
				if err != nil {
					panic(err)
				}
			}
			tagNode := redBlackTree.getRBNode(tt.nodeKye)
			brNode := tagNode.getBrother()
			if brNode == nil {
				return
			}
			assert.Equal(t, tt.want, brNode.Key)

		})
	}
}

func TestRBNode_getGrandParent(t *testing.T) {
	tests := []struct {
		name    string
		k       []int
		nodeKye int
		want    int
	}{
		{
			name: "nil",
			k:    nil,
		},
		{
			name:    "no-grandpa",
			nodeKye: 1,
			k:       []int{1},
		},
		{
			name:    "no-grandpa",
			nodeKye: 1,
			k:       []int{1, 2},
		},
		{
			name:    "have grandpa",
			k:       []int{1, 2, 3, 4},
			nodeKye: 4,
			want:    2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redBlackTree := NewRBTree[int, string](compare())
			for i := 0; i < len(tt.k); i++ {
				err := redBlackTree.Add(&RBNode[int, string]{
					Key: tt.k[i],
				})
				if err != nil {
					panic(err)
				}
			}
			tagNode := redBlackTree.getRBNode(tt.nodeKye)
			brNode := tagNode.getGrandParent()
			if brNode == nil {
				return
			}
			assert.Equal(t, tt.want, brNode.Key)

		})
	}
}

func TestRBNode_getUncle(t *testing.T) {
	tests := []struct {
		name    string
		k       []int
		nodeKye int
		want    int
	}{
		{
			name: "nil",
			k:    nil,
		},
		{
			name:    "no-uncle",
			nodeKye: 1,
			k:       []int{1},
		},
		{
			name:    "no-uncle",
			nodeKye: 1,
			k:       []int{1, 2},
		},
		{
			name:    "have uncle",
			k:       []int{1, 2, 3, 4},
			nodeKye: 4,
			want:    1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redBlackTree := NewRBTree[int, string](compare())
			for i := 0; i < len(tt.k); i++ {
				err := redBlackTree.Add(&RBNode[int, string]{
					Key: tt.k[i],
				})
				if err != nil {
					panic(err)
				}
			}
			tagNode := redBlackTree.getRBNode(tt.nodeKye)
			brNode := tagNode.getUncle()
			if brNode == nil {
				return
			}
			assert.Equal(t, tt.want, brNode.Key)

		})
	}
}

func TestRBTree_Root(t *testing.T) {
	tests := []struct {
		name    string
		rbTree  *RBTree[int, int]
		k       []int
		rootKey int
		size    int
	}{
		{
			name:   "nil",
			rbTree: nil,
			size:   0,
		},
		{
			name:   "new RBTree",
			rbTree: NewRBTree[int, int](compare()),
			size:   0,
		},
		{
			name:    "one",
			k:       []int{1},
			rbTree:  NewRBTree[int, int](compare()),
			rootKey: 1,
			size:    1,
		},
		{
			name:    "case1",
			k:       []int{1, 2, 3, 4},
			rootKey: 2,
			size:    4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redBlackTree := NewRBTree[int, string](compare())
			for i := 0; i < len(tt.k); i++ {
				err := redBlackTree.Add(&RBNode[int, string]{
					Key: tt.k[i],
				})
				if err != nil {
					panic(err)
				}
			}
			assert.Equal(t, tt.size, redBlackTree.Size())
			root := redBlackTree.Root()
			if root == nil {
				return
			}
			assert.Equal(t, tt.rootKey, redBlackTree.Root().Key)

		})
	}
}

func Test_nodeCheck(t *testing.T) {
	IsRedBlackTreeCase := []struct {
		name  string
		node  *RBNode[int, string]
		count int
		want  bool
	}{
		{
			name:  "nil",
			node:  nil,
			count: 0,
			want:  true,
		},
		{
			name:  "node-nil",
			node:  nil,
			count: 0,
			want:  true,
		},
		{
			name:  "root",
			node:  &RBNode[int, string]{left: nil, right: nil, color: Black},
			count: 1,
			want:  true,
		},
		//			 root(b)
		//			/
		//		   a(r)
		{
			name: "root-oneChild",
			node: &RBNode[int, string]{
				left: &RBNode[int, string]{
					right: nil,
					left:  nil,
					color: Red,
				},
				right: nil,
				color: Black,
			},
			count: 1,
			want:  true,
		},
		//			 root(b)
		//			/	    \
		//		   a(r)	     b(b)
		{
			name: "root-twoChild",
			node: &RBNode[int, string]{
				left: &RBNode[int, string]{
					right: nil,
					left:  nil,
					color: Red,
				},
				right: &RBNode[int, string]{
					right: nil,
					left:  nil,
					color: Black,
				},
				color: Black,
			},
			count: 1,
			want:  false,
		},
		//			 root(b)
		//			/	    \
		//		   a(b)	     b(b)
		//		 /  \        /    \
		//      nil  c(r)    d(r)   nil
		//           / \     / \
		//          nil nil nil nil
		{
			name: "blackNodeNotSame",
			node: &RBNode[int, string]{
				left: &RBNode[int, string]{
					right: &RBNode[int, string]{
						right: nil,
						left:  nil,
						color: Red,
					},
					left:  nil,
					color: Black,
				},
				right: &RBNode[int, string]{
					right: nil,
					left: &RBNode[int, string]{
						right: nil,
						left:  nil,
						color: Red,
					},
					color: Black,
				},
				color: Black,
			},
			count: 2,
			want:  true,
		},
		{
			name: "root-grandson",
			node: &RBNode[int, string]{
				parent: nil,
				Key:    7,
				left: &RBNode[int, string]{
					Key:   5,
					color: Black,
					left: &RBNode[int, string]{
						Key:   4,
						color: Red,
					},
					right: &RBNode[int, string]{
						Key:   6,
						color: Red,
					},
				},
				right: &RBNode[int, string]{
					Key:   10,
					color: Red,
					left: &RBNode[int, string]{
						Key:   9,
						color: Black,
						left: &RBNode[int, string]{
							Key:   8,
							color: Red,
						},
					},
					right: &RBNode[int, string]{
						Key:   12,
						color: Black,
						left: &RBNode[int, string]{
							Key:   11,
							color: Red,
						},
					},
				},
				color: Black,
			},
			count: 2,
			want:  true,
		},
	}
	for _, tt := range IsRedBlackTreeCase {
		t.Run(tt.name, func(t *testing.T) {
			num := 0
			res := nodeCheck[int](tt.node, tt.count, num)
			assert.Equal(t, tt.want, res)

		})
	}
}

func TestRBNode_SetValue(t *testing.T) {
	tcase := []struct {
		name     string
		node     *RBNode[int, int]
		value    int
		wantNode *RBNode[int, int]
	}{
		{
			name:     "nil",
			node:     nil,
			value:    1,
			wantNode: nil,
		},
		{
			name: "new node",
			node: &RBNode[int, int]{
				Key:    1,
				Value:  0,
				left:   nil,
				right:  nil,
				parent: nil,
				color:  Red,
			},
			value: 1,
			wantNode: &RBNode[int, int]{
				Key:    1,
				Value:  1,
				left:   nil,
				right:  nil,
				parent: nil,
				color:  Red,
			},
		},
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			tt.node.SetValue(tt.value)
			assert.Equal(t, tt.wantNode, tt.node)
		})
	}
}

func TestRBTree_getSuccessor(t *testing.T) {
	tests := []struct {
		name      string
		k         []int
		successor int
		wantKey   int
	}{
		{
			name:      "nil-successor",
			k:         nil,
			successor: 8,
		},
		{
			name:      "have no successor",
			k:         []int{2},
			successor: 2,
		},
		{
			name:      "have right successor",
			k:         []int{5, 4, 6, 3, 2},
			successor: 3,
			wantKey:   4,
		},
		{
			name:      "have no-right successor",
			k:         []int{5, 4, 6, 3, 2},
			successor: 4,
			wantKey:   5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redBlackTree := NewRBTree[int, string](compare())
			for i := 0; i < len(tt.k); i++ {
				err := redBlackTree.Add(&RBNode[int, string]{
					Key: tt.k[i],
				})
				if err != nil {
					return
				}
			}
			tagNode := redBlackTree.getRBNode(tt.successor)
			successorNode := redBlackTree.getSuccessor(tagNode)
			if successorNode == nil {
				return
			}
			assert.Equal(t, tt.wantKey, successorNode.Key)
		})
	}
}
