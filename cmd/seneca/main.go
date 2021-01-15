package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

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

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fs := core.Filesystem{Root: usr.HomeDir}

	if arg == "letters" || arg == "l" {

		app.ReadPaper(fs)

	} else if arg == "raw" {

		url, doi, err := rawPrompt()
		if err != nil {
			log.Fatal(err)
		}

		paper, err := app.AddPaperRaw(url, doi, fs)
		if err != nil {
			log.Fatal(err)
		}

		err = confirmPrompt()
		if err != nil {
			log.Fatal(err)
		} else {
			fs.ReadPaper(paper)
		}

	} else {

		url := os.Args[1]
		paper, err := app.AddPaper(url, fs)
		if err != nil {
			log.Fatal(err)
		}

		err = confirmPrompt()
		if err != nil {
			log.Fatal(err)
		} else {
			fs.ReadPaper(paper)
		}
	}

}

// Utility to prompt for raw paper input
func rawPrompt() (string, string, error) {

	var url, doi string

	prompt := promptui.Prompt{
		Label: "url",
	}

	url, err := prompt.Run()

	if err != nil {
		return "", "", err

	} else {

		prompt := promptui.Prompt{
			Label: "doi",
		}

		doi, err = prompt.Run()
		if err != nil {
			return "", "", err
		} else {
			return url, doi, nil
		}
	}

}

// Utility to present confirm prompt
func confirmPrompt() error {

	prompt := promptui.Prompt{
		Label:     "access",
		IsConfirm: true,
	}

	_, err := prompt.Run()
	return err

}
