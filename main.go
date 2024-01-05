package main

import (
	"financial-monitoring/gold"
	"fmt"
)

func main() {
	goldService := gold.NewService()
	price, err := goldService.FetchGoldPrice()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	fmt.Printf("Current gold price: %d", price)

	panic("Implementing!!!!")
}
