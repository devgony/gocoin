package main

import "fmt"

func plus(a ...int) int {
	var total int
	for index, item := range a {
		total += item
	}
	return total
}

func main() {
	result := plus(2, 3, 4, 5, 6)
	fmt.Println(result)
}
