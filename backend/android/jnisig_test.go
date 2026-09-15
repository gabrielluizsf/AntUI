package android_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

// A JNI call has to be made with the return type its signature declares.
//
// This is the one class of mistake in this whole tree that nothing catches.
// It is not a compile error — the signature is a string and the call is a
// method name, and Go has no idea they are meant to agree. It is not a
// runtime error either: calling a method that returns an int as though it
// returned nothing works, on this ABI, on this VM, today. The value is left
// where nobody looks and the frame happens to unwind the same way.
//
// It is undefined behaviour, and the only thing that says so is the VM's own
// CheckJNI, which is off unless somebody turns it on:
//
//	adb shell setprop debug.checkjni 1 && adb shell stop && adb shell start
//
// With it on, the process aborts with "the return type of CallVoidMethodA
// does not match int android.speech.tts.TextToSpeech.
// setOnUtteranceProgressListener". That is how the one this test was written
// for was found, and turning CheckJNI on is worth doing now and then for the
// mistakes a static check cannot see. This catches the mismatch that a
// static check can, on every ordinary run.
func TestJNICallsMatchTheirSignatures(t *testing.T) {
	// What each call promises the method returns. jni.Sig's first argument
	// is the return type, so the two have to agree.
	expect := map[string]string{
		"CallVoid": "TVoid", "CallStaticVoid": "TVoid",
		"InvokeVoid": "TVoid", "StaticVoid": "TVoid",

		"CallBool": "TBool", "CallStaticBool": "TBool", "InvokeBool": "TBool",
		"CallInt": "TInt", "CallStaticInt": "TInt", "InvokeInt": "TInt",
		"StaticIntOf": "TInt", "ConstantInt": "TInt",
		"CallLong": "TLong", "CallStaticLong": "TLong", "InvokeLong": "TLong",
		"CallFloat": "TFloat", "CallStaticFloat": "TFloat", "InvokeFloat": "TFloat",
		"CallDouble": "TDouble", "CallStaticDouble": "TDouble",

		// A constructor's signature ends in V, whatever it constructs: the
		// object comes back from NewObject and not from the descriptor.
		"Make": "TVoid",
	}
	// These return an object of some kind, and anything that is not a
	// primitive will do.
	object := map[string]bool{
		"Call": true, "Invoke": true, "Static": true,
		"CallObject": true, "CallStaticObject": true,
	}
	// The primitives, by the names jni gives them. Anything else is an
	// object of some kind — a class, an array, or java/lang/Object.
	primitives := map[string]bool{
		"jni.TVoid": true, "jni.TBool": true, "jni.TByte": true,
		"jni.TChar": true, "jni.TShort": true, "jni.TInt": true,
		"jni.TLong": true, "jni.TFloat": true, "jni.TDouble": true,
	}
	// And these read a String out, whatever the declared class is.
	str := map[string]bool{
		"InvokeString": true, "CallStaticString": true, "StaticString": true,
		"CallString": true,
	}

	fset := token.NewFileSet()
	checked := 0
	err := filepath.WalkDir("..", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}
		slash := filepath.ToSlash(path)
		if !strings.Contains(slash, "/android/") {
			return nil
		}
		// The jni package defines these; it does not call them with
		// signatures of its own.
		if strings.Contains(slash, "/android/jni/") {
			return nil
		}
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Logf("could not parse %s: %v", path, err)
			return nil
		}

		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			name := sel.Sel.Name
			want, primitive := expect[name]
			if !primitive && !object[name] && !str[name] {
				return true
			}
			sig := findSig(call.Args)
			if sig == nil || len(sig.Args) == 0 {
				return true // a signature built somewhere else; nothing to compare
			}
			checked++

			got := render(sig.Args[0])
			switch {
			case primitive:
				if got != "jni."+want {
					t.Errorf("%s: %s with a signature that returns %s",
						at(fset, call.Pos()), name, got)
				}
			case object[name]:
				if primitives[got] {
					t.Errorf("%s: %s with a signature that returns the primitive %s",
						at(fset, call.Pos()), name, got)
				}
			case str[name]:
				if !strings.Contains(got, "java/lang/String") && !strings.Contains(got, "TString") {
					t.Errorf("%s: %s with a signature that returns %s",
						at(fset, call.Pos()), name, got)
				}
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if checked == 0 {
		t.Fatal("no calls were checked at all; the search is broken, not the code")
	}
	t.Logf("%d JNI calls checked against their signatures", checked)
}

// findSig is the jni.Sig(...) among a call's arguments, if there is one
// written out there rather than built elsewhere.
func findSig(args []ast.Expr) *ast.CallExpr {
	for _, a := range args {
		call, ok := a.(*ast.CallExpr)
		if !ok {
			continue
		}
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "Sig" {
			return call
		}
	}
	return nil
}

// render writes an expression back out, enough to compare two type names by.
func render(e ast.Expr) string {
	switch v := e.(type) {
	case *ast.SelectorExpr:
		return render(v.X) + "." + v.Sel.Name
	case *ast.Ident:
		return v.Name
	case *ast.CallExpr:
		out := render(v.Fun) + "("
		for i, a := range v.Args {
			if i > 0 {
				out += ", "
			}
			out += render(a)
		}
		return out + ")"
	case *ast.BasicLit:
		return v.Value
	}
	return "?"
}

func at(fset *token.FileSet, pos token.Pos) string {
	p := fset.Position(pos)
	return filepath.ToSlash(p.Filename) + ":" + itoa(p.Line)
}
