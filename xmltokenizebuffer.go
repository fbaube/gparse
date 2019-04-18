package gparse

import (
	"encoding/xml"
	"io"
	S "strings"

	"github.com/pkg/errors"
)

// XmlTokenizeBuffer takes a string, so we can assume that we can
// discard it after use cos the caller has another copy of it.
// To be safe, it copies every token using `xml.CopyToken(T)`.
func XmlTokenizeBuffer(s string) (xtokens []xml.Token, err error) {
	var e error
	var T, TT xml.Token
	xtokens = make([]xml.Token, 0, 100)

	r := S.NewReader(s)
	var parser = xml.NewDecoder(r)
	// Strict mode does not enforce XML namespace requirements. In parti-
	// cular it does not reject namespace tags that use undefined prefixes.
	// Such tags are recorded with the unknown prefix as the namespace URL.
	parser.Strict = false
	// When Strict == false, AutoClose is a set of elements to consider
	// closed immediately after they are opened, regardless of whether
	// an end element is present. For example, <br/>.
	// TODO: Add anything for LwDITA ?
	parser.AutoClose = xml.HTMLAutoClose
	// Entity can map non-standard entity names to string replacements.
	// The parser is preloaded with the following standard XML mappings,
	// whether or not they are also provided in the actual map content:
	//	"lt": "<", "gt": ">", "amp": "&", "apos": "'", "quot": `"`
	// NOTE It doesn't do parameter entities, and we havnt necessarily
	// parsed any entities at all yet, so don't bother trying to use this.
	// NOTE If you dump all these, you find that there's a zillion of'em.
	parser.Entity = xml.HTMLEntity

	for {
		T, e = parser.Token()
		// fmt.Printf("%+v, %s, %s \n", token, tokenType, tokenString)
		if e == io.EOF {
			if ExtraInfo {
				println("==> xml tokenization eof ok")
			}
			break
		}
		if e != nil {
			return xtokens, errors.Wrap(e, "gxml.MakeXmlTokensFromContent.tokenization")
		}
		TT = xml.CopyToken(T)
		xtokens = append(xtokens, TT)
	}
	return xtokens, nil
}
