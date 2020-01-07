package gparse

import (
	"encoding/xml"
	"fmt"
	"path"
	FP "path/filepath"
	S "strings"
	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
)

// XmlCatalog represents a parsed XML catalog file.
type XmlCatalog struct {
	XMLName xml.Name `xml:"catalog"̀`
	// "public" or "system"
	Prefer       string        `xml:"prefer,attr"`
	XmlPublicIDs []XmlPublicID `xml:"public"`
	// We do this so we can peel off the directory path
	FU.AbsFilePathParts
}

// XmlPublicID representa a line item from a parsed XML catalog file.
type XmlPublicID struct {
	XMLName  xml.Name `xml:"public"`
	PublicID string   `xml:"publicId,attr"`
	// The short form, as declared in the catalog file.
	// TODO It's a string, but let's hope it gets decoded OK.
	SystemID `xml:"uri,attr"`
	// The long form, as resolved
	FU.AbsFilePath
	error // in case an entry barfs

	// "+" or "-"
	Registration string
	// If not OASIS then could be any one of a myriad of others.
	IsOasis bool
	// "OASIS" or something else
	Organization string
	// PTClass is typically "DTD" (filename.dtd)
	// or "ELEMENTS" (filename.mod)
	PublicTextClass string
	// PTDesc is the distinguishing string, e.g.
	// PUBLIC "-//OASIS//DTD (PublicTextDesc)//EN".
	// Sometimes it can have an optional embedded
	// version number, such as "DITA 1.3".
	PublicTextDesc string
}

func NewXmlPublicID(s string) (*XmlPublicID, error) {
	// println("NewXmlPublicID:", s)
	if s == "" {
		return nil, nil
	}
	if SU.IsXmlQuoted(s) {
		s = SU.MustXmlUnquote(s)
	}
	// -//OASIS//DTD LIGHTWEIGHT DITA Topic//EN
	var ss []string
	ss = S.Split(s, "/")
	// fmt.Printf("(DD) (%d) %#v \n", len(ss), ss)
	ss = SU.DeleteEmptyStrings(ss)
	// {"-", "OASIS", "DTD LIGHTWEIGHT DITA Topic", "EN"}
	// fmt.Printf("(DD) (%d) %#v \n", len(ss), ss)
	if len(ss) != 4 || ss[0] != "-" || ss[3] != "EN" {
		return nil, fmt.Errorf("Malformed Public ID<" + s + ">")
	}
	p := new(XmlPublicID)
	p.Organization = ss[1]
	p.IsOasis = ("OASIS" == p.Organization)
	p.PublicTextClass, p.PublicTextDesc = SU.SplitOffFirstWord(ss[2])
	// ilog.Printf("PubID|%s|%s|%s|\n", p.Organization, p.PTClass, p.PTDesc)
	// fmt.Printf("(DD:pPID) PubID<%s|%s|%s>\n", p.Organization, p.PTClass, p.PTDesc)
	return p, nil
}

func (p *XmlCatalog) GetByPublicID(s string) *XmlPublicID {
	if s == "" {
		return nil
	}
	for _, entry := range p.XmlPublicIDs {
		if s == entry.PublicID {
			return &entry
		}
	}
	return nil
}

