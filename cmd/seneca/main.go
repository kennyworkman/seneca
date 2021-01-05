package main

import (
	"github.com/kennyworkman/seneca/pkg/app"
	"github.com/kennyworkman/seneca/pkg/core"
)

func main() {

	// CLI flags.
	// url := os.Args[1]

	fs := core.Filesystem{Root: "/Users/kenny/seneca"}
	// app.AddPaper(url, fs)
	app.ReadPaper(fs)
}
