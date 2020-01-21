package gparse

var ExtraInfo bool

// DTDtypeFileExtensions are for content guessing.
var DTDtypeFileExtensions = []string{".dtd", ".mod", ".ent"}

// MarkdownFileExtensions are for content guessing.
// What might also work is ".m*d*"
var MarkdownFileExtensions = []string{".md", ".mdown", ".markdown", ".mkdn"}

type SystemID string
type PublicID string
