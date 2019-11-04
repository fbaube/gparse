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
