package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
	// typeNames   = flag.String("type", "", "comma-separated list of type names; must be set")
	// output      = flag.String("output", "", "output file name; default srcdir/<type>_string.go")
	// trimprefix  = flag.String("trimprefix", "", "trim the `prefix` from the generated constant names")
	// linecomment = flag.Bool("linecomment", false, "use line comment text as printed text when present")
	// buildTags   = flag.String("tags", "", "comma-separated list of build tags to apply")
	_replace = flag.Bool("replace", false, "")
	_rename  = flag.String("rename", "", "")
	_entity  = flag.String("entity", "", "")
	_use     = flag.String("use", "", "")
	_repo    = flag.String("repo", "", "")
	_unix    = flag.Bool("unix", false, "transfers fields which's tag has suffix '_time' from string to int64 if true")
	_debug   = flag.Bool("v", false, "")
	_help    = flag.Bool("help", false, "")
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
	if *_help {
		flag.Usage()
		return
	}

	if len(*_entity) == 0 && len(*_use) == 0 && len(*_repo) == 0 {
		flag.Usage()
		requireNoError(errors.New("entity/use/repo at least one param provide"))
	}

	if *_debug {
		println()
		println("replace:", *_replace)
		println("rename:", *_rename)
		println("entity:", *_entity)
		println("use:", *_use)
		println("repo:", *_repo)
		println("unix:", *_unix)
		println("debug:", *_debug)
		println()
	}

	dir, err := getDir()
	requireNoError(err, "get dir")

	pt := PayloadTransformer{
		Dir: dir,
	}

	err = pt.Init()
	requireNoError(err, "initialize payload transformer")

	entity := pt.GetEntity()
	println("struct name:", entity.StructName)
	println("struct:", entity.Struct)
	println("method name:", entity.MethodName)
	println("method:", entity.Method)

	fp := FileParser{
		Dir:              dir,
		Entity:           strings.Trim(strings.Trim(*_entity, "\""), "'"),
		EntityStruct:     pt.GetEntity(),
		UseCase:          strings.Trim(strings.Trim(*_use, "\""), "'"),
		UseCaseStruct:    pt.GetUseCase(),
		Repository:       strings.Trim(strings.Trim(*_repo, "\""), "'"),
		RepositoryStruct: pt.GetRepository(),
	}

	err = fp.ParseFile()
	requireNoError(err, "parse file with file parser")

	fp.InsertStruct()
	err = fp.SaveFile()
	requireNoError(err, "save file with file parser")
}

func getDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", errors.Errorf("get word dir, err: %+v", err)
	}

	dir := filepath.Join(cwd, os.Getenv("GOFILE"))
	return dir, nil
}

func requireNoError(err error, msg ...string) {
	if err != nil {
		if len(msg) == 0 || len(msg[0]) == 0 {
			log.Fatal(err)
		}
		log.Fatalf("%s, err: %+v", msg[0], err)
	}
}

func requireNotNil(a any, msg ...string) {
	if a == nil {
		if len(msg) == 0 || len(msg[0]) == 0 {
			log.Fatalf("nil struct (%T)", a)
		}
		log.Fatalf("nil struct (%T), err: %s", a, msg[0])
	}
}

type PayloadTransformer struct {
	Dir       string
	PrintAst  bool
	structure Structure
	isReq     bool
}

func (pt *PayloadTransformer) Init() error {

	fset, tf, f := pt.parsingFile(pt.Dir)
	targetStruct, err := pt.getTargetStruct(tf, f)
	if err != nil {
		return errors.Wrap(err, "get target struct")
	}

	structName, isReq, handledStruct, err := pt.handleStruct(targetStruct)
	if err != nil {
		return errors.Wrap(err, "handle struct")
	}

	if pt.PrintAst {
		ast.Print(fset, handledStruct)
		println("isReq:", isReq)
	}

	var cache []byte
	buffer := bytes.NewBuffer(cache)
	err = format.Node(buffer, token.NewFileSet(), handledStruct)
	if err != nil {
		errors.Errorf("format node, err: %+v", err)
	}

	pt.structure = Structure{
		StructName: structName,
		Struct:     buffer.String(),
	}
	pt.isReq = isReq
	return nil
}

