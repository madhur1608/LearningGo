package main

import "fmt"

func fact(n int) int {
	if n == 0 {
		return 1
	}
	return n * fact(n-1)  // 7 * 7-1 * 7-2 * 7-3 * 7-4 * 7-5 * 7-6 == 5040
}

func main() {
	fmt.Println(fact(7))


	var fib func(n int) int

	fib = func(n int) int {
		if n < 2 {
			return n
		}
		return fib(n-1) + fib(n-2)  // 7th fibonacci number = 13
	}
	fmt.Println(fib(7))
}