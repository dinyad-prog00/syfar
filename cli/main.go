package main

import (
	"fmt"
	"os"
	"syfar/cli/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "syfar: Error\n%v\n", err)
		os.Exit(2)
	}
}
