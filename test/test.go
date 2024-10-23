package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Item struct {
	index int
	name  string
}

func main() {
	fmt.Println("Please Choose: ")

	var items []Item
	items = append(items, Item{
		index: 1,
		name:  "apple",
	})
	items = append(items, Item{
		index: 2,
		name:  "banana",
	})
	items = append(items, Item{
		index: 3,
		name:  "cherry",
	})

	for _, item := range items {
		fmt.Printf("{%d} : %s\n", item.index, item.name)
	}

	var chose string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter item index: ")
	if scanner.Scan() {
		chose = scanner.Text()
	}

	// Convert the chosen index to an integer
	chosenIndex, err := strconv.Atoi(chose)
	if err != nil {
		fmt.Println("Invalid input. Please enter a number.")
		return
	}

	for _, item := range items {
		if chosenIndex == item.index {
			fmt.Printf("You chose %s\n", item.name)
			break
		}
	}
}
