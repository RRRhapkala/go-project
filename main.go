package main

import "fmt"

func meshupSlices(one, two []int) {
	if len(one) == len(two) {
		for i := 0; i < len(one); {
			fmt.Println(one, "\n", two)
			variable := one[i]
			one[i] = two[i]
			two[i] = variable
			i += 2
		}
	}

}

func main() {
	fmt.Println("Go Project")
	meshupSlices([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 2})
}
