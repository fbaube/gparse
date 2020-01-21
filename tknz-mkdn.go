package gparse

var MNdTypes = []string{"nil", "Blk", "Inl", "Doc"}

// var theSRC []byte // string
// var NdKdNms []string
/*
var TheSourceBfr []byte
var TheSourceAfr []byte
var TheReader text.Reader
var r RRR.Renderer
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
