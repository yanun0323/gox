package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func mainBackup() {
	setupLog()

	if *_help {
		flag.Usage()
		return
	}

	if len(*_targetFile) == 0 {
		flag.Usage()
		requireNoError(errors.New("entity/use/repo at least one param provide"))
	}

	if *_debug {
		println()
		println("\t", "replace", "=", *_replace)
		println("\t", "name", "=", *_targetFile)
		println()
	}

	environmentPrint()

	wd, fd := getDir()
	// Create a FileSet to work with
	fset := token.NewFileSet()
	// Parse the file and create an AST
	file, err := parser.ParseFile(fset, fd, nil, parser.ParseComments)
	requireNoError(err, "parse file set")

	requireNoError(ast.Print(fset, file), "print ast")

	for _, d := range file.Decls {
		println(d.Pos())
	}

	return

	var buf strings.Builder
	requireNoError(format.Node(&buf, fset, file), "format node")

	targetFilePath := filepath.Join(wd, *_targetFile)
	requireNoError(os.MkdirAll(filepath.Dir(targetFilePath), 0o777), "mkdir all")

	f, err := os.Create(targetFilePath)
	requireNoError(err, "open target file")
	defer f.Close()

	_, err = f.WriteString(buf.String())
	requireNoError(err, "write string into file")

	return

	return

	loader := SourceStructLoader{
		Dir: fd,
	}
	err = loader.Load()
	requireNoError(err, "initialize current structure loader")

	structure := loader.GetStruct()
	if *_debug {
		println()
		println("dir:", fd)
		println("struct name:", structure.InterfaceName)
		println("struct type:", structure.GetStructType())
		println("struct:", structure.Interface)
		println("package:", currentPackage())
		println()
	}

	if t := structure.GetStructType(); t != "interface" {
		requireNoError(fmt.Errorf("unsupported type %s, this command only works for `interface`", t))
	}

	pkg := currentPackage()

	var generator *Generator
	switch pkg {
	case "usecase":
		generator = NewGenerator(pkg, "use", *_targetFile, *_replace, &structure, _usecasePathFn)
	case "repository":
		generator = NewGenerator(pkg, "repo", *_targetFile, *_replace, &structure, _repositoryPathFn)
	default:
		log.Println("unsupported package name:", pkg)
	}

	if *_debug {
		generator.DebugPrint()
	}

	err = generator.Save(*_targetFile)
	requireNoError(err, "save generated file")
}
