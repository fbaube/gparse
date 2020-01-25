package gparse

import (
	"fmt"
	S "strings"
	PU "github.com/fbaube/parseutils"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// DoGTokens_html turns every `MkdnToken` Markdown token into a
// `GToken`. It's pretty simple, because no tree building is done yet.
//
func DoGTokens_html(pCPR *PU.ConcreteParseResults_html) ([]*GToken, error) {
	var NL []*html.Node
	var DL []int
	var p *GToken
	var gTokens = make([]*GToken,0)
	var NT html.NodeType // 1..3
	var datom atom.Atom
	// var NKi int

	NL = pCPR.NodeList
	DL = pCPR.NodeDepths

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
		NT = n.Type
		datom = n.DataAtom
		ds := S.TrimSpace(datom.String())
		Ds := S.TrimSpace(n.Data)

		if ds == "" && NT != html.DocumentNode {
			continue
		}
		/*
		fmt.Printf("html: NT<%d/%s> datom<%s> Data<%s> NS<%s> \n",
			S.TrimSpace(datom.String()), S.TrimSpace(n.Data), n.Namespace)
			// and Attr []Attribute
		*/
		s := fmt.Sprintf("L%d%s (%d:%s)  ",
			p.Depth, S.Repeat("  ", p.Depth-1), NT, PU.NTstring(NT))
		if ds == Ds {
			if ds == "" && NT != html.DocumentNode { // Now handled above! 
				s += "SKIP "
				} else {
					s += fmt.Sprintf("dd<%s> ", ds)
				}
		} else if ds == "" {
			s += fmt.Sprintf("Data<%s> ", Ds)
		} else {
			s += fmt.Sprintf("datom<%s> Data<%s> ",	ds, Ds)
		}
		if n.Namespace != "" {
			s += fmt.Sprintf("NS<%s> ", n.Namespace)
		}
		println(s)
	}
	return gTokens, nil
}

		/*
		if i == 0 {
			// Has to be Document token
			if NK != ast.KindDocument {
				panic("NOT DOC")
			}
			// println("=DOC=DOC=:", litter.Sdump(MT))
		}
		* /
		switch NT { // ast.NodeKind
			/*
			    Type      NodeType
							ErrorNode NodeType = iota
    					TextNode
    					DocumentNode
    					ElementNode
    					CommentNode
    					DoctypeNode
			    DataAtom  atom.Atom
							integer codes (a.k.a. atoms) for a fixed set of common HTML
							strings: tag names and attribute keys like "p" and "id".
			    Data      string
			    Namespace string
			    Attr      []Attribute
			* /
		case html.ElementNode:
						// A StartElement has a Name (GName) and Attributes (GAtt's)
						p.TTType = "SE"
						// type xml.StartElement struct { Name Name ; Attr []Attr }
						xTag := xml.CopyToken(xt).(xml.StartElement)
						p.GName = GName(xTag.Name)
						p.GName.FixNS()
						// println("SE:", pGT.GName.String())
						if p.GName.Space == NS_XML {
							p.GName.Space = "xml:"
						}
						for _, A := range xTag.Attr {
							if A.Name.Space == NS_XML {
								// println("TODO check name.local: newgtoken/L36 xml:" + A.Name.Local)
								A.Name.Space = "xml:"
							}
							a := GAtt(A)
							// aa := &a
							p.GAtts = append(p.GAtts, a) // aa)
						}
						p.Keyword = ""
						p.Otherwords = ""
						// fmt.Printf("<!--Start-Tag--> %s \n", outGT.Echo())
						gTokens = append(gTokens, p)
						continue

/*

					case html.EndElement:
						// An EndElement has a Name (GName).
						p.TTType = "EE"
						// type xml.EndElement struct { Name Name }
						xTag := xml.CopyToken(xt).(xml.EndElement)
						p.GName = GName(xTag.Name)
						if p.GName.Space == NS_XML {
							p.GName.Space = "xml:"
						}
						p.Keyword = ""
						p.Otherwords = ""
						// fmt.Printf("<!--End-Tagnt--> %s \n", outGT.Echo())
						gTokens = append(gTokens, p)
						continue

					case html.ProcInst:
						p.TTType = "PI"
						// type xml.ProcInst struct { Target string ; Inst []byte }
						xTag := xt.(xml.ProcInst)
						p.Keyword = S.TrimSpace(xTag.Target)
						p.Otherwords = S.TrimSpace(string(xTag.Inst))
						// fmt.Printf("<!--ProcInstr--> <?%s %s?> \n",
						// 	outGT.Keyword, outGT.Otherwords)
						gTokens = append(gTokens, p)
						continue

* /

					case html.CharData:
						// type CharData []byte
						p.TTType = "CD"
						bb := []byte(xml.CopyToken(xt).(xml.CharData))
						s := S.TrimSpace(string(bb))
						// pGT.Keyword remains ""
						p.Otherwords = s
						if s == "" {
							// ilog.Printf("PCDATA is all whitespace: \n")
							// DO NOTHING
							// NOTE This may do weird things to elements
							// that have text content models.
							// println("WARNING: Got an all-whitespace xml.CharData")
							continue
						}
						// fmt.Printf("<!--Char-Data--> %s \n", outGT.Otherwords)
						gTokens = append(gTokens, p)
						continue

					case html.CommentNode:
						// type Comment []byte
						p.TTType = "Cmt"
						// pGT.Keyword remains ""
						p.Otherwords = S.TrimSpace(string([]byte(xt.(xml.Comment))))
						// fmt.Printf("<!-- Comment --> <!-- %s --> \n", outGT.Otherwords)
						gTokens = append(gTokens, p)
						continue

					case html.Directive: // type Directive []byte
						p.TTType = "Dir"
						s := S.TrimSpace(string([]byte(xt.(xml.Directive))))
						p.Keyword, p.Otherwords = SU.SplitOffFirstWord(s)
						// fmt.Printf("<!--Directive--> <!%s %s> \n",
						// 	outGT.Keyword, outGT.Otherwo rds)
						gTokens = append(gTokens, p)
						continue

					default:
						p.TTType = "ERR"
						println(fmt.Sprintf("Unrecognized xml.Token type<%T> for: %+v", xt, xt))
						panic("OOPS")
						// continue
					}
				}
			}

/*
				case ast.KindAutoLink:
					// https://github.github.com/gfm/#autolinks
					// Autolinks are absolute URIs and email addresses btwn < and >.
					// They are parsed as links, with the link target reused as the link label.
					p.NodeKind = "KindAutoLink"
					p.DitaTag = "xref"
					p.HtmlTag = "a@href"
					n2 := n.(*ast.AutoLink)
					fmt.Printf("AutoLink: %+v \n", *n2)
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
					fmt.Printf("Blockquote: %+v \n", *n2)
					// type Blockquote struct {
					//   BaseBlock
					// }
					// w.WriteString("<blockquote>\n")
				case ast.KindCodeBlock:
					p.NodeKind = "KindCodeBlock"
					p.DitaTag = "?pre+?code"
					p.HtmlTag = "pre+code"
					n2 := n.(*ast.CodeBlock)
					fmt.Printf("CodeBlock: %+v \n", *n2)
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
					fmt.Printf("Emphasis: %+v \n", *n2)
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
					fmt.Printf("FencedCodeBlock: %+v \n", *n2)
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
					fmt.Printf("HTMLBlock: %+v \n", *n2)
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
					fmt.Printf("Heading: %+v \n", *n2)
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
					fmt.Printf("Image: %+v \n", *n2)
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
					fmt.Printf("Link: %+v \n", *n2)
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
					fmt.Printf("List: %+v \n", *n2)
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
					fmt.Printf("ListItem: %+v \n", *n2)
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
					fmt.Printf("RawHTML: %+v \n", *n2)
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
					fmt.Printf("Text: %+v \n", *n2)
					// // sDump = litter.Sdump(*n2)
					// type Text struct {
					//   BaseInline
					//   Segment is a position in a source text.
					//   Segment textm.Segment
					//   flags uint8
					// }
					segment := n2.Segment
					// p.NodeText = fmt.Sprintf("KindText:\n | %s", string(TheReader.Value(segment)))
					p.NodeText = /* fmt.Sprintf("KindText:\n | %s", * / string(pCPR.Reader.Value(segment)) //)
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
					* /
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

*/
