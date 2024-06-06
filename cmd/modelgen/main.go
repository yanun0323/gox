package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

const _commandName = "modelgen"

var (
	_replace   = flag.Bool("replace", false, "replace all structure and method if there's already a same structure")
	_debug     = flag.Bool("v", false, "show debug information")
	_help      = flag.Bool("h", false, "show command help")
	_helpAll   = flag.Bool("help", false, "show command help detail")
	_noComment = flag.Bool("nc", false, "do not generate comment")

	_p   = flag.String("p", "", "file name to generate payload structure")
	_pj  = flag.Bool("pj", false, "remain/generate struct json tag")
	_pr  = flag.Bool("pr", false, "replace payload structure and method if there's already a specified same structure (no working when provides -replace)")
	_p2e = flag.Bool("p2e", false, "payload to entity: generate payload method 'ToEntity'")
	_p2r = flag.Bool("p2r", false, "payload to repository: generate payload method 'ToRepository'")
	_p2u = flag.Bool("p2u", false, "payload to usecase: generate payload method 'ToUseCase'")
	_pfe = flag.Bool("pfe", false, "payload from entity: generate func 'NewXXXFromEntity'")
	_pfr = flag.Bool("pfr", false, "payload from repository: generate func 'NewXXXFromRepository'")
	_pfu = flag.Bool("pfu", false, "payload from usecase: generate func 'NewXXXFromUseCase'")
	_pu  = flag.Bool("pu", false, "set fields which's field name has suffix 'Time' from 'string' to 'int64'")
	_pt  = flag.Bool("pt", false, "set fields which's field name has suffix 'Time' from 'int64' to 'string'")

	_e   = flag.String("e", "", "file name to generate entity structure")
	_ej  = flag.Bool("ej", false, "remain/generate struct json tag")
	_er  = flag.Bool("er", false, "replace entity structure and method if there's already a specified same structure (no working when provides -replace)")
	_e2p = flag.Bool("e2p", false, "entity to payload: generate entity method 'ToPayload'")
	_e2r = flag.Bool("e2r", false, "entity to repository: generate entity method 'ToRepository'")
	_e2u = flag.Bool("e2u", false, "entity to usecase: generate entity method 'ToUseCase'")
	_efp = flag.Bool("efp", false, "entity from payload: generate func 'NewXXXFromPayload'")
	_efr = flag.Bool("efr", false, "entity from repository: generate func 'NewXXXFromRepository'")
	_efu = flag.Bool("efu", false, "entity from usecase: generate func 'NewXXXFromUseCase'")
	_eu  = flag.Bool("eu", false, "set fields which's field name has suffix 'Time' from 'string' to 'int64'")
	_et  = flag.Bool("et", false, "set fields which's field name has suffix 'Time' from 'int64' to 'string'")

	_r   = flag.String("r", "", "file name to generate repository structure")
	_rj  = flag.Bool("rj", false, "remain/generate struct json tag")
	_rr  = flag.Bool("rr", false, "replace repository structure and method if there's already a specified same structure (no working when provides -replace)")
	_r2p = flag.Bool("r2p", false, "repository to payload: generate repository method 'ToPayload'")
	_r2e = flag.Bool("r2e", false, "repository to entity: generate repository method 'ToEntity'")
	_r2u = flag.Bool("r2u", false, "repository to usecase: generate repository method 'ToUseCase'")
	_rfp = flag.Bool("rfp", false, "repository from payload: generate func 'NewXXXFromPayload'")
	_rfe = flag.Bool("rfe", false, "repository from entity: generate func 'NewXXXFromEntity'")
	_rfu = flag.Bool("rfu", false, "repository from usecase: generate func 'NewXXXFromUseCase'")
	_ru  = flag.Bool("ru", false, "set fields which's field name has suffix 'Time' from 'string' to 'int64'")
	_rt  = flag.Bool("rt", false, "set fields which's field name has suffix 'Time' from 'int64' to 'string'")

	_u   = flag.String("u", "", "file name to generate usecase structure")
	_uj  = flag.Bool("uj", false, "remain/generate struct json tag")
	_ur  = flag.Bool("ur", false, "replace usecase structure and method if there's already a specified same structure (no working when provides -replace)")
	_u2p = flag.Bool("u2p", false, "usecase to payload: generate usecase method 'ToPayload'")
	_u2e = flag.Bool("u2e", false, "usecase to entity: generate usecase method 'ToEntity'")
	_u2r = flag.Bool("u2r", false, "usecase to repository: generate usecase method 'ToRepository'")
	_ufp = flag.Bool("ufp", false, "usecase from payload: generate func 'NewXXXFromPayload'")
	_ufe = flag.Bool("ufe", false, "usecase from entity: generate func 'NewXXXFromEntity'")
	_ufr = flag.Bool("ufr", false, "usecase from repository: generate func 'NewXXXFromRepository'")
	_uu  = flag.Bool("uu", false, "set fields which's field name has suffix 'Time' from 'string' to 'int64'")
	_ut  = flag.Bool("ut", false, "set fields which's field name has suffix 'Time' from 'int64' to 'string'")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "%s: 根據定義的結構 package, 生成程式碼到對應位置: payload/entity/usecase/repository\n", _commandName)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-h\t顯示用法\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-help\t顯示用法詳細資訊\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-nc\t生成 Struct 及 Method 時, 不生成 Comment\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-replace\t強制取代目標現有的 Struct 及 Method (所有目標, 相當於 -pr -er -ur -rr)\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t{$}\t代表要產的目標類型, 用以下字母替換:\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t\tp = payload \n")
	fmt.Fprintf(os.Stderr, "\t\te = entity\n")
	fmt.Fprintf(os.Stderr, "\t\tu = usecase\n")
	fmt.Fprintf(os.Stderr, "\t\tr = repository\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-{$}\t目標檔名(不需要路徑)\t\t\t\t-p=member.go\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-{$}j\t保留/生成 目標結構的 json tag\t\t\t-pj\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-{$}r\t強制取代目標現有的 Struct 及 Method\t\t-pr\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-{$}u\t將目標 Time 結尾的 string 欄位轉換為 int64\t-pu\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-{$}t\t將目標 Time 結尾的 int64 欄位轉換為 string\t-pt\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-{$}2{$}\t指定目標生成 ToPayload/ToUseCase/ToEntity/ToRepository Method \t\t\t-p2u (payload 生成 ToUseCase Method)\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-{$}f{$}\t指定目標生成 FromPayload/FromUseCase/FromEntity/FromRepository Function \t-pfu (payload package 生成 NewFromUseCase Function)\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t範例:\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t//go:generate %s -p=member.go -p2u -pu -pj\n", _commandName)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t\t-p=member.go\t生成程式碼在 payload member.go 檔案內\n")
	fmt.Fprintf(os.Stderr, "\t\t-p2u\t\t生成 ToUseCase Method\n")
	fmt.Fprintf(os.Stderr, "\t\t-pu\t\t將 time 結尾的 string 欄位轉換為 int64\n")
	fmt.Fprintf(os.Stderr, "\t\t-pj\t\t生成 json tag\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t//go:generate %s -u=member.go -ufr\n", _commandName)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t\t-u=member.go\t生成程式碼在 usecase member.go 檔案內\n")
	fmt.Fprintf(os.Stderr, "\t\t-ufr\t\t生成 NewFromUsecase Function\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t//go:generate %s -p2e -p2u -e=member.go -e2u -eu -u=member_use.go -uu\n", _commandName)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t\t-p2e\t\t在 payload 生成 ToEntity Method\n")
	fmt.Fprintf(os.Stderr, "\t\t-p2u\t\t在 payload 生成 ToUseCase Method\n")
	fmt.Fprintf(os.Stderr, "\t\t-e=member.go\t生成程式碼在 entity member.go 檔案內\n")
	fmt.Fprintf(os.Stderr, "\t\t-e2u\t\t在 entity 生成 ToUseCase Method\n")
	fmt.Fprintf(os.Stderr, "\t\t-eu\t\t將 time 結尾的 string 欄位轉換為 int64\n")
	fmt.Fprintf(os.Stderr, "\t\t-u=use.go\t生成程式碼在 usecase  use.go 檔案內\n")
	fmt.Fprintf(os.Stderr, "\t\t-uu\t\t將 time 結尾的 string 欄位轉換為 int64\n")
	fmt.Fprintf(os.Stderr, "\n")
}

