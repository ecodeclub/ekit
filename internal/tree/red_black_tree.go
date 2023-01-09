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

import "github.com/gotomicro/ekit"

const (
	Red   = false
	Black = true
)

type RedBlackTree[T any] struct {
	root    *TreeNode[T]
	compare ekit.Comparator[T]
}

type TreeNode[T any] struct {
	color               bool
	key                 T
	left, right, parent *TreeNode[T]
}

// NewRedBlackTree 构建红黑树
func NewRedBlackTree[T any](compare ekit.Comparator[T]) *RedBlackTree[T] {
	return &RedBlackTree[T]{
		compare: compare,
		root:    nil,
	}
}

// Add 增加节点
func (redBlackTree *RedBlackTree[T]) Add(node *TreeNode[T]) {
	redBlackTree.addNode(node)
}

// Delete 删除节点
func (redBlackTree *RedBlackTree[T]) Delete(key T) {
	node := redBlackTree.getTreeNode(key)
	if node == nil {
		return
	}
	redBlackTree.deleteNode(node)
}

// Find 查找节点
func (redBlackTree *RedBlackTree[T]) Find(key T) *TreeNode[T] {
	return redBlackTree.getTreeNode(key)
}

func (redBlackTree *RedBlackTree[T]) addNode(node *TreeNode[T]) {
	t := redBlackTree.root
	if t == nil {
		redBlackTree.root = node
		return
	}
	cmp := 0
	parent := &TreeNode[T]{}
	for t != nil {
		parent = t
		cmp = redBlackTree.compare(node.key, t.key)
		if cmp < 0 {
			t = t.left
		} else if cmp > 0 {
			t = t.right
		} else {
			// 重复节点,重新覆盖还是不处理
			t = node
			return
		}
	}
	tempNode := &TreeNode[T]{
		key:    node.key,
		parent: parent,
	}
	if cmp < 0 {
		parent.left = tempNode
	} else {
		parent.right = tempNode
	}
	redBlackTree.fixAfterAdd(tempNode)
}

func (redBlackTree *RedBlackTree[T]) deleteNode(node *TreeNode[T]) {
	// node 有2个子节点,需要选一个替换原有的位置
	if node.left != nil && node.right != nil {
		s := redBlackTree.successor(node)
		node = s
	}
	replacement := &TreeNode[T]{}
	if node.left != nil {
		replacement = node.left
	} else {
		replacement = node.right
	}
	if replacement != nil {
		replacement.parent = node.parent
	}
	if replacement != nil {
		replacement.parent = node.parent
		if node.parent == nil {
			redBlackTree.root = replacement
		} else if node == node.parent.left {
			node.parent.left = replacement
		} else {
			node.parent.right = replacement
		}
		node.left = nil
		node.right = nil
		node.parent = nil
		if node.color {
			redBlackTree.fixAfterDelete(node)
		}
	} else if node.parent == nil {
		redBlackTree.root = nil
	} else {
		if node.color {
			redBlackTree.fixAfterDelete(node)
		}
		if node.parent != nil {
			if node == node.parent.left {
				node.parent.left = nil
			} else if node == node.parent.right {
				node.parent.right = nil
			}
			node.parent = nil
		}
	}
}

func (redBlackTree *RedBlackTree[T]) successor(node *TreeNode[T]) *TreeNode[T] {
	if node == nil {
		return nil
	} else if node.right != nil {
		p := node.right
		for p.left != nil {
			p = p.left
		}
		return p
	} else {
		p := node.parent
		ch := node
		for p != nil && ch != p.right {
			ch = p
			p = p.parent
		}
		return p
	}

}

func (redBlackTree *RedBlackTree[T]) getTreeNode(key T) *TreeNode[T] {
	node := redBlackTree.root
	for node != nil {
		cmp := redBlackTree.compare(key, node.key)
		if cmp < 0 {
			node = node.left
		} else if cmp > 0 {
			node = node.right
		} else {
			return node
		}
	}
	return nil
}

