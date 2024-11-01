package main

import (
	"fmt"
	"os"
)

func nonmain() {
	fmt.Println("os.Exit()")
	os.Exit(0)
}

func main() {
	fmt.Println("os.Exit()")
	os.Exit(0) // want "call to os.Exit in main"
}
