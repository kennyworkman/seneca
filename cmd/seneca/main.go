package main

import (
	"fmt"
	"os"

	"github.com/kennyworkman/seneca/pkg/app"
	"github.com/kennyworkman/seneca/pkg/core"
)

func main() {

	if len(os.Args) > 2 {
		fmt.Printf("Too many args")
		return
	}

	arg := os.Args[1]

	fs := core.Filesystem{Root: "/Users/kenny/seneca"}
	if arg == "letters" || arg == "l" {
		app.ReadPaper(fs)
	} else {
		url := os.Args[1]
		app.AddPaper(url, fs)
	}

}
