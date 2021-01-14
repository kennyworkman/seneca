package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kennyworkman/seneca/pkg/app"
	"github.com/kennyworkman/seneca/pkg/core"
	"github.com/manifoldco/promptui"
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

		prompt := promptui.Prompt{
			Label:     "access",
			IsConfirm: true,
		}

		_, err = prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		} else {
			fs.ReadPaper(paper)
		}
	}

}
