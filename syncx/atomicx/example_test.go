// Copyright 2021 gotomicro
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

package atomicx

import "fmt"

func ExampleNewValue() {
	val := NewValue[int]()
	data := val.Load()
	fmt.Println(data)
	// Output:
	// 0
}

func ExampleNewValueOf() {
	val := NewValueOf[int](123)
	data := val.Load()
	fmt.Println(data)
	// Output:
	// 123
}

func ExampleValue_Load() {
	val := NewValueOf[int](123)
	data := val.Load()
	fmt.Println(data)
	// Output:
	// 123
}

func ExampleValue_Store() {
	val := NewValueOf[int](123)
	data := val.Load()
	fmt.Println(data)
	val.Store(456)
	data = val.Load()
	fmt.Println(data)
	// Output:
	// 123
	// 456
}
