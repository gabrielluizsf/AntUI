package shim

// SourceLevel is the Java version the shim is compiled to. It is old on
// purpose: nothing here needs anything newer, and the older the class file
// the wider the range of build tools that will dex it.
const SourceLevel = "8"
