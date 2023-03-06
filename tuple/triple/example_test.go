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

package ekit_test

import (
	"fmt"
	. "github.com/gotomicro/ekit/tuple/triple"
)

func ExampleTriple_Copy() {
	triple := Triple{
		First:  1,
		Second: "one",
		Third:  "second"} // <1,"one","second">

	triple = triple.Copy(
		Triple{Third: "first"})

	fmt.Println(triple.ToString())

	// Output: <1,one,first>
}

func ExampleTriple_ToList() {
	triple := Triple{
		First:  1,
		Second: "one"}
	fmt.Println(triple.ToList())

	//Output: [1 one]
}

func ExampleTriple_ToString() {
	triple := Triple{
		First:  1,
		Second: "one"}
	fmt.Println(triple.ToString())

	//Output: <1,one>
}
