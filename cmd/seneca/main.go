package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func main() {

	// CLI flags.
	url := os.Args[1]
	err := addPaper(url)
	if err != nil {
		log.Fatal(err)
	}

}

func addPaper(url string) error {

	// Turn url into new doi labeled dir
	mirrors := getSciHubURLs()
	pdfURI, doi := getPDFSource(mirrors, url)

	resp, err := http.Get(pdfURI)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// DOI syntax - https://www.doi.org/doi_handbook/2_Numbering.html#2.2
	splitDOI := strings.Split(doi, "/")
	prefix, postfix := splitDOI[0], splitDOI[1]

	// Create prefix dir if not exist

	// Create postfix dir if not exist
	out, err := os.Create(postfix)
	if err != nil {
		return err
	}
	defer out.Close()

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	print(n, prefix)

	return nil
}

func readPaper(doi string) error {
	return nil
}

func searchPaper(query string) error {
	return nil
}

// https://pkg.go.dev/golang.org/x/net/html#example-Parse
func getSciHubURLs() []string {
	resp, err := http.Get("https://sci-hub.now.sh")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	tokenized, err := html.Parse(strings.NewReader(string(body)))

	// Return list sorted topologically
	var f func(*html.Node)
	var mirrors []string
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" && strings.Contains(a.Val, "https://sci-hub") {
					mirrors = append(mirrors, a.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(tokenized)
	return mirrors
}

func getPDFSource(mirrors []string, url string) (string, string) {
	candidate := mirrors[0] + url

	resp, err := http.Get(candidate)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	// find doi in body
	doi := findDOI(body)

	tokenized, err := html.Parse(strings.NewReader(string(body)))
	var iframeSearch func(*html.Node)
	var pdfURI string
	iframeSearch = func(n *html.Node) {
		if n.Data == "iframe" {
			for _, iframe := range n.Attr {
				if iframe.Key == "src" {
					pdfURI = iframe.Val
					return
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			iframeSearch(c)
		}

	}

	iframeSearch(tokenized)

	if pdfURI != "" {
		pdfURI = "https:" + pdfURI
	}

	return pdfURI, doi
}

func findDOI(body []byte) string {
	reader := bytes.NewReader(body)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {

		line := scanner.Text()
		if strings.Contains(line, "var doi") {
			// Pull the DOI out of js variable assignment.
			// var doi = '<doi>';
			r, _ := regexp.Compile("'(.*?)'")
			return strings.Trim(r.FindString(line), "'")
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return ""
}

type Paper struct {
	DOI string
}
