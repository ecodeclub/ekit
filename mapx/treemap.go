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

package mapx

import (
	"errors"
	"github.com/gotomicro/ekit"
)

const (
	Red   = false
	Black = true
)

var (
	errTreeMapComparatorIsNull = errors.New("ekit: Comparator不能为nil")
	errTreeMapKeyNotFound      = errors.New("ekit: TreeMap未找到指定Key")
)

// TreeMap 是基于红黑树实现的Map
// 需要注意TreeMap是有序的所以必须传入比较器
// compare	比较器
// root	根节点
type TreeMap[Key ekit.RealNumber, Val any] struct {
	compare ekit.Comparator[Key]
	root    *treeNode[Key, Val]
	size    int
}

//	treeNode TreeMap上的节点，基于红黑树设计
//	values 节点值,可以为nil
//	key 节点必须为可比较类型
//	left 左子节点
//	right  右子节点
//  parent 父节点
//  color 着色 Black：true,Red :false
type treeNode[Key ekit.RealNumber, Val any] struct {
	values Val
	key    Key
	left   *treeNode[Key, Val]
	right  *treeNode[Key, Val]
	parent *treeNode[Key, Val]
	color  bool
}

func (node *treeNode[Key, Val]) setValue(val Val) {
	node.values = val
}

// BuildTreeMap TreeMap构造方法
// 支持传入compare比较器，并根据传入的m构建TreeMap
// 需注意比较器compare不能为nil
func BuildTreeMap[Key ekit.RealNumber, Val any](compare ekit.Comparator[Key], m map[Key]Val) (*TreeMap[Key, Val], error) {
	treeMap := NewTreeMap[Key, Val]()
	if err := treeMap.SetComparable(compare); err != nil {
		return nil, err
	}
	if err := treeMap.PutAll(m); err != nil {
		return nil, err
	}
	return treeMap, nil
}

// NewTreeMap TreeMap构造方法,创建一个的TreeMap
// 需注意比较器默认为ekit.ComparatorRealNumber
func NewTreeMap[Key ekit.RealNumber, Val any]() *TreeMap[Key, Val] {
	treeMap := &TreeMap[Key, Val]{
		compare: ekit.ComparatorRealNumber[Key],
	}
	return treeMap
}

func (treeMap *TreeMap[Key, Val]) SetComparable(compare ekit.Comparator[Key]) error {
	if compare != nil {
		treeMap.compare = compare
		return nil
	} else {
		return errTreeMapComparatorIsNull
	}
}

// PutAll 将一个可比较Key的map传入TreeMap
// 需注意如果map中的Key已存在TreeMap将被替换
// 错误：TreeMap中比较器为nil将会返回error
func (treeMap *TreeMap[Key, Val]) PutAll(m map[Key]Val) error {
	if len(m) != 0 {
		return nil
	}
	keys, values := KeysValues[Key, Val](m)
	for i := 0; i < len(keys); i++ {
		err := treeMap.Put(keys[i], values[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// Put 在TreeMap插入指定Key的节点
// 需注意如果TreeMap已存在该Key那么原值会被替换
// 错误：
// TreeMap中比较器为nil将会返回error
func (treeMap *TreeMap[Key, Val]) Put(k Key, v Val) error {
	if treeMap.compare == nil {
		return errTreeMapComparatorIsNull
	}
	t := treeMap.root
	if t == nil {
		treeMap.root = &treeNode[Key, Val]{
			key:    k,
			values: v,
		}
		treeMap.size++
		return nil
	}
	cmp := 0
	parent := &treeNode[Key, Val]{}

	for t != nil {
		parent = t
		cmp = treeMap.compare(k, t.key)
		if cmp < 0 {
			t = t.left
		} else if cmp > 0 {
			t = t.right
		} else {
			t.setValue(v)
			return nil
		}
	}

	node := &treeNode[Key, Val]{
		key:    k,
		values: v,
		parent: parent,
	}
	if cmp < 0 {
		parent.left = node
	} else {
		parent.right = node
	}
	treeMap.fixAfterPut(node)
	treeMap.size++

	return nil
}

// Get 在TreeMap找到指定Key的节点,返回Val
// 错误：
// TreeMap未找到指定Key将会返回error
// TreeMap中比较器为nil将会返回error
func (treeMap *TreeMap[Key, Val]) Get(k Key) (Val, error) {
	var v Val
	if treeMap.compare == nil {
		return v, errTreeMapComparatorIsNull
	}
	node := treeMap.root
	for node != nil {
		cmp := treeMap.compare(k, node.key)
		if cmp < 0 {
			node = node.left
		} else if cmp > 0 {
			node = node.right
		} else {
			return node.values, nil
		}
	}
	return v, errTreeMapKeyNotFound
}

// fixAfterPut 着色旋转
func (treeMap *TreeMap[Key, Val]) fixAfterPut(x *treeNode[Key, Val]) {
	x.color = Red
	for x != nil && x != treeMap.root && !x.parent.color {
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
					treeMap.rotateLeft(x)
				}
				x.parentOf().setColor(Black)
				x.parentOf().parentOf().setColor(Red)
				treeMap.rotateRight(x.parentOf().parentOf())
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
					treeMap.rotateRight(x)
				}
				x.parentOf().setColor(Black)
				x.parentOf().parentOf().setColor(Red)
				treeMap.rotateLeft(x.parentOf().parentOf())
			}
		}
	}
	treeMap.root.setColor(Black)
}

func (node *treeNode[Key, Val]) colorOf() bool {
	if node == nil {
		//空的叶子节点均为黑色
		return Black
	}
	return node.color
}

func (node *treeNode[Key, Val]) setColor(color bool) {
	if node == nil {
		return
	}
	node.color = color
}

func (node *treeNode[Key, Val]) parentOf() *treeNode[Key, Val] {
	if node == nil {
		return nil
	}
	return node.parent
}

func (node *treeNode[Key, Val]) leftOf() *treeNode[Key, Val] {
	if node == nil {
		return nil
	}
	return node.left
}

func (node *treeNode[Key, Val]) rightOf() *treeNode[Key, Val] {
	if node == nil {
		return nil
	}
	return node.right
}

// rotateLeft 左旋转
func (treeMap *TreeMap[Key, Val]) rotateLeft(node *treeNode[Key, Val]) {
	if node != nil {
		r := node.right
		node.right = r.left
		if r.left != nil {
			r.left.parent = node
		}
		r.parent = node.parent
		if node.parent == nil {
			treeMap.root = r
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
func (treeMap *TreeMap[Key, Val]) rotateRight(node *treeNode[Key, Val]) {
	if node != nil {
		l := node.left
		node.left = l.right
		if l.right != nil {
			l.right.parent = node
		}
		l.parent = node.parent
		if node.parent == nil {
			treeMap.root = l
		} else if node.parent.right == node {
			node.parent.right = l
		} else {
			node.parent.left = l
		}
		l.right = node
		node.parent = l
	}
}
