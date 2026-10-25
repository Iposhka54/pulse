package main

import (
	"fmt"
	"os"

	"github.com/Iposhka54/pulse/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
