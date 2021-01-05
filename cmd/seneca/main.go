package main

import (
	"os"

	"github.com/kennyworkman/seneca/pkg/app"
	"github.com/kennyworkman/seneca/pkg/core"
)

func main() {

	// CLI flags.
	url := os.Args[1]

	fs := core.Filesystem{Root: "/Users/kenny/seneca"}
	app.ReadPaper(url, fs)
}
