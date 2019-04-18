package gparse

// Please refer to the Annotated XML Specification
// http://www.xml.com/axml/testaxml.htm
//
// This file requires some explanation. When considering DOCTYPES, the cases
// that most interest us all involve OASIS ExternalID's for DITA and LwDITA,
// and these are mostly PUBLIC ExternalID's. This yield several requirements.
//
// When parsing an XML document, we have a DOCTYPE declaration like this:
// <!DOCTYPE topic PUBLIC "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN" "lw-topic.dtd">
// "topic" is a root element ("RootTag") and the rest is an ExternalID.
//
// When parsing a DTD, we're sposta use an XML Catalog, but instead we will
// take advantage of an OASIS documentation convention to look for a string
// like one of these in a comment:
// .dtd: [PUBLIC] "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN"
// .mod: [PUBLIC] "-//OASIS//ELEMENTS LIGHTWEIGHT DITA Topic//EN"
// Each is an ExternalID.
//
// Therefore our "strings of interest" are PUBLIC ID's, which have the form
// PUBLIC "-//OASIS//(PTClass) (PTDescription)//EN" ["optional SystemLiteral"]
// where PT = "Public Text", PTClass is "DTD" or "ELEMENTS", and PTDesc will
// tell us about DITA or LwDITA (or possibly even something else).
//
// To be more exact, an ExternalID is one of these:
// SYSTEM SystemLiteral (or)
// PUBLIC PubidLiteral SystemLiteral
// where
// SystemLiteral is a quoted string containing anything except quote characters.
// PuidLiteral contains CR, LF, " ", alphanumeric, and any of [-'()+,./:=?;!*#@$_%].

// XmlInfo represents and stores :ul:
// :: XML preamble ("<?xml ...")
// :: DOCTYPE declaration
// :: ?? File input options (taken from TODO) (see Go xml.Decoder)
// :: XML content type: DTD, (single) RootTag, Fragments (i.e. multiple
// "root" tags), MixedContent (a subtype of RootTag), INVALID
// -ul-
// The latest LwDITA values (2017.08):
// PUBLIC "-//OASIS//DTD LIGHTWEIGHT DITA (Map|Topic)//EN"
//
type XmlFileMeta struct {
	XmlContype
	// XmlPreamble (i.e. <?xml ...?> ) might be empty but it's not worth the
	// trouble to make it a pointer and anyways if the file does not have one
	// then just copy in the default preamble defined in the Golang library.
	XmlPreamble
	// TagDefCt is for DTD-type files (.dtd, .mod)
	TagDefCt int // Nr of <!ELEMENT ...>
	// XmlDoctype is non-nil IFF a DOCTYPE directive was encountered
	*XmlDoctype
	// CmtdOasisPublicID *ParsedPublicID

	// RootTagIndex int  // Or some sort of pointer to the tree
	// If RootTagCt is >1, the content is a Fragment
	RootTagCt int // If
}

func (xfm XmlFileMeta) Echo() string {
	return xfm.XmlPreamble.Echo() + xfm.XmlDoctype.Echo()
}

func (xfm XmlFileMeta) String() string {
	var s1, s2 = "[NO_xmlPreamble]", "[NO_xmlDoctype]"
	if "" != xfm.XmlPreamble.Raw {
		s1 = xfm.XmlPreamble.String()
	}
	if xfm.XmlDoctype != nil && "" != xfm.XmlDoctype.raw {
		s2 = xfm.XmlDoctype.String()
	}
	return s1 + s2
}
