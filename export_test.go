package goquery

import (
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestTextWithBr(t *testing.T) {
	link := "https://www.roadshowchina.cn/Home/Meet/detail.html?mid=8170"
	rsp, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", rsp.StatusCode, rsp.Status)
	}

	// Load the HTML document
	doc, err := NewDocumentFromReader(rsp.Body)
	if err != nil {
		log.Fatal(err)
	}
	content := TextWithBr(doc.Find("body").Nodes)

	log.Printf("content:[%v]", content)
}

func TestTextWithTag(t *testing.T) {
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
	// content = strings.TrimPrefix(content, TextWithTag(doc.Find(".ct_hqimg").Nodes))
	// content = strings.TrimPrefix(content, TextWithTag(doc.Find("blockquote").Nodes))
	// content = strings.ReplaceAll(content, "<p>投资B站却巨亏？又一私募爆雷！B站暴涨150%，基金却跌剩3成！募集资金投向了哪儿？是否涉嫌违规承诺？</p>", "")

	log.Printf("content:[%v]", content)
}

func TestTextWithAllTag(t *testing.T) {
	res, err := http.Get("https://finance.sina.com.cn/wm/2020-05-21/doc-iircuyvi4252783.shtml")
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
	content := TextWithAllTag(doc.Find(".article").Nodes)
	content = strings.TrimPrefix(content, TextWithTag(doc.Find(".ct_hqimg").Nodes))
	content = strings.TrimPrefix(content, TextWithTag(doc.Find("blockquote").Nodes))
	content = strings.ReplaceAll(content, "<p>投资B站却巨亏？又一私募爆雷！B站暴涨150%，基金却跌剩3成！募集资金投向了哪儿？是否涉嫌违规承诺？</p>", "")
	content = strings.TrimPrefix(content, "\n")
	log.Printf("[%v]", len("\n"))
	log.Printf("[%v]", content[:1] == "\n")
	log.Printf("[%v]", content[:1])
	log.Printf("content:[%v]", content)
}

func TestTextSimple(t *testing.T) {
	link := "http://finance.china.com.cn/stock/usstock/20200319/5226061.shtml"
	res, err := http.Get(link)
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
	content := TextWithTag(doc.Find("p").Nodes)
	log.Printf("content:[%v]", content)
}
