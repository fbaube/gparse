package gparse

// DispFmtgType specifies the "rendering context".
type DispFmtgType string

// DispFmtgTypes is CDATA, ID/REF, etc., plus a reserved/TBD entry for "enum".
// NOTE that these strings are used in comments thruout this package.
var DispFmtgTypes = []DispFmtgType{
	"nilerror",
	"ROOT",  // Document root
	"BLCK",  //
	"INLN",  //
	"NONE",  //
}

func (DFT DispFmtgType) LongForm() string {
	switch DFT {
	case "ROOT":
		return "Doc-root"
	case "BLCK":
		return "BlockCtx"
	case "INLN":
		return "InlineCtx"
	case "NONE":
		return "N/A-None"
	}
	return "DispFmtgType-LongForm-ERROR"
}
