package mapx

import (
	"errors"
	"github.com/gotomicro/ekit"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTreeMap_Put(t *testing.T) {
	tests := []struct {
		name    string
		k       []int
		v       []string
		treeMap *TreeMap[int, string]
		wantKey []int
		wantVal []string
		wantErr error
	}{
		{
			name:    "nil",
			treeMap: NewTreeMap[int, string](nil),
			k:       []int{0},
			v:       []string{"0"},
			wantErr: errors.New("ekit: Comparator不能为nil"),
		},
		{
			name:    "single",
			treeMap: NewTreeMap[int, string](compare()),
			k:       []int{0},
			v:       []string{"0"},
			wantKey: []int{0},
			wantVal: []string{"0"},
			wantErr: nil,
		},
		{
			name:    "multiple",
			treeMap: NewTreeMap[int, string](compare()),
			k:       []int{0, 1, 2},
			v:       []string{"0", "1", "2"},
			wantKey: []int{0, 1, 2},
			wantVal: []string{"0", "1", "2"},
			wantErr: nil,
		},
		{
			name:    "same",
			treeMap: NewTreeMap[int, string](compare()),
			k:       []int{0, 0, 1, 2, 2, 3},
			v:       []string{"0", "999", "1", "998", "2", "3"},
			wantKey: []int{0, 1, 2, 3},
			wantVal: []string{"999", "1", "2", "3"},
			wantErr: nil,
		},
		{
			name:    "same",
			treeMap: NewTreeMap[int, string](compare()),
			k:       []int{0, 0},
			v:       []string{"0", "999"},
			wantKey: []int{0},
			wantVal: []string{"999"},
			wantErr: nil,
		},
		{
			name:    "disorder",
			treeMap: NewTreeMap[int, string](compare()),
			k:       []int{1, 2, 0, 3, 5, 4},
			v:       []string{"1", "2", "0", "3", "5", "4"},
			wantKey: []int{0, 1, 2, 3, 4, 5},
			wantVal: []string{"0", "1", "2", "3", "4", "5"},
			wantErr: nil,
		},
		{
			name:    "disorder-same",
			treeMap: NewTreeMap[int, string](compare()),
			k:       []int{1, 2, 2, 0, 3, 3},
			v:       []string{"1", "2", "998", "0", "3", "997"},
			wantKey: []int{0, 1, 2, 3},
			wantVal: []string{"0", "1", "998", "997"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < len(tt.k); i++ {
				err := tt.treeMap.Put(tt.k[i], tt.v[i])
				if err != nil {
					assert.Equal(t, tt.wantErr, err)
					return
				}
			}
			keys, val := tt.treeMap.keyValue()
			assert.ElementsMatch(t, tt.wantKey, keys)
			assert.ElementsMatch(t, tt.wantVal, val)
		})
	}
}

type kv[Key comparable, Val any] struct {
	ks   []Key
	vals []Val
}

func (treeMap *TreeMap[int, string]) keyValue() ([]int, []string) {
	treeNode := treeMap.root
	m := &kv[int, string]{}

	midOrder(treeNode, m)
	return m.ks, m.vals
}
func midOrder[Key comparable, Val any](node *treeNode[Key, Val], m *kv[Key, Val]) {
	//先遍历左子树
	if node.left != nil {
		midOrder(node.left, m)
	}
	//再遍历自己
	if node != nil {
		m.ks = append(m.ks, node.key)
		m.vals = append(m.vals, node.values)
	}
	//最后遍历右子树
	if node.right != nil {
		midOrder(node.right, m)
	}
}

//func TestTreeMap_PutAll(t *testing.T) {
//	type fields struct {
//		compare ekit.Comparator
//		root    *treeNode[Key, Val]
//		size    int
//	}
//	type args struct {
//		m map[Key]Val
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr assert.ErrorAssertionFunc
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			treeMap := &TreeMap{
//				compare: tt.fields.compare,
//				root:    tt.fields.root,
//				size:    tt.fields.size,
//			}
//			tt.wantErr(t, treeMap.PutAll(tt.args.m), fmt.Sprintf("PutAll(%v)", tt.args.m))
//		})
//	}
//}

//
//func Test_treeNode_setValue(t *testing.T) {
//	type fields struct {
//		values Val
//		key    Key
//		left   *treeNode[Key, Val]
//		right  *treeNode[Key, Val]
//		parent *treeNode[Key, Val]
//		color  bool
//	}
//	type args struct {
//		val any
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			node := &treeNode{
//				values: tt.fields.values,
//				key:    tt.fields.key,
//				left:   tt.fields.left,
//				right:  tt.fields.right,
//				parent: tt.fields.parent,
//				color:  tt.fields.color,
//			}
//			node.setValue(tt.args.val)
//		})
//	}
//}

func compare() ekit.Comparator[int] {
	return func(a, b int) int {
		if a < b {
			return -1
		}
		if a == b {
			return 0
		}
		return 1
	}
}
