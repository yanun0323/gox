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
	_help    = flag.Bool("h", false, "show command help")
	_helpAll = flag.Bool("help", false, "show command help detail")

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
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "%s:\n", _commandName)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-h\t顯示用法\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-help\t顯示用法詳細資訊\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t{$}\t代表要產的目標類型, 用以下字母替換:\n")
	fmt.Fprintf(os.Stderr, "\t\tp=payload e=entity u=usecase r=repository\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-{$}\t目標檔名(不需要路徑)\t\t\t\t-p=member.go\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-{$}k\t保留目標 Struct Tag\t\t\t\t-pk\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-{$}r\t強制取代目標現有的 Struct 及 Method\t\t-pr\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-{$}u\t將目標 Time 結尾的 string 欄位轉換為 int64\t-pu\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-{$}t\t將目標 Time 結尾的 int64 欄位轉換為 string\t-pt\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-{$}2{$}\t指定目標生成 ToUseCase/ToEntity/ToRepository Method \t-p2u (payload 生成 ToUseCase 的方法)\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t範例:\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t//go:generate %s -p=member.go -p2u -pu\n", _commandName)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t生成 payload, 生成在 member.go 內, 生成 ToUseCase 的方法, 將 time 結尾的 string 欄位轉換為 int64\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t//go:generate %s -p2e -p2u -e=member.go -e2u -eu -u=member_use.go -uu\n", _commandName)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t在 payload 生成 ToEntity & ToUseCase 的方法\n")
	fmt.Fprintf(os.Stderr, "\t生成 entity, 生成在 member.go 內, 生成 ToUseCase 的方法, 將 time 結尾的 string 欄位轉換為 int64\n")
	fmt.Fprintf(os.Stderr, "\t生成 usecase, 生成在 member_use.go 內, 將 time 結尾的 string 欄位轉換為 int64\n")
	fmt.Fprintf(os.Stderr, "\n")
}

func main() {
	setupLog()

	if *_helpAll {
		flag.PrintDefaults()
		flag.Usage()
		return
	}

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

	generator.ProvideSourceStructure(structure, filename)
	generator.Gen()
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
