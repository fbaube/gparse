package gparse

// This file: About the "xml" namespace

// TODO NS is specified, inherited, default

// NS_XML is the XML namespace.
var NS_XML = "http://www.w3.org/XML/1998/namespace"

// NS_XML is the OASIS namespace for XML catalogs.
var NS_OASIS_XML_CATALOG = "urn:oasis:names:tc:entity:xmlns:xml:catalog"

// WARNING: Go has lotsa XML namespace problems:
// https://github.com/golang/go/issues/13400#issuecomment-162459219

// XML_NS_Recognized is recognized values in the XML namespace.
var XML_NS_Recognized = []string{
	// `lang` identifies the human language used in the
	// scope of the element to which it's attached.
	"lang",
	// `space (default|preserve)´ says whether white space
	// to be considered as significant in the scope of the
	// element to which it's attached.
	"space",
	// The XML Base spec (Second edition) describes a facility,
	// like HTML BASE, for defining base URIs for parts of XML
	// documents. It defines xml:base, and describes in detail
	// how to use it in processing relative URI references.
	"base",
	// The xml:id spec defines xml:id, which is of type "ID".
	// independently of any DTD or schema.
	"id",
	// "Father" denotes Jon Bosak, the chair of the original XML WG:
	// "In appreciation for his vision, leadership and dedication
	// the W3C XML Plenary today 10.02.2000 reserves for Jon Bosak
	// in perpetuity the XML name "xml:Father".
	"Father",
}

// More info:
// - https://www.w3.org/2001/xml.xsd
// - https://www.w3.org/TR/xmlbase/

// NOTE also: other names beginning 'xml'.
// The XML spec reserves all names beginning with the letters 'xml'
// in any combination of upper & lower -case for use by the W3C.
// To date three such names have been given definitions — although
// these names are not in the XML namespace, they are listed here
// as a convenience:
// - "xml": See the XML declaration and the XML namespace prefix.
// - "xmlns": See namespace declarations.
// - xml-stylesheet: See the xml-stylesheet processing instruction.

// NOTE about "xml:id": https://www.w3.org/TR/xml-id/ :
// The normalized value of the attribute is an NCName
// IAW "Namespaces in XML Recommendation" i.e. NCName
// for XML 1.0 (or NCName for XML 1.1).
//
// D.1 With DTD Validation:
// DTD authors are encouraged to use xml:id attributes when
// providing identifiers for elements declared in their DTDS.
// The following DTD fragment illustrates a sample declaration
// for the xml:id attribute:
//
// <!ATTLIST  someElement  xml:id  ID  #IMPLIED >
//
// DTD authors are encouraged to declare attributes named xml:id
// with the type ID. A document that uses xml:id attributes that
// have a declared type other than ID will always generate xml:id errors.
//
// Consumers of documents validated using properly declared xml:id
// attributes can recognize IDs through the attribute type property.
//
// Name productions:
// NCName	 ::=  Name - (Char* ':' Char*)  // An XML Name, minus the ":"
//   Name  ::=  NameStartChar (NameChar)*
// NameStartChar ::= ":" | [A-Z] | "_" | [a-z] | [#xC0-#xD6] | [#xD8-#xF6] | [#xF8-#x2FF] | [#x370-#x37D] | [#x37F-#x1FFF] | [#x200C-#x200D] | [#x2070-#x218F] | [#x2C00-#x2FEF] | [#x3001-#xD7FF] | [#xF900-#xFDCF] | [#xFDF0-#xFFFD] | [#x10000-#xEFFFF]
//      NameChar ::= NameStartChar | "-" | "." | [0-9] | #xB7 | [#x0300-#x036F] | [#x203F-#x2040]
//
// https://www.w3.org/TR/xml-id/#dt-xml-id-error
// Definition: An xml:id error is a non-fatal error that occurs when an xml:id
// processor finds that a document has violated the constraints of this spec.
//
// https://www.w3.org/TR/xml-names/
