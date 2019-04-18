package gparse

// This file: Generic Golang XML Entities ("ents")

// GEnt is a generic XML Entity. Either
// - an Internal Parsed General entity:   <!ENTITY foo "bar">  OR
// - an Internal Parsed Parameter entity: <!ENTITY % foo "bar">
//
// "bar" may not use any of:
// '&', '%', '"', '%Name;', '&Name;', Unicode char ref. (Sez who ?)
//
type GEnt struct {
	// e.g. "foo"
	NameOnly string
	// including "%|&" and ";" i.e. "&foo;" or "%foo;"
	NameAsRef string
	// true if parameter entity, false if general entity
	TypeIsParm bool
	// "%" if parameter entity, "&" if general entity
	RefChar string
	// External entities only (PUBLIC, SYSTEM)
	IsSystemID bool
	IsPublicID bool
	ID         string
	URI        string
	Fullpath   string

	TheRest string
}

func (ge GEnt) String() string {
	return "[GEnt]"
}
