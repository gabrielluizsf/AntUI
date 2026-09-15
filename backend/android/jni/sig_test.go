package jni

import "testing"

func TestSig(t *testing.T) {
	for _, c := range []struct {
		got, want string
	}{
		{Sig(TVoid), "()V"},
		{Sig(TVoid, TInt), "(I)V"},
		{Sig(TVoid, TInt, TString), "(ILjava/lang/String;)V"},
		{Sig(TString), "()Ljava/lang/String;"},
		{Sig(TBool, TArray(TString)), "([Ljava/lang/String;)Z"},
		{Sig(TArray(TInt), TLong, TDouble), "(JD)[I"},
		{Sig(TClass("android/net/Uri"), TString), "(Ljava/lang/String;)Landroid/net/Uri;"},
		// A constructor: void, and the name is what makes it one.
		{Sig(TVoid, TContext), "(Landroid/content/Context;)V"},
	} {
		if c.got != c.want {
			t.Errorf("got %q, want %q", c.got, c.want)
		}
	}
}

// A class name is written with slashes in a signature and with dots
// everywhere else, and the two are used in different halves of the same API.
// Taking a dotted name and correcting it is worth more than refusing it: the
// signature that results from getting it wrong finds no method, and the
// error says nothing about why.
func TestTClassTakesEitherSpelling(t *testing.T) {
	want := "Landroid/net/Uri;"
	if got := TClass("android/net/Uri"); got != want {
		t.Errorf("slashes gave %q", got)
	}
	if got := TClass("android.net.Uri"); got != want {
		t.Errorf("dots gave %q, want %q", got, want)
	}
}

func TestTArrayNests(t *testing.T) {
	if got := TArray(TArray(TByte)); got != "[[B" {
		t.Errorf("an array of arrays of bytes is %q, want [[B", got)
	}
	if got := TArray(TClass("java/lang/String")); got != "[Ljava/lang/String;" {
		t.Errorf("an array of strings is %q", got)
	}
}
