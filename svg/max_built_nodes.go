package svg

// maxBuiltNodes is how many nodes one drawing may build before the rest of it
// is left out. A `<use>` copies what it names, and a `<use>` of a group of
// `<use>`s copies all of them: ten of them a level deep and six levels down is
// a million shapes out of a file of a few hundred bytes, the same trick as
// the XML billion laughs played through instancing (librsvg's
// CVE-2019-20446). A real drawing — a map, a chart, an icon sheet — is
// thousands of nodes; one past this many is built to run the machine out of
// memory rather than to be looked at.
const maxBuiltNodes = 200_000