// fixAfterAdd 插入时着色旋转
func (redBlackTree *RedBlackTree[T]) fixAfterAdd(x *TreeNode[T]) {
	x.color = Red
	for x != nil && x != redBlackTree.root && !x.parent.color {
		if x.parentOf() == x.parentOf().parentOf().leftOf() {
			y := x.parentOf().parentOf().rightOf()
			if !y.colorOf() {
				x.parent.setColor(Black)
				y.setColor(Black)
				x.parentOf().parentOf().setColor(Red)
				x = x.parentOf().parentOf()
			} else {
				if x == x.parentOf().parentOf().rightOf() {
					x = x.parentOf()
					redBlackTree.rotateLeft(x)
				}
				x.parentOf().setColor(Black)
				x.parentOf().parentOf().setColor(Red)
				redBlackTree.rotateRight(x.parentOf().parentOf())
			}
		} else {
			y := x.parentOf().parentOf().leftOf()
			if !y.colorOf() {
				x.parentOf().setColor(Black)
				y.setColor(Black)
				x.parentOf().parentOf().setColor(Red)
				x = x.parentOf().parentOf()
			} else {
				if x == x.parentOf().leftOf() {
					x = x.parentOf()
					redBlackTree.rotateRight(x)
				}
				x.parentOf().setColor(Black)
				x.parentOf().parentOf().setColor(Red)
				redBlackTree.rotateLeft(x.parentOf().parentOf())
			}
		}
	}
	redBlackTree.root.setColor(Black)
}

// fixAfterDelete 删除时着色旋转
func (redBlackTree *RedBlackTree[T]) fixAfterDelete(x *TreeNode[T]) {
	for x != redBlackTree.root && x.colorOf() {
		if x == x.parent.leftOf() {
			sib := x.parentOf().rightOf()
			if !sib.colorOf() {
				sib.setColor(Black)
				sib.parentOf().setColor(Red)
				redBlackTree.rotateLeft(x.parentOf())
				sib = x.parentOf().rightOf()
			}
			if sib.leftOf().colorOf() && sib.rightOf().color {
				sib.setColor(Red)
				x = x.parentOf()
			} else {
				if sib.rightOf().colorOf() {
					sib.leftOf().setColor(Black)
					sib.setColor(Red)
					redBlackTree.rotateRight(sib)
					sib = x.parentOf().rightOf()
				}
				sib.setColor(x.parentOf().colorOf())
				x.parentOf().setColor(Black)
				sib.rightOf().setColor(Black)
				redBlackTree.rotateLeft(x.parentOf())
				x = redBlackTree.root
			}

		} else {
			sib := x.parentOf().leftOf()
			if !sib.colorOf() {
				sib.setColor(Black)
				x.parentOf().setColor(Red)
				redBlackTree.rotateRight(x.parentOf())
				sib = x.parentOf().leftOf()
			}
			if sib.rightOf().colorOf() && sib.leftOf().colorOf() {
				sib.setColor(Red)
				x = x.parentOf()
			} else {
				if sib.leftOf().colorOf() {
					sib.rightOf().setColor(Black)
					sib.setColor(Red)
					redBlackTree.rotateLeft(sib)
					sib = x.parentOf().leftOf()
				}
				sib.setColor(x.parentOf().colorOf())
				x.parentOf().setColor(Black)
				sib.leftOf().setColor(Black)
				redBlackTree.rotateRight(x.parentOf())
				x = redBlackTree.root
			}
		}
	}
	x.setColor(Black)
}

// rotateLeft 左旋转
func (redBlackTree *RedBlackTree[T]) rotateLeft(node *TreeNode[T]) {
	if node != nil {
		r := node.right
		node.right = r.left
		if r.left != nil {
			r.left.parent = node
		}
		r.parent = node.parent
		if node.parent == nil {
			redBlackTree.root = r
		} else if node.parent.left == node {
			node.parent.left = r
		} else {
			node.parent.right = r
		}
		r.left = node
		node.parent = r
	}
}

// rotateRight 右旋转
func (redBlackTree *RedBlackTree[T]) rotateRight(node *TreeNode[T]) {
	if node != nil {
		l := node.left
		node.left = l.right
		if l.right != nil {
			l.right.parent = node
		}
		l.parent = node.parent
		if node.parent == nil {
			redBlackTree.root = l
		} else if node.parent.right == node {
			node.parent.right = l
		} else {
			node.parent.left = l
		}
		l.right = node
		node.parent = l
	}
}

func (node *TreeNode[T]) colorOf() bool {
	if node == nil {
		return Black
	}
	return node.color
}

func (node *TreeNode[T]) setColor(color bool) {
	if node == nil {
		return
	}
	node.color = color
}

func (node *TreeNode[T]) parentOf() *TreeNode[T] {
	if node == nil {
		return nil
	}
	return node.parent
}

func (node *TreeNode[T]) leftOf() *TreeNode[T] {
	if node == nil {
		return nil
	}
	return node.left
}

func (node *TreeNode[T]) rightOf() *TreeNode[T] {
	if node == nil {
		return nil
	}
	return node.right
}
