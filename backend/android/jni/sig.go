package jni

import "strings"

// The type letters a JNI signature is built from. They are the JVM's own
// spelling, which is the only spelling: a signature is not parsed leniently,
// and one wrong letter finds no method rather than the wrong one — which is
// the good case. The bad case is a signature that matches something else.
//
// This file has no build tag on purpose. Building a signature is string
// work, it is where the mistakes are, and it can be tested on any machine.
const (
	TVoid   = "V"
	TBool   = "Z"
	TByte   = "B"
	TChar   = "C"
	TShort  = "S"
	TInt    = "I"
	TLong   = "J"
	TFloat  = "F"
	TDouble = "D"

	// The three classes that come up in almost every signature.
	TString  = "Ljava/lang/String;"
	TObject  = "Ljava/lang/Object;"
	TContext = "Landroid/content/Context;"
)

// TClass is the type of a class, named with slashes: TClass("android/net/Uri")
// is "Landroid/net/Uri;". A name given with dots is taken to be a mistake and
// converted, because the two spellings are used in different places and
// getting them the wrong way round is the most common way to write a
// signature that finds nothing.
func TClass(name string) string {
	return "L" + strings.ReplaceAll(name, ".", "/") + ";"
}

// TArray is the type of an array of something: TArray(TInt) is "[I".
func TArray(of string) string { return "[" + of }

// Sig builds a method signature from what it returns and what it takes:
//
//	Sig(TVoid, TInt, TString)            →  "(ILjava/lang/String;)V"
//	Sig(TClass("android/net/Uri"))       →  "()Landroid/net/Uri;"
//	Sig(TBool, TArray(TString))          →  "([Ljava/lang/String;)Z"
//
// The return type comes first because that is the order it is spoken in,
// even though the signature writes it last.
func Sig(ret string, args ...string) string {
	var b strings.Builder
	b.Grow(len(ret) + 2 + 16*len(args))
	b.WriteByte('(')
	for _, a := range args {
		b.WriteString(a)
	}
	b.WriteByte(')')
	b.WriteString(ret)
	return b.String()
}
