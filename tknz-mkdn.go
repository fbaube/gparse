package gparse

import (
	"github.com/sanity-io/litter"
	GM "github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	RRR "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

// var theSRC []byte // string
// var NdKdNms []string
var TheSourceBfr []byte
var TheSourceAfr []byte
var TheReader text.Reader
var r RRR.Renderer

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

/* Variables
   var DefinitionList = &definitionList{}
      use PHP Markdown Extra Definition lists. (DT, : DD, DD, <br>)
   var Footnote = &footnote{}
       use PHP Markdown Extra Footnotes. fn.[^1] ... [^1]: The fn.
   var GFM = &gfm{}
       provides Github Flavored markdown functionalities.
   var Linkify = &linkify{}
       parse text that seems like a URL.
   var Strikethrough = &strikethrough{}
       use strikethru expressions like '~~text~~' .
   var Table = &table{}
       use GFM tables .
   var TaskList = &taskList{}
       use GFM task lists. [ ] [x]
   var Typographer = &typographer{}
       replace punctuations with typographic entities.
*/

/*
   An Option interface is a functional option type for the Parser.
   A ParseOption is a functional option type for the Parser.Parse.

   WithAutoHeadingID    enables custom heading IDs and autogen'd heading ids.
   WithHeadingAttribute enables custom heading attributes.
   WithASTTransformers  lets you add ASTTransformers to the parser.
   WithAttribute        enables custom attributes (BUT, only hdgs).
   WithBlockParsers  lets you add BlockParsers  to the parser.
   WithInlineParsers lets you add InlineParsers to the parser.
   WithOption          sets an arbitrary option to the parser.
   WithParagraphTransformers  lets you add PT's to the parser.
   WithContext                lets you override a default context.
   WithWorkers sets the number of inline parsing workers (goroutines).
*/

/*
   HTML Renderer options
   html.WithWriter(html.Writer) write contents to an io.Writer.
   html.WithHardWraps - Render new lines as <br>.
   html.WithXHTML     - Render as XHTML.
   html.WithUnsafe    - By default, GM does not render raw HTMLs and potentially
       dangerous links. With this option, GM renders these contents as-is.
*/

var pMTokzn *MkdnTokenization

/*
type MarkdownAST struct {
	ast.Node
}
func (p MarkdownAST) Echo() string {
	return "MKDN ECHO" // p.Node.String()
}
func (p MarkdownAST) EchoTo(w io.Writer) {
	w.Write([]byte(p.Echo()))
}
func (p MarkdownAST) String() string {
	// return p.Node.String()
	return fmt.Sprintf("%+v", p)
}
func (p MarkdownAST) DumpTo(w io.Writer) {
	w.Write([]byte(p.String()))
}
*/

// MkdnTokenizeBuffer takes a string and returns both the tree produced by
// the parser AND a flat list of the Nodes.
func MkdnTokenizeBuffer(s string) (*MkdnTokenization, error) {
	var GoldMarkDown GM.Markdown
	// var REND renderer.Renderer

	GoldMarkDown = GM.New(
		GM.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			extension.Strikethrough,
			extension.Table,
		),
		GM.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		GM.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithWriter(html.DefaultWriter), // os.Stdout),
		),
	)
	var TheParser parser.Parser
	var TheParseTree ast.Node
	TheSourceBfr = []byte(s) // p.CheckedPath.Raw)
	println("TheSource:", litter.Sdump(s))
	// r = GM.DefaultRenderer() // GoldMarkDown.Renderer().(html.Renderer)
	r = GoldMarkDown.Renderer()
	TheParser = GoldMarkDown.Parser()
	TheReader = text.NewReader(TheSourceBfr)
	TheParseTree = TheParser.Parse(TheReader)
	TheSourceAfr = TheReader.Source()
	pMTokzn = new(MkdnTokenization)
	pMTokzn.TreeRootNode = TheParseTree
	/*
		println("======================================")
		println("Markdown TREE DUMP:")
		// ParseTree.Dump(TheSource, 0)
		ast.DumpHelper(TheParseTree, TheSourceAfr, 0, nil, nil)
	*/
	println("======================================")
	println("======================================")

	TheReader.ResetPosition()
	return pMTokzn, nil
}

func (pMT *MkdnTokenization) MkdnNodeListFromAST() {
	pMT.ListNodes = make([]ast.Node, 0) // MarkupStringer, 0)
	pMT.ListNodeDepths = make([]int, 0)
	e := ast.Walk(pMT.TreeRootNode, walkerFuncMkdnGathertreenodes)
	if e != nil {
		panic(e)
	}
}
	// println("======================================")

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


var MNdTypes = []string{"nil", "Blk", "Inl", "Doc"}

var MWalkLevel int

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

func MStrattributesFromAttributes(atts []ast.Attribute) []strattribute {
	var stratts []strattribute
	for _, attr := range atts {
		// litter.Dump(attr)
		// if ok,_ := []uint8{
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
}
