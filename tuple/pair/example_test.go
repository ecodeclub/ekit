package ekit_test

import (
	"fmt"
	. "github.com/gotomicro/ekit/tuple/pair"
)

func ExamplePair_Copy() {
	pair := Pair{
		First:  1,
		Second: "two"} // <1,"two">
	pair = pair.Copy(
		Pair{Second: "one"})
	fmt.Println(pair.ToString())

	// Output: <1,one>
}

func ExamplePair_ToList() {
	pair := Pair{
		First:  1,
		Second: "one"}
	fmt.Println(pair.ToList())

	//Output: [1 one]
}

func ExamplePair_ToString() {
	pair := Pair{
		First:  1,
		Second: "one"}
	fmt.Println(pair.ToString())

	//Output: <1,one>
}
