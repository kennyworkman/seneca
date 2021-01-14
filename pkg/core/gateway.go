package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Gateway interface {
	GetPaper(paper *Paper) error
	AddPaper(paper *Paper) error
	DeletePaper(paper *Paper) error
}

type Filesystem struct {
	Root string
}

func (f Filesystem) GetRoot() string {
	return filepath.Join(f.Root, "/.seneca")
}

// Wraps grep to search through .txt representation of pdfs without parsing results
func (f Filesystem) RawSearch(query string) ([]byte, error) {
	root := filepath.Join(f.Root, "/.seneca")
	grep := exec.Command("grep", "-ir", "--include", "*.txt", query, root)
	return grep.Output()
}

// Runs grep on filesystem and derives Paper objects from result
func (f Filesystem) SearchAndParse(query string) ([]*Paper, error) {

	root := filepath.Join(f.Root, "/.seneca")
	grep := exec.Command("grep", "-lir", "--include", "*.txt", query, root)

	filesText, err := grep.Output()
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(string(filesText)), "\n")
	papers := []*Paper{}
	for _, file := range files {
		paper := txtFilePathtoPaper(file)
		// Hacky to conform to Field over Method interface of promptui.
		head, err := paper.GetHead(f)
		if err != nil {
			return nil, err
		}

		grep, err := paper.GetGrep(query, f)
		if err != nil {
			return nil, err
		}
		paper.Head = head
		paper.Grep += grep
		papers = append(papers, paper)
	}

	return papers, nil
}

// Utility to convert .txt file path to Paper object
func txtFilePathtoPaper(path string) *Paper {
	splitPath := strings.Split(path, "/")
	prefix, postfix := splitPath[len(splitPath)-2], splitPath[len(splitPath)-1]
	return NewPaper(prefix, postfix[:len(postfix)-4])
}

func (f Filesystem) AddPaper(paper *Paper) error {

	if f.ExistsPaper(paper) {
		fmt.Printf("\nPaper already exists in seneca")
		return nil
	}

	prefix, _ := paper.splitDOI()

	// Mkdir is safe. Won't overwrite existing folders/files.
	os.MkdirAll(filepath.Join(f.GetRoot(), prefix), 0700)

	pdfFile := paper.pdfFile()
	noteFile := paper.noteFile()
	os.Chdir(filepath.Join(f.GetRoot(), prefix))

	pdf, err := os.Create(pdfFile)
	if err != nil {
		return err
	}
	defer pdf.Close()

	note, err := os.Create(noteFile)
	if err != nil {
		return err
	}
	defer note.Close()

	_, err = pdf.Write(paper.RawData)
	if err != nil {
		return err
	}

	// Create txt representation of pdf for grep.
	_, err = exec.Command("pdftotext", pdfFile).Output()
	if err != nil {
		return err
	}

	head, err := exec.Command("head", paper.txtFile()).Output()
	if err != nil {
		return err
	}

	// Add all note metadata + boilerplate here.
	head = append(head, []byte("\n----")...)
	_, err = note.Write(head)
	if err != nil {
		return err
	}

	return nil
}

func (f Filesystem) ExistsPaper(paper *Paper) bool {
	prefix, _ := paper.splitDOI()

	// Mkdir is safe. Won't overwrite existing folders/files.
	os.MkdirAll(filepath.Join(f.GetRoot(), prefix), 0700)

	pdfFile := paper.pdfFile()
	os.Chdir(filepath.Join(f.GetRoot(), prefix))

	if _, err := os.Stat(pdfFile); err == nil {
		return true
	}
	return false
}

func (f Filesystem) ReadPaper(paper *Paper) error {

	prefix, _ := paper.splitDOI()
	os.Chdir(filepath.Join(f.GetRoot(), prefix))

	// Write current date as a header
	note, err := os.OpenFile(paper.noteFile(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer note.Close()
	dateHeader := "\n\n## " + time.Now().Format("Mon Jan 2")
	_, err = note.Write([]byte(dateHeader))
	if err != nil {
		return err
	}

	// Open pdf
	exec.Command("zathura", paper.pdfFile()).Start()

	// Emulate terminal and open vim
	// https://stackoverflow.com/questions/21513321/how-to-start-vim-from-go
	vim := exec.Command("vim", paper.noteFile())
	vim.Stdin = os.Stdin
	vim.Stdout = os.Stdout
	err = vim.Run()
	if err != nil {
		return err
	}

	return nil

}

func (f Filesystem) DeletePaper(paper *Paper) error {
	return nil
}

func (f Filesystem) fileHead(prefix, postfix string) (string, error) {
	root := filepath.Join(f.Root, "/.seneca")
	path := filepath.Join(root, prefix, postfix+".txt")

	head, err := exec.Command("head", path).Output()
	if err != nil {
		return "", err
	}
	return string(head), nil
}

func (f Filesystem) fileGrep(query, prefix, postfix string) (string, error) {
	root := filepath.Join(f.Root, "/.seneca")
	path := filepath.Join(root, prefix, postfix+".txt")

	grepOut, err := exec.Command("grep", "-m", "10", "-in", query, path).Output()
	if err != nil {
		return "", err
	}
	return string(grepOut), nil
}
