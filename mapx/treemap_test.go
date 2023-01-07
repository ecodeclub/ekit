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
			k:       []int{1, 3, 2, 0, 2, 3},
			v:       []string{"1", "2", "998", "0", "3", "997"},
			wantKey: []int{0, 1, 2, 3},
			wantVal: []string{"0", "1", "3", "997"},
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

func TestTreeMap_PutAll(t *testing.T) {
	tests := []struct {
		name    string
		m       map[int]string
		treeMap *TreeMap[int, string]
		wantKey []int
		wantVal []string
		wantErr error
	}{
		{
			name:    "nil",
			treeMap: NewTreeMap[int, string](compare()),
			m:       nil,
			wantKey: nil,
			wantVal: nil,
			wantErr: nil,
		},
		{
			name:    "empty",
			treeMap: NewTreeMap[int, string](compare()),
			m:       map[int]string{},
			wantKey: nil,
			wantVal: nil,
			wantErr: nil,
		},
		{
			name:    "single",
			treeMap: NewTreeMap[int, string](compare()),
			m: map[int]string{
				0: "0",
			},
			wantKey: []int{0},
			wantVal: []string{"0"},
			wantErr: nil,
		},
		{
			name:    "multiple",
			treeMap: NewTreeMap[int, string](compare()),
			m: map[int]string{
				0: "0",
				1: "1",
				2: "2",
			},
			wantKey: []int{0, 1, 2},
			wantVal: []string{"0", "1", "2"},
			wantErr: nil,
		},
		{
			name:    "disorder",
			treeMap: NewTreeMap[int, string](compare()),
			m: map[int]string{
				1: "1",
				2: "2",
				0: "0",
				3: "3",
				5: "5",
				4: "4",
			},
			wantKey: []int{0, 1, 2, 3, 4, 5},
			wantVal: []string{"0", "1", "2", "3", "4", "5"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.treeMap.PutAll(tt.m)
			if err != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}
			keys, val := tt.treeMap.keyValue()
			assert.ElementsMatch(t, tt.wantKey, keys)
			assert.ElementsMatch(t, tt.wantVal, val)
		})
	}
}

func TestTreeMap_Get(t *testing.T) {
	var tests = []struct {
		name    string
		m       map[int]int
		findKey int
		treeMap *TreeMap[int, int]
		wantVal int
		wantErr error
	}{
		{
			name:    "empty-TreeMap",
			treeMap: NewTreeMap[int, int](compare()),
			m:       map[int]int{},
			findKey: 0,
			wantVal: 0,
			wantErr: errors.New("ekit: TreeMap未找到指定Key"),
		},
		{
			name: "compare-nil",
			treeMap: &TreeMap[int, int]{
				compare: nil,
			},
			m:       map[int]int{},
			findKey: 0,
			wantVal: 0,
			wantErr: errors.New("ekit: Comparator不能为nil"),
		},
		{
			name:    "find",
			treeMap: NewTreeMap[int, int](compare()),
			m: map[int]int{
				1: 1,
				2: 2,
				0: 0,
				3: 3,
				5: 5,
				4: 4,
			},
			findKey: 2,
			wantVal: 2,
			wantErr: nil,
		},
		{
			name:    "not-find",
			treeMap: NewTreeMap[int, int](compare()),
			m: map[int]int{
				1: 1,
				2: 2,
				0: 0,
				3: 3,
				5: 5,
				4: 4,
			},
			findKey: 6,
			wantVal: 0,
			wantErr: errors.New("ekit: TreeMap未找到指定Key"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.treeMap.compare == nil {
				val, err := tt.treeMap.Get(tt.findKey)
				if err != nil {
					assert.Equal(t, tt.wantErr, err)
					assert.Equal(t, tt.wantVal, val)
					return
				}
			} else {
				tt.treeMap.PutAll(tt.m)
				val, err := tt.treeMap.Get(tt.findKey)
				assert.Equal(t, tt.wantErr, err)
				assert.Equal(t, tt.wantVal, val)
			}

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
	if treeMap.size > 0 {
		midOrder(treeNode, m)
		return m.ks, m.vals
	}
	return nil, nil
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

func compare() ekit.Comparator[int] {
	return ekit.ComparatorRealNumber[int]
}
