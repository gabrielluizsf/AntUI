package sdk

// JDK is a Java compiler, needed only to rebuild the Java shim. Building an
// app does not need one: the shim ships already compiled.
type JDK struct {
	Dir   string
	Javac string
}
