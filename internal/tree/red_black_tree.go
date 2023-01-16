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

	"github.com/gotomicro/ekit"
)

const (
	Red   = false
	Black = true
)

var (
	ErrRBTreeSameRBNode = errors.New("ekit: RBTree不能添加重复节点Key")
	ErrRBTreeNotRBNode  = errors.New("ekit: RBTree不存在节点Key")
	// errRBTreeCantRepaceNil = errors.New("ekit: RBTree不能将节点替换为nil")
)

type RBTree[K any, V any] struct {
	root    *RBNode[K, V]
	compare ekit.Comparator[K]
	size    int
}

func (rb *RBTree[K, V]) Size() int {
	if rb == nil {
		return 0
	}
	return rb.size
}

func (rb *RBTree[K, V]) Root() *RBNode[K, V] {
	if rb == nil {
		return nil
	}
	return rb.root
}

type RBNode[K any, V any] struct {
	color               bool
	key                 K
	value               V
	left, right, parent *RBNode[K, V]
}

func (node *RBNode[K, V]) setValue(v V) {
	if node == nil {
		return
	}
	node.value = v
}
func (node *RBNode[K, V]) Left() *RBNode[K, V] {
	return node.getLeft()
}

func (node *RBNode[K, V]) Right() *RBNode[K, V] {
	return node.getRight()
}

func (node *RBNode[K, V]) Parent() *RBNode[K, V] {
	return node.getParent()
}

func (node *RBNode[K, V]) Key() K {
	return node.key
}

func (node *RBNode[K, V]) Value() V {
	return node.value
}

// NewRBTree 构建红黑树
func NewRBTree[K any, V any](compare ekit.Comparator[K]) *RBTree[K, V] {
	return &RBTree[K, V]{
		compare: compare,
		root:    nil,
	}
}

func newRBNode[K any, V any](key K, value V) *RBNode[K, V] {
	return &RBNode[K, V]{
		key:    key,
		value:  value,
		color:  Red,
		left:   nil,
		right:  nil,
		parent: nil,
	}
}

// Add 增加节点
func (rb *RBTree[K, V]) Add(key K, value V) error {
	return rb.addNode(newRBNode(key, value))
}

// Delete 删除节点
func (rb *RBTree[K, V]) Delete(key K) {
	if node := rb.findRBNode(key); node != nil {
		rb.deleteNode(node)
	}
}

// Find 查找节点
func (rb *RBTree[K, V]) Find(key K) *RBNode[K, V] {
	return rb.findRBNode(key)
}
func (rb *RBTree[K, V]) SetValue(key K, value V) error {
	node := rb.Find(key)
	if node == nil {
		return ErrRBTreeNotRBNode
	}
	node.setValue(value)
	return nil
}

// addNode 插入新节点
func (rb *RBTree[K, V]) addNode(node *RBNode[K, V]) error {
	var fixNode *RBNode[K, V]
	if rb.root == nil {
		rb.root = newRBNode[K, V](node.key, node.value)
		fixNode = rb.root
	} else {
		t := rb.root
		cmp := 0
		parent := &RBNode[K, V]{}
		for t != nil {
			parent = t
			cmp = rb.compare(node.key, t.key)
			if cmp < 0 {
				t = t.left
			} else if cmp > 0 {
				t = t.right
			} else if cmp == 0 {
				return ErrRBTreeSameRBNode
			}
		}
		fixNode = &RBNode[K, V]{
			key:    node.key,
			parent: parent,
			value:  node.value,
			color:  Red,
		}
		if cmp < 0 {
			parent.left = fixNode
		} else {
			parent.right = fixNode
		}
	}
	rb.size++
	rb.fixAfterAdd(fixNode)
	return nil
}

