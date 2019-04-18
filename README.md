Package gparse processes Golang markup language tokens: initially
XML, but with support for other flavors of LwDITA markup RSN (HTML5,
Markdown-XP). 

Files in this directory use Markdown, so use `godoc2md` on'em.

This package makes own-versions of Golang XML structures so that
they get sensible new names and handy methods.

Two shortened names (Att *Attribute*, Elm *Element*) keep code readable.

### About XML content, including mixed content

When working with XML we can generally distinguish three types of files:
- Record-oriented XML data - expressed using XML elements
- Natural language XML documents - also expressed using XML elements,
and known as **mixed content**
- Validation rules - generally expressed as XSD, RNG, or DTD. It is
interesting to note that DTDs actually obey the same fundamental XML
syntax rules as the other two types (record-oriented, mixed content);
the typical DTD file extensions (`.dtd .mod`) are helpful to humans
but are not strictly required as a signal to a parser that fully
understands the syntax of all three XML file types.

That being said, this package can superficially digest many directives
(e.g. ELEMENT, ATTLIST, ENTITY) but does not yet (at this level)
completely parse them, or act on them (by performing transclusion).
