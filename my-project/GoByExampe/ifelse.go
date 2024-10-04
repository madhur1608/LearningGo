package main

import "fmt"

func main() {
	if 7 % 2 == 0 {						//Here’s a basic example.
		fmt.Println("7 is even number")
	} else {
		fmt.Println("7 is odd number")
	}

	if 8 % 4 == 0 {						//You can have an if statement without an else.
		fmt.Println("8 is divisible by 4")
	}

	if 8 % 2 == 0|| 7 % 2 == 0 {		//You can have an if statement without an else.
		fmt.Println("either 8 or 7 are even")
	}

	if num := 9; num < 0 {				//A statement can precede conditionals; any variables declared in this statement are available in the current and all subsequent branches.
		fmt.Println(num, "is negative")
	} else if num < 10 {
		fmt.Println(num, "has 1 digit")
	} else {
		fmt.Println(num, "has multiple digits")
	}
}

// Note that you don’t need parentheses around conditions in Go, but that the braces are required.