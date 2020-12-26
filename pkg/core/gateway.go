package core

import (
	"os"
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

func (f Filesystem) GetPaper(paper *Paper) error {
	return nil
}

func (f Filesystem) AddPaper(paper *Paper) error {
	prefix, postfix := paper.splitDOI()

	// Mkdir is safe.
	os.Mkdir(filepath.Join(f.Root, prefix), 0700)
	os.Mkdir(filepath.Join(f.Root, prefix, postfix), 0700)

	// create notes buffer and pdf file

	pdfFile := paper.pdfFile()
	noteFile := paper.noteFile()
	os.Chdir(filepath.Join(f.Root, prefix, postfix))

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

	// n, err := io.Copy()
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (f Filesystem) DeletePaper(paper *Paper) error {
	return nil
}

type Database struct {
}
