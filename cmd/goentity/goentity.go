package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var (
	// typeNames   = flag.String("type", "", "comma-separated list of type names; must be set")
	// output      = flag.String("output", "", "output file name; default srcdir/<type>_string.go")
	// trimprefix  = flag.String("trimprefix", "", "trim the `prefix` from the generated constant names")
	// linecomment = flag.Bool("linecomment", false, "use line comment text as printed text when present")
	// buildTags   = flag.String("tags", "", "comma-separated list of build tags to apply")
	dir  = flag.String("dir", "", "output dir destination")
	unix = flag.Bool("unix", false, "transfers fields which's tag has suffix '_time' from string to int64 if true")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of goentity:\n")
	fmt.Fprintf(os.Stderr, "\tgoentity [flags] -type T [directory]\n")
	fmt.Fprintf(os.Stderr, "\tgoentity [flags] -type T files... # Must be a single package\n")
	fmt.Fprintf(os.Stderr, "For more information, see:\n")
	fmt.Fprintf(os.Stderr, "\thttps://pkg.go.dev/golang.org/x/tools/cmd/goentity\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("goentity: ")
	flag.Usage = Usage
	flag.Parse()
	if len(*dir) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	cwd, err := os.Getwd()
	requireNoError(err, "get work dir")

	path := filepath.Join(cwd, os.Getenv("GOFILE"))
	log.Default().Printf("current file (%s)", path)

	fset := token.NewFileSet()
	_, err = parser.ParseFile(fset, path, nil, parser.AllErrors)
	requireNoError(err, "parse file")

	// ast.Print(fset, f)

	pos, err := strconv.Atoi(os.Getenv("GOLINE"))
	requireNoError(err, "get file line")

	obj := fset.File(token.Pos(pos))
	println("fset.base: ", fset.Base())

	println("obj.name:", obj.Name())
	println("obj.line count:", obj.LineCount())
	println(obj.Size())
	println(obj.Base())
	for _, l := range obj.Lines() {
		print("line:", l)
	}

	ast.Inspect(fset, func(n ast.Node) bool {

	})
}

func requireNoError(err error, msg ...string) {
	if err != nil {
		if len(msg) == 0 || len(msg[0]) == 0 {
			log.Fatal(err)
		}
		log.Fatalf("%s, err: %+v", msg[0], err)
	}
}

func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}
