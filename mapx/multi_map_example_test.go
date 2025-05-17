package mapx

import (
	"fmt"
	"sort"

	"github.com/ecodeclub/ekit/tuple/pair"
)

func ExampleMultiMap_Iterate() {
	multiMap := NewMultiHashMap[testStringData, int](0)
	arr := make([]pair.Pair[string, int], 0)
	multiMap.Put(testStringData{data: "hello"}, 1)
	multiMap.Put(testStringData{data: "world"}, 2)
	multiMap.Put(testStringData{data: "world"}, 3)
	multiMap.Put(testStringData{data: "world"}, 4)
	multiMap.Put(testStringData{data: "ekit"}, 3)

	multiMap.Iterate(
		func(key testStringData, val int) bool {
			arr = append(arr, pair.NewPair(key.data, val))
			return true
		})

	sort.Slice(arr, func(i, j int) bool {
		if arr[i].Key == arr[j].Key {
			return arr[i].Value < arr[j].Value
		}
		return arr[i].Key < arr[j].Key
	})

	for _, pa := range arr {
		fmt.Println(pa.Key, pa.Value)
	}
	// Output:
	// ekit 3
	// hello 1
	// world 2
	// world 3
	// world 4
}
