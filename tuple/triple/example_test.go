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
