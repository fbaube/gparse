package gparse

import (
	S "strings"

	"github.com/sanity-io/litter"
	"github.com/yuin/goldmark/ast"
)

// GTokznFromMkdnTokens turns every `MkdnToken` Markdown token into a
// `GToken`. It's pretty simple, because no tree building is done yet.
//
func GTokznFromMkdnTokens(mtokens []MarkupStringer/*MkdnToken*/) (gtokens GTokenization, err error) {
	var MS MarkupStringer
	var MT MkdnToken
	var pMT *MkdnToken
	var pGT *GToken
	var isNotFirstToken = false

	// First dump them all in an indented tree
	println("======================================")
	println("MkdnTokens TREE DUMP:")
	for _, pp := range mtokens {
		p := pp.(MkdnToken)
		var pfx = S.Repeat("  ", p.NodeDepth-1)
		println(pfx, p.NodeType[:1], S.TrimPrefix(p.NodeKind, "Kind"),
			p.DitaTag, p.HtmlTag, p.NodeText)
	}

	for _, MS = range mtokens {
		MT = MS.(MkdnToken)
		pGT = new(GToken)
		pGT.BaseToken = *pMT

		// Is it the very first token ?
		if !isNotFirstToken {
			// Has to be Document token
			if MT.NodeKindEnum != ast.KindDocument {
				panic("NOT DOC")
			}
			println("=DOC=DOC=:", litter.Sdump(MT))
		}
		isNotFirstToken = true

		/* Fields:
		NodeDepth    int // from node walker
		NodeType     string
		NodeKind     string
		NodeKindEnum ast.NodeKind
		NodeKindInt  int
		// NodeText is the text of the MD node,
		//  and it is not present for all nodes.
		NodeText string
		// DitaTag and HtmlTag are the equivalent LwDITA and (X)HTML tags,
		// possibly with an attribute specified too. sDitaTag is authoritative;
		// sHtmlTag is provided mainly as an aid to understanding the code.
		DitaTag, HtmlTag string
		NodeNumeric   int // Headings, Emphasis, ...?
		*/

		switch MT.NodeKindEnum {

		case ast.KindAutoLink:
		case ast.KindBlockquote:
		case ast.KindCodeBlock:
		case ast.KindCodeSpan:
		case ast.KindDocument:
		case ast.KindEmphasis:
		case ast.KindFencedCodeBlock:
		case ast.KindHTMLBlock:
		case ast.KindHeading:
		case ast.KindImage:
		case ast.KindLink:
		case ast.KindList:
		case ast.KindListItem:
		case ast.KindParagraph:
		case ast.KindRawHTML:
		case ast.KindText:
		case ast.KindTextBlock:
		case ast.KindThematicBreak:
			/*
				case xml.StartElement:
					// A StartElement has a Name (GName) and Attributes (GAtt's)
					pGT.GTagTokType = "SE"
					// type xml.StartElement struct { Name Name ; Attr []Attr }
					xTag := xml.CopyToken(MT).(xml.StartElement)
					pGT.GName = GName(xTag.Name)
					// println("SE:", pGT.GName.String())
					if pGT.GName.Space == NS_XML {
						pGT.GName.Space = "xml"
					}
					for _, A := range xTag.Attr {
						if A.Name.Space == NS_XML {
							// println("TODO check name.local: newgtoken/L36 xml:" + A.Name.Local)
							A.Name.Space = "xml"
						}
						a := GAtt(A)
						// aa := &a
						pGT.GAttList = append(pGT.GAttList, a) // aa)
					}
					pGT.Keyword = ""
					pGT.Otherwords = ""
					// fmt.Printf("<!--Start-Tag--> %s \n", outGT.Echo())
					gtokens = append(gtokens, pGT)
					continue

				case xml.EndElement:
					// An EndElement has a Name (GName).
					pGT.GTagTokType = "EE"
					// type xml.EndElement struct { Name Name }
					xTag := xml.CopyToken(MT).(xml.EndElement)
					pGT.GName = GName(xTag.Name)
					if pGT.GName.Space == NS_XML {
						pGT.GName.Space = "xml"
					}
					pGT.Keyword = ""
					pGT.Otherwords = ""
					// fmt.Printf("<!--End-Tagnt--> %s \n", outGT.Echo())
					gtokens = append(gtokens, pGT)
					continue

				case xml.ProcInst:
					pGT.GTagTokType = "PI"
					// type xml.ProcInst struct { Target string ; Inst []byte }
					xTag := MT.(xml.ProcInst)
					pGT.Keyword = S.TrimSpace(xTag.Target)
					pGT.Otherwords = S.TrimSpace(string(xTag.Inst))
					// fmt.Printf("<!--ProcInstr--> <?%s %s?> \n",
					// 	outGT.Keyword, outGT.Otherwords)
					gtokens = append(gtokens, pGT)
					continue

				case xml.CharData:
					// type CharData []byte
					pGT.GTagTokType = "CD"
					bb := []byte(xml.CopyToken(MT).(xml.CharData))
					s := S.TrimSpace(string(bb))
					// pGT.Keyword remains ""
					pGT.Otherwords = s
					if s == "" {
						// ilog.Printf("PCDATA is all whitespace: \n")
						// DO NOTHING
						// NOTE This may do weird things to elements
						// that have text content models.
						// println("WARNING: Got an all-whitespace xml.CharData")
						continue
					}
					// fmt.Printf("<!--Char-Data--> %s \n", outGT.Otherwords)
					gtokens = append(gtokens, pGT)
					continue

				case xml.Comment:
					// type Comment []byte
					pGT.GTagTokType = "Cmt"
					// pGT.Keyword remains ""
					pGT.Otherwords = S.TrimSpace(string([]byte(MT.(xml.Comment))))
					// fmt.Printf("<!-- Comment --> <!-- %s --> \n", outGT.Otherwords)
					gtokens = append(gtokens, pGT)
					continue

				case xml.Directive: // type Directive []byte
					pGT.GTagTokType = "Dir"
					s := S.TrimSpace(string([]byte(MT.(xml.Directive))))
					pGT.Keyword, pGT.Otherwords = SU.SplitOffFirstWord(s)
					// fmt.Printf("<!--Directive--> <!%s %s> \n",
					// 	outGT.Keyword, outGT.Otherwo rds)
					gtokens = append(gtokens, pGT)
					continue

				default:
					pGT.GTagTokType = "ERR"
					println(fmt.Sprintf("Unrecognized xml.Token type<%T> for: %+v", MT, MT))
					panic("OOPS")
					// continue
			*/
		}
	}
	return gtokens, nil
}
