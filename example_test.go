package goquery

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
)

// This example scrapes the reviews shown on the home page of metalsucks.net.
func TestExample(t *testing.T) {
	// Request the HTML page.
	// res, err := http.Get("https://finance.sina.com.cn/stock/stockptd/2020-05-13/doc-iircuyvi2920672.shtml")
	res, err := http.Get("https://finance.sina.com.cn/stock/jsy/2020-05-14/doc-iircuyvi3071692.shtml")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(strings.TrimPrefix(TextWithTag(doc.Find(".article").Nodes), TextWithTag(doc.Find(".ct_hqimg").Nodes)))
	// log.Println(strings.TrimPrefix(TextWithTag(doc.Find(".article").Nodes), TextWithTag(doc.Find("blockquote").Nodes)))
	log.Println(TextWithTag(doc.Find(".article").Nodes))
	log.Println()
}

// This example shows how to use NewDocumentFromReader from a file.
func ExampleNewDocumentFromReader_file() {
	// create from a file
	f, err := os.Open("some/file.html")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	doc, err := NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}
	// use the goquery document...
	_ = doc.Find("h1")
}

// This example shows how to use NewDocumentFromReader from a string.
func ExampleNewDocumentFromReader_string() {
	// create from a string
	data := `
<html>
	<head>
		<title>My document</title>
	</head>
	<body>
		<h1>Header</h1>
	</body>
</html>`

	doc, err := NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	header := doc.Find("h1").Text()
	fmt.Println(header)

	// Output: Header
}
