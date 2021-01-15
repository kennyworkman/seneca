package app

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/kennyworkman/seneca/pkg/core"
	"golang.org/x/net/html"
)

// Adds a paper represented by URL to seneca.
func AddPaper(url string, fs core.Filesystem) (*core.Paper, error) {

	// Easy to recover the DOI from sci-hub's landing html.
	// Need DOI for persistence scheme.
	// Therefore need hardcoded recovery of DOI and pdf uri from:
	//     * bioarxiv
	//     * arxiv

	var doi, pdfURI string
	var err error
	if strings.Contains(url, "arxiv.org") {

		fmt.Printf("\nRecognized arxiv paper.")

		// Recover DOI from arxiv URL
		splitURL := strings.Split(url, "/")
		doi = "arxiv/" + splitURL[len(splitURL)-1]

		pdfURI = strings.Replace(url, "abs", "pdf", 1) + ".pdf"

	} else if strings.Contains(url, "biorxiv.org") {

		fmt.Printf("\nRecognized biorxiv paper.")

		// Recover DOI from biorxiv URL
		splitURL := strings.Split(url, "/")
		doi = splitURL[len(splitURL)-2] + "/" + splitURL[len(splitURL)-1]

		pdfURI = url + ".pdf"

	} else {

		// Check list of mirrors for sci-hub thats up.
		// Scrape from there.
		mirrors := getSciHubURLs()
		if len(mirrors) == 0 {
			return nil, fmt.Errorf("No scihub mirrors found.")
		}
		fmt.Printf("\nAttempting download from sci-hub mirrors.")
		pdfURI, doi, err = getPDFSource(mirrors, url)
		if err != nil {
			return nil, err
		}
	}

	fmt.Printf("\nFetching paper binary...")
	client := &http.Client{}
	req, err := http.NewRequest("GET", pdfURI, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	pdfBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	paper := &core.Paper{DOI: doi, RawData: pdfBytes}
	fs.AddPaper(paper)

	return paper, nil
}

func AddPaperRaw(url, doi string, fs core.Filesystem) (*core.Paper, error) {

	fmt.Printf("\nFetching paper binary...")
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	pdfBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	paper := &core.Paper{DOI: doi, RawData: pdfBytes}
	fs.AddPaper(paper)

	return paper, nil
}

// https://pkg.go.dev/golang.org/x/net/html#example-Parse
// Retrieve list of mirrors from url
func getSciHubURLs() (mirrors []string) {
	resp, err := http.Get("https://sci-hub.now.sh")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	tokenized, err := html.Parse(strings.NewReader(string(body)))

	// Return list sorted topologically
	var f func(*html.Node)
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

// TODO: Tries from list of mirrors until accessible pdf URI and doi successfully pulled.
func getPDFSource(mirrors []string, url string) (pdfURI string, doi string, err error) {
	candidate := mirrors[0] + url

	resp, err := http.Get(candidate)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	// find doi in body
	doi = findDOI(body)

	tokenized, err := html.Parse(strings.NewReader(string(body)))
	var iframeSearch func(*html.Node)
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

	if len(pdfURI) == 0 {
		return "", "", fmt.Errorf("Unable to recover paper with scihub")
	}

	if pdfURI[:6] != "https:" {
		pdfURI = "https:" + pdfURI
	}

	return pdfURI, doi, nil
}

// TODO: helper method to parse sci-hub page for doi
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
