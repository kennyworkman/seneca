package core

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	grep := exec.Command("grep", "-lr", "--include", "*.txt", query, root)

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
		paper.Detail = head
		paper.Detail += grep
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

	// Add all note metadata + boilerplate here.
	_, err = note.Write([]byte("metadata"))
	if err != nil {
		return err
	}

	// Create txt representation of pdf for grep.
	_, err = exec.Command("pdftotext", pdfFile).Output()
	if err != nil {
		return err
	}

	return nil
}

func (f Filesystem) GetPaper(paper *Paper) error {
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