// deleteNode 红黑树删除方法
// 删除分两步,第一步取出后继节点,第二部着色旋转
// 取后继节点
// case1:node左右非空子节点,通过getSuccessor获取后继节点
// case2:node左右只有一个非空子节点
// case3:node左右均为空节点
// 着色旋转
// case1:当删除节点非空且为黑色时,会违反红黑树任何路径黑节点个数相同的约束,所以需要重新平衡
// case2:当删除红色节点时,不会破坏任何约束,所以不需要平衡
func (rb *RBTree[K, V]) deleteNode(node *RBNode[K, V]) {
	// node左右非空,取后继节点
	if node.left != nil && node.right != nil {
		s := rb.findSuccessor(node)
		node.key = s.key
		node.value = s.value
		node = s
	}
	var replacement *RBNode[K, V]
	// node节点只有一个非空子节点
	if node.left != nil {
		replacement = node.left
	} else {
		replacement = node.right
	}
	if replacement != nil {
		replacement.parent = node.parent
		if node.parent == nil {
			rb.root = replacement
		} else if node == node.parent.left {
			node.parent.left = replacement
		} else {
			node.parent.right = replacement
		}
		node.left = nil
		node.right = nil
		node.parent = nil
		if node.color {
			rb.fixAfterDelete(replacement)
		}
	} else if node.parent == nil {
		// 如果node节点无父节点,说明node为root节点
		rb.root = nil
	} else {
		// node子节点均为空
		if node.color {
			rb.fixAfterDelete(node)
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
	rb.size--
}

// findSuccessor 寻找后继节点
// case1: node节点存在右子节点,则右子树的最小节点是node的后继节点
// case2: node节点不存在右子节点,则其第一个为左节点的祖先的父节点为node的后继节点
func (rb *RBTree[K, V]) findSuccessor(node *RBNode[K, V]) *RBNode[K, V] {
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
		for p != nil && ch == p.right {
			ch = p
			p = p.parent
		}
		return p
	}

}

func (rb *RBTree[K, V]) findRBNode(key K) *RBNode[K, V] {
	node := rb.root
	for node != nil {
		cmp := rb.compare(key, node.key)
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
// 如果是空节点、root节点、父节点是黑无需构建
// 可分为3种情况
// fixUncleRed 叔叔节点是红色右节点
// fixAddLeftBlack 叔叔节点是黑色右节点
// fixAddRightBlack 叔叔节点是黑色左节点
func (rb *RBTree[K, V]) fixAfterAdd(x *RBNode[K, V]) {
	x.color = Red
	for x != nil && x != rb.root && !x.getParent().getColor() {
		y := x.getUncle()
		if !y.getColor() {
			x = rb.fixUncleRed(x, y)
			continue
		}
		if x.getParent() == x.getGrandParent().getLeft() {
			x = rb.fixAddLeftBlack(x)
			continue
		}
		x = rb.fixAddRightBlack(x)
	}
	rb.root.setColor(Black)
}

// fixAddLeftRed 叔叔节点是红色右节点，由于不能存在连续红色节点,此时祖父节点x.getParent().getParent()必为黑。另x为红所以叔父节点需要变黑，祖父变红，此时红黑树完成
//
//							  b(b)                    b(r)
//							/		\				/		\
//						  a(r)        y(r)  ->   a(b)        y(b)
//						/   \       /  \         /   \       /  \
//		            x(r)    nil   nil  nil    x (r) nil   nil  nil
//	             	/  \                      /  \
//	            	nil nil                   nil nil
func (rb *RBTree[K, V]) fixUncleRed(x *RBNode[K, V], y *RBNode[K, V]) *RBNode[K, V] {
	x.getParent().setColor(Black)
	y.setColor(Black)
	x.getGrandParent().setColor(Red)
	x = x.getGrandParent()
	return x
}

// fixAddLeftBlack 叔叔节点是黑色右节点.x节点是父节点左节点,执行左旋，此时x节点变为原x节点的父节点a,也就是左子节点。的接着将x的父节点和爷爷节点的颜色对换。然后对爷爷节点进行右旋转,此时红黑树完成
// 如果x为左节点则跳过左旋操作
//
//							  b(b)                    b(b)                b(r)
//							/		\				/		\            /   \
//						  a(r)        y(b)  ->   a(r)        y(b)  ->  a(b)   y(b)
//						/   \       /  \         /   \       /  \      /  \    /  \
//		               nil   x (r) nil  nil      x(r) nil  nil  nil   x(r) nil nil nil
//	           		 		 /  \               /  \                  / \
//	           		 		nil nil             nil nil              nil nil
func (rb *RBTree[K, V]) fixAddLeftBlack(x *RBNode[K, V]) *RBNode[K, V] {
	if x == x.getParent().getRight() {
		x = x.getParent()
		rb.rotateLeft(x)
	}
	x.getParent().setColor(Black)
	x.getGrandParent().setColor(Red)
	rb.rotateRight(x.getGrandParent())
	return x
}

// fixAddRightBlack 叔叔节点是黑色左节点.x节点是父节点右节点,执行右旋，此时x节点变为原x节点的父节点a,也就是右子节点。接着将x的父节点和爷爷节点的颜色对换。然后对爷爷节点进行右旋转,此时红黑树完成
// 如果x为右节点则跳过右旋操作
//
//							  b(b)                    b(b)                b(r)
//							/		\				/		\            /   \
//						  y(b)       a(r)  ->   y(b)        a(r)  ->  y(b)     a(b)
//						/   \       /  \         /   \       /  \      /  \    /  \
//		               nil   nil x(r)  nil      nil nil  nil  x(r)   nil nil  nil  x(r)
//	           		 		      /  \                         /  \               /  \
//	           		 		      nil nil                    nil nil              nil nil
func (rb *RBTree[K, V]) fixAddRightBlack(x *RBNode[K, V]) *RBNode[K, V] {
	if x == x.getParent().getLeft() {
		x = x.getParent()
		rb.rotateRight(x)
	}
	x.getParent().setColor(Black)
	x.getGrandParent().setColor(Red)
	rb.rotateLeft(x.getGrandParent())
	return x
}

// fixAfterDelete 删除时着色旋转
// 根据x是节点位置分为fixAfterDeleteLeft,fixAfterDeleteRight两种情况
func (rb *RBTree[K, V]) fixAfterDelete(x *RBNode[K, V]) {
	for x != rb.root && x.getColor() {
		if x == x.parent.getLeft() {
			x = rb.fixAfterDeleteLeft(x)
		} else {
			x = rb.fixAfterDeleteRight(x)
		}
	}
	x.setColor(Black)
}

// fixAfterDeleteLeft 处理x为左子节点时的平衡处理
func (rb *RBTree[K, V]) fixAfterDeleteLeft(x *RBNode[K, V]) *RBNode[K, V] {
	sib := x.getParent().Right()
	if !sib.getColor() {
		sib.setColor(Black)
		sib.getParent().setColor(Red)
		rb.rotateLeft(x.getParent())
		sib = x.getParent().getRight()
	}
	if sib.getLeft().getColor() && sib.getRight().getColor() {
		sib.setColor(Red)
		x = x.getParent()
	} else {
		if sib.getRight().getColor() {
			sib.getLeft().setColor(Black)
			sib.setColor(Red)
			rb.rotateRight(sib)
			sib = x.getParent().getRight()
		}
		sib.setColor(x.getParent().getColor())
		x.getParent().setColor(Black)
		sib.getRight().setColor(Black)
		rb.rotateLeft(x.getParent())
		x = rb.root
	}
	return x
}

// fixAfterDeleteRight 处理x为右子节点时的平衡处理
func (rb *RBTree[K, V]) fixAfterDeleteRight(x *RBNode[K, V]) *RBNode[K, V] {
	sib := x.getParent().Left()
	if !sib.getColor() {
		sib.setColor(Black)
		x.getParent().setColor(Red)
		rb.rotateRight(x.getParent())
		sib = x.getBrother()
	}
	if sib.getRight().getColor() && sib.getLeft().getColor() {
		sib.setColor(Red)
		x = x.getParent()
	} else {
		if sib.getLeft().getColor() {
			sib.getRight().setColor(Black)
			sib.setColor(Red)
			rb.rotateLeft(sib)
			sib = x.getParent().getLeft()
		}
		sib.setColor(x.getParent().getColor())
		x.getParent().setColor(Black)
		sib.getLeft().setColor(Black)
		rb.rotateRight(x.getParent())
		x = rb.root
	}
	return x
}

// rotateLeft 左旋转
//
//							  b                    a
//							/	\				  /	  \
//						  c       a  ->    		 b     y
//								 / \            /  \
//		                     	x    y     		c	x

func (rb *RBTree[K, V]) rotateLeft(node *RBNode[K, V]) {
	if node == nil || node.getRight() == nil {
		return
	}
	r := node.right
	node.right = r.left
	if r.left != nil {
		r.left.parent = node
	}
	r.parent = node.parent
	if node.parent == nil {
		rb.root = r
	} else if node.parent.left == node {
		node.parent.left = r
	} else {
		node.parent.right = r
	}
	r.left = node
	node.parent = r

}

// rotateRight 右旋转
//
//						  b                    c
//						/	\				  /	  \
//					  c       a  ->    		 x     b
//					 /	\	                       / \
//	                 x  y  	     	 	  	       y  a
func (rb *RBTree[K, V]) rotateRight(node *RBNode[K, V]) {
	if node == nil || node.getLeft() == nil {
		return
	}
	l := node.left
	node.left = l.right
	if l.right != nil {
		l.right.parent = node
	}
	l.parent = node.parent
	if node.parent == nil {
		rb.root = l
	} else if node.parent.right == node {
		node.parent.right = l
	} else {
		node.parent.left = l
	}
	l.right = node
	node.parent = l

}

func (node *RBNode[K, V]) getColor() bool {
	if node == nil {
		return Black
	}
	return node.color
}

func (node *RBNode[K, V]) setColor(color bool) {
	if node == nil {
		return
	}
	node.color = color
}

func (node *RBNode[K, V]) getParent() *RBNode[K, V] {
	if node == nil {
		return nil
	}
	return node.parent
}

func (node *RBNode[K, V]) getLeft() *RBNode[K, V] {
	if node == nil {
		return nil
	}
	return node.left
}

func (node *RBNode[K, V]) getRight() *RBNode[K, V] {
	if node == nil {
		return nil
	}
	return node.right
}

func (node *RBNode[K, V]) getUncle() *RBNode[K, V] {
	if node == nil {
		return nil
	}
	return node.getParent().getBrother()
}
func (node *RBNode[K, V]) getGrandParent() *RBNode[K, V] {
	if node == nil {
		return nil
	}
	return node.getParent().getParent()
}
func (node *RBNode[K, V]) getBrother() *RBNode[K, V] {
	if node == nil {
		return nil
	}
	if node == node.getParent().getLeft() {
		return node.getParent().getRight()
	}
	return node.getParent().getLeft()
}
