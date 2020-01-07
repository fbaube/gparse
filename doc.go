// Package gparse processes markup language tokens, primarily supporting the
// three formats of LwDITA: XDITA (XML), HDITA (HTML5), MDITA (Markdown-XP).
//
// *Note: This package's inline documentation uses Markdown, so for best
// results, use (godoc2mcd)[], like so: `godoc2md `.*
//
// *Terminology: Instead of the term **parse tree**, this package uses
// the term **CST** (Concrete Syntax Tree), to contrast & compare to
// AST (Abstract Syntax Tree). See for example this introduction
// (on Wikipedia)[https://en.wikipedia.org/wiki/Abstract_syntax_tree].*
//
// This package uses (`yuin/goldmark`)[https://github.com/yuin/goldmark]
// as its Markdown parser for several reasons but mainly because work on
// v2 of the (`BlackFriday`)[https://github.com/russross/blackfriday]
// Markdown parser seems to have stalled. Also because goldmark already
// creates a CST (which goldmark calls an AST) whereas Blackfriday does
// not yet. Also goldmark is Commonmark-compliant, but this does not
// guarantee compliance with LwDITA MDITA and MDITA-XP.
//
// ### Technical Approach
//
// This package makes its own new types from Go stdlib XML structures, so that
// they get sensible new names and handy methods while retaining an ability to
// be type-cast back to the Golang stdlib equivalents.
//
// Short names (Att for *Attribute*, Elm for *Element*, Doc for *Document*)
// keep code readable.
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
