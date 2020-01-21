package gparse

var HNdTypes = []string{"nil", "Blk", "Inl", "Doc"}

// NewXmlItems (incl. ENTs, IDs, etc.)
// New GRefs, other link types
//
// XDITA / HDITA / MDITA:
// <xref> / <a href> / [link](/URI "title")
// <image>2 / <img> / ![alt text for an image](images/ image_name.jpg)
// <keydef> / <div data-class="keydef"> / MDITA-XP <div data- class="keydef"> in HDITA syntax
// <topicref> / <a href> inside a <li> / [link](/URI "title") inside a list item
// <media-source> / <source>
// @href
// @id
// @conref / @data-conref
// @keys   / @data-keys
// @keyref / @data-keyref

// println("processmkdn.go/Process:\n", p.CheckedPath.Raw)

// var e error
// var mdRoot *BF.Node
// var yamlFrontmatter []string
