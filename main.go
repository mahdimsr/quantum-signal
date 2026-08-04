package main

import (
	"fmt"
	"quantum-signal/scheduler"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file : %s", err)
	}

	fmt.Println("📶 Quantum signal started...")

	scheduler.Start(5)
}
