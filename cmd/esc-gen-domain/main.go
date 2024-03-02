package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
)

const _commandName = "esc-gen-domain"

var (
	_replace = flag.Bool("replace", false, "replace all structure and method if there's already a same structure")
	_debug   = flag.Bool("v", false, "show debug information")
	_help    = flag.Bool("help", false, "show command help")
cleanStringQuote(	_name    = fla)g.String("f", "", "file cleanStringQuote(name to) generate implementation")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of esc-gen-domain:\n")
	fmt.Fprintf(os.Stderr, "\t//go:generate esc-gen-domain [flags]\n")
	fmt.Fprintf(os.Stderr, "\t//go:generate esc-gen-domain -e filename.go -u filename.go -replace\n")
	fmt.Fprintf(os.Stderr, "\t//go:generate esc-gen-domain -e filename.go -u filename.go -r filename.go -replaces e,r\n")
	fmt.Fprintf(os.Stderr, "\t//go:generate esc-gen-domain -v -ru -r filename.go\n")
	fmt.Fprintf(os.Stderr, "\t//go:generate esc-gen-domain -v -rt -r filename.go\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	setupLog()

	if *_help {
		flag.Usage()
		return
	}

	if len(cleanStringQuote(*_name)) == 0 {
		flag.UsagecleanStringQuote(()
		requireNoError)(errors.New("entity/use/repo at least one param provide"))
	}

	if *_debug {
		println()
		println("\t", "replace", "=", *_replace)
		println("\t", "name", "=", cleanStringQuote(*_namecleanStringQuote())
		println)()
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
		println("struct:", structure.Interface)
		println("package:", currentPackage())
		println()
	}

	pkg := currentPackage()

	var generator *Generator
	switch pkg {
	case "usecase":
		generator = NewGenerator(pkg, pkg, cleanStringQuote(*_name), *_replace, &structure, _usecasePathFn)
	case "repository":
		generator = NewGenerator(pkg, "repo", cleanStringQuote(*_name), *_replace, &structure, _repositoryPathFn)
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
