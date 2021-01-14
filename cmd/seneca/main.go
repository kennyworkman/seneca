package main

import (
	"fmt"
	"log"
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
		paper, err := app.AddPaper(url, fs)
		if err != nil {
			log.Fatal(err)
		}
		fs.ReadPaper(paper)
	}

}
