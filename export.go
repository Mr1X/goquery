package goquery

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// TextWithTag query text and retain some tag like <p>,<img>
func TextWithTag(nodes []*html.Node) string {
	var buf bytes.Buffer

	// Slightly optimized vs calling Each: no single selection object created
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type != html.TextNode {
			switch n.Data {
			case "p":
				buf.WriteString("<p>")
			case "img":
				for _, v := range n.Attr {
					if v.Key == "src" {
						buf.WriteString(fmt.Sprintf(`<img src="%v">`, strings.TrimPrefix(v.Val, "//")))
					}
				}
			}

		} else {
			buf.WriteString(n.Data)
		}

		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		if n.Type != html.TextNode {
			switch n.Data {
			case "p":
				buf.WriteString("</p>")
			}
		}
	}
	for _, n := range nodes {
		f(n)
	}

	return buf.String()
}
