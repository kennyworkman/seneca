package core

import "strings"

type Paper struct {
	DOI     string
	RawData []byte
}

type PaperRepository interface {
	GetPaper(paper *Paper) error
	AddPaper(paper *Paper) error
	DeletePaper(paper *Paper) error
}

// DOI syntax - https://www.doi.org/doi_handbook/2_Numbering.html#2.2
func (p Paper) splitDOI() (string, string) {
	splitDOI := strings.Split(p.DOI, "/")
	return splitDOI[0], splitDOI[1]
}

func (p Paper) pdfFile() string {
	_, postfix := p.splitDOI()
	return postfix + ".pdf"
}

func (p Paper) noteFile() string {
	_, postfix := p.splitDOI()
	return postfix + "_notes.md"
}
