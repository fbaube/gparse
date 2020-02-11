package gparse

// TTType specifies the type an input token or markup tag. Values are based
// on the tokens output'd by the stdlib `xml.Decoder`, with some additions.
type TTType string

// GTagTokTypes is CDATA, ID/REF, etc., plus a reserved/TBD entry for "enum".
// NOTE that these strings are used in comments thruout this package.
var TTTypes = []TTType{
	"nilerror",
	"SE",  // StartElement  // could be:  "<s>"
	"EE",  // EndElement    // could be: "</s>" // Not used in GTokens & GTrees
	"SC",  // SelfClosingTag// could be: "<s/>" // Usage is unclear / TBD
	"CD",  // CDATA         // could be: "s"
	"PI",  // Proc. Instr.  // could be: "<?s?>"
	"Cmt", // XML comment   // could be: "<!--s-->"
	"Dir", // XML directive // could be: "<!s>"
	// The following are actually DIRECTIVE SUBTYPES, but they
	// are put in this list so that they can be assigned freely.
	"DOCTYPE",
	"ELEMENT",
	"ATTLIST",
	"ENTITY",
	"NOTATION",
	// The following are TBD.
	"ID",
	"IDREF",
	"ENUM",
}

func (TT TTType) LongForm() string {
	switch TT {
	case "SE":
		return "Start-Tag"
	case "EE":
		return "End-Tagnt"
	case "CD":
		return "Char-Data"
	case "Cmt":
		return " Comment "
	case "PI":
		return "ProcInstr"
	case "Dir":
		return "Directive"
	case "SC":
		return "SelfClose"
	case "Doc":
		return "DocuStart"
	}
	return string(TT) 
}
