package mapx

import "github.com/gotomicro/ekit"

// TreeMap 是基于红黑树实现的Map
// 需要注意TreeMap是有序的所以必须传入比较器
// 如果没有传入那么Key必须要实现Comparator
type TreeMap[Key comparable, Val any] struct {
	compare ekit.Comparator[Key]
	root    Node[Key, Val]
	size    int
}

type Node[Key comparable, Val any] struct {
	values Val
	key    Key
	left   *Node[Key, Val]
	right  *Node[Key, Val]
	parent *Node[Key, Val]
}

func NewTreeMapWithComparator[Key comparable, Val any](compare ekit.Comparator[Key]) TreeMap[Key, Val] {
	return TreeMap[Key, Val]{
		compare: compare,
	}
}
func NewTreeMapWithMap() {

}
func NewTreeMap() {

}

func PutAll[Key comparable, Val any](m map[Key]Val) {

}

func Put[Key comparable, Val any](m map[Key]Val) {

}

func Get[Key comparable, Val any](m map[Key]Val) {

}
