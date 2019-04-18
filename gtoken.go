package gparse

// This file: Structures for Generic Golang Tokens.
// They are based on struct `xml.Token` returned by the Golang XML parser
// but have been generalized to be usable for other LwDITA formats.

import (
	"encoding/xml"
	"io"
)

// GToken is meant to simplify & unify tokenisation across LwDITA's three
// supported input formats: XDITA XML, HDITA HTML5, and MD-XP Markdown.
// It also serves to represent all the various kinds of XML Directives,
// including DTDs(!).
//
// To do this, the tokens produced by each parsing API are reduced to
// their essentials:
// - token type (defined by the enumeration `GTagTokType`)
// - token text (tag name or non-tag text content)
// - tag attributes
//
// NOTE that XML Directives are later "normalized", but that's another story.
//
type GToken struct {
	// Keep the original token around, just in case.
	// Note that this `xml.Token` (or the entire `GToken`) might be erased in
	// later processing, if (for example) it is a CDATA that has only whitespace.
	xml.Token
	// GTagTokType enumerates the types of struct `GToken` and also the types of
	// struct `GTag`, which are a strict superset. Therefore the two structs use
	// a shared "type" enumeration.
	//
	// NOTE that "EE" (`EndElement`) is OK for a `GToken.Type` but (probably)
	// not for a `GTag.Type`, cos the existence of a matching `EndElement` for
	// every `StartElement` should be assumed (but need not actually be present)
	// in a valid `GTree`.
	GTagTokType
	// GName is for "SE" and "EE" *only*
	GName
	// GAttList is for "SE" *only*
	GAttList
	// Keyword is for ProcInst "PI" and Directive "Dir", *only*
	Keyword string
	// Otherwords is for all *except* "SE" and "EE"
	Otherwords string
}

// Echo implements Markupper.
func (T GToken) Echo() string {
	// var s string
	switch T.GTagTokType {

	case "SE":
		return "<" + T.GName.Echo() + T.GAttList.Echo() + ">"

	case "EE":
		return "</" + T.GName.Echo() + ">"

	case "SC":
		// panic("gparse.echo.L61.SC!")
		logerr.Println("Bogus token <SC>")
		return "ERR!"

	case "CD":
		return T.Otherwords

	case "PI":
		return "<?" + T.Keyword + " " + T.Otherwords + "?>"

	case "Cmt":
		return "<!-- " + T.Otherwords + " -->"

	case "DIR":
	default: // Directive subtypes, after Directives have been normalized
		return "<!" + T.Keyword + " " + T.Otherwords + ">"
	}
	return "<!-- ?! GToken.ERR ?! -->"
}

// EchoTo implements Markupper.
func (T GToken) EchoTo(w io.Writer) {
	w.Write([]byte(T.Echo()))
}

// String implements Markupper.
func (T GToken) String() string {
	return ("<!--" + T.GTagTokType.LongForm() + "-->  " + T.Echo())
}

// String implements Markupper.
func (T GToken) DumpTo(w io.Writer) {
	w.Write([]byte(T.String()))
}
