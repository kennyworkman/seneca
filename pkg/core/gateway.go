package core

import (
	"os"
	"os/exec"
	"path/filepath"
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
	grep := exec.Command("grep", "-r", "--include", "*.txt", query, root)
	return grep.Output()
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

	// todo add note boilerplate
	_, err = note.Write([]byte("metadata"))
	if err != nil {
		return err
	}

	// Create txt representation of pdf for grep

	return nil
}

func (f Filesystem) GetPaper(paper *Paper) error {
	return nil
}

func (f Filesystem) DeletePaper(paper *Paper) error {
	return nil
}

type Database struct {
}
