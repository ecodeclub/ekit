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

package slice

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		want []string
	}{
		{
			name: "src nil",
			want: []string{},
		},
		{
			name: "src empty",
			src:  []int{},
			want: []string{},
		},
		{
			name: "src has element",
			src:  []int{1, 2, 3},
			want: []string{"1", "2", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Map(tt.src, func(idx int, src int) string {
				return strconv.Itoa(src)
			})
			assert.Equal(t, res, tt.want)
		})
	}
}

func ExampleMap() {
	src := []int{1, 2, 3}
	dst := Map(src, func(idx int, src int) string {
		return strconv.Itoa(src)
	})
	fmt.Println(dst)
	// Output: [1 2 3]
}

func TestFilterMap(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		want []string
	}{
		{
			name: "src nil",
			want: []string{},
		},
		{
			name: "src empty",
			src:  []int{},
			want: []string{},
		},
		{
			name: "src has element",
			src:  []int{1, -2, 3},
			want: []string{"1", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := FilterMap(tt.src, func(idx int, src int) (string, bool) {
				return strconv.Itoa(src), src >= 0
			})
			assert.Equal(t, res, tt.want)
		})
	}
}

func ExampleFilterMap() {
	src := []int{1, -2, 3}
	dst := FilterMap[int, string](src, func(idx int, src int) (string, bool) {
		return strconv.Itoa(src), src >= 0
	})
	fmt.Println(dst)
	// Output: [1 3]
}

func TestToMapV(t *testing.T) {
	t.Run("integer-string to map[int]int", func(t *testing.T) {
		elements := []string{"1", "2", "3", "4", "5"}
		resMap := ToMapV(elements, func(str string) (int, int) {
			num, _ := strconv.Atoi(str)
			return num, num
		})
		epectedMap := map[int]int{
			1: 1,
			2: 2,
			3: 3,
			4: 4,
			5: 5,
		}
		assert.Equal(t, epectedMap, resMap)
	})
	t.Run("struct<string, string, int> to map[string]struct<string, string, int>", func(t *testing.T) {
		type eleType struct {
			A string
			B string
			C int
		}
		elements := []eleType{
			{
				A: "a",
				B: "b",
				C: 1,
			},
			{
				A: "c",
				B: "d",
				C: 2,
			},
		}
		resMap := ToMapV(elements, func(ele eleType) (string, eleType) {
			return ele.A, ele
		})
		epectedMap := map[string]eleType{
			"a": {
				A: "a",
				B: "b",
				C: 1,
			},
			"c": {
				A: "c",
				B: "d",
				C: 2,
			},
		}
		assert.Equal(t, epectedMap, resMap)
	})

	t.Run("struct<string, string, int> to map[string]struct<string, string, int>, 重复的key", func(t *testing.T) {
		type eleType struct {
			A string
			B string
			C int
		}
		elements := []eleType{
			{
				A: "a",
				B: "b",
				C: 1,
			},
			{
				A: "c",
				B: "d",
				C: 2,
			},
			{
				A: "a",
				B: "d",
				C: 3,
			},
		}
		resMap := ToMapV(elements, func(ele eleType) (string, eleType) {
			return ele.A, ele
		})
		epectedMap := map[string]eleType{
			"a": {
				A: "a",
				B: "d",
				C: 3,
			},
			"c": {
				A: "c",
				B: "d",
				C: 2,
			},
		}
		assert.Equal(t, epectedMap, resMap)
	})

	t.Run("传入nil slice,返回空map", func(t *testing.T) {
		var elements []string = nil
		resMap := ToMapV(elements, func(str string) (int, int) {
			num, _ := strconv.Atoi(str)
			return num, num
		})
		epectedMap := make(map[int]int)
		assert.Equal(t, epectedMap, resMap)
	})
}

func TestToMap(t *testing.T) {
	t.Run("integer-string to map[int]string", func(t *testing.T) {
		elements := []string{"1", "2", "3", "4", "5"}
		resMap := ToMap(elements, func(str string) int {
			num, _ := strconv.Atoi(str)
			return num
		})
		epectedMap := map[int]string{
			1: "1",
			2: "2",
			3: "3",
			4: "4",
			5: "5",
		}
		assert.Equal(t, epectedMap, resMap)
	})
	t.Run("struct<string, string, int> to map[string]struct<string, string, int>", func(t *testing.T) {
		type eleType struct {
			A string
			B string
			C int
		}
		elements := []eleType{
			{
				A: "a",
				B: "b",
				C: 1,
			},
			{
				A: "c",
				B: "d",
				C: 2,
			},
		}
		resMap := ToMap(elements, func(ele eleType) string {
			return ele.A
		})
		epectedMap := map[string]eleType{
			"a": {
				A: "a",
				B: "b",
				C: 1,
			},
			"c": {
				A: "c",
				B: "d",
				C: 2,
			},
		}
		assert.Equal(t, epectedMap, resMap)
	})

	t.Run("struct<string, string, int> to map[string]struct<string, string, int>, 重复的key", func(t *testing.T) {
		type eleType struct {
			A string
			B string
			C int
		}
		elements := []eleType{
			{
				A: "a",
				B: "b",
				C: 1,
			},
			{
				A: "c",
				B: "d",
				C: 2,
			},
		}
		resMap := ToMap(elements, func(ele eleType) string {
			return ele.A
		})
		epectedMap := map[string]eleType{
			"a": {
				A: "a",
				B: "b",
				C: 1,
			},
			"c": {
				A: "c",
				B: "d",
				C: 2,
			},
		}
		assert.Equal(t, epectedMap, resMap)
	})

	t.Run("传入nil slice,返回空map", func(t *testing.T) {
		var elements []string = nil
		resMap := ToMap(elements, func(str string) int {
			num, _ := strconv.Atoi(str)
			return num
		})
		epectedMap := make(map[int]string)
		assert.Equal(t, epectedMap, resMap)
	})
}

func ExampleToMap() {
	elements := []string{"1", "2", "3", "4", "5"}
	resMap := ToMap(elements, func(str string) int {
		num, _ := strconv.Atoi(str)
		return num
	})
	fmt.Println(resMap)
	// Output: map[1:1 2:2 3:3 4:4 5:5]
}

func ExampleToMapV() {
	type eleType struct {
		A string
		B string
		C int
	}
	type eleTypeOut struct {
		A string
		B string
	}
	elements := []eleType{
		{
			A: "a",
			B: "b",
			C: 1,
		},
		{
			A: "c",
			B: "d",
			C: 2,
		},
	}
	resMap := ToMapV(elements, func(ele eleType) (string, eleTypeOut) {
		return ele.A, eleTypeOut{
			A: ele.A,
			B: ele.B,
		}
	})
	fmt.Println(resMap)
	// Output: map[a:{a b} c:{c d}]
}
