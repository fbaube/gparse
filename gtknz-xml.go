package gparse

import (
	"encoding/xml"
	"fmt"
	S "strings"
	// "net/http"
	// "github.com/yuin/goldmark/ast"
	SU "github.com/fbaube/stringutils"
)

// GTokznFromXmlTokens is TBS.
func GTokznFromXmlTokens(xtokens []xml.Token) (gtokens GTokenization, err error) {
	// NOTE Returns (`nil,nil`) if the token is valid but useless, and
	// can be skipped, such as an `xml.CharData` that is all whitespace.
	var XT xml.Token
	var pGT *GToken

	for _, XT = range xtokens {
		pGT = new(GToken)
		pGT.BaseToken = XT

		switch XT.(type) {

		case xml.StartElement:
			// A StartElement has a Name (GName) and Attributes (GAtt's)
			pGT.TTType = "SE"
			// type xml.StartElement struct { Name Name ; Attr []Attr }
			xTag := xml.CopyToken(XT).(xml.StartElement)
			pGT.GName = GName(xTag.Name)
			pGT.GName.FixNS()
			// println("SE:", pGT.GName.String())
			if pGT.GName.Space == NS_XML {
				pGT.GName.Space = "xml:"
			}
			for _, A := range xTag.Attr {
				if A.Name.Space == NS_XML {
					// println("TODO check name.local: newgtoken/L36 xml:" + A.Name.Local)
					A.Name.Space = "xml:"
				}
				a := GAtt(A)
				// aa := &a
				pGT.GAtts = append(pGT.GAtts, a) // aa)
			}
			pGT.Keyword = ""
			pGT.Otherwords = ""
			// fmt.Printf("<!--Start-Tag--> %s \n", outGT.Echo())
			gtokens = append(gtokens, pGT)
			continue

		case xml.EndElement:
			// An EndElement has a Name (GName).
			pGT.TTType = "EE"
			// type xml.EndElement struct { Name Name }
			xTag := xml.CopyToken(XT).(xml.EndElement)
			pGT.GName = GName(xTag.Name)
			if pGT.GName.Space == NS_XML {
				pGT.GName.Space = "xml:"
			}
			pGT.Keyword = ""
			pGT.Otherwords = ""
			// fmt.Printf("<!--End-Tagnt--> %s \n", outGT.Echo())
			gtokens = append(gtokens, pGT)
			continue

		case xml.ProcInst:
			pGT.TTType = "PI"
			// type xml.ProcInst struct { Target string ; Inst []byte }
			xTag := XT.(xml.ProcInst)
			pGT.Keyword = S.TrimSpace(xTag.Target)
			pGT.Otherwords = S.TrimSpace(string(xTag.Inst))
			// fmt.Printf("<!--ProcInstr--> <?%s %s?> \n",
			// 	outGT.Keyword, outGT.Otherwords)
			gtokens = append(gtokens, pGT)
			continue

		case xml.CharData:
			// type CharData []byte
			pGT.TTType = "CD"
			bb := []byte(xml.CopyToken(XT).(xml.CharData))
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
			pGT.TTType = "Cmt"
			// pGT.Keyword remains ""
			pGT.Otherwords = S.TrimSpace(string([]byte(XT.(xml.Comment))))
			// fmt.Printf("<!-- Comment --> <!-- %s --> \n", outGT.Otherwords)
			gtokens = append(gtokens, pGT)
			continue

		case xml.Directive: // type Directive []byte
			pGT.TTType = "Dir"
			s := S.TrimSpace(string([]byte(XT.(xml.Directive))))
			pGT.Keyword, pGT.Otherwords = SU.SplitOffFirstWord(s)
			// fmt.Printf("<!--Directive--> <!%s %s> \n",
			// 	outGT.Keyword, outGT.Otherwo rds)
			gtokens = append(gtokens, pGT)
			continue

		default:
			pGT.TTType = "ERR"
			println(fmt.Sprintf("Unrecognized xml.Token type<%T> for: %+v", XT, XT))
			panic("OOPS")
			// continue
		}
	}
	return gtokens, nil
}
