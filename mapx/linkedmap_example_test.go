package mapx

import (
	"fmt"
	"sort"
)

func ExampleLinkedMap_Iterate() {
	linkedMap := NewLinkedHashMap[testStringData, int](0)
	strArr := make([]string, 0)
	linkedMap.Put(testStringData{data: "hello"}, 1)
	linkedMap.Put(testStringData{data: "world"}, 2)
	linkedMap.Put(testStringData{data: "ekit"}, 3)

	linkedMap.Iterate(
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