type Structure struct {
	StructName string
	Struct     string
	MethodName string
	Method     string
}

const _methodTemplate = `func (%s *%s) %s() *%s {
	return &%s{
		%s
	}
}`

func (pt PayloadTransformer) GetEntity() Structure {
	// !isReq -> empty method
	if !pt.isReq {
		return pt.structure
	}

	// isReq -> ToUse
	pt.structure.MethodName = "ToUseCase"
	short := pt.getShortName(pt.structure.StructName)
	fields := pt.getFields(pt.structure)
	content := make([]string, 0, len(fields))
	for _, f := range fields {
		content = append(content, fmt.Sprintf("%s: usecase.%s,", f, f))
	}

	pt.structure.Method = fmt.Sprintf(_methodTemplate,
		short,
		pt.structure.StructName,
		"ToUseCase",
		"usecase."+pt.structure.StructName,
		"usecase."+pt.structure.StructName,
		strings.Join(content, "\n"),
	)

	formatted, err := format.Source([]byte(pt.structure.Method))
	if err == nil {
		pt.structure.Method = string(formatted)
	}

	if *_debug {
		println("entity:", pt.structure.Method)
	}

	return pt.structure
}

func (pt PayloadTransformer) GetUseCase() Structure {
	// isReq -> ToRepository
	if pt.isReq {
		pt.structure.MethodName = "ToRepository"
		short := pt.getShortName(pt.structure.StructName)
		fields := pt.getFields(pt.structure)
		content := make([]string, 0, len(fields))
		for _, f := range fields {
			content = append(content, fmt.Sprintf("%s: repository.%s,", f, f))
		}

		pt.structure.Method = fmt.Sprintf(_methodTemplate,
			short,
			pt.structure.StructName,
			"ToRepository",
			"repository."+pt.structure.StructName,
			"repository."+pt.structure.StructName,
			strings.Join(content, "\n"),
		)

		formatted, err := format.Source([]byte(pt.structure.Method))
		if err == nil {
			pt.structure.Method = string(formatted)
		}

		if *_debug {
			println("use case:", pt.structure.Method)
		}
	} else {
		// !isReq -> ToPayload
		pt.structure.MethodName = "ToPayload"
		short := pt.getShortName(pt.structure.StructName)
		fields := pt.getFields(pt.structure)
		content := make([]string, 0, len(fields))
		for _, f := range fields {
			content = append(content, fmt.Sprintf("%s: payload.%s,", f, f))
		}

		pt.structure.Method = fmt.Sprintf(_methodTemplate,
			short,
			pt.structure.StructName,
			"ToPayload",
			"payload."+pt.structure.StructName,
			"payload."+pt.structure.StructName,
			strings.Join(content, "\n"),
		)

		formatted, err := format.Source([]byte(pt.structure.Method))
		if err == nil {
			pt.structure.Method = string(formatted)
		}

		if *_debug {
			println("use case:", pt.structure.Method)
		}
	}

	return pt.structure
}

func (pt PayloadTransformer) GetRepository() Structure {
	// isReq -> empty method
	if pt.isReq {
		return pt.structure
	}
	// !isReq -> ToUseCase
	pt.structure.MethodName = "ToUseCase"
	short := pt.getShortName(pt.structure.StructName)
	fields := pt.getFields(pt.structure)
	content := make([]string, 0, len(fields))
	for _, f := range fields {
		content = append(content, fmt.Sprintf("%s: usecase.%s,", f, f))
	}

	pt.structure.Method = fmt.Sprintf(_methodTemplate,
		short,
		pt.structure.StructName,
		"ToUseCase",
		"usecase."+pt.structure.StructName,
		"usecase."+pt.structure.StructName,
		strings.Join(content, "\n"),
	)

	formatted, err := format.Source([]byte(pt.structure.Method))
	if err == nil {
		pt.structure.Method = string(formatted)
	}

	if *_debug {
		println("entity:", pt.structure.Method)
	}

	return pt.structure
}

