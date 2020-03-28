package main

import (
	"fmt"

	"github.com/nestorhq/nestor/internal/config"
)

func main() {
	config, err := config.ReadConfig("nestor.yml")
	if err != nil {
		panic(err)
	}
	fmt.Printf("config: %v\n", config)
}
