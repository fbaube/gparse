package gparse

import (
	"log"
	"os"
)

var logerr *log.Logger

// ExtraInfo is a very coarse flag for constraining output when multiple files
// are being processed.
var ExtraInfo bool

// Ldate       // the date in the local time zone: 2009/01/23
// Ltime       // the time in the local time zone: 01:23:23
// Lmicrosec's // microsecond resolution: 01:23:23.123123.  assumes Ltime.
// Llongfile   // full file name and line number: /a/b/c/d.go:23
// Lshortfile  // basic file name and line number: d.go:23 ; overrides Llongfile
// LUTC        // if Ldate or Ltime is set, use UTC, not local time zone

func init() {
	logerr = log.New(os.Stderr, "ERR:gtoken> ", log.Lshortfile)
}

// XmlContype describes the top-down XML structure. Its use is TBS.
type XmlContype string

// XmlContypes, maybe DTDmod should be DTDelms.
var XmlContypes = []XmlContype{"Unknown", "DTD", "DTDmod", "DTDent",
	"RootTagData", "RootTagMixedContent", "MultipleRootTags", "INVALID"}

// DTDtypeFileExtensions are for content guessing.
var DTDtypeFileExtensions = []string{".dtd", ".mod", ".ent"}

// MarkdownFileExtensions are for content guessing.
// What might also work is ".m*d*"
var MarkdownFileExtensions = []string{".md", ".mdown", ".markdown", ".mkdn"}

type SystemID string
type PublicID string
