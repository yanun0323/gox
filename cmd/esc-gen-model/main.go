package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
)

const _commandName = "esc-gen-model"

var (
	_replace = flag.Bool("replace", false, "replace all structure and method if there's already a same structure")
	_debug   = flag.Bool("v", false, "show debug information")
	_help    = flag.Bool("help", false, "show command help")

	_p   = flag.String("p", "", "file name to generate payload structure")
	_pk  = flag.Bool("pk", false, "keep struct field tags")
	_pr  = flag.Bool("pr", false, "replace payload structure and method if there's already a specified same structure (no working when provides -replace)")
	_p2e = flag.Bool("p2e", false, "payload to entity: generate payload method 'ToEntity'")
	_p2r = flag.Bool("p2r", false, "payload to repository: generate payload method 'ToRepository'")
	_p2u = flag.Bool("p2u", false, "payload to usecase: generate payload method 'ToUseCase'")
	_pu  = flag.Bool("pu", false, "set fields which's field name has suffix 'Time' from 'string' to 'int64'")
	_pt  = flag.Bool("pt", false, "set fields which's field name has suffix 'Time' from 'int64' to 'string'")

	_e   = flag.String("e", "", "file name to generate entity structure")
	_ek  = flag.Bool("ek", false, "keep struct field tags")
	_er  = flag.Bool("er", false, "replace entity structure and method if there's already a specified same structure (no working when provides -replace)")
	_e2p = flag.Bool("e2p", false, "entity to payload: generate entity method 'ToPayload'")
	_e2r = flag.Bool("e2r", false, "entity to repository: generate entity method 'ToRepository'")
	_e2u = flag.Bool("e2u", false, "entity to usecase: generate entity method 'ToUseCase'")
	_eu  = flag.Bool("eu", false, "set fields which's field name has suffix 'Time' from 'string' to 'int64'")
	_et  = flag.Bool("et", false, "set fields which's field name has suffix 'Time' from 'int64' to 'string'")

	_r   = flag.String("r", "", "file name to generate repository structure")
	_rk  = flag.Bool("rk", false, "keep struct field tags")
	_rr  = flag.Bool("rr", false, "replace repository structure and method if there's already a specified same structure (no working when provides -replace)")
	_r2p = flag.Bool("r2p", false, "repository to payload: generate repository method 'ToPayload'")
	_r2e = flag.Bool("r2e", false, "repository to entity: generate repository method 'ToEntity'")
	_r2u = flag.Bool("r2u", false, "repository to usecase: generate repository method 'ToUseCase'")
	_ru  = flag.Bool("ru", false, "set fields which's field name has suffix 'Time' from 'string' to 'int64'")
	_rt  = flag.Bool("rt", false, "set fields which's field name has suffix 'Time' from 'int64' to 'string'")

	_u   = flag.String("u", "", "file name to generate usecase structure")
	_uk  = flag.Bool("uk", false, "keep struct field tags")
	_ur  = flag.Bool("ur", false, "replace usecase structure and method if there's already a specified same structure (no working when provides -replace)")
	_u2p = flag.Bool("u2p", false, "usecase to payload: generate usecase method 'ToPayload'")
	_u2e = flag.Bool("u2e", false, "usecase to entity: generate usecase method 'ToEntity'")
	_u2r = flag.Bool("u2r", false, "usecase to repository: generate usecase method 'ToRepository'")
	_uu  = flag.Bool("uu", false, "set fields which's field name has suffix 'Time' from 'string' to 'int64'")
	_ut  = flag.Bool("ut", false, "set fields which's field name has suffix 'Time' from 'int64' to 'string'")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of esc-gen-model:\n")
	fmt.Fprintf(os.Stderr, "\t//go:generate esc-gen-model [flags]\n")
	fmt.Fprintf(os.Stderr, "\t//go:generate esc-gen-model -e filename.go -u filename.go -replace\n")
	fmt.Fprintf(os.Stderr, "\t//go:generate esc-gen-model -e filename.go -u filename.go -r filename.go -replaces e,r\n")
	fmt.Fprintf(os.Stderr, "\t//go:generate esc-gen-model -v -ru -r filename.go\n")
	fmt.Fprintf(os.Stderr, "\t//go:generate esc-gen-model -v -rt -r filename.go\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	setupLog()

	if *_help {
		flag.Usage()
		return
	}

	if len(*_e) == 0 && len(*_u) == 0 && len(*_r) == 0 && len(*_p) == 0 {
		flag.Usage()
		requireNoError(errors.New("entity/use/repo at least one param provide"))
	}

	if *_debug {
		println()
		println("\t", "replace", "=", *_replace)
		println("\t", "entity", "=", *_e)
		println("\t", "use", "=", *_u)
		println("\t", "repo", "=", *_r)
		println("\t", "debug", "=", *_debug)
		println("\t", "pkg", "=", currentPackage())
		println()
	}

	if *_replace {
		*_pr = true
		*_er = true
		*_rr = true
		*_ur = true
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
		println("struct name:", structure.StructName)
		println("struct:", structure.Struct)
		println("package:", currentPackage())
		println()
	}

	generator := Generator{
		Payload:    NewElement(_payload, cleanStringQuote(*_p), *_pr, false, *_p2e, *_p2r, *_p2u, *_pu, *_pt, *_pk, _payloadPathFn),
		Entity:     NewElement(_entity, cleanStringQuote(*_e), *_er, *_e2p, false, *_e2r, *_e2u, *_eu, *_et, *_ek, _entityPathFn),
		Repository: NewElement(_repository, cleanStringQuote(*_r), *_rr, *_r2p, *_r2e, false, *_r2u, *_ru, *_rt, *_rk, _repositoryPathFn),
		Usecase:    NewElement(_usecase, cleanStringQuote(*_u), *_ur, *_u2p, *_u2e, *_u2r, false, *_uu, *_ut, *_uk, _usecasePathFn),
	}
	filename := os.Getenv("GOFILE")
	println("filename", filename)
	generator.ProvideSourceStructure(structure, filename)
	generator.Gen()
	generator.DebugPrint()

	err = generator.Save(internalPath)
	requireNoError(err, "save generated file")
}

func setupLog() {
	log.SetFlags(0)
	log.SetPrefix(_commandName + ": ")
	flag.Usage = Usage
	flag.Parse()
}