func (PayloadTransformer) getShortName(name string) string {
	cache := make([]byte, 0, 10)
	for _, char := range name {
		if char >= 'A' && char <= 'Z' {
			cache = append(cache, byte(char))
		}
	}
	return string(bytes.ToLower(cache))
}

func (PayloadTransformer) getFields(st Structure) []string {
	fields := make([]string, 0, 10)
	rows := strings.Split(st.Struct, "\n")
	rows = rows[1:]
	rows = rows[:len(rows)-1]
	for _, row := range rows {
		f := strings.Split(strings.TrimSpace(row), " ")[0]
		if len(f) != 0 {
			fields = append(fields, f)
		}
	}
	return fields
}

func (PayloadTransformer) parsingFile(dir string) (*token.FileSet, *token.File, *ast.File) {
	tfs := token.NewFileSet()
	f, err := parser.ParseFile(tfs, dir, nil, parser.AllErrors)
	requireNoError(err, "parse file")

	tf := tfs.File(f.Pos())
	return tfs, tf, f
}

func (PayloadTransformer) getTargetStruct(tf *token.File, f *ast.File) (*ast.GenDecl, error) {
	pos, err := strconv.Atoi(os.Getenv("GOLINE"))
	requireNoError(err, "get file line")
	targetPos := tf.LineStart(pos + 1)
	for _, d := range f.Decls {
		g, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}

		if g.TokPos == targetPos {
			return g, nil
		}
	}
	return nil, errors.New("missing target struct")
}

func (pt PayloadTransformer) handleStruct(st *ast.GenDecl) (string, bool, *ast.GenDecl, error) {
	targetStruct, err := pt.parsingStruct(st)
	if err != nil {
		return "", false, nil, errors.Wrap(err, "parsing struct")
	}

	isReq := pt.renameAndIsRequest(targetStruct)

	if *_unix {
		structContainer, ok := targetStruct.Type.(*ast.StructType)
		if !ok {
			return "", false, nil, errors.New("transfer to struct type")
		}

		requireNotNil(structContainer.Fields, "struct container fields")
		requireNotNil(structContainer.Fields.List, "struct container fields list")

		println("list count:", len(structContainer.Fields.List))
		for i, field := range structContainer.Fields.List {
			requireNotNil(field, "list field")
			requireNotNil(field.Tag, "list field tag")
			if field.Tag == nil {
				continue
			}

			if !strings.Contains(field.Tag.Value, "_time\"") && !strings.Contains(field.Tag.Value, "_at\"") {
				continue
			}

			fieldType, ok := field.Type.(*ast.Ident)
			if !ok {
				return "", false, nil, errors.Errorf("type error at field (%d)", i)
			}

			if strings.EqualFold(fieldType.Name, "string") {
				fieldType.Name = "int64"
			}
		}
	}

	return targetStruct.Name.Name, isReq, st, nil
}

func (PayloadTransformer) parsingStruct(st *ast.GenDecl) (*ast.TypeSpec, error) {
	for _, field := range st.Specs {
		fd, ok := field.(*ast.TypeSpec)
		if ok {
			return fd, nil
		}
	}
	return nil, errors.New("type spec not found")
}

func (pt PayloadTransformer) renameAndIsRequest(target *ast.TypeSpec) bool {
	requireNotNil(target, "target struct")
	requireNotNil(target.Name, "target struct name")

	defer func() {
		if len(*_rename) != 0 {
			target.Name.Name = *_rename
		} else {
		}
	}()

	if pt.cutSuffix(target.Name, "Req") {
		return true
	}

	if pt.cutSuffix(target.Name, "Request") {
		return true
	}

	if pt.cutSuffix(target.Name, "Res") {
		return true
	}

	if pt.cutSuffix(target.Name, "Resp") {
		return false
	}

	if pt.cutSuffix(target.Name, "Response") {
		return false
	}

	return false
}

func (PayloadTransformer) cutSuffix(name *ast.Ident, suffix string) bool {
	if strings.HasSuffix(name.Name, suffix) {
		before, _ := strings.CutSuffix(name.Name, suffix)
		name.Name = before
		return true
	}
	return false
}
