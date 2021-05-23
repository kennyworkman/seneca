package app

import (
	"errors"
	"fmt"
	"log"

	"github.com/kennyworkman/seneca/pkg/core"
	"github.com/manifoldco/promptui"
)

// Retrieve a paper and associated letter (note buffer)
func ReadPaper(fs core.Filesystem) {

	// Interactive terminal list of results from search query
	// https://stackoverflow.com/questions/41037870/go-exec-command-run-command-which-contains-pipe
	validate := func(input string) error {

		// filesystem search ?, using root
		res, _ := fs.RawSearch(input)

		if len(res) == 0 {
			return errors.New("No results")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "grep",
		Validate: validate,
	}

	query, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	} else {

		res, err := fs.SearchAndParse(query)
		if err != nil {
			log.Fatal(err)
		}

		templates := &promptui.SelectTemplates{
			Inactive: "{{ .DOI }}",
			Active:   " {{ .DOI | cyan }}",
			Details:  "\n{{ .Head | blue }}\n{{ .Grep | red}}",
		}

		// Retrieve / open associated paper and note buffer
		selectPrompt := promptui.Select{
			Label:     "Select Paper",
			Items:     res,
			Templates: templates,
		}

		selectIdx, _, err := selectPrompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		fs.ReadPaper(res[selectIdx])
	}
}
