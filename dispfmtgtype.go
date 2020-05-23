package gparse

// DispFmtgType (display formatting type) specifies the "rendering context".
type DispFmtgType string

// DispFmtgTypes specifies how an element fits into layout. 
var DispFmtgTypes = []DispFmtgType{
	"nilerror",
	"ROOT",  // Document root
	"BLCK",  //
	"INLN",  //
	"NONE",  //
}

// LongForm returns a marginally-more-user-frenly description.
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
