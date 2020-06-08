package goquery

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/axgle/mahonia"
	"github.com/go-xweb/log"
	"golang.org/x/net/html"
)

var (
	classEnableMap = map[string]bool{
		// https://tieba.baidu.com/p/6728038381
		"d_post_content j_d_post_content  clearfix": true,
		// http://bbs.tianya.cn/post-stocks-2177796-1.shtml
		"bbs-content clearfix": true,
	}
	classDisableMap = map[string]bool{
		// http://mini.eastday.com/a/200515112914295.html
		"guess_like clear-fix bottom_cns": true,
		"footerDetail":                    true,
		"phone-code-box":                  true,
		"other-login-title":               true,
		"bind-form":                       true,
		"report-box ie_fixed":             true,
		"left_chunk ordernum1":            true,
		"mask":                            true,
		// http://finance.china.com.cn/stock/usstock/20200319/5226061.shtml
		"fr navr":                  true,
		"bottom":                   true,
		"AboutUs":                  true,
		"Map":                      true,
		"Contact":                  true,
		"code":                     true,
		"footer w1000 wauto hauto": true,
		// https://baijiahao.baidu.com/s?id=1666726079754462244&wfr=spider&for=pc
		"article-desc clearfix": true,
		// https://www.zhitongcaijing.com/content/detail/301770.html
		"float-r":  true,
		"detail-r": true,
		// https://www.a963.com/news/106521.shtml
		"right_title clearfix up_2": true,
		// https://www.prnasia.com/releases/
		"card-text-wrap": true,
		// https://tieba.baidu.com/p/6728038381
		"search_form":                  true,
		"sign_tip_bdwrap clearfix":     true,
		"nav_wrap nav_wrap_add_border": true,
		"p_favthr_main":                true,
		// http://bbs.tianya.cn/post-stocks-2177796-1.shtml
		"atl-location clearfix": true,
		"read-menu cf":          true,
		"action-tyf-shang":      true,
		"action-tyf-zan":        true,
		"post-div clearfix":     true,
		"foot":                  true,
		// http://mapp.jrj.com.cn/news/usstock/2020/05/27174029780301.shtml
		"time":            true,
		"suggest-reading": true,
		"shareBoxIcon":    true,
		// https://haokan.baidu.com/v?pd=wisenatural&vid=4438594192845039018
		"header-right float-right": true,
		"land-bottom":              true,
		"videoinfo-text clearfix":  true,
		// https://www.iqiyi.com/v_19rrlo3yvo.html
		"popup-con": true,
	}
)

// TextWithBr query text only br
func TextWithBr(nodes []*html.Node) string {
	var buf bytes.Buffer

	pnodes := TextSelectP(nodes)

	var f func(*html.Node)
	f = func(n *html.Node) {
		switch n.Type {
		case html.TextNode:
			data := strings.TrimSpace(n.Data)
			if data == "" {
				return
			}
			buf.WriteString(data)
			buf.WriteString("<br>")
			buf.WriteString("<br>")
			return
		case html.ElementNode:
			switch n.Data {
			case "img":
				for _, v := range n.Attr {
					if v.Key == "src" {
						buf.WriteString(fmt.Sprintf(`<img src="%v">`, v.Val))
						buf.WriteString("<br>")
						buf.WriteString("<br>")
					}
				}
			}
		}

		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	for _, n := range pnodes {
		f(n)
	}

	return buf.String()
}

//src为要转换的字符串
func coverGBKToUTF8(src string) string {
	// 网上搜有说要调用translate函数的，实测不用
	return mahonia.NewDecoder("gbk").ConvertString(src)
}

// TextSelectP select p nodes
func TextSelectP(nodes []*html.Node) []*html.Node {

	pnodes := []*html.Node{}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "p", "font":
				for _, v := range n.Attr {
					switch v.Key {
					case "class", "id":
						// log.Infof(v.Val)
						_, enable := classEnableMap[v.Val]
						if enable {
							pnodes = append(pnodes, n)
							return
						}
						_, disable := classDisableMap[v.Val]
						if disable {
							return
						}
					}

				}

				pnodes = append(pnodes, n)
				return
			case "div":
				for _, v := range n.Attr {
					switch v.Key {
					case "class", "id":
						// log.Infof(v.Val)
						_, enable := classEnableMap[v.Val]
						if enable {
							pnodes = append(pnodes, n)
							return
						}
						_, disable := classDisableMap[v.Val]
						if disable {
							return
						}
					}

				}
			case "figure", "footer":
				return
			}

		}
		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	for _, n := range nodes {
		f(n)
	}

	return pnodes
}

// TextOnlyPImg query text and retain some tag like <p>,<img>
func TextOnlyPImg(nodes []*html.Node) string {
	var buf bytes.Buffer

	var f func(*html.Node)
	f = func(n *html.Node) {
		data := ""

		log.Infof("n.Type:%v n.Data:[%v]", n.Type, n.Data)

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
			case "style", "script":
				return
			}
		case html.TextNode:
			data = strings.TrimSpace(n.Data)
			buf.WriteString(data)
		default:
			return
		}

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
