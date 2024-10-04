package main

import (
	"fmt"
)

func main() {
	i := 1 				//The most basic type, with a single condition.
	for i <= 3 {
		fmt.Println(i)
		i += 1
	}

	for j := 0; j < 3; j++ {	//A classic initial/condition/after for loop.
		fmt.Println(j)
	}

	for i := range 3 {			//Another way of accomplishing the basic “do this N times” iteration is range over an integer.
		fmt.Println("range", i)
	}

	for {						//for without a condition will loop repeatedly until you break out of the loop or return from the enclosing function.
		fmt.Println("loop")
		break
	}

	for n := range 6 {			//You can also continue to the next iteration of the loop.
		if n % 2 == 0 {
			continue
		}
		fmt.Println(n)
	}
}