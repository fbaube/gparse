package gparse

import (
	"encoding/xml"
	"fmt"
	S "strings"
	// "net/http"
	// "github.com/yuin/goldmark/ast"
	SU "github.com/fbaube/stringutils"
	PU "github.com/fbaube/parseutils"
)

// DoGTokens_xml is TBS.
func DoGTokens_xml(pCPR *PU.ConcreteParseResults_xml) ([]*GToken, error) {
	var XTs []xml.Token
	var xt    xml.Token
	var p       *GToken
	var i int
	var gTokens = make([]*GToken,0)
	var gDepths = make([]int, 0)
	var iDepth  = 1 // current depth
	var prDpth  int // depth for printing
	var canSkip bool

	if pCPR.NodeDepths != nil {
		panic("XML tokens already have depths")
	}
	pCPR.NodeDepths = make([]int, 0)
	XTs = pCPR.NodeList

	for i, xt = range XTs {
		p = new(GToken)
		p.BaseToken = xt
		prDpth  = iDepth
		canSkip = false

		switch xt.(type) {

		case xml.StartElement:
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
				p.GAtts = append(p.GAtts, a)
			}
			p.Keyword    = ""
			p.Otherwords = ""
			// fmt.Printf("<!--Start-Tag--> %s \n", outGT.Echo())
			iDepth++

		case xml.EndElement:
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
			iDepth--
			canSkip = true

		case xml.ProcInst:
			p.TTType = "PI"
			// type xml.ProcInst struct { Target string ; Inst []byte }
			xTag := xt.(xml.ProcInst)
			p.Keyword = S.TrimSpace(xTag.Target)
			p.Otherwords = S.TrimSpace(string(xTag.Inst))
			// fmt.Printf("<!--ProcInstr--> <?%s %s?> \n",
			// 	outGT.Keyword, outGT.Otherwords)

		case xml.CharData:
			// type CharData []byte
			p.TTType = "CD"
			bb := []byte(xml.CopyToken(xt).(xml.CharData))
			s := S.TrimSpace(string(bb))
			// pGT.Keyword remains ""
			p.Otherwords = s
			if s == "" {
				canSkip = true
				p.Depth = prDpth
				gTokens = append(gTokens, nil)
				gDepths = append(gDepths, prDpth)
				continue
				// ilog.Printf("PCDATA is all whitespace: \n")
				// DO NOTHING
				// NOTE This may do weird things to elements
				// that have text content models.
				// println("WARNING: Got an all-whitespace xml.CharData")
			}
			// } else {
				// fmt.Printf("<!--Char-Data--> %s \n", outGT.Otherwords)

		case xml.Comment:
			// type Comment []byte
			p.TTType = "Cmt"
			// pGT.Keyword remains ""
			p.Otherwords = S.TrimSpace(string([]byte(xt.(xml.Comment))))
			// fmt.Printf("<!-- Comment --> <!-- %s --> \n", outGT.Otherwords)

		case xml.Directive: // type Directive []byte
			p.TTType = "Dir"
			s := S.TrimSpace(string([]byte(xt.(xml.Directive))))
			p.Keyword, p.Otherwords = SU.SplitOffFirstWord(s)
			// fmt.Printf("<!--Directive--> <!%s %s> \n",
			// 	outGT.Keyword, outGT.Otherwo rds)

		default:
			p.TTType = "ERR"
			println(fmt.Sprintf("Unrecognized xml.Token type<%T> for: %+v", xt, xt))
			panic("OOPS")
			// continue
		}

		p.Depth = prDpth
		gTokens = append(gTokens, p)
		gDepths = append(gDepths, prDpth)

		if p != nil {
			sCS := ""
			if canSkip {
				sCS = "(canSkip?)"
			} else {
				fmt.Printf("[%02d:L%d]%s (%s) %s %s \n",
					i, prDpth, S.Repeat("  ", prDpth), p.TTType, p.Echo(), sCS)
			}
		}
	}
	pCPR.NodeDepths = gDepths
	return gTokens, nil
}
