package main

import (
	"log"
	"os"

	"github.com/kennyworkman/seneca/pkg/app"
)

func main() {

	// CLI flags.
	url := os.Args[1]
	err := app.AddPaper(url)
	if err != nil {
		log.Fatal(err)
	}

}
