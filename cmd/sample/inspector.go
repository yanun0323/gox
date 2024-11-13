package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func main() {
	fmt.Printf("Running %s go on %s\n", os.Args[0], os.Getenv("GOFILE"))

	cwd, fd := getDir()
	fmt.Printf("  cwd = %s\n", cwd)
	fmt.Printf("  os.Args = %#v\n", os.Args)

	fset, _, f := getAstFile(fd)
	printAst(fset, f)

	for _, ev := range []string{"GOARCH", "GOOS", "GOFILE", "GOLINE", "GOPACKAGE", "DOLLAR"} {
		println("\t", ev, "=", os.Getenv(ev))
	}
}

func getDir() (currentDirectory string, currentFile string) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fileDir := filepath.Join(cwd, os.Getenv("GOFILE"))
	return cwd, fileDir
}

func getAstFile(filePath string) (*token.FileSet, *token.File, *ast.File) {
	tfs := token.NewFileSet()
	f, err := parser.ParseFile(tfs, filePath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	tf := tfs.File(f.Pos())
	return tfs, tf, f
}

func printAst(fset *token.FileSet, f *ast.File) {
	_ = ast.Print(fset, f)
}
