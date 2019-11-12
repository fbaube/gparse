package gparse

import (
  "encoding/xml"
  "golang.org/x/net/html"
  "github.com/yuin/goldmark/ast"
)

type XmlTokenization struct {
  Tokens []xml.Token
}

// MkdnTokenization is a bit dodgy cos `ast.Node` is an interface, not a struct.
type MkdnTokenization struct {
  TreeRootNode ast.Node
  ListNodes  []ast.Node
  ListNodeDepths  []int
}

type HtmlTokenization struct {
  TreeRootNodeP *html.Node
  ListNodesP []*html.Node
  ListNodeDepths []int
}

func (hNodes HtmlTokenization) GetAllByAnyTag(ss []string) []*html.Node {
	if ss == nil || len(ss) == 0 {
		return nil
	}
  var ret []*html.Node
  ret = make([]*html.Node, 0)
	for _, p := range hNodes.ListNodesP {

  }
  return ret
}

// GetAllByTag returns a slice of `*html.Node`. It checks the basic tag only,
// not any namespace. Note that these tag lookup func's default to searching
// the `ListNodesP`, not the tree of `Node`s.
func (hT HtmlTokenization) GetAllByTag(s string) []*html.Node {
	if s == "" {
		return nil
	}
  var ret []*html.Node
	ret = make([]*html.Node , 0)
	for _, p := range hT.ListNodesP {

  }
  return ret
}

// GetAllByTag returns a slice of `ast.Node`. It checks the basic tag only,
// not any namespace. Note that these tag lookup func's default to searching
// the `ListNodesP`, not the tree of `Node`s.
func (mNodes MkdnTokenization) GetAllByTag(s string) []ast.Node {
	if s == "" {
		return nil
	}
  var ret []ast.Node
	ret = make([]ast.Node,0)
	for _, p := range mNodes.ListNodes {
  }
  return ret
}
