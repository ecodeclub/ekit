package mapx

import (
	"fmt"
	"hash/crc32"
	"sort"
)

type testStringData struct {
	data string
}

func (ts testStringData) Code() uint64 {
	return uint64(crc32.ChecksumIEEE([]byte(ts.data)))
}

func (ts testStringData) Equals(other any) bool {
	otherv, ok := other.(testStringData)
	if !ok {
		return false
	}
	return ts.data == otherv.data
}

func ExampleHashMap_Iterate() {
	hashMap := NewHashMap[testStringData, int](0)
	strArr := make([]string, 0)
	hashMap.Put(testStringData{data: "hello"}, 1)
	hashMap.Put(testStringData{data: "world"}, 2)
	hashMap.Put(testStringData{data: "ekit"}, 3)

	hashMap.Iterate(
		func(key testStringData, val int) bool {
			strArr = append(strArr, key.data)
			return true
		})

	sort.Strings(strArr)
	for _, s := range strArr {
		fmt.Println(s)
	}

	// Output:
	// ekit
	// hello
	// world
}
