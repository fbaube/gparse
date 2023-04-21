package gparse

// We want to present a single interface for tokens produced by parsing.
// This means wrapping the tokens we get from the various parsing libraries.
// Let's start by examining our raw materials. Note that we do not attach
// tree-structure information to tokens; this is done via `GTag`s in a `GTree`.
//
// "XML": [XU.XToken] (and specifically, [encoding/xml.StartElement]):
//   type StartElement struct {
//     Name Name
//     Attr []Attr
//   }
//
// "MKDN": `github.com/yuin/goldmark/ast/Node`(excluding tree-related stuff):
//   Type() NodeType (RBIN, see below)
//   Kind() NodeKind (the various rich text tags)
//   Text(source []byte) []byte
//   Attributes() []Attribute
//
// "HTML": `net/html`:
//   type Token struct {
//     Type     TokenType
//     DataAtom atom.Atom
//     Data     string
//     Attr     []Attribute
//   }
//
// Fields we definitely need in EACH wrapped token pretype:
// - `GName` (i.e. the tag, and can include NS = Namespace)
// - `GAtts` (to generalize XML & HTML attributes and Markdown properties;
// make it a map so that there are no accidental type conversions)
// - `Depth` (assumed to start at 0 for Doc/Root, but the exact starting
// index is not critical unless we mix token families)
// - `TTType` (i.e. Token/Tag Type)
// - `DispFmtgType` (Display Formatting Type: Root, Blck, Inln, None)
// - `Keyword` & `Otherwords` (for XML stuff and other in-content meta-info)
//
// So each token pretype is wrapped in order to deliver these fields.