// NewXmlCatalogFromFile is a convenience function that reads in the
// file and then processes the file contents. It is not clear what the
// constraints on the path are (but a relative path should work okay).
func NewXmlCatalogFromFile(fpath string) (pXC *XmlCatalog, err error) {
	if fpath == "" {
		return nil, nil
	}
	var CP *FU.CheckedContent
	CP = FU.NewCheckedContentFromPath(fpath)
	if CP.GetError() != nil {
		println("==> Can't read catalog file:", fpath, ", reason:", CP.Error())
		return nil, fmt.Errorf("gparse.NewXmlCatalog.ReadFile<%s>: %w",
			fpath, CP.GetError())
	}

	// ==============================

	var xtokens []xml.Token
	var e error
	xtokens, e = XmlTokenizeBuffer(CP.Raw)
	if e != nil {
		return nil, fmt.Errorf("XmlTokenizeBuffer: %w", e)
	}
	var gtokzn GTokenization
	gtokzn, e = GTokznFromXmlTokens(xtokens)
	if e != nil {
		return nil, fmt.Errorf("gtoken.NewGtokznFromXmlTokens: %w", e)
	}
	var gktnRoot *GToken
	var gtknEntries GTokenization
	gktnRoot = gtokzn.GetFirstByTag("catalog")
	gtknEntries = gtokzn.GetAllByTag("public")
	if gktnRoot == nil {
		panic("No <catalog> root elm")
	}
	pXC = new(XmlCatalog)
	pXC.XMLName = xml.Name(gktnRoot.GName)
	pXC.Prefer = gktnRoot.GetAttVal("prefer")
	pXC.XmlPublicIDs = make([]XmlPublicID, 0)
	// We do this so we can peel off the directory path
	// pXC.FileFullName
	for _, GT := range gtknEntries {
		// println("  CAT-ENTRY:", GT.Echo()) // entry.GAttList.Echo())
		pID, e := NewXmlPublicIDfromGToken(GT)
		// NOTE Gotta fix the filepath
		pID.AbsFilePath = FU.AbsFilePath(
			FU.AbsWRT(string(pID.AbsFilePath), FP.Dir(string(fpath))))
		if e != nil {
			panic(e)
		}
		if pID == nil {
			println("Got NIL from:", GT.Echo())
		}
		pXC.XmlPublicIDs = append(pXC.XmlPublicIDs, *pID)
	}

	// ==============================

	// NOTE: The following code is UGLY and needs to be FIXED.
	pXC.AbsFilePathParts = *FU.AbsFP(fpath).NewAbsPathParts()
	fileDir := path.Dir(pXC.AbsFilePathParts.Echo())
	println("XML catalog fileDir:", fileDir)
	for _, entry := range pXC.XmlPublicIDs {
		println("  Entry's AbsFilePath:", /* FIXME MU.Tilded */ (entry.AbsFilePath.S()))
		// entry.AbsFilePath = FU.AbsFilePath(path.Join(fileDir, entry.AbsFilePath.S()))
		entry.AbsFilePath = FU.AbsFilePath(fileDir + FU.PathSep + string(entry.AbsFilePath))
	}
	ok := pXC.Validate()
	if !ok {
		panic("BAD CAT")
	}
	// println("==> Processed XML catalog at:", pXC.FileFullName.String())
	// println("TODO: insert file path for catalog file")
	return pXC, nil
}

func NewXmlPublicIDfromGToken(pT *GToken) (pID *XmlPublicID, err error) {
	if pT == nil {
		return nil, nil
	}
	println("L.173 NEW_XmlPublicIDfromGToken:", pT.Echo())
	fmt.Printf("L.174 GT: %+v \n", *pT)
	NS := pT.GName.Space
	if NS != "" && NS != NS_OASIS_XML_CATALOG {
		panic("XML catalog entry has bad NS: " + NS)
	}
	println("L.179 Space:", pT.GName.Space, "/ Local:", pT.GName.Local)
	attPid := pT.GAtts.GetAttVal("publicId")
	attUri := pT.GAtts.GetAttVal("uri")
	if attPid == "" && attUri == "" {
		println("Empty GToken for Public ID!")
		return nil, nil
	}
	println("L.186 attPid is:", attPid)
	println("L.187 attUri is:", attUri)

	// -//OASIS//DTD LIGHTWEIGHT DITA Topic//EN
	var ss []string
	ss = S.Split(attPid, "/")
	// fmt.Printf("(DD) (%d) %#v \n", len(ss), ss)
	ss = SU.DeleteEmptyStrings(ss)
	// {"-", "OASIS", "DTD LIGHTWEIGHT DITA Topic", "EN"}
	// fmt.Printf("(DD:PIDss) (%d) %v \n", len(ss), ss)
	if len(ss) != 4 || ss[0] != "-" || ss[3] != "EN" {
		return nil, fmt.Errorf("Malformed Public ID<%s>", attPid)
	}
	pID = new(XmlPublicID)
	// NOTE DANGER This is probably relative not absolute,
	// and has to be fixed by the caller
	pID.PublicID = attPid
	pID.SystemID = SystemID(attUri)
	pID.AbsFilePath = FU.AbsFilePath(attUri)
	pID.Organization = ss[1]
	pID.IsOasis = ("OASIS" == pID.Organization)
	pID.PublicTextClass, pID.PublicTextDesc = SU.SplitOffFirstWord(ss[2])
	// ilog.Printf("PubID|%s|%s|%s|\n", p.Organization, p.PTClass, p.PTDesc)
	// fmt.Printf("(DD:pPID) PubID<%s|%s|%s> AFP<%s>\n",
	//  	pID.Organization, pID.PublicTextClass,
	//		pID.PublicTextDesc, pID.AbsFilePath)
	return pID, nil
}

