package tree

import (
	"fmt"
	"github.com/gotomicro/ekit"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRedBlackTree(t *testing.T) {
	tests := []struct {
		name    string
		compare ekit.Comparator[int]
		wantV   *RedBlackTree[int]
	}{
		{
			name:    "int",
			compare: compare(),
			wantV: &RedBlackTree[int]{
				compare: compare(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redBlackTree := NewRedBlackTree[int](compare())
			assert.ElementsMatch(t, tt.wantV, redBlackTree)
		})
	}
}
func compare() ekit.Comparator[int] {
	return ekit.ComparatorRealNumber[int]
}

// isRedBlackTree 检测是否满足红黑树
func isRedBlackTree[T any](root *TreeNode[T]) bool {
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
func nodeCheck[T any](node *TreeNode[T], count int, num int) bool {
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
		node *TreeNode[int]
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
			node: &TreeNode[int]{left: nil, right: nil, color: Black},
			want: true,
		},
		//			 root(b)
		//			/
		//		   a(b)
		{
			name: "root-oneChild",
			node: &TreeNode[int]{
				left: &TreeNode[int]{
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
			node: &TreeNode[int]{
				left: &TreeNode[int]{
					right: nil,
					left:  nil,
					color: Red,
				},
				right: &TreeNode[int]{
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
			node: &TreeNode[int]{
				left: &TreeNode[int]{
					right: &TreeNode[int]{
						right: nil,
						left:  nil,
						color: Red,
					},
					left:  nil,
					color: Black,
				},
				right: &TreeNode[int]{
					right: nil,
					left: &TreeNode[int]{
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
			node: &TreeNode[int]{
				parent: nil,
				left: &TreeNode[int]{
					parent: &TreeNode[int]{
						color: Black,
					},
					right: &TreeNode[int]{
						right: nil,
						left:  nil,
						color: Red,
						parent: &TreeNode[int]{
							color: Red,
						},
					},
					left:  nil,
					color: Red,
				},
				right: &TreeNode[int]{
					parent: &TreeNode[int]{
						color: Black,
					},
					right: nil,
					left: &TreeNode[int]{
						right: nil,
						left:  nil,
						color: Red,
						parent: &TreeNode[int]{
							color: Red,
						},
					},
					color: Red,
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
				redBlackTree.Add(&TreeNode[int]{
					key: tt.k[i],
				})
			}
			res := isRedBlackTree[int](redBlackTree.root)
			assert.Equal(t, tt.want, res)

		})
	}
}

func TestRedBlackTree_Delete(t *testing.T) {}

func TestRedBlackTree_Find(t *testing.T) {}

func TestRedBlackTree_addNode(t *testing.T) {}

func TestRedBlackTree_deleteNode(t *testing.T) {}

func TestRedBlackTree_fixAfterAdd(t *testing.T) {}

func TestRedBlackTree_fixAfterDelete(t *testing.T) {}

func TestRedBlackTree_getTreeNode(t *testing.T) {}

func TestRedBlackTree_rotateLeft(t *testing.T) {}

func TestRedBlackTree_rotateRight(t *testing.T) {}

func TestRedBlackTree_successor(t *testing.T) {}

func Test_treeNode_colorOf(t *testing.T) {}

func Test_treeNode_leftOf(t *testing.T) {}

func Test_treeNode_parentOf(t *testing.T) {}

func Test_treeNode_rightOf(t *testing.T) {}

func Test_treeNode_setColor(t *testing.T) {}
