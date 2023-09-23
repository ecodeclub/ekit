package tree

import (
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/internal/tree"
	"github.com/stretchr/testify/assert"
	"testing"
)

func compare() ekit.Comparator[int] {
	return ekit.ComparatorRealNumber[int]
}

func TestNewRBTree(t *testing.T) {
	testCases := []struct {
		name    string
		compare ekit.Comparator[int]
		wantErr error
	}{
		{
			name:    "compare is nil",
			compare: nil,
			wantErr: ErrRBTreeComparatorIsNull,
		},
		{
			name:    "compare is ok",
			compare: compare(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewRBTree[int, string](tc.compare)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestRBTree_Add(t *testing.T) {
	testCases := []struct {
		name    string
		compare ekit.Comparator[int]
		key     int
		value   string
		wantErr error
	}{
		{
			name:    "no err is ok",
			compare: compare(),
			key:     1,
			value:   "value1",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rbTree, _ := NewRBTree[int, string](tc.compare)
			err := rbTree.Add(tc.key, tc.value)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestRBTree_Delete(t *testing.T) {
	testCases := []struct {
		name     string
		compare  ekit.Comparator[int]
		key      int
		wantBool bool
	}{
		{
			name:     "no err is ok",
			compare:  compare(),
			key:      1,
			wantBool: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rbTree, _ := NewRBTree[int, string](tc.compare)
			_, resultBool := rbTree.Delete(tc.key)
			assert.Equal(t, tc.wantBool, resultBool)
		})
	}
}

func TestRBTree_Set(t *testing.T) {
	testCases := []struct {
		name    string
		compare ekit.Comparator[int]
		key     int
		value   string
		wantErr error
	}{
		{
			name:    "no err is ok",
			compare: compare(),
			key:     1,
			value:   "value1",
			wantErr: tree.ErrRBTreeNotRBNode,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rbTree, _ := NewRBTree[int, string](tc.compare)
			err := rbTree.Set(tc.key, tc.value)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestRBTree_Find(t *testing.T) {
	testCases := []struct {
		name    string
		compare ekit.Comparator[int]
		key     int
		wantErr error
	}{
		{
			name:    "no err is ok",
			compare: compare(),
			key:     1,
			wantErr: tree.ErrRBTreeNotRBNode,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rbTree, _ := NewRBTree[int, string](tc.compare)
			_, err := rbTree.Find(tc.key)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestRBTree_Size(t *testing.T) {
	testCases := []struct {
		name     string
		compare  ekit.Comparator[int]
		wantSize int
	}{
		{
			name:     "no err is ok",
			compare:  compare(),
			wantSize: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rbTree, _ := NewRBTree[int, string](tc.compare)
			size := rbTree.Size()
			assert.Equal(t, tc.wantSize, size)
		})
	}
}

func TestRBTree_KeyValues(t *testing.T) {
	testCases := []struct {
		name    string
		compare ekit.Comparator[int]
	}{
		{
			name:    "no err is ok",
			compare: compare(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rbTree, _ := NewRBTree[int, string](tc.compare)
			keys, values := rbTree.KeyValues()
			assert.Equal(t, 0, len(keys))
			assert.Equal(t, 0, len(values))
		})
	}
}
