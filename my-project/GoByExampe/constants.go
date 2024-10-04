package main

import (
	"fmt"
	"math"
)

const s string = "constants" //const declares a constant value.

func main() {
	fmt.Println(s)

	const n = 500000000 //A const statement can appear anywhere a var statement can.

	const d = 3e20 / n  //Constant expressions perform arithmetic with arbitrary precision.

	fmt.Println(d)

	fmt.Println(int64(d)) //A numeric constant has no type until it’s given one, such as by an explicit conversion.

	fmt.Println(math.Sin(n)) //A number can be given a type by using it in a context that requires one, such as a variable assignment or function call. For example, here math.Sin expects a float64.
}