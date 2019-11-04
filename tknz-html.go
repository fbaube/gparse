package gparse

import (
	S "strings"
	"golang.org/x/net/html"
)

/*
// Package ast defines AST nodes that represent markdown elements.
package ast

import (
	"bytes"
	"fmt"
	"strings"

	textm "github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)
*/

// NewXmlItems (incl. ENTs, IDs, etc.)
// New GRefs, other link types
//
// XDITA / HDITA / MDITA:
// <xref> / <a href> / [link](/URI "title")
// <image>2 / <img> / ![alt text for an image](images/ image_name.jpg)
// <keydef> / <div data-class="keydef"> / MDITA-XP <div data- class="keydef"> in HDITA syntax
// <topicref> / <a href> inside a <li> / [link](/URI "title") inside a list item
// <media-source> / <source>
// @href
// @id
// @conref / @data-conref
// @keys   / @data-keys
// @keyref / @data-keyref

// println("processmkdn.go/Process:\n", p.CheckedPath.Raw)

// var e error
// var mdRoot *BF.Node
// var yamlFrontmatter []string

/*
var TheHTokens []MarkupStringer // *HtmlToken

type HtmlAST struct {
	html.Node
}
func (p HtmlAST) Echo() string {
	return "HTML ECHO" // p.Node.String()
}
func (p HtmlAST) EchoTo(w io.Writer) {
	w.Write([]byte(p.Echo()))
}
func (p HtmlAST) String() string {
	// return p.Node.String()
	return fmt.Sprintf("%+v", p)
}
func (p HtmlAST) DumpTo(w io.Writer) {
	w.Write([]byte(p.String()))
}
*/

var pHTokzn *HtmlTokenization

// HtmlTokenizeBuffer takes a string and returns both the tree produced by
// the parser AND a flat list of the Nodes.
func HtmlTokenizeBuffer(s string) (*HtmlTokenization, error) {
	var e error
	pHTokzn = new(HtmlTokenization)
	pHTokzn.TreeRootNodeP, e = html.Parse(S.NewReader(s))
	if e != nil { panic(e) }

	/*
		println("======================================")
		println("Markdown TREE DUMP:")
		// ParseTree.Dump(TheSource, 0)
		ast.DumpHelper(TheParseTree, TheSourceAfr, 0, nil, nil)
	*/
	println("======================================")
	println("======================================")

	pHTokzn.ListNodesP = make([]*html.Node, 0)
	println("Html TREE WALK:")
	// FIXME !!!!!!!!!!!! _ = HtmlWalk(pRoot, walkerFuncHtmlGathertreenodes)
	println("======================================")

	/*
		// fmt.Printf("    FM: %+v \n", yamlFrontmatter)
		if p.Header != nil {
			println("--> YAML frontmatter:\n", p.HedRaw, "---")
		}
	*/
	/*
		// fmt.Printf("    MD: %+v \n", *mdRoot)
		println("==BEG== DumpNode:BF:Root")
		// FIXME gparse.DumpBFnode(mdRoot, 0)
		println("==MID== DumpNode:BF:Root")
		NormalizeTextLeaves(mdRoot)
		// FIXME gparse.DumpBFnode(mdRoot, 0)
		println("==END== DumpNode:BF:Root")
	*/
	return pHTokzn, nil
}

var HNdTypes = []string{"nil", "Blk", "Inl", "Doc"}

var HWalkLevel int

// A BaseBlock struct implements the Node interface.
/*
type BaseBlock struct {
	BaseNode
	blankPreviousLines bool
	lines              *textm.Segments
}
type baseLink struct {
  BaseInline
  // Destination is a destination(URL) of this link.
  Destination []byte
  // Title is a title of this link.
  Title []byte
}
*/

func (pHT *HtmlTokenization) HtmlNodeListFromAST() {
	pHT.ListNodesP = make([]*html.Node, 0) // MarkupStringer, 0)
	pHT.ListNodeDepths = make([]int, 0)
	HtmlWalk(pHT.TreeRootNodeP, walkerFuncHtmlGathertreenodes)
}

func HStrattributesFromAttributes(atts []html.Attribute) []strattribute {
	//?? var stratts []strattribute
	for _, attr := range atts {
		println("HtmlAttr:", "NS", attr.Namespace, "Key", attr.Key, "Val", attr.Val)
		// litter.Dump(attr)
		// if ok,_ := []uint8{
		/* =================================
		strattr := new(strattribute)
		strattr.Name = string(attr.Name)
		switch attr.Value.(type) {
		case []uint8:
			strattr.Value = string(attr.Value.([]uint8))
		case [][]uint8:
			strattr.Value = ""
			var bbbb [][]byte
			var bb []byte
			bbbb = attr.Value.([][]byte)
			for _, bb = range bbbb {
				strattr.Value += string(bb) // attr.Value.([]uint8))
			}
		}
		stratts = append(stratts, *strattr)
	}
	return stratts
	*/
	}
return nil
}

// Walker is a function that will be called when Walk find a
// new node.
// entering is set true before walks children, false after walked children.
// If Walker returns error, Walk function immediately stop walking.
type HtmlWalker func(n *html.Node, entering bool)

// HtmlWalk walks a AST tree by the depth first search algorithm.
func HtmlWalk(n *html.Node, walker HtmlWalker) {
	walker(n, true)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		HtmlWalk(c, walker)
	}
	walker(n, false)
}
