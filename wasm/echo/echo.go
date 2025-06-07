package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	out := make([]string, 0, len(os.Args)-1)
	for _, arg := range os.Args[1:] {
		out = append(out, fmt.Sprint(arg))
	}

	fmt.Printf("%s\n", strings.Join(out, " "))
}
