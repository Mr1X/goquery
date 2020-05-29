package goquery

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/go-xweb/log"
	"golang.org/x/net/html"
)

// TextOnlyPImg query text and retain some tag like <p>,<img>
func TextOnlyPImg(nodes []*html.Node) string {
	var buf bytes.Buffer

	var f func(*html.Node)
	f = func(n *html.Node) {
		data := ""
		switch n.Type {
		case html.ElementNode:
			switch n.Data {
			case "p":
				buf.WriteString("<" + n.Data + ">")
			case "img":
				for _, v := range n.Attr {
					if v.Key == "src" {
						buf.WriteString(fmt.Sprintf(`<img src="%v">`, v.Val))
					}
				}
			}
		case html.TextNode:
			data = strings.TrimSpace(n.Data)
			buf.WriteString(data)
		default:
			return
		}

		log.Debugf("n.Type:%v n.Data:[%v]", n.Type, data)

		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		if n.Type == html.ElementNode {
			switch n.Data {
			case "p":
				buf.WriteString("</" + n.Data + ">")
			}
		}
	}
	for _, n := range nodes {
		f(n)
	}

	return buf.String()
}

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

		switch n.Data {
		case "style", "script", "iframe":
			return
		}

		data := strings.TrimSpace(n.Data)

		log.Debugf("n.Type:%v n.Data:[%v]", n.Type, data)

		if n.Type == html.ElementNode {
			switch n.Data {
			case "p", "font", "ul", "li", "tr", "td":
				buf.WriteString("<" + n.Data + ">")
			case "img":
				for _, v := range n.Attr {
					if v.Key == "src" {
						buf.WriteString(fmt.Sprintf(`<img src="%v">`, v.Val))
					}
				}
			}

		}
		if n.Type == html.TextNode {
			buf.WriteString(data)
		}

		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		if n.Type == html.ElementNode {
			switch n.Data {
			case "p", "font":
				buf.WriteString("</" + n.Data + ">")
			}
		}
	}
	for _, n := range nodes {
		f(n)
	}

	return buf.String()
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

// TextSimple query text
func TextSimple(nodes []*html.Node) string {
	var buf bytes.Buffer

	// Slightly optimized vs calling Each: no single selection object created
	var f func(*html.Node)
	f = func(n *html.Node) {
		// 注释的部分
		if n.Type == html.CommentNode {
			return
		}

		switch n.Parent.Data {
		case "style", "script":
			return
		}

		data := strings.TrimSpace(n.Data)

		log.Infof("n:%v", n)
		log.Infof("n.Type:%v n.DataAtom:%v n.Data:[%v] n.Namespace:%v", n.Type, n.DataAtom, data, n.Namespace)

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
