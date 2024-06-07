package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strconv"
)

type SourceStructLoader struct {
	Dir       string
	PrintAst  bool
	structure Structure
}

func (pt *SourceStructLoader) Load() error {
	fset, tf, f := pt.parsingFile(pt.Dir)
	targetStruct, err := pt.getTargetStruct(tf, f)
	if err != nil {
		return fmt.Errorf("get target struct, err: %w", err)
	}

	structName, handledStruct, err := pt.handleStruct(targetStruct)
	if err != nil {
		return fmt.Errorf("handle struct, err: %w", err)
	}

	if pt.PrintAst {
		ast.Print(fset, handledStruct)
	}

	var cache []byte
	buffer := bytes.NewBuffer(cache)
	err = format.Node(buffer, token.NewFileSet(), handledStruct)
	if err != nil {
		fmt.Errorf("format node, err: %+v", err)
	}

	pt.structure = Structure{
		StructName: structName,
		Struct:     buffer.String(),
	}
	return nil
}

func (pt SourceStructLoader) GetStruct() Structure {
	return pt.structure
}

func (SourceStructLoader) parsingFile(dir string) (*token.FileSet, *token.File, *ast.File) {
	tfs := token.NewFileSet()
	f, err := parser.ParseFile(tfs, dir, nil, parser.AllErrors)
	requireNoError(err, "parse file")

	tf := tfs.File(f.Pos())
	return tfs, tf, f
}

func (loader SourceStructLoader) getTargetStruct(tf *token.File, f *ast.File) (*ast.GenDecl, error) {
	line, err := strconv.Atoi(os.Getenv("GOLINE"))
	requireNoError(err, "get file line")

	for ; line < tf.LineCount(); line++ {
		targetPos := tf.LineStart(line)
		for _, d := range f.Decls {
			g, ok := d.(*ast.GenDecl)
			if !ok {
				continue
			}

			if g.TokPos == targetPos {
				return g, nil
			}
		}
	}

	return nil, errors.New("missing target struct")
}

func (SourceStructLoader) isComment(tf *token.File, f *ast.File, line int) bool {
	l := line
	targetPos := tf.LineStart(l)
	for l < tf.LineCount() {
		for _, c := range f.Comments {
			println("comment pos:", c.Pos(), "target pos:", targetPos)
			if c.Pos() == targetPos {
				return true
			}
		}
		l += 1
		targetPos = tf.LineStart(l)
	}
	return false
}

func (pt SourceStructLoader) handleStruct(st *ast.GenDecl) (string, *ast.GenDecl, error) {
	targetStruct, err := pt.parsingStruct(st)
	if err != nil {
		return "", nil, fmt.Errorf("parsing struct, err: %w", err)
	}

	return targetStruct.Name.Name, st, nil
}

func (SourceStructLoader) parsingStruct(st *ast.GenDecl) (*ast.TypeSpec, error) {
	for _, field := range st.Specs {
		fd, ok := field.(*ast.TypeSpec)
		if ok {
			return fd, nil
		}
	}
	return nil, errors.New("type spec not found")
}
