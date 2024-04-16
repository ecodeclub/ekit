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

package list

import (
	"testing"

	"github.com/ecodeclub/ekit"
	"github.com/stretchr/testify/assert"
)

func TestNewSkipList(t *testing.T) {
	testCases := []struct {
		name      string
		compare   ekit.Comparator[int]
		wantSlice []int
	}{
		{
			name:      "new skip list",
			compare:   ekit.ComparatorRealNumber[int],
			wantSlice: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sl := NewSkipList(tc.compare)
			assert.Equal(t, tc.wantSlice, sl.AsSlice())
		})
	}
}

func TestSkipList_AsSlice(t *testing.T) {
	testCases := []struct {
		name      string
		compare   ekit.Comparator[int]
		wantSlice []int
	}{
		{
			name:      "no err is ok",
			compare:   ekit.ComparatorRealNumber[int],
			wantSlice: []int{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sl := NewSkipList[int](tc.compare)
			assert.Equal(t, tc.wantSlice, sl.AsSlice())
		})
	}
}

func TestSkipList_Cap(t *testing.T) {
	testCases := []struct {
		name     string
		compare  ekit.Comparator[int]
		wantSize int
	}{
		{
			name:     "no err is ok",
			compare:  ekit.ComparatorRealNumber[int],
			wantSize: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sl := NewSkipList[int](tc.compare)
			assert.Equal(t, tc.wantSize, sl.Cap())
		})
	}
}

func TestSkipList_DeleteElement(t *testing.T) {
	testCases := []struct {
		name     string
		compare  ekit.Comparator[int]
		value    int
		wantBool bool
	}{
		{
			name:     "no err is ok",
			compare:  ekit.ComparatorRealNumber[int],
			value:    1,
			wantBool: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sl := NewSkipList[int](tc.compare)
			ok := sl.DeleteElement(tc.value)
			assert.Equal(t, tc.wantBool, ok)
		})
	}
}

func TestSkipList_Insert(t *testing.T) {
	testCases := []struct {
		name      string
		compare   ekit.Comparator[int]
		key       int
		wantSlice []int
	}{
		{
			name:      "no err is ok",
			compare:   ekit.ComparatorRealNumber[int],
			key:       1,
			wantSlice: []int{1},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sl := NewSkipList[int](tc.compare)
			sl.Insert(tc.key)
			assert.Equal(t, tc.wantSlice, sl.AsSlice())
		})
	}
}

func TestSkipList_Len(t *testing.T) {
	testCases := []struct {
		name     string
		compare  ekit.Comparator[int]
		wantSize int
	}{
		{
			name:     "no err is ok",
			compare:  ekit.ComparatorRealNumber[int],
			wantSize: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sl := NewSkipList[int](tc.compare)
			assert.Equal(t, tc.wantSize, sl.Len())
		})
	}
}

func TestSkipList_Search(t *testing.T) {
	testCases := []struct {
		name     string
		compare  ekit.Comparator[int]
		value    int
		wantBool bool
	}{
		{
			name:     "no err is ok",
			compare:  ekit.ComparatorRealNumber[int],
			value:    1,
			wantBool: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sl := NewSkipList[int](tc.compare)
			ok := sl.Search(tc.value)
			assert.Equal(t, tc.wantBool, ok)
		})
	}
}
