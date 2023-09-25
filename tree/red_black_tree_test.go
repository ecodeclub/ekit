// Copyright 2021 ecodeclub
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
	"testing"

	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/internal/tree"
	"github.com/stretchr/testify/assert"
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
			wantErr: errRBTreeComparatorIsNull,
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
