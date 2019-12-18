package goexamples

import (
	"fmt"
	"math"
)

// func main() {
// 	//goValues();
// 	//goVariables();
// 	goConstants()

// }
func GoValues() {
	fmt.Println("hai")
	fmt.Println("1+1 =", 1+1)
	fmt.Println(true && false)
	fmt.Println(true || false)
	fmt.Println(!true)
}
func goVariables() {
	var a int = 1
	fmt.Println(a)
	var b = 2
	fmt.Println(b)
	var c, d int = 3, 4
	fmt.Println(c, d)
	e := "sai"
	fmt.Println(e)
	var f string = "everest"
	fmt.Println(f)
	var g bool = true
	fmt.Println(g)

}
func goConstants() {
	// A `const` statement can appear anywhere a `var`
	// statement can.
	const n = 500000000

	// Constant expressions perform arithmetic with
	// arbitrary precision.
	const d = 3e20 / n
	fmt.Println(d)

	// A numeric constant has no type until it's given
	// one, such as by an explicit conversion.
	fmt.Println(int64(d))

	// A number can be given a type by using it in a
	// context that requires one, such as a variable
	// assignment or function call. For example, here
	// `math.Sin` expects a `float64`.
	fmt.Println(math.Sin(n))
}
