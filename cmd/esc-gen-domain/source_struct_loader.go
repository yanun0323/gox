package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

type SourceStructLoader struct {
	Dir      string
	PrintAst bool
	inter    Interface
}

func (pt *SourceStructLoader) Load() error {
	fset, tf, f := pt.parsingFile(pt.Dir)
	targetInterface, err := pt.getTargetInterface(tf, f)
	if err != nil {
		return errors.Wrap(err, "get target interface")
	}

	interfaceName, handledInterface, err := pt.handleInterface(targetInterface)
	if err != nil {
		return errors.Wrap(err, "handle interface")
	}

	if pt.PrintAst {
		ast.Print(fset, handledInterface)
	}

	var cache []byte
	buffer := bytes.NewBuffer(cache)
	err = format.Node(buffer, token.NewFileSet(), handledInterface)
	if err != nil {
		errors.Errorf("format node, err: %+v", err)
	}

	pt.inter = Interface{
		InterfaceName: interfaceName,
		Interface:     buffer.String(),
	}
	return nil
}

func (pt SourceStructLoader) GetStruct() Interface {
	return pt.inter
}

func (SourceStructLoader) parsingFile(dir string) (*token.FileSet, *token.File, *ast.File) {
	tfs := token.NewFileSet()
	f, err := parser.ParseFile(tfs, dir, nil, parser.AllErrors)
	requireNoError(err, "parse file")

	tf := tfs.File(f.Pos())
	return tfs, tf, f
}

func (SourceStructLoader) getTargetInterface(tf *token.File, f *ast.File) (*ast.GenDecl, error) {
	pos, err := strconv.Atoi(os.Getenv("GOLINE"))
	requireNoError(err, "get file line")

	targetPos := tf.LineStart(pos + 1)
	for _, d := range f.Decls {
		g, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}

		if g.TokPos == targetPos {
			return g, nil
		}
	}
	return nil, errors.New("missing target interface")
}

func (pt SourceStructLoader) handleInterface(st *ast.GenDecl) (string, *ast.GenDecl, error) {
	targetInterface, err := pt.parsingInterface(st)
	if err != nil {
		return "", nil, errors.Wrap(err, "parsing interface")
	}

	return targetInterface.Name.Name, st, nil
}

func (SourceStructLoader) parsingInterface(st *ast.GenDecl) (*ast.TypeSpec, error) {
	for _, field := range st.Specs {
		fd, ok := field.(*ast.TypeSpec)
		if ok {
			return fd, nil
		}
	}
	return nil, errors.New("type spec not found")
}
