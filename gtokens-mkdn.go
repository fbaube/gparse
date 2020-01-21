package gparse

import (
	PU "github.com/fbaube/parseutils"
	"github.com/yuin/goldmark/ast"
)

/*
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
*/

// DoGTokens_mkdn turns every `MkdnToken` Markdown token into a
// `GToken`. It's pretty simple, because no tree building is done yet.
//
func DoGTokens_mkdn(pCPR *PU.ConcreteParseResults_mkdn) ([]*GToken, error) {
	var NL []ast.Node
	var DL []int
	var p *GToken
	var gTokens = make([]*GToken,0)
	var NT ast.NodeType // 1..3
	var NK ast.NodeKind
	var NKi int

	NL = pCPR.NodeList
	DL = pCPR.NodeDepths
	/*
	var MS MarkupStringer
	var MT MkdnToken
	var pMT *MkdnToken
	*/

	// First dump them all in an indented tree
	/*
	println("======================================")
	println("MkdnTokens TREE DUMP:")
	for _, pp := range mtokens {
		p := pp.(MkdnToken)
		var pfx = S.Repeat("  ", p.NodeDepth-1)
		println(pfx, p.NodeType[:1], S.TrimPrefix(p.NodeKind, "Kind"),
			p.DitaTag, p.HtmlTag, p.NodeText)
	}
	*/
	for i, n := range NL {
		p = new(GToken)
		p.BaseToken = n
		p.Depth = DL[i]
		NT = n.Type()
		NK = n.Kind()
		NKi = int(NK)
		if i == 0 {
			// Has to be Document token
			if NK != ast.KindDocument {
				panic("NOT DOC")
			}
			// println("=DOC=DOC=:", litter.Sdump(MT))
		}
		/*
		Fields:
		NodeDepth    int // from node walker
		NodeType     string
		NodeKind     string
		NodeKindEnum ast.NodeKind
		NodeKindInt  int
		// NodeText is the text of the MD node,
		//  and it is not present for all nodes.
		NodeText string
		*/

		switch NK { // ast.NodeKind

				case ast.KindAutoLink:
					// https://github.github.com/gfm/#autolinks
					// Autolinks are absolute URIs and email addresses btwn < and >.
					// They are parsed as links, with the link target reused as the link label.
					p.NodeKind = "KindAutoLink"
					p.DitaTag = "xref"
					p.HtmlTag = "a@href"
					n2 := n.(*ast.AutoLink)
					// type AutoLink struct {
					//   BaseInline
					//   Type is a type of this autolink.
					//   AutoLinkType AutoLinkType
					//   Protocol specified a protocol of the link.
					//   Protocol []byte
					//   value *Text
					// }
					// w.WriteString(`<a href="`)
					// url := n.URL(source)
					// label := n.Label(source)
					// if n.AutoLinkType == ast.AutoLinkEmail &&
					//    !bytes.HasPrefix(bytes.ToLower(url), []byte("mailto:")) {
					//   w.WriteString("mailto:")
					// }
					// w.Write(util.EscapeHTML(util.URLEscape(url, false)))
					// w.WriteString(`">`)
					// w.Write(util.EscapeHTML(label))
					// w.WriteString(`</a>`)
				case ast.KindBlockquote:
					p.NodeKind = "KindBlockquote"
					p.DitaTag = "?blockquote"
					p.HtmlTag = "blockquote"
					n2 := n.(*ast.Blockquote)
					// type Blockquote struct {
					//   BaseBlock
					// }
					// w.WriteString("<blockquote>\n")
				case ast.KindCodeBlock:
					p.NodeKind = "KindCodeBlock"
					p.DitaTag = "?pre+?code"
					p.HtmlTag = "pre+code"
					n2 := n.(*ast.CodeBlock)
					// type CodeBlock struct {
					//   BaseBlock
					// }
					// w.WriteString("<pre><code>")
					// r.writeLines(w, source, n)
				case ast.KindCodeSpan:
					p.NodeKind = "KindCodeSpan"
					p.DitaTag = "?code"
					p.HtmlTag = "code"
					// // n2 := n.(*ast.CodeSpan)
					// // sDump = litter.Sdump(*n2)
					// type CodeSpan struct {
					//   BaseInline
					// }
					// w.WriteString("<code>")
					// for c := n.FirstChild(); c != nil; c = c.NextSibling() {
					//   segment := c.(*ast.Text).Segment
					//   value := segment.Value(source)
					//   if bytes.HasSuffix(value, []byte("\n")) {
					//     r.Writer.RawWrite(w, value[:len(value)-1])
					//     if c != n.LastChild() {
					//       r.Writer.RawWrite(w, []byte(" "))
					//     }
					//   } else {
					//     r.Writer.RawWrite(w, value)
				case ast.KindDocument:
					// Note that metadata comes btwn this
					// start-of-document tag and the content ("body").
					p.NodeKind = "KindDocument"
					p.DitaTag = "topic"
					p.HtmlTag = "html"
				case ast.KindEmphasis:
					p.NodeKind = "KindEmphasis"
					// iLevel 2 | iLevel 1
					p.DitaTag = "b|i"
					p.HtmlTag = "strong|em"
					n2 := n.(*ast.Emphasis)
					p.NodeNumeric = n2.Level
					// type Emphasis struct {
					//   BaseInline
					//   Level is a level of the emphasis.
					//   Level int
					// }
					// tag := "em"
					// if n.Level == 2 {
					//   tag = "strong"
					// }
					// if entering {
					//   w.WriteByte('<')
					//   w.WriteString(tag)
					//   w.WriteByte('>')
				case ast.KindFencedCodeBlock:
					p.NodeKind = "KindFencedCodeBlock"
					p.DitaTag = "?code"
					p.HtmlTag = "code"
					n2 := n.(*ast.FencedCodeBlock)
					// type FencedCodeBlock struct {
					//   BaseBlock
					//   Info returns a info text of this fenced code block.
					//   Info *Text
					//   language []byte
					// }
					// w.WriteString("<pre><code")
					// language := n.Language(source)
					// if language != nil {
					//   w.WriteString(" class=\"language-")
					//   r.Writer.Write(w, language)
				case ast.KindHTMLBlock:
					p.NodeKind = "KindHTMLBlock"
					p.DitaTag = "?htmlblock"
					p.HtmlTag = "?htmlblock"
					n2 := n.(*ast.HTMLBlock)
					// type HTMLBlock struct {
					//   BaseBlock
					//   Type is a type of this html block.
					//   HTMLBlockType HTMLBlockType
					//   ClosureLine is a line that closes this html block.
					//   ClosureLine textm.Segment
					// }
					// if r.Unsafe {
					//   l := n.Lines().Len()
					//   for i := 0; i < l; i++ {
					//     line := n.Lines().At(i)
					//     w.Write(line.Value(source))
					//   }
					// } else {
					//   w.WriteString("<!-- raw HTML omitted -->\n")
				case ast.KindHeading:
					p.NodeKind = "KindHeading"
					p.DitaTag = "?"
					p.HtmlTag = "h%d"
					n2 := n.(*ast.Heading)
					p.NodeNumeric = n2.Level
					// type Heading struct {
					//   BaseBlock
					//   Level returns a level of this heading.
					//   This value is between 1 and 6.
					//   Level int
					// }
				// w.WriteString("<h")
				// w.WriteByte("0123456"[n.Level])
				case ast.KindImage:
					p.NodeKind = "KindImage"
					p.DitaTag = "image"
					p.HtmlTag = "img"
					n2 := n.(*ast.Image)
					// type Image struct {
					//   baseLink
					// }
					// w.WriteString("<img src=\"")
					// if r.Unsafe || !IsDangerousURL(n.Destination) {
					//   w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
					// }
					// w.WriteString(`" alt="`)
					// w.Write(n.Text(source))
					// w.WriteByte('"')
					// if n.Title != nil {
					//   w.WriteString(` title="`)
					//   r.Writer.Write(w, n.Title)
					//   w.WriteByte('"')
					// }
					// if r.XHTML {
					//   w.WriteString(" />")
					// } else {
					//   w.WriteString(">")
				case ast.KindLink:
					p.NodeKind = "KindLink"
					p.DitaTag = "xref"
					p.HtmlTag = "a@href"
					n2 := n.(*ast.Link)
					// type Link struct {
					//   baseLink
					// }
					// w.WriteString("<a href=\"")
					// if r.Unsafe || !IsDangerousURL(n.Destination) {
					//   w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
					// }
					// w.WriteByte('"')
					// if n.Title != nil {
					//   w.WriteString(` title="`)
					//   r.Writer.Write(w, n.Title)
					//   w.WriteByte('"')
					// }
					// w.WriteByte('>')
				case ast.KindList:
					p.NodeKind = "KindList"
					n2 := n.(*ast.List)
					if n2.IsOrdered() {
						p.DitaTag = "ol"
						p.HtmlTag = "ol"
					} else {
						p.DitaTag = "ul"
						p.HtmlTag = "ul"
					}
					// type List struct {
					//   BaseBlock
					//   Marker is a markar character like '-', '+', ')' and '.'.
					//   Marker byte
					//   IsTight is a true if this list is a 'tight' list.
					//   See https://spec.commonmark.org/0.29/#loose for details.
					//   IsTight bool
					//   Start is an initial number of this ordered list.
					//   If this list is not an ordered list, Start is 0.
					//   Start int
					// }
					// tag := "ul"
					// if n.IsOrdered() {
					//   tag = "ol"
					// }
					// w.WriteByte('<')
					// w.WriteString(tag)
					// if n.IsOrdered() && n.Start != 1 {
					//   fmt.Fprintf(w, " start=\"%d\">\n", n.Start)
					// } else {
					//   w.WriteString(">\n")
				case ast.KindListItem:
					p.NodeKind = "KindListItem"
					n2 := n.(*ast.ListItem)
					p.DitaTag = "li"
					p.HtmlTag = "li"
					// type ListItem struct {
					//   BaseBlock
					//   Offset is an offset potision of this item.
					//   Offset int
					// }
					// w.WriteString("<li>")
					// fc := n.FirstChild()
					// if fc != nil {
					//   if _, ok := fc.(*ast.TextBlock); !ok {
					//     w.WriteByte('\n')
				case ast.KindParagraph:
					p.NodeKind = "KindParagraph"
					p.DitaTag = "p"
					p.HtmlTag = "p"
					// // n2 := n.(*ast.Paragraph)
					// // sDump = litter.Sdump(*n2)
					// type Paragraph struct {
					//   BaseBlock
					// }
					// w.WriteString("<p>")
				case ast.KindRawHTML:
					p.NodeKind = "KindRawHTML"
					p.DitaTag = "?rawhtml"
					p.HtmlTag = "?rawhtml"
					n2 := n.(*ast.RawHTML)
					// type RawHTML struct {
					//   BaseInline
					//   Segments *textm.Segments
					// }
					// if r.Unsafe {
					// n := node.(*ast.RawHTML)
					// l := n.Segments.Len()
					// for i := 0; i < l; i++ {
					//   segment := n.Segments.At(i)
					//   w.Write(segment.Value(source))
					// }
				case ast.KindText:
					p.NodeKind = "KindText"
					n2 := n.(*ast.Text)
					p.DitaTag = "?text"
					p.HtmlTag = "?text"
					// // sDump = litter.Sdump(*n2)
					// type Text struct {
					//   BaseInline
					//   Segment is a position in a source text.
					//   Segment textm.Segment
					//   flags uint8
					// }
					segment := n2.Segment
					// p.NodeText = fmt.Sprintf("KindText:\n | %s", string(TheReader.Value(segment)))
					p.NodeText = /* fmt.Sprintf("KindText:\n | %s", */ string(pCPR.Reader.Value(segment)) //)
					/*
						if n.IsRaw() {
							r.Writer.RawWrite(w, segment.Value(TheSource))
						} else {
							r.Writer.Write(w, segment.Value(TheSource))
							if n.HardLineBreak() || (n.SoftLineBreak() && r.HardWraps) {
								if r.XHTML {
									w.WriteString("<br />\n")
								} else {
									w.WriteString("<br>\n")
								}
							} else if n.SoftLineBreak() {
								w.WriteByte('\n')
							}
						}
					*/
				case ast.KindTextBlock:
					p.NodeKind = "KindTextBlock"
					p.DitaTag = "?textblock"
					p.HtmlTag = "?textblock"
					// // n2 := n.(*ast.TextBlock)
					// // sDump = litter.Sdump(*n2)
					// type TextBlock struct {
					//   BaseBlock
					// }
					// if _, ok := n.NextSibling().(ast.Node); ok && n.FirstChild() != nil {
					//   w.WriteByte('\n')
				case ast.KindThematicBreak:
					p.NodeKind = "KindThematicBreak"
					p.DitaTag = "hr"
					p.HtmlTag = "hr"
					// type ThemanticBreak struct {
					//   BaseBlock
					// }
					// if r.XHTML {
					//   w.WriteString("<hr />\n")
					// } else {
					//   w.WriteString("<hr>\n")
				default:
					p.NodeKind = "KindUNK"
					p.DitaTag = "UNK"
					p.HtmlTag = "UNK"
				}
		}
	return gTokens, nil
}