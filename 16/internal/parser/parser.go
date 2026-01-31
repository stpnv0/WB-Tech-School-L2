package parser

import (
	"bytes"
	"io"

	"golang.org/x/net/html"
)

type LinkInfo struct {
	URL       string
	LinkType  string // "page", "css", "image"
	Tag       string // "a", "link"
	Attribute string // "href", "src"
}

type TagConfig struct {
	Attributes []string
	LinkType   string
}

var tagConfigs = map[string]TagConfig{
	"a":      {[]string{"href"}, "page"},
	"link":   {[]string{"href"}, "css"},
	"script": {[]string{"src"}, "js"},
	"img":    {[]string{"src", "srcset"}, "image"},
}

type ParseResult struct {
	Links        []LinkInfo // all links with metadata
	HTMLContent  []byte     // original
	ModifiedHTML []byte     // modified
}

func Parse(HTMLContent []byte, urlToLocalMap map[string]string) (*ParseResult, error) {
	li, err := extractLinks(HTMLContent)
	if err != nil {
		return nil, err
	}

	// Если нет маппинга - возвращаем оригинальный HTML без замены ссылок
	var modifiedHTML []byte
	if urlToLocalMap != nil && len(urlToLocalMap) > 0 {
		modifiedHTML, err = replaceLinks(HTMLContent, urlToLocalMap)
		if err != nil {
			return nil, err
		}
	} else {
		modifiedHTML = HTMLContent
	}

	return &ParseResult{
		Links:        li,
		HTMLContent:  HTMLContent,
		ModifiedHTML: modifiedHTML,
	}, nil
}

func extractLinks(HTML []byte) ([]LinkInfo, error) {
	var links []LinkInfo

	r := bytes.NewReader(HTML)
	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			if z.Err() == io.EOF {
				break
			}
			return nil, z.Err()
		}

		if tt == html.StartTagToken || tt == html.SelfClosingTagToken {
			token := z.Token()
			if cfg, ok := tagConfigs[token.Data]; ok {
				for _, attr := range token.Attr {
					if contains(cfg.Attributes, attr.Key) {
						if attr.Val == "" {
							continue
						}

						links = append(links, LinkInfo{
							URL:       attr.Val,
							LinkType:  cfg.LinkType,
							Tag:       token.Data,
							Attribute: attr.Key,
						})
					}
				}
			}
		}
	}

	return links, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}

	return false
}

func replaceLinks(data []byte, urlToLocalMap map[string]string) ([]byte, error) {
	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	replaceLinksInNode(doc, urlToLocalMap)

	var buf bytes.Buffer
	err = html.Render(&buf, doc)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func replaceLinksInNode(n *html.Node, urlToLocalMap map[string]string) {
	if n.Type == html.ElementNode {
		for i, attr := range n.Attr {
			if attr.Key == "href" || attr.Key == "src" {
				if localPath, ok := urlToLocalMap[attr.Val]; ok {
					n.Attr[i].Val = localPath
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		replaceLinksInNode(c, urlToLocalMap)
	}
}
