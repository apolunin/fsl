package main

import (
	"fmt"
	"os"

	"github.com/apolunin/fsl/engine"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("usage: fslengine <file1> <file2> ... <file_n>")
		return
	}

	if err := engine.Run(os.Args[1:]...); err != nil {
		fmt.Printf("error: %v", err)
	}
}
