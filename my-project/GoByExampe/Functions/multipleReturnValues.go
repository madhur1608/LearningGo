package main

func vals() (int, int) {
	return 3, 7
}

func main() {
	a, b := vals()

	println(a)
	println(b)

	_, c := vals()

	println(c)
}