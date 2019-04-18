package gparse

import (
	"encoding/xml"
	"fmt"
	// "errors"
	"io/ioutil"
	"path"
	FP "path/filepath"
	S "strings"

	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	"github.com/pkg/errors"
)

// XmlCatalog represents a parsed XML catalog file.
type XmlCatalog struct {
	XMLName xml.Name `xml:"catalog"̀`
	// "public" or "system"
	Prefer       string        `xml:"prefer,attr"`
	XmlPublicIDs []XmlPublicID `xml:"public"`
	// We do this so we can peel off the directory path
	FU.FileFullName
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
		return nil, errors.New("Malformed Public ID<" + s + ">")
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

// NewXmlCatalogFromFile is a convenience function that reads in the file
// // and then passes the buffered file contents to the next function, below.
func NewXmlCatalogFromFile(fpath string) (pXC *XmlCatalog, err error) {
	if fpath == "" {
		return nil, nil
	}
	bb, e := ioutil.ReadFile(string(fpath))
	if e != nil {
		println("==> Can't read catalog file<", fpath, ">, reason:", e.Error())
		return nil, errors.Wrapf(e, "gxml.NewXmlCatalog.ReadFile<%s>", fpath)
	}

	// ==============================

	// gtokzn, e := GTokenizeXmlBuffer(theContent)
	var xtokens []xml.Token
	xtokens, e = XmlTokenizeBuffer(string(bb))
	if e != nil {
		return nil, errors.Wrap(e, "XmlTokenizeBuffer")
	}
	// G-Tokenize
	var gtokzn GTokenization
	gtokzn, e = MakeFromXmlTokens(xtokens)
	if e != nil {
		return nil, errors.Wrap(e, "gtoken.MakeFromXmlTokens")
	}
	var root *GToken
	var entries GTokenization
	root = gtokzn.GetFirstByTag("catalog")
	entries = gtokzn.GetAllByTag("public")
	if root == nil {
		panic("No root elm in XML catalog")
	}
	pXC = new(XmlCatalog)
	pXC.XMLName = xml.Name(root.GName)
	pXC.Prefer = root.GetAttVal("prefer")
	pXC.XmlPublicIDs = make([]XmlPublicID, 0)
	// We do this so we can peel off the directory path
	// pXC.FileFullName
	for _, GT := range entries {
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

	pXC.FileFullName = *FU.NewFileFullName(FU.RelFilePath(fpath))
	fileDir := path.Dir(pXC.FileFullName.Echo())
	println("XML catalog fileDir:", fileDir)
	for _, entry := range pXC.XmlPublicIDs {
		println("  Entry AbsFilePath:", string(entry.AbsFilePath))
		// entry.AbsFilePath = FU.AbsFilePath(path.Join(fileDir, entry.AbsFilePath.S()))
		entry.AbsFilePath = FU.AbsFilePath(fileDir + FU.PathSep + string(entry.AbsFilePath))
	}
	ok := pXC.ValidateCatalog()
	if !ok {
		panic("BAD CAT")
	}
	// println("==> Processed XML catalog at:", pXC.FileFullName.String())
	// println("TODO: insert file path for catalog file")
	return pXC, nil
}

/*
// NewXmlCatalogFromBuffer takes a string, not an io.Reader.
// Ya gotta figger it'll be useful at some point or another.
func NewXmlCatalogFromBuffer(inString FU.FileContent) (pXC *XmlCatalog, err error) {
	if inString == "" {
		return nil, nil
	}
	// println("NewXmlCatalogFromBuffer:<< ", string(inString), ">>")
	var pCat = new(XmlCatalog)
	bb := []byte(inString)
	e := xml.NewDecoder(bytes.NewReader(bb)).Decode(pCat)
	if e != nil {
		println("==> Can't process catalog file:", e.Error())
		return nil, errors.Wrap(e, "gxml.NewXmlCatalog.xmlDecode")
	}
	// print(pCat.String())
	println("parsed:<<", pCat.String(), ">>")
	// TODO TODO for each failed ID, use the func that knows the slashes
	return pCat, nil
}
*/

func NewXmlPublicIDfromGToken(pT *GToken) (pID *XmlPublicID, err error) {
	if pT == nil {
		return nil, nil
	}
	// println("NEW_XmlPublicIDfromGToken:", pT.Echo())
	// fmt.Printf("GT: %+v \n", *pT)
	NS := pT.GName.Space
	if NS != "" && NS != NS_OASIS_XML_CATALOG {
		panic("XML catalog entry has bad NS: " + NS)
	}
	// println("Space:", pT.GName.Space, "/ Local:", pT.GName.Local)
	attPid := pT.GAttList.GetAttVal("publicId")
	attUri := pT.GAttList.GetAttVal("uri")
	if attPid == "" && attUri == "" {
		println("Empty GToken for Public ID!")
		return nil, nil
	}
	// println("attPid is:", attPid)
	// println("attUri is:", attUri)

	// -//OASIS//DTD LIGHTWEIGHT DITA Topic//EN
	var ss []string
	ss = S.Split(attPid, "/")
	// fmt.Printf("(DD) (%d) %#v \n", len(ss), ss)
	ss = SU.DeleteEmptyStrings(ss)
	// {"-", "OASIS", "DTD LIGHTWEIGHT DITA Topic", "EN"}
	// fmt.Printf("(DD:PIDss) (%d) %v \n", len(ss), ss)
	if len(ss) != 4 || ss[0] != "-" || ss[3] != "EN" {
		return nil, errors.New("Malformed Public ID<" + attPid + ">")
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
func (p *XmlCatalog) ValidateCatalog() (retval bool) {
	retval = true
	for i, pEntry := range p.XmlPublicIDs {
		if "" == pEntry.PublicID {
			println("OOPS:", pEntry.String())
			panic(fmt.Sprintf("Missing Public ID in catalog entry[%d]: %s",
				i, p.FileFullName.String()))
		}
		var abspath FU.AbsFilePath
		abspath = p.FileFullName.DirPath.Append(FU.RelFilePath(pEntry.SystemID))
		pIF, e := FU.NewInputFile(FU.RelFilePath(abspath)) // downcast
		if e != nil {
			fmt.Printf("==> Catalog<%s>: Bad System ID / URI <%s> for Public ID <%s> \n",
				p.FileFullName.String(), pEntry.SystemID, pEntry.PublicID)
			retval = false
			continue
		}
		// NOTE The loop variable "entry" is by value, not reference !
		// entry.FilePath = FU.FilePath(pIF.FileFullName.String())
		p.XmlPublicIDs[i].AbsFilePath = FU.AbsFilePath(pIF.FileFullName.String())

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
			pEntry.error = errors.New("Malformed Public ID<" + s + ">")
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
