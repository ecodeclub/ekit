package mapx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMyHashMap(t *testing.T) {
	testKV := []struct {
		key testData
		val any
	}{
		{
			key: testData{
				id: 1,
			},
			val: 1,
		},
		{
			key: testData{
				id: 2,
			},
			val: 2,
		},
		{
			key: testData{
				id: 3,
			},
			val: 3,
		},
		{
			key: testData{
				id: 11,
			},
			val: 11,
		},
	}
	myhashmap := NewHashMap[testData](10)
	for _, kv := range testKV {
		myhashmap.Put(kv.key, kv.val)
	}
	wantHashMap := MyHashMap[testData]{
		hashmap: map[uint64]*Node{
			1: &Node{
				key:   testData{id: 1},
				value: 1,
				next: &Node{
					key:   testData{id: 11},
					value: 11,
				},
			},
			2: NewNode(NewTestData(2), 2),
			3: NewNode(NewTestData(3), 3),
		},
	}

	assert.Equal(t, wantHashMap.hashmap, myhashmap.hashmap)
	getTestCases := []struct {
		name    string
		key     testData
		wantVal any
		isFound bool
	}{
		{
			name: "get normal val",
			key: testData{
				id: 1,
			},
			wantVal: 1,
			isFound: true,
		},
		{
			name: "hash conflicts",
			key: testData{
				id: 11,
			},
			wantVal: 11,
			isFound: true,
		},
		{
			name: "hash not Found",
			key: testData{
				id: 8,
			},
			isFound: false,
		},
		{
			name: "val not Found",
			key: testData{
				id: 21,
			},
			isFound: false,
		},
	}
	for _, tc := range getTestCases {
		t.Run(tc.name, func(t *testing.T) {
			val, ok := myhashmap.Get(tc.key)
			assert.Equal(t, tc.isFound, ok)
			if !ok {
				return
			}
			assert.Equal(t, tc.wantVal, val)
		})
	}

}

type testData struct {
	id int
}

func (t testData) Code() uint64 {
	hash := t.id % 10
	return uint64(hash)
}

func (t testData) Comparable(key any) bool {
	val, ok := key.(testData)
	if !ok {
		return false
	}
	if t.id != val.id {
		return false
	}
	return true
}

func NewTestData(id int) testData {
	return testData{
		id: id,
	}
}
