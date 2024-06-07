package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const _commandName = "domaingen"

var (
	_replace    = flag.Bool("replace", false, "replace all structure and method if there's already a same structure")
	_debug      = flag.Bool("v", false, "show debug information")
	_help       = flag.Bool("h", false, "show command help")
	_targetFile = flag.String("target", "", "target file name to generate implementation")
	_structName = flag.String("struct", "", "target implementation structure name")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "%s: 根據定義的介面 package, 生成程式碼到對應位置\n", _commandName)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-h\t\t顯示用法\n")
	fmt.Fprintf(os.Stderr, "\t-target\t\t目標檔案名稱\t\t\t-target=../../usecase/member_usecase.go\n")
	fmt.Fprintf(os.Stderr, "\t-replace\t強制取代目標相同名稱的 Method\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t範例:\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t//go:generate %s -target=../../usecase/member.go -replace\n", _commandName)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t\t-f=member.go\t生成程式碼在 usecase/repository member.go 檔案內\n")
	fmt.Fprintf(os.Stderr, "\t\t-replace\t強制取代目標相同名稱的 Method\n")
	fmt.Fprintf(os.Stderr, "\n")
}

func main() {
	setupLog()

	if *_help {
		flag.Usage()
		return
	}

	checkTargetFile()
	debugPrint()
	environmentPrint()

	_, fd := getDir()

	helper, err := NewAstHelper(fd)
	requireNoError(err, "new ast helper")

	line, _ := strconv.Atoi(os.Getenv("GOLINE"))
	pos := helper.tf.LineStart(line)

	var decl *ast.GenDecl
	for i, d := range helper.f.Decls[1:] {
		println(i)
		gd, ok := d.(*ast.GenDecl)
		if ok && pos >= d.Pos() || pos <= d.End() {
			decl = gd
			break
		}
	}

	if decl == nil {
		requireNoError(errors.New("interface not found"))
	}

	for _, sp := range decl.Specs {
		s, ok := sp.(*ast.TypeSpec)
		if !ok {
			continue
		}

		fmt.Printf("%+v\n", *s)
	}

	helper.PrintAst()

	// if err := SaveAst(fset, f, filepath.Join(wd, *_targetFile)); err != nil {
	// 	requireNoError(err)
	// }
}

func setupLog() {
	log.SetFlags(0)
	log.SetPrefix(_commandName + ": ")
	flag.Usage = Usage
	flag.Parse()
}

func checkTargetFile() {
	if len(*_targetFile) == 0 {
		flag.Usage()
		requireNoError(errors.New("entity/use/repo at least one param provide"))
	}
}

func debugPrint() {
	if *_debug {
		println()
		println("\t", "replace", "=", *_replace)
		println("\t", "name", "=", *_targetFile)
		println()
	}
}

func getDir() (currentDirectory string, currentFile string) {
	cwd, err := os.Getwd()
	requireNoError(err, "Getwd")

	fileDir := filepath.Join(cwd, os.Getenv("GOFILE"))
	return cwd, fileDir
}
