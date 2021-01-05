package app

import (
	"errors"
	"fmt"

	"github.com/kennyworkman/seneca/pkg/core"
	"github.com/manifoldco/promptui"
)

// Retrieve a paper and associated letter (note buffer)
func ReadPaper(query string, fs core.Filesystem) {

	// Interactive terminal list of results from search query
	// https://stackoverflow.com/questions/41037870/go-exec-command-run-command-which-contains-pipe
	validate := func(input string) error {

		// filesystem search ?, using root
		res, _ := fs.RawSearch(input)

		fmt.Printf("\n%+v\n", string(res))

		if len(res) == 0 {
			return errors.New("No results")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "grep",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)

	// Retrieve / open associated paper and note buffer
}
