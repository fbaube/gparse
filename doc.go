// Package gparse processes Golang markup language tokens, primarily
// supporting the three formats of LwDITA: XDITA (XML), HDITA (HTML5),
// MDITA (Markdown-XP).
//
// This package uses the BlackFriday (v2) Markdown parser for two reasons:
// v2 has been updated to produce ASTs, and v1 is used by Hugo. Conformance
// with LwDITA MDITA and MDITA-XP is not guaranteed - not at all.
//
// Golang comments in this directory use Markdown, so use `godoc2md` on'em.
//
// ### Technical Approach
//
// This package makes its own versions of Golang XML structures so that they
// get sensible new names and handy methods, while retaining type-cast-ability.
//
// Short names (Att for *Attribute*, Elm for *Element*) keep code readable.
//
// This code *should* work with XML namespaces, but this is completely untested.
//
// ### Method naming
//
// - `NewFoo(..)` always allocates new memory and returns a pointer.
// - `Echo()` echoes an object back in source XML form, but normalized.
// - `String()` outputs a human-friendly form useful for development
// and debugging but usually indigestible to an XML parser.
//
// ### About XML content, including mixed content
//
// When working with XML we can generally distinguish three types of files:
// - Record-oriented XML data - expressed using XML elements
// - Natural language XML documents - also expressed using XML elements,
// and known as **mixed content**
// - Validation rules - generally expressed as XSD, RNG, or DTD. It is
// interesting to note that DTDs actually obey the same fundamental XML
// syntax rules as the other two types (record-oriented, mixed content);
// the typical DTD file extensions (`.dtd .mod`) are helpful to humans
// but are not strictly required as a signal to a parser that fully
// understands the syntax of all three XML file types (and all types
// of XML entities)
//
// That being said, this package can superficially digest many directives
// (e.g. ELEMENT, ATTLIST, ENTITY) but does not yet (at this level)
// completely parse them, or act on them (by performing transclusion).
//
package gparse
