package main

import "fmt"

func main() {
	fmt.Println("Welcome to the Quiz Game!!")
	var name string
	fmt.Printf("Enter your name: ")
	fmt.Scan(&name)

	fmt.Printf("Hello %v, Welcome to the game!\n", name)

	fmt.Print("Enter your age: ")
	var age uint
    fmt.Scan(&age)

	fmt.Println(age >= 10)

	if age >= 10{
		fmt.Println("Yay you can play!!")
	}else{
		fmt.Println("You cannot play!!")
		return
	}

	score := 0
	num_questions := 2

	fmt.Printf("What is better, RTX 3080 or RTX 3090? ")
	var answer1 string
	var answer2 string
	fmt.Scan(&answer1, &answer2)

	if answer1 + " " + answer2  == "RTX 3090" || answer1 + " " + answer2  == "rtx 3090"{
		fmt.Println("Correct!!")
		score++
	}else{
		fmt.Println("Incorrect!!")
	}

	fmt.Printf("How many cores does the Ryzen 9 3900x have? ")
	var cores uint
	fmt.Scan(&cores)
	if cores == 12 {
		fmt.Println("Correct!!")
		score++
	}else{
		fmt.Println("Incorrect!!")
	}
	fmt.Printf("Your score is %v out %v\n", score, num_questions)

	percent := (float64(score) / float64(num_questions)) * 100
	
	fmt.Printf("Your scored %v%%.", percent)
}
