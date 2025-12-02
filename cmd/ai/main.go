package main

import (
	"ai/internal/config"
	"fmt"
)

func main() {
	cfg, err := config.Load("./config.yaml")
	if err != nil {
		fmt.Println("err", err)
		return
	}
	fmt.Printf("%s", cfg)
}
