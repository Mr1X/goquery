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
	res, err := http.Get("https://finance.sina.com.cn/stock/stockptd/2020-05-13/doc-iircuyvi2920672.shtml")
	// res, err := http.Get("https://finance.sina.com.cn/stock/jsy/2020-05-14/doc-iircuyvi3071692.shtml")
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
	content := TextWithTag(doc.Find(".article").Nodes)
	content = strings.TrimPrefix(content, TextWithTag(doc.Find(".ct_hqimg").Nodes))
	content = strings.TrimPrefix(content, TextWithTag(doc.Find("blockquote").Nodes))
	content = strings.ReplaceAll(content, "<p>投资B站却巨亏？又一私募爆雷！B站暴涨150%，基金却跌剩3成！募集资金投向了哪儿？是否涉嫌违规承诺？</p>", "")
	content = strings.TrimPrefix(content, "\n")
	log.Printf("[%v]", len("\n"))
	log.Printf("[%v]", content[:1] == "\n")
	log.Printf("[%v]", content[:1])
	log.Printf("content:[%v]", content)
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
