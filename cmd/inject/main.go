package main

import (
	"fmt"
	"os"
)

func main() {
	if wd, err := os.Getwd(); err != nil {
		fmt.Print(err)
	} else {
		fmt.Print(wd)
	}
}
