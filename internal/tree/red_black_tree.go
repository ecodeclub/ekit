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

type RBTree[T any] struct {
	root    *RBNode[T]
	compare ekit.Comparator[T]
}

type RBNode[T any] struct {
	color               bool
	key                 T
	left, right, parent *RBNode[T]
}

// NewRedBlackTree 构建红黑树
func NewRedBlackTree[T any](compare ekit.Comparator[T]) *RBTree[T] {
	return &RBTree[T]{
		compare: compare,
		root:    nil,
	}
}

// Add 增加节点
func (redBlackTree *RBTree[T]) Add(node *RBNode[T]) {
	redBlackTree.addNode(node)
}

// Delete 删除节点
func (redBlackTree *RBTree[T]) Delete(key T) {
	node := redBlackTree.getRBNode(key)
	if node == nil {
		return
	}
	redBlackTree.deleteNode(node)
}

// Find 查找节点
func (redBlackTree *RBTree[T]) Find(key T) *RBNode[T] {
	return redBlackTree.getRBNode(key)
}

func (redBlackTree *RBTree[T]) addNode(node *RBNode[T]) {
	t := redBlackTree.root
	if t == nil {
		redBlackTree.root = node
		return
	}
	cmp := 0
	parent := &RBNode[T]{}
	for t != nil {
		parent = t
		cmp = redBlackTree.compare(node.key, t.key)
		if cmp < 0 {
			t = t.left
		} else if cmp > 0 {
			t = t.right
		} else {
			return
		}
	}
	tempNode := &RBNode[T]{
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

// deleteNode 红黑树删除方法
// 删除分两步,第一步取出后继节点,第二部着色旋转
// 取后继节点
// case1:node左右非空子节点,通过successor获取后继节点
// case2:node左右只有一个非空子节点
// case3:node左右均为空节点
// 着色旋转
// case1:当删除节点非空且为黑色时,会违反红黑树任何路径黑节点个数相同的约束,所以需要重新平衡
// case2:当删除红色节点时,不会破坏任何约束,所以不需要平衡
func (redBlackTree *RBTree[T]) deleteNode(node *RBNode[T]) {
	// node左右非空,取后继节点
	if node.left != nil && node.right != nil {
		s := redBlackTree.successor(node)
		node.key = s.key
		node = s
	}
	replacement := &RBNode[T]{}
	// node节点只有一个非空子节点
	if node.left != nil {
		replacement = node.left
	} else {
		replacement = node.right
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
			redBlackTree.fixAfterDelete(replacement)
		}
	} else if node.parent == nil {
		// 如果node节点无父节点,说明node为root节点
		redBlackTree.root = nil
	} else {
		// node子节点均为空
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

// successor 寻找后继节点
// case1: node节点存在右子节点,则右子树的最小节点是node的后继节点
// case2: node节点不存在右子节点,则其第一个为左节点的祖先的父节点为node的后继节点
func (redBlackTree *RBTree[T]) successor(node *RBNode[T]) *RBNode[T] {
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

func (redBlackTree *RBTree[T]) getRBNode(key T) *RBNode[T] {
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
// 如果是空节点、root节点、父节点是黑无需构建
// 可分为4种情况
// fixUncleRed 叔叔节点是红色右节点
// fixAddLeftBlack 叔叔节点是黑色右节点
// fixAddRightBlack 叔叔节点是黑色左节点
func (redBlackTree *RBTree[T]) fixAfterAdd(x *RBNode[T]) {
	x.color = Red
	for x != nil && x != redBlackTree.root && !x.parent.color {
		if x.parentOf() == x.parentOf().parentOf().leftOf() {
			// y是叔节点
			y := x.parentOf().parentOf().rightOf()
			if !y.colorOf() {
				x = redBlackTree.fixUncleRed(x, y)
			} else {
				x = redBlackTree.fixAddLeftBlack(x)
			}
		} else {
			y := x.parentOf().parentOf().leftOf()
			if !y.colorOf() {
				x = redBlackTree.fixUncleRed(x, y)
			} else {
				x = redBlackTree.fixAddRightBlack(x)
			}
		}
	}
	redBlackTree.root.setColor(Black)
}

// fixAddLeftRed 叔叔节点是红色右节点，由于不能存在连续红色节点,此时祖父节点x.parentOf().parentOf()必为黑。另x为红所以叔父节点需要变黑，祖父变红，此时红黑树完成
//
//							  b(b)                    b(r)
//							/		\				/		\
//						  a(r)        y(r)  ->   a(b)        y(b)
//						/   \       /  \         /   \       /  \
//		            x(r)    nil   nil  nil    x (r) nil   nil  nil
//	             	/  \                      /  \
//	            	nil nil                   nil nil
func (redBlackTree *RBTree[T]) fixUncleRed(x *RBNode[T], y *RBNode[T]) *RBNode[T] {
	x.parentOf().setColor(Black)
	y.setColor(Black)
	x.parentOf().parentOf().setColor(Red)
	x = x.parentOf().parentOf()
	return x
}

// fixAddLeftBlack 叔叔节点是黑色右节点.x节点是父节点右节点,执行左旋，此时x节点变为原x节点的父节点a,也就是左子节点。的接着将x的父节点和爷爷节点的颜色对换。然后对爷爷节点进行右旋转,此时红黑树完成
// 如果x为左节点则跳过左旋操作
//
//							  b(b)                    b(b)                b(r)
//							/		\				/		\            /   \
//						  a(r)        y(b)  ->   a(r)        y(b)  ->  a(b)   y(b)
//						/   \       /  \         /   \       /  \      /  \    /  \
//		               nil   x (r) nil  nil      x(r) nil  nil  nil   x(r) nil nil nil
//	           		 		 /  \               /  \                  / \
//	           		 		nil nil             nil nil              nil nil
func (redBlackTree *RBTree[T]) fixAddLeftBlack(x *RBNode[T]) *RBNode[T] {
	if x == x.parentOf().rightOf() {
		x = x.parentOf()
		redBlackTree.rotateLeft(x)
	}
	x.parentOf().setColor(Black)
	x.parentOf().parentOf().setColor(Red)
	redBlackTree.rotateRight(x.parentOf().parentOf())
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
func (redBlackTree *RBTree[T]) fixAddRightBlack(x *RBNode[T]) *RBNode[T] {
	if x == x.parentOf().leftOf() {
		x = x.parentOf()
		redBlackTree.rotateRight(x)
	}
	x.parentOf().setColor(Black)
	x.parentOf().parentOf().setColor(Red)
	redBlackTree.rotateLeft(x.parentOf().parentOf())
	return x
}

// fixAfterDelete 删除时着色旋转
// 根据x是节点位置分为fixAfterDeleteLeft,fixAfterDeleteRight两种情况
func (redBlackTree *RBTree[T]) fixAfterDelete(x *RBNode[T]) {
	for x != redBlackTree.root && x.colorOf() {
		if x == x.parent.leftOf() {
			x = redBlackTree.fixAfterDeleteLeft(x)
		} else {
			x = redBlackTree.fixAfterDeleteRight(x)
		}
	}
	x.setColor(Black)
}

// fixAfterDeleteLeft 处理x为左子节点时的平衡处理
func (redBlackTree *RBTree[T]) fixAfterDeleteLeft(x *RBNode[T]) *RBNode[T] {
	sib := x.parentOf().rightOf()
	if !sib.colorOf() {
		sib.setColor(Black)
		sib.parentOf().setColor(Red)
		redBlackTree.rotateLeft(x.parentOf())
		sib = x.parentOf().rightOf()
	}
	if sib.leftOf().colorOf() && sib.rightOf().colorOf() {
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
	return x
}

// fixAfterDeleteRight 处理x为右子节点时的平衡处理
func (redBlackTree *RBTree[T]) fixAfterDeleteRight(x *RBNode[T]) *RBNode[T] {
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
	return x
}

// rotateLeft 左旋转
//
//							  b                    a
//							/	\				  /	  \
//						  c       a  ->    		 b     y
//								 / \            /  \
//		                     	x    y     		c	x

func (redBlackTree *RBTree[T]) rotateLeft(node *RBNode[T]) {
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
//
//						  b                    c
//						/	\				  /	  \
//					  c       a  ->    		 x     b
//					 /	\	                       / \
//	                 x  y  	     	 	  	       y  a
func (redBlackTree *RBTree[T]) rotateRight(node *RBNode[T]) {
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

func (node *RBNode[T]) colorOf() bool {
	if node == nil {
		return Black
	}
	return node.color
}

func (node *RBNode[T]) setColor(color bool) {
	if node == nil {
		return
	}
	node.color = color
}

func (node *RBNode[T]) parentOf() *RBNode[T] {
	if node == nil {
		return nil
	}
	return node.parent
}

func (node *RBNode[T]) leftOf() *RBNode[T] {
	if node == nil {
		return nil
	}
	return node.left
}

func (node *RBNode[T]) rightOf() *RBNode[T] {
	if node == nil {
		return nil
	}
	return node.right
}
