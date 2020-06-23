package gparse

import (
	"encoding/xml"
	"fmt"
	// "path"
	FP "path/filepath"
	S "strings"
	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	PU "github.com/fbaube/parseutils"
	XM "github.com/fbaube/xmlmodels"
	"github.com/fbaube/gtoken"
)

// XmlCatalogRecord represents a parsed XML catalog file, at the top level.
type XmlCatalogRecord struct {
	XMLName xml.Name `xml:"catalog"̀`
	// "public" or "system"
	Prefer  string   `xml:"prefer,attr"`
	XmlPublicIDsubrecords []XM.PIDSIDcatalogFileRecord `xml:"public"`
	// We do this so we can peel off the directory path
	FU.AbsFilePath
}

func (p *XmlCatalogRecord) GetByPublicID(s string) *XM.PIDSIDcatalogFileRecord {
	if s == "" {
		return nil
	}
	for _, entry := range p.XmlPublicIDsubrecords {
		if s == string(entry.XmlPublicID) {
			return &entry
		}
	}
	return nil
}

// NewXmlCatalogRecordFromFile is a convenience function that reads in the
// file and then processes the file contents. It is not clear what the
// constraints on the path are (but a relative path should work okay).
func NewXmlCatalogRecordFromFile(fpath string) (pXC *XmlCatalogRecord, err error) {
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

	var pCPR *PU.ConcreteParseResults_xml
	var e error
	pCPR, e = PU.GetParseResults_xml(CP.Raw)
	if e != nil {
		return nil, fmt.Errorf("gparse.xml.parseResults: %w", e)
	}
	var GTs []*gtoken.GToken
	GTs, e = gtoken.DoGTokens_xml(pCPR)
	if e != nil {
		return nil, fmt.Errorf("gparse.xml.GTokens: %w", e)
	}
	var gktnRoot *gtoken.GToken
	var gtknEntries []*gtoken.GToken
	gktnRoot    = gtoken.GetFirstByTag(GTs, "catalog")
	gtknEntries = gtoken.GetAllByTag(GTs, "public")
	if gktnRoot == nil {
		panic("No <catalog> root elm")
	}
	pXC = new(XmlCatalogRecord)
	pXC.XMLName = xml.Name(gktnRoot.GName)
	pXC.Prefer = gktnRoot.GetAttVal("prefer")
	pXC.XmlPublicIDsubrecords = make([]XM.PIDSIDcatalogFileRecord, 0)

	for _, GT := range gtknEntries {
		// println("  CAT-ENTRY:", GT.Echo()) // entry.GAttList.Echo())
		pID, e := NewSIDPIDcatalogRecordfromGToken(GT)
		// NOTE Gotta fix the filepath
		pID.AbsFilePath = // FU.AbsFilePath(
			FU.AbsWRT(string(pID.AbsFilePath), FP.Dir(string(fpath))) // )
		if e != nil {
			panic(e)
		}
		if pID == nil {
			println("Got NIL from:", GT.Echo())
		}
		pXC.XmlPublicIDsubrecords = append(pXC.XmlPublicIDsubrecords, *pID)
	}

	// ==============================

	// NOTE The following code is UGLY and needs to be FIXED.
	fileDir := pXC.AbsFilePath.DirPath()
	println("XML catalog fileDir:", fileDir)
	for _, entry := range pXC.XmlPublicIDsubrecords {
		println("  Entry's AbsFilePath:", /* FIXME:60 MU.Tilded*/ (entry.AbsFilePath))
		entry.AbsFilePath =  fileDir.S() + FU.PathSep + string(entry.AbsFilePath)
	}
	ok := pXC.Validate()
	if !ok {
		panic("BAD CAT")
	}
	// println("==> Processed XML catalog at:", pXC.FileFullName.String())
	// println("TODO: insert file path for catalog file")
	return pXC, nil
}

func NewSIDPIDcatalogRecordfromGToken(pT *gtoken.GToken) (pID *XM.PIDSIDcatalogFileRecord, err error) {
	if pT == nil {
		return nil, nil
	}
	println("L.173 NEW_XmlPublicIDfromGToken:", pT.Echo())
	fmt.Printf("L.174 GT: %+v \n", *pT)
	NS := pT.GName.Space
	if NS != "" && NS != XM.NS_OASIS_XML_CATALOG {
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
	pID = new(XM.PIDSIDcatalogFileRecord)
	// NOTE DANGER This is probably relative not absolute,
	// and has to be fixed by the caller
	pID.XmlPublicID = XM.XmlPublicID(attPid)
	pID.XmlSystemID = XM.XmlSystemID(attUri)
	pID.AbsFilePath = attUri // FU.AbsFilePath(attUri)
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
func (p *XmlCatalogRecord) Validate() (retval bool) {
	retval = true
	for i, pEntry := range p.XmlPublicIDsubrecords {
		if "" == pEntry.XmlPublicID {
			println("OOPS:", pEntry.String())
			panic(fmt.Sprintf("Missing Public ID in catalog entry[%d]: %s",
				i, p.AbsFilePath)) // Parts.String()))
		}
		var abspath FU.AbsFilePath
		abspath = p.AbsFilePath.DirPath().Append(string(pEntry.XmlSystemID))
		// pIF, e := FU.NewInputFile(FU.RelFilePath(abspath)) // downcast
		pIF := FU.NewCheckedContentFromPath(abspath.S())
		// (&FU.CheckedPath{RelFilePath: abspath.AsRelFP()}).Resolve() //.Check()
		if !pIF.IsOkayFile() { // pIF.PathType() != "FILE" { // e != nil {
			fmt.Printf("==> Catalog<%s>: Bad System ID / URI <%s> for Public ID <%s> \n",
				p.AbsFilePath, pEntry.XmlSystemID, pEntry.XmlPublicID)
			retval = false
			continue
		}
		// NOTE The loop variable "entry" is by value, not reference !
		// entry.FilePath = FU.FilePath(pIF.FileFullName.String())
		p.XmlPublicIDsubrecords[i].AbsFilePath =
			/*FU.AbsFilePath(*/ pIF.AbsFilePath//)

		// Now do some fancy parsing of the Public ID
		var s = string(pEntry.XmlPublicID)
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
			pEntry.Err = fmt.Errorf("Malformed Public ID: %s", s)
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