// Validate validates an XML catalog. It checks that the listed files exist
// and that the IDs (as strings that are not parsed yet) are well-formed.
// It assumes that the catalog has already been loaded from an XML catalog
// file on-disk. The return value is false if _any_ entry fails to load,
// but also each entry has its own error field.
func (p *XmlCatalog) Validate() (retval bool) {
	retval = true
	for i, pEntry := range p.XmlPublicIDs {
		if "" == pEntry.PublicID {
			println("OOPS:", pEntry.String())
			panic(fmt.Sprintf("Missing Public ID in catalog entry[%d]: %s",
				i, p.AbsFilePathParts.String()))
		}
		var abspath FU.AbsFilePath
		abspath = p.AbsFilePathParts.DirPath.Append(string(pEntry.SystemID))
		// pIF, e := FU.NewInputFile(FU.RelFilePath(abspath)) // downcast
		pIF := FU.NewCheckedContentFromPath(abspath.S())
		// (&FU.CheckedPath{RelFilePath: abspath.AsRelFP()}).Resolve() //.Check()
		if !pIF.IsOkayFile() { // pIF.PathType() != "FILE" { // e != nil {
			fmt.Printf("==> Catalog<%s>: Bad System ID / URI <%s> for Public ID <%s> \n",
				p.AbsFilePathParts.String(), pEntry.SystemID, pEntry.PublicID)
			retval = false
			continue
		}
		// NOTE The loop variable "entry" is by value, not reference !
		// entry.FilePath = FU.FilePath(pIF.FileFullName.String())
		p.XmlPublicIDs[i].AbsFilePath = FU.AbsFilePath(pIF.AbsFilePathParts.String())

		// Now do some fancy parsing of the Public ID
		var s = pEntry.PublicID
		if SU.IsXmlQuoted(s) {
			s = SU.MustXmlUnquote(s)
		}
		// -//OASIS//DTD LIGHTWEIGHT DITA Topic//EN
		var ss []string
		ss = S.Split(s, "/")
		// fmt.Printf("(DD) (%d) %#v \n", len(ss), ss)
		ss = SU.DeleteEmptyStrings(ss)
		// {"-", "OASIS", "DTD LIGHTWEIGHT DITA Topic", "EN"}
		// fmt.Printf("(DD) (%d) %#v \n", len(ss), ss)
		if len(ss) != 4 || ss[0] != "-" || ss[3] != "EN" {
			retval = false
			pEntry.error = fmt.Errorf("Malformed Public ID: %s", s)
			continue
		}
		pEntry.Organization = ss[1]
		pEntry.IsOasis = ("OASIS" == pEntry.Organization)
		pEntry.PublicTextClass, pEntry.PublicTextDesc = SU.SplitOffFirstWord(ss[2])
		// ilog.Printf("PubID|%s|%s|%s|\n", p.Organization, p.PTClass, p.PTDesc)
		// fmt.Printf("(DD:pPID) PubID<%s|%s|%s>\n", p.Organization, p.PTClass, p.PTDesc)
	}
	return true
}

// String returns the public ID _unquoted_.
// <!DOCTYPE topic "-//OASIS//DTD LIGHTWEIGHT DITA Topic//EN">
func (p XmlPublicID) String() string {
	return fmt.Sprintf("%s//%s//%s %s//EN",
		p.Registration, p.Organization, p.PublicTextClass, p.PublicTextDesc)
}