func main() {
	setupLog()
	return
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
		println("struct type:", structure.GetStructType())
		println("struct:", structure.Struct)
		println("package:", currentPackage())
		println()
	}

	if t := structure.GetStructType(); t != "struct" {
		requireNoError(fmt.Errorf("unsupported type %s, this command only works for `struct`", t))
	}

	generator := Generator{
		Payload: NewElement(ElementParam{
			pkg:            _payload,
			filename:       cleanStringQuote(*_p),
			replace:        *_pr,
			toPayload:      false,
			toEntity:       *_p2e,
			toRepository:   *_p2r,
			toUseCase:      *_p2u,
			fromPayload:    false,
			fromEntity:     *_pfe,
			fromRepository: *_pfr,
			fromUseCase:    *_pfu,
			unix:           *_pu,
			timestamp:      *_pt,
			genJson:        *_pj,
			pathFunc:       _payloadPathFn,
		}),
		Entity: NewElement(ElementParam{
			pkg:            _entity,
			filename:       cleanStringQuote(*_e),
			replace:        *_er,
			toPayload:      *_e2p,
			toEntity:       false,
			toRepository:   *_e2r,
			toUseCase:      *_e2u,
			fromPayload:    *_efp,
			fromEntity:     false,
			fromRepository: *_efr,
			fromUseCase:    *_efu,
			unix:           *_eu,
			timestamp:      *_et,
			genJson:        *_ej,
			pathFunc:       _entityPathFn,
		}),
		Repository: NewElement(ElementParam{
			pkg:            _repository,
			filename:       cleanStringQuote(*_r),
			replace:        *_rr,
			toPayload:      *_r2p,
			toEntity:       *_r2e,
			toRepository:   false,
			toUseCase:      *_r2u,
			fromPayload:    *_rfp,
			fromEntity:     *_rfe,
			fromRepository: false,
			fromUseCase:    *_rfu,
			unix:           *_ru,
			timestamp:      *_rt,
			genJson:        *_rj,
			pathFunc:       _repositoryPathFn,
		}),
		Usecase: NewElement(ElementParam{
			pkg:            _usecase,
			filename:       cleanStringQuote(*_u),
			replace:        *_ur,
			toPayload:      *_u2p,
			toEntity:       *_u2e,
			toRepository:   *_u2r,
			toUseCase:      false,
			fromPayload:    *_ufp,
			fromEntity:     *_ufe,
			fromRepository: *_ufr,
			fromUseCase:    false,
			unix:           *_uu,
			timestamp:      *_ut,
			genJson:        *_uj,
			pathFunc:       _usecasePathFn,
		}),
	}

	filename := os.Getenv("GOFILE")

	generator.ProvideSourceStructure(structure, filename)
	generator.Gen()
	if *_debug {
		generator.DebugPrint()
	}

	err = generator.Save(internalPath, *_noComment)
	requireNoError(err, "save generated file")
}

func setupLog() {
	log.SetFlags(0)
	log.SetPrefix(_commandName + ": ")
	flag.Usage = Usage
	flag.Parse()
}
