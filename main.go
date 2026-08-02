package main

import (
	"fmt"
	"quantum-signal/scheduler"
)

func main() {
	fmt.Println("📶 Quantum signal started...")

	scheduler.Start(5)
}
