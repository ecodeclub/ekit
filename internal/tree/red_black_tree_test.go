package tree

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gotomicro/ekit"
	"github.com/stretchr/testify/assert"
)

func TestNewRedBlackTree(t *testing.T) {
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
			redBlackTree := NewRedBlackTree[int](compare())
			assert.Equal(t, tt.wantV, isRedBlackTree[int](redBlackTree.root))
		})
	}
}
func compare() ekit.Comparator[int] {
	return ekit.ComparatorRealNumber[int]
}

// isRedBlackTree 检测是否满足红黑树
func isRedBlackTree[T any](root *RBNode[T]) bool {
	// 检测节点是否黑色
	if !root.colorOf() {
		return false
	}
	// count 取最左树的黑色节点作为对照
	count := 0
	num := 0
	node := root
	for node != nil {
		if node.color {
			count++
		}
		node = node.leftOf()
	}
	return nodeCheck[T](root, count, num)
}

// nodeCheck 节点检测
// 1、是否有连续的红色节点
// 2、每条路径的黑色节点是否一致
func nodeCheck[T any](node *RBNode[T], count int, num int) bool {
	if node == nil {
		return true
	}
	if !node.colorOf() && !node.parent.colorOf() {
		fmt.Println("存在连续红色节点")
		return false
	}
	if node.colorOf() {
		num++
	}
	if node.leftOf() == nil && node.rightOf() == nil {
		if num != count {
			fmt.Println("黑色节点数量不一致")
			return false
		}
	}
	return nodeCheck(node.left, count, num) && nodeCheck(node.right, count, num)
}

func TestRedBlackTree_Add(t *testing.T) {
	isRedBlackTreeCase := []struct {
		name string
		node *RBNode[int]
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
			node: &RBNode[int]{left: nil, right: nil, color: Black},
			want: true,
		},
		//			 root(b)
		//			/
		//		   a(b)
		{
			name: "root-oneChild",
			node: &RBNode[int]{
				left: &RBNode[int]{
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
			node: &RBNode[int]{
				left: &RBNode[int]{
					right: nil,
					left:  nil,
					color: Red,
				},
				right: &RBNode[int]{
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
			node: &RBNode[int]{
				left: &RBNode[int]{
					right: &RBNode[int]{
						right: nil,
						left:  nil,
						color: Red,
					},
					left:  nil,
					color: Black,
				},
				right: &RBNode[int]{
					right: nil,
					left: &RBNode[int]{
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
		//			 root(b)
		//			/	    \
		//		   a(r)	     b(r)
		//		 /  \        /    \
		//      nil  c(r)    d(r)   nil
		//           / \     / \
		//          nil nil nil nil
		{
			name: "root-grandson",
			node: &RBNode[int]{
				parent: nil,
				key:    7,
				left: &RBNode[int]{
					key:   5,
					color: Black,
					left: &RBNode[int]{
						key:   4,
						color: Red,
					},
					right: &RBNode[int]{
						key:   6,
						color: Red,
					},
				},
				right: &RBNode[int]{
					key:   10,
					color: Red,
					left: &RBNode[int]{
						key:   9,
						color: Black,
						left: &RBNode[int]{
							key:   8,
							color: Red,
						},
					},
					right: &RBNode[int]{
						key:   12,
						color: Black,
						left: &RBNode[int]{
							key:   11,
							color: Red,
						},
					},
				},
				color: Black,
			},
			want: false,
		},
	}
	for _, tt := range isRedBlackTreeCase {
		t.Run(tt.name, func(t *testing.T) {
			res := isRedBlackTree[int](tt.node)
			assert.Equal(t, tt.want, res)

		})
	}
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
			wantErr: nil,
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
			redBlackTree := NewRedBlackTree[int](compare())
			for i := 0; i < len(tt.k); i++ {
				redBlackTree.Add(&RBNode[int]{
					key: tt.k[i],
				})
			}
			res := isRedBlackTree[int](redBlackTree.root)
			assert.Equal(t, tt.want, res)

		})
	}
}

func TestRedBlackTree_Delete(t *testing.T) {
	tcase := []struct {
		name   string
		delKey int
		key    []int
		want   bool
	}{
		{
			name:   "nil",
			delKey: 0,
			key:    nil,
			want:   true,
		},
		{
			name:   "node-empty",
			delKey: 0,
			key:    []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			want:   true,
		},
		{
			name:   "左右非空子节点,删除节点为黑色",
			delKey: 11,
			key:    []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			want:   true,
		},
		{
			name:   "左右只有一个非空子节点,删除节点为黑色",
			delKey: 11,
			key:    []int{4, 5, 6, 7, 8, 9, 11, 12},
			want:   true,
		},
		{
			name:   "左右均为空节点,删除节点为黑色",
			delKey: 12,
			key:    []int{4, 5, 6, 7, 8, 9, 12},
			want:   true,
		}, {
			name:   "左右非空子节点,删除节点为红色",
			delKey: 5,
			key:    []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			want:   true,
		},
		// 此状态无法构造出正确的红黑树
		// {
		//	name:   "左右只有一个非空子节点,删除节点为红色",
		//	delKey: 5,
		//	key:    []int{4, 5, 6, 7, 8, 9, 11, 12},
		//	want:   true,
		// },
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			rbTree := NewRedBlackTree[int](compare())
			for i := 0; i < len(tt.key); i++ {
				node := &RBNode[int]{
					key: tt.key[i],
				}
				rbTree.Add(node)
			}
			assert.Equal(t, tt.want, isRedBlackTree[int](rbTree.root))
			rbTree.Delete(tt.delKey)
			assert.Equal(t, tt.want, isRedBlackTree[int](rbTree.root))

		})
	}
}

func TestRedBlackTree_Find(t *testing.T) {
	tcase := []struct {
		name      string
		findKey   int
		key       []int
		wantKey   int
		wantError error
	}{
		{
			name:      "nil",
			findKey:   0,
			key:       nil,
			wantError: errors.New("未找到0节点"),
		},
		{
			name:      "node-empty",
			findKey:   0,
			key:       []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			wantError: errors.New("未找到0节点"),
		},
		{
			name:    "find",
			findKey: 11,
			key:     []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			wantKey: 11,
		}, {
			name:    "find",
			findKey: 12,
			key:     []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			wantKey: 12,
		}, {
			name:    "find",
			findKey: 7,
			key:     []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			wantKey: 7,
		},
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			rbTree := NewRedBlackTree[int](compare())
			for i := 0; i < len(tt.key); i++ {
				node := &RBNode[int]{
					key: tt.key[i],
				}
				rbTree.Add(node)
			}
			assert.Equal(t, true, isRedBlackTree[int](rbTree.root))
			findNode := rbTree.Find(tt.findKey)
			if findNode == nil {
				assert.Equal(t, tt.wantError, errors.New("未找到0节点"))
			} else {
				assert.Equal(t, tt.wantKey, findNode.key)
			}
		})
	}
}

