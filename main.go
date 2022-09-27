package main

import (
	"fmt"
	"os"
	"string_search/search"
)

func usage(Args []string) {
	fmt.Println(Args[0], " <string to lookup>")
	os.Exit(0)
}

func main() {
	if len(os.Args) != 2 {
		usage(os.Args)
	}

	search.Run(os.Args[1])
}
