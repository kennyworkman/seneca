package core

import "strings"

type Paper struct {
	DOI string
	// Head results (first few lines) on this paper
	Head string
	// Grep results for some query
	Grep    string
	RawData []byte
}

type PaperRepository interface {
	GetPaper(paper *Paper) error
	AddPaper(paper *Paper) error
	DeletePaper(paper *Paper) error
}

// Creates Paper object from DOI
func NewPaper(doiPrefix, doiPostfix string) *Paper {
	return &Paper{DOI: doiPrefix + "/" + doiPostfix}
}

// DOI syntax - https://www.doi.org/doi_handbook/2_Numbering.html#2.2
// Returns a (prefix, postfix) pair.
// The prefix can itself have multiple nested identifiers in pathological cases:
//	ie. 10.1093/bioinformatics/btz781 -> (10.1093/bioinformatics, btz781).
func (p Paper) splitDOI() (string, string) {
	splitDOI := strings.Split(p.DOI, "/")
	return strings.Join(splitDOI[:len(splitDOI)-1], "/"), splitDOI[len(splitDOI)-1]
}

func (p Paper) pdfFile() string {
	_, postfix := p.splitDOI()
	return postfix + ".pdf"
}

func (p Paper) txtFile() string {
	_, postfix := p.splitDOI()
	return postfix + ".txt"
}

func (p Paper) noteFile() string {
	_, postfix := p.splitDOI()
	return postfix + "_notes.md"
}

// Returns first several lines of text representation.
func (p Paper) GetHead(fs Filesystem) (string, error) {
	prefix, postfix := p.splitDOI()
	return fs.fileHead(prefix, postfix)
}

// Returns query grepped on this paper alone.
func (p Paper) GetGrep(query string, fs Filesystem) (string, error) {
	prefix, postfix := p.splitDOI()
	return fs.fileGrep(query, prefix, postfix)
}