func TestRedBlackTree_addNode(t *testing.T) {
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
			wantErr: nil,
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
			redBlackTree := NewRedBlackTree[int](compare())
			for i := 0; i < len(tt.k); i++ {
				redBlackTree.addNode(&RBNode[int]{
					key: tt.k[i],
				})
			}
			res := isRedBlackTree[int](redBlackTree.root)
			assert.Equal(t, tt.want, res)

		})
	}
}

func TestRedBlackTree_deleteNode(t *testing.T) {
	tcase := []struct {
		name      string
		delKey    int
		key       []int
		want      bool
		wantError error
	}{
		{
			name:      "nil",
			delKey:    0,
			key:       nil,
			wantError: errors.New("未找到节点0"),
		},
		{
			name:      "node-empty",
			delKey:    0,
			key:       []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			wantError: errors.New("未找到节点0"),
		},
		{
			name:   "左右非空子节点,删除节点为黑色",
			delKey: 11,
			key:    []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			want:   true,
		},
		{
			name:   "左右只有一个非空子节点,删除节点为黑色",
			delKey: 11,
			key:    []int{4, 5, 6, 7, 8, 9, 11, 12},
			want:   true,
		},
		{
			name:   "左右均为空节点,删除节点为黑色",
			delKey: 12,
			key:    []int{4, 5, 6, 7, 8, 9, 12},
			want:   true,
		}, {
			name:   "左右非空子节点,删除节点为红色",
			delKey: 5,
			key:    []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			want:   true,
		},
		// 此状态无法构造出正确的红黑树
		// {
		//	name:   "左右只有一个非空子节点,删除节点为红色",
		//	delKey: 5,
		//	key:    []int{4, 5, 6, 7, 8, 9, 11, 12},
		//	want:   true,
		// },
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			rbTree := NewRedBlackTree[int](compare())
			for i := 0; i < len(tt.key); i++ {
				node := &RBNode[int]{
					key: tt.key[i],
				}
				rbTree.Add(node)
			}
			delNode := rbTree.getRBNode(tt.delKey)
			if delNode == nil {
				assert.Equal(t, tt.wantError, errors.New("未找到节点0"))
			} else {
				rbTree.deleteNode(delNode)
				assert.Equal(t, tt.want, isRedBlackTree[int](rbTree.root))
			}
		})
	}
}

func TestRedBlackTree_fixAfterAdd(t *testing.T) {}

func TestRedBlackTree_fixAfterDelete(t *testing.T) {}

func TestRedBlackTree_getTreeNode(t *testing.T) {
	tcase := []struct {
		name      string
		findKey   int
		key       []int
		wantKey   int
		wantError error
	}{
		{
			name:      "nil",
			findKey:   0,
			key:       nil,
			wantError: errors.New("未找到0节点"),
		},
		{
			name:      "node-empty",
			findKey:   0,
			key:       []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			wantError: errors.New("未找到0节点"),
		},
		{
			name:    "find",
			findKey: 11,
			key:     []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			wantKey: 11,
		}, {
			name:    "find",
			findKey: 12,
			key:     []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			wantKey: 12,
		}, {
			name:    "find",
			findKey: 7,
			key:     []int{4, 5, 6, 7, 8, 9, 10, 11, 12},
			wantKey: 7,
		},
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			rbTree := NewRedBlackTree[int](compare())
			for i := 0; i < len(tt.key); i++ {
				node := &RBNode[int]{
					key: tt.key[i],
				}
				rbTree.Add(node)
			}
			assert.Equal(t, true, isRedBlackTree[int](rbTree.root))
			findNode := rbTree.getRBNode(tt.findKey)
			if findNode == nil {
				assert.Equal(t, tt.wantError, errors.New("未找到0节点"))
			} else {
				assert.Equal(t, tt.wantKey, findNode.key)
			}
		})
	}
}

func TestRedBlackTree_rotateLeft(t *testing.T) {

}

func TestRedBlackTree_rotateRight(t *testing.T) {}

func TestRedBlackTree_successor(t *testing.T) {}

func Test_treeNode_colorOf(t *testing.T) {}

func Test_treeNode_leftOf(t *testing.T) {}

func Test_treeNode_parentOf(t *testing.T) {}

func Test_treeNode_rightOf(t *testing.T) {}

func Test_treeNode_setColor(t *testing.T) {}
