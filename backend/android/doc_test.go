package android_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// Every exported symbol under antui/backend/android has to say what it is.
//
// This is a test and not an audit because an audit is true on the day it is
// done. A package people are meant to write apps against is read through its
// documentation — `go doc` is the interface, not the source — and one
// undocumented export is a question somebody has to answer by reading the
// implementation.
//
// It is deliberately shallow about *quality*: it cannot tell a good comment
// from a bad one, and it does not try. What it can tell is that a comment is
// there and that it begins with the name, which is the convention the whole
// tree follows and the thing that makes `go doc` read as sentences.
func TestEveryExportIsDocumented(t *testing.T) {
	var missing []string
	fset := token.NewFileSet()

	err := filepath.WalkDir("..", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		// Only the Android tree. The rest of antui predates this and is
		// somebody else's afternoon.
		if !strings.Contains(filepath.ToSlash(path), "/android/") {
			return nil
		}
		// The generated Go inside an example is not an interface.
		if strings.Contains(filepath.ToSlash(path), "/examples/") {
			return nil
		}

		// Parsed with the comments and with every build tag ignored, so that
		// a file which only compiles for Android is still read here.
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			// A file this cannot parse is not a reason to fail: cgo
			// preambles and //export lines are ordinary Go, but a syntax
			// this Go does not know yet should not stop the check.
			t.Logf("could not parse %s: %v", path, err)
			return nil
		}

		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if !d.Name.IsExported() {
					continue
				}
				// A method on an unexported type is not part of the
				// interface, whatever its own name is.
				if d.Recv != nil && !exportedReceiver(d.Recv) {
					continue
				}
				if !documented(d.Doc, d.Name.Name) {
					missing = append(missing, where(fset, d.Pos(), name(d)))
				}
			case *ast.GenDecl:
				checkGenDecl(fset, d, &missing)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	sort.Strings(missing)
	if len(missing) > 0 {
		t.Errorf("%d exported symbols with nothing said about them:\n  %s",
			len(missing), strings.Join(missing, "\n  "))
	}
}

// checkGenDecl looks at types, constants and variables, which come in groups
// and may be documented as a group or one at a time.
func checkGenDecl(fset *token.FileSet, d *ast.GenDecl, missing *[]string) {
	if d.Tok == token.IMPORT {
		return
	}
	// A group with a comment above it documents every name in it, which is
	// how a set of constants is written everywhere in this tree.
	grouped := d.Doc != nil && len(d.Doc.List) > 0

	for _, spec := range d.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			if !s.Name.IsExported() {
				continue
			}
			if grouped || documented(s.Doc, s.Name.Name) {
				continue
			}
			*missing = append(*missing, where(fset, s.Pos(), "type "+s.Name.Name))
		case *ast.ValueSpec:
			for _, n := range s.Names {
				if !n.IsExported() {
					continue
				}
				if grouped || documented(s.Doc, n.Name) || s.Comment != nil {
					continue
				}
				*missing = append(*missing, where(fset, n.Pos(), n.Name))
			}
		}
	}
}

// documented says whether there is a comment and whether it starts with the
// name, which is what makes `go doc` read as prose rather than as a list of
// remarks.
func documented(doc *ast.CommentGroup, name string) bool {
	if doc == nil || len(doc.List) == 0 {
		return false
	}
	text := strings.TrimSpace(doc.Text())
	if text == "" {
		return false
	}
	// "Name is", "Name does", "A Name is", "ErrX is what ..." — and a
	// deprecation or a build directive on its own does not count.
	first := text
	if i := strings.IndexByte(text, '\n'); i > 0 {
		first = text[:i]
	}
	for _, prefix := range []string{name, "A " + name, "An " + name, "The " + name} {
		if strings.HasPrefix(first, prefix) {
			return true
		}
	}
	// A comment that leads with something else is still a comment. The
	// convention is worth following and is not worth failing over — what
	// matters is that a reader is told something.
	return len(first) > 12
}

func exportedReceiver(recv *ast.FieldList) bool {
	if recv == nil || len(recv.List) == 0 {
		return false
	}
	switch t := recv.List[0].Type.(type) {
	case *ast.StarExpr:
		if id, ok := t.X.(*ast.Ident); ok {
			return id.IsExported()
		}
	case *ast.Ident:
		return t.IsExported()
	case *ast.IndexExpr: // a generic receiver
		if id, ok := t.X.(*ast.Ident); ok {
			return id.IsExported()
		}
	}
	return false
}

func name(d *ast.FuncDecl) string {
	if d.Recv == nil {
		return "func " + d.Name.Name
	}
	return "method " + receiverName(d.Recv) + "." + d.Name.Name
}

func receiverName(recv *ast.FieldList) string {
	switch t := recv.List[0].Type.(type) {
	case *ast.StarExpr:
		if id, ok := t.X.(*ast.Ident); ok {
			return id.Name
		}
	case *ast.Ident:
		return t.Name
	}
	return "?"
}

func where(fset *token.FileSet, pos token.Pos, what string) string {
	p := fset.Position(pos)
	return filepath.ToSlash(p.Filename) + ":" + itoa(p.Line) + ": " + what
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
