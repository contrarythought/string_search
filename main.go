package main

import (
	"fmt"
	"os"
	"string_search/search"
)

func usage(Args []string) {
	fmt.Println(Args[0], " [starting search path] [string to lookup]")
	os.Exit(0)
}

func main() {
	if len(os.Args) != 3 {
		usage(os.Args)
	}

	files := search.Run(os.Args[1:])

	fmt.Println(files.Files)
}
