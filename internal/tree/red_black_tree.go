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
	root    *treeNode[T]
	compare ekit.Comparator[T]
}

type treeNode[T any] struct {
	color               bool
	key                 T
	left, right, parent *treeNode[T]
}

// NewRedBlackTree 构建红黑树
func NewRedBlackTree[T any](compare ekit.Comparator[T]) *RedBlackTree[T] {
	return &RedBlackTree[T]{
		compare: compare,
		root:    nil,
	}
}

// Add 增加节点
func (redBlackTree *RedBlackTree[T]) Add(node *treeNode[T]) {
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
func (redBlackTree *RedBlackTree[T]) Find(key T) *treeNode[T] {
	return redBlackTree.getTreeNode(key)
}

func (redBlackTree *RedBlackTree[T]) addNode(node *treeNode[T]) {
	t := redBlackTree.root
	if t == nil {
		redBlackTree.root = node
		return
	}
	cmp := 0
	parent := &treeNode[T]{}
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
		}
	}
	tempNode := &treeNode[T]{
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

func (redBlackTree *RedBlackTree[T]) deleteNode(node *treeNode[T]) {

	redBlackTree.fixAfterDelete(node)
}

func (redBlackTree *RedBlackTree[T]) getTreeNode(key T) *treeNode[T] {
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
func (redBlackTree *RedBlackTree[T]) fixAfterAdd(x *treeNode[T]) {
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
func (redBlackTree *RedBlackTree[T]) fixAfterDelete(x *treeNode[T]) {

}

// rotateLeft 左旋转
func (redBlackTree *RedBlackTree[T]) rotateLeft(node *treeNode[T]) {
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
func (redBlackTree *RedBlackTree[T]) rotateRight(node *treeNode[T]) {
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

func (node *treeNode[T]) colorOf() bool {
	if node == nil {
		return Black
	}
	return node.color
}

func (node *treeNode[T]) setColor(color bool) {
	if node == nil {
		return
	}
	node.color = color
}

func (node *treeNode[T]) parentOf() *treeNode[T] {
	if node == nil {
		return nil
	}
	return node.parent
}

func (node *treeNode[T]) leftOf() *treeNode[T] {
	if node == nil {
		return nil
	}
	return node.left
}

func (node *treeNode[T]) rightOf() *treeNode[T] {
	if node == nil {
		return nil
	}
	return node.right
}
