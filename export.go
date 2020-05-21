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
		// 注释的部分
		if n.Type == html.CommentNode {
			return
		}

		if n.Parent.Data == "style" {
			return
		}

		data := strings.TrimSpace(n.Data)

		// log.Infof("n.Type:%v n.Data:[%v]", n.Type, data)

		if n.Type != html.TextNode {
			switch n.Data {
			case "p":
				buf.WriteString("\n<p>")
			case "img":
				for _, v := range n.Attr {
					if v.Key == "src" {
						// buf.WriteString(fmt.Sprintf(`<img src="%v">`, strings.TrimPrefix(v.Val, "//")))
						img := v.Val
						if !strings.HasPrefix(img, "http") {
							img = "//" + img
						}
						buf.WriteString(fmt.Sprintf(`<img src="%v">`, img))
					}
				}
			}

		} else {
			// buf.WriteString(n.Data)
			buf.WriteString(data)
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

	return strings.TrimPrefix(buf.String(), "\n")
}

// TextWithAllTag query text and retain all tag
func TextWithAllTag(nodes []*html.Node) string {
	var buf bytes.Buffer

	// Slightly optimized vs calling Each: no single selection object created
	var f func(*html.Node)
	f = func(n *html.Node) {
		// 注释的部分
		if n.Type == html.CommentNode {
			return
		}

		if n.Parent.Data == "style" {
			return
		}

		data := strings.TrimSpace(n.Data)

		// log.Infof("n.Type:%v n.Data:[%v]", n.Type, data)

		if n.Type != html.TextNode {
			switch n.Data {
			case "img":
				for _, v := range n.Attr {
					if v.Key == "src" {
						img := v.Val
						if strings.HasPrefix(img, "//") {
							img = strings.Replace(img, "//", "https://", 1)
						}
						buf.WriteString(fmt.Sprintf(`<img src="%v">`, img))
					}
				}
			default:
				buf.WriteString("<" + n.Data + ">")
			}

		} else {
			buf.WriteString(data)
		}

		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		if n.Type != html.TextNode {
			switch n.Data {
			case "img":
				return
			default:
				buf.WriteString("</" + n.Data + ">")
			}
		}
	}
	for _, n := range nodes {
		f(n)
	}

	return buf.String()
}
