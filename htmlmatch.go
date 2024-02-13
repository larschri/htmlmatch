package htmlmatch

import (
	"bytes"
	"io"
	"strings"

	"golang.org/x/net/html"
)

func elementMatch(full, pattern *html.Node) bool {
	if full.Type != pattern.Type {
		return false
	}
	if pattern.Type == html.TextNode && strings.HasPrefix(pattern.Data, "substring:") {
		if !strings.Contains(full.Data, strings.TrimPrefix(pattern.Data, "substring:")) {
			return false
		}
	} else if pattern.Type == html.TextNode && strings.HasPrefix(pattern.Data, "verbatim:") {
		if full.Data != strings.TrimPrefix(pattern.Data, "verbatim:") {
			return false
		}
	} else {
		if strings.TrimSpace(full.Data) != strings.TrimSpace(pattern.Data) {
			return false
		}
	}
outer:
	for i := range pattern.Attr {
		for j := range full.Attr {
			if pattern.Attr[i] == full.Attr[j] {
				continue outer
			}
		}
		return false
	}
	return true
}

func skipWhitespace(n *html.Node) *html.Node {
	if n != nil && n.Type == html.TextNode && strings.TrimSpace(n.Data) == "" {
		return n.NextSibling
	}
	return n
}

func containsTree(full, pattern *html.Node) *html.Node {
	pattern = skipWhitespace(pattern)
	for full != nil && pattern != nil {
		if elementMatch(full, pattern) && containsTree(full.FirstChild, pattern.FirstChild) == nil {
			pattern = skipWhitespace(pattern.NextSibling)
		} else {
			pattern = containsTree(full.FirstChild, pattern)
		}
		full = full.NextSibling
	}
	return pattern
}

// ContainsTree returns true if pattern is contained in the full tree. Every
// node and attribute in pattern must be present in full in the same order and
// parent-child structure. However, full can contain other attributes and nodes
// in between those that matches pattern.
func ContainsTree(full, pattern *html.Node) bool {
	return containsTree(full, pattern) == nil
}

func parseVerbatim(parent *html.Node, tkz *html.Tokenizer) (html.Token, error) {
	for {
		tkz.Next()
		token := tkz.Token()
		node := &html.Node{
			Parent:      parent,
			PrevSibling: parent.LastChild,
			DataAtom:    token.DataAtom,
			Data:        token.Data,
			Attr:        token.Attr,
		}
		switch tkz.Token().Type {
		case html.ErrorToken:
			return token, tkz.Err()
		case html.TextToken:
			node.Type = html.TextNode
		case html.StartTagToken:
			node.Type = html.ElementNode
			end, err := parseVerbatim(node, tkz)
			if err != nil {
				return end, err
			}
		case html.EndTagToken:
			return token, nil
		case html.SelfClosingTagToken:
			node.Type = html.ElementNode
		case html.CommentToken:
			continue
		case html.DoctypeToken:
			continue
		}
		if parent.LastChild != nil {
			parent.LastChild.NextSibling = node
		} else {
			parent.FirstChild = node
		}
		parent.LastChild = node
	}
}

// ParseVerbatim parses the given content into *html.Node. The result is a
// verbatim translation of the input content, without any additions or
// removals. This differs from
// https://pkg.go.dev/golang.org/x/net/html, which is HTML5 compatible.
func ParseVerbatim(content string) (*html.Node, error) {
	tokenizer := html.NewTokenizer(bytes.NewBuffer([]byte(content)))
	var node html.Node
	if _, err := parseVerbatim(&node, tokenizer); err != nil && err != io.EOF {
		return nil, err
	}

	return node.FirstChild, nil
}

// MustParseVerbatim will panics if parsing fails.
func MustParseVerbatim(content string) *html.Node {
	n, err := ParseVerbatim(content)
	if err != nil {
		panic(err)
	}
	return n
}
