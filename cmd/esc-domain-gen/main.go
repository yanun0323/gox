package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
)

const _commandName = "esc-domain-gen"

var (
	_replace = flag.Bool("replace", false, "replace all structure and method if there's already a same structure")
	_debug   = flag.Bool("v", false, "show debug information")
	_help    = flag.Bool("h", false, "show command help")
	_name    = flag.String("f", "", "file name to generate implementation")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "%s: 根據定義的介面 package, 生成程式碼到對應位置: usecase/repository\n", _commandName)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-h\t\t顯示用法\n")
	fmt.Fprintf(os.Stderr, "\t-replace\t強制取代目標重複的 Method\n")
	fmt.Fprintf(os.Stderr, "\t-f\t\t目標檔案名稱\t\t-f=member_usecase.go\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t範例:\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t//go:generate %s -f=member.go -replace\n", _commandName)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t到 usecase/repository member.go 生成程式碼, 強制取代目標重複的 Method\n")
	fmt.Fprintf(os.Stderr, "\n")
}

func main() {
	setupLog()

	if *_help {
		flag.Usage()
		return
	}

	if len(*_name) == 0 {
		flag.Usage()
		requireNoError(errors.New("entity/use/repo at least one param provide"))
	}

	if *_debug {
		println()
		println("\t", "replace", "=", *_replace)
		println("\t", "name", "=", *_name)
		println()
	}

	dir, err := getDir()
	requireNoError(err, "get dir")

	internalPath, err := findInternalPath(dir)
	requireNoError(err, "find internal path")

	loader := SourceStructLoader{
		Dir: dir,
	}
	err = loader.Load()
	requireNoError(err, "initialize current structure loader")

	structure := loader.GetStruct()
	if *_debug {
		println()
		println("dir:", dir)
		println("internal path:", internalPath)
		println("struct name:", structure.InterfaceName)
		println("struct type:", structure.GetStructType())
		println("struct:", structure.Interface)
		println("package:", currentPackage())
		println()
	}

	if t := structure.GetStructType(); t != "interface" {
		requireNoError(errors.Errorf("unsupported type %s, this command only works for `interface`", t))
	}

	pkg := currentPackage()

	var generator *Generator
	switch pkg {
	case "usecase":
		generator = NewGenerator(pkg, "use", *_name, *_replace, &structure, _usecasePathFn)
	case "repository":
		generator = NewGenerator(pkg, "repo", *_name, *_replace, &structure, _repositoryPathFn)
	default:
		log.Println("unsupported package name:", pkg)
	}

	if *_debug {
		generator.DebugPrint()
	}

	err = generator.Save(internalPath)
	requireNoError(err, "save generated file")
}

func setupLog() {
	log.SetFlags(0)
	log.SetPrefix(_commandName + ": ")
	flag.Usage = Usage
	flag.Parse()
}
