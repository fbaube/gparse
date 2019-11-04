package gparse

import "fmt"

// ExtraInfo is a very coarse flag for constraining output when multiple files
// are being processed.
var ExtraInfo bool

type strattribute struct{ Name, Value string }

func (SA strattribute) String() string {
	return fmt.Sprintf("%s<%s> ", SA.Name, SA.Value)
}

/*
var logerr *log.Logger
func init() {
	logerr = log.New(os.Stderr, "ERR:gtoken> ", log.Lshortfile)
}
*/

// DTDtypeFileExtensions are for content guessing.
var DTDtypeFileExtensions = []string{".dtd", ".mod", ".ent"}

// MarkdownFileExtensions are for content guessing.
// What might also work is ".m*d*"
var MarkdownFileExtensions = []string{".md", ".mdown", ".markdown", ".mkdn"}

type SystemID string
type PublicID string
