package gparse

import (
	"io"
	"github.com/yuin/goldmark/ast"
)

// MkdnToken is the output of the parser `github.com/yuin/goldmark`.
// Note that the parser delivers them in a tree structure but that
// then we organise them into a flat list for convenience.
// These are handled differently from XML tokens, which are returned
// by the Go XML parser, and have a very simple interface but then
// have to be typecast into a particular token type. For Markdown
// tokens, we store the parser's token but then try to condense the
// token's data into a uniform set of fields.
// Note that a MkdnToken is converted into a GToken, which has these fields:
// - GTagTokType
// - GName ("SE" and "EE" *only*)
// - GAttList ("SE" *only*)
// - Keyword string (ProcInst "PI" and Directive "Dir", *only*)
// - Otherwords string (for all *except* "SE" and "EE")
type MkdnToken struct {
	ast.Node
	NodeDepth    int // from node walker
	NodeType     string // "nil", "Blk", "Inl", "Doc"
	NodeKind     string // the many rich text tags
	NodeKindEnum ast.NodeKind
	NodeKindInt  int
	// NodeText is the text of the MD node,
	//  and it is not present for all nodes.
	NodeText string
	// DitaTag and HtmlTag are the equivalent LwDITA and (X)HTML tags,
	// possibly with an attribute specified too. sDitaTag is authoritative;
	// sHtmlTag is provided mainly as an aid to understanding the code.
	DitaTag, HtmlTag string
	NodeNumeric      int // Headings, Emphasis, ...?
}

func (p MkdnToken) Echo() string {
	return "MKDN_TOKEN_ECHO"
}
func (p MkdnToken) EchoTo(w io.Writer){
	w.Write([]byte(p.Echo()))
}
func (p MkdnToken) String() string {
	return "MKDN_TOKEN_STRING"
}
func (p MkdnToken) DumpTo(w io.Writer) {
	w.Write([]byte(p.String()))
}
