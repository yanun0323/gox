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

const _commandName = "domaingen"

var (
	_replace    = flag.Bool("replace", false, "replace all structure and method if there's already a same structure")
	_debug      = flag.Bool("v", false, "show debug information")
	_help       = flag.Bool("h", false, "show command help")
	_targetFile = flag.String("target", "", "target file name to generate implementation")
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

	dir, err := getDir()
	requireNoError(err, "get dir")

	cwd, err := os.Getwd()
	requireNoError(err, "get word dir")

	// Create a FileSet to work with
	fset := token.NewFileSet()
	// Parse the file and create an AST
	file, err := parser.ParseFile(fset, dir, nil, parser.ParseComments)
	requireNoError(err, "parse file set")

	requireNoError(ast.Print(fset, file), "print ast")

	var buf strings.Builder
	requireNoError(format.Node(&buf, fset, file), "format node")

	path := filepath.Join(cwd, *_targetFile)
	println(cwd)
	println(*_targetFile)
	println(path)

	requireNoError(os.MkdirAll(filepath.Dir(path), 0o777), "mkdir all")
	f, err := os.Create(path)
	requireNoError(err, "open target file")
	defer f.Close()

	_, err = f.WriteString(buf.String())
	requireNoError(err, "write string into file")

	println(buf.String())

	return

	ast.Inspect(file, func(n ast.Node) bool {
		// Find Function Call Statements
		switch n := n.(type) {
		case *ast.Comment:
			println("Comment", n.Pos(), "~", n.End())
		case *ast.CommentGroup:
			println("CommentGroup", n.Pos(), "~", n.End())

			// Expressions
		case *ast.Field:
			println("Field", n.Pos(), "~", n.End())
		case *ast.FieldList:
			println("FieldList", n.Pos(), "~", n.End())
		case *ast.BadExpr:
			println("BAD:", n.Pos(), "~", n.End())
		case *ast.Ident:
			println("Ident:", n.Name, n.Pos(), "~", n.End())
		case *ast.Ellipsis:
			println("Ellipsis:", n.Elt, n.Pos(), "~", n.End())
		case *ast.BasicLit:
			println("BasicLit:", n.Value, n.Kind.String(), n.Pos(), "~", n.End())
		case *ast.FuncLit:
			println("Func:", n.Type.Func, n.Pos(), "~", n.End())
		case *ast.CompositeLit:
			println("Composite:", n.Type, n.Pos(), "~", n.End())

			// Expression
		case *ast.ParenExpr:
			println("Paren:", n.X, n.Pos(), "~", n.End())
		case *ast.SelectorExpr:
			println("SELECTOR", n.X, n.Pos(), "~", n.End())
		case *ast.IndexExpr:
			println("Index:", n.X, n.Pos(), "~", n.End())
		case *ast.IndexListExpr:
			println("IndexList:", n.X, n.Pos(), "~", n.End())
		case *ast.SliceExpr:
			println("Slice:", n.X, n.Pos(), "~", n.End())
		case *ast.TypeAssertExpr:
			println("TypeAssert:", n.X, n.Pos(), "~", n.End())
		case *ast.CallExpr:
			println("Call:", n.Fun, n.Pos(), "~", n.End())
		case *ast.StarExpr:
			println("Star:", n.X, n.Pos(), "~", n.End())
		case *ast.UnaryExpr:
			println("Unary:", n.X, n.Pos(), "~", n.End())
		case *ast.BinaryExpr:
			println("Binary:", n.X, n.Pos(), "~", n.End())
		case *ast.KeyValueExpr:
			println("KeyValue:", n.Key, n.Value, n.Pos(), "~", n.End())

			// TYPE
		case *ast.ArrayType:
			println("Array:", n.Elt, n.Len, n.Pos(), "~", n.End())
		case *ast.StructType:
			println("Struct:", n.Fields, n.Pos(), "~", n.End())
		case *ast.FuncType:
			println("FuncType", n.Pos(), "~", n.End())
		case *ast.InterfaceType:
			println("InterfaceType:", n.Methods, n.Pos(), "~", n.End())
		case *ast.MapType:
			println("MapType", n.Key, n.Value, n.Pos(), "~", n.End())
		case *ast.ChanType:
			println("ChanType", n.Value, n.Pos(), "~", n.End())

			// Statements
		case *ast.BadStmt:
			println("BadStmt", n.Pos(), "~", n.End())
		case *ast.DeclStmt:
			println("DeclStmt", n.Pos(), "~", n.End())
		case *ast.EmptyStmt:
			println("EmptyStmt", n.Pos(), "~", n.End())
		case *ast.LabeledStmt:
			println("LabeledStmt", n.Pos(), "~", n.End())
		case *ast.ExprStmt:
			println("ExprStmt", n.Pos(), "~", n.End())
		case *ast.SendStmt:
			println("SendStmt", n.Pos(), "~", n.End())
		case *ast.IncDecStmt:
			println("IncDecStmt", n.Pos(), "~", n.End())
		case *ast.AssignStmt:
			println("AssignStmt", n.Pos(), "~", n.End())
		case *ast.GoStmt:
			println("GoStmt", n.Pos(), "~", n.End())
		case *ast.DeferStmt:
			println("DeferStmt", n.Pos(), "~", n.End())
		case *ast.ReturnStmt:
			println("ReturnStmt", n.Pos(), "~", n.End())
		case *ast.BranchStmt:
			println("BranchStmt", n.Pos(), "~", n.End())
		case *ast.BlockStmt:
			println("BlockStmt", n.Pos(), "~", n.End())
		case *ast.IfStmt:
			println("IfStmt", n.Pos(), "~", n.End())
		case *ast.CaseClause:
			println("CaseClause", n.Pos(), "~", n.End())
		case *ast.SwitchStmt:
			println("SwitchStmt", n.Pos(), "~", n.End())
		case *ast.TypeSwitchStmt:
			println("TypeSwitchStmt", n.Pos(), "~", n.End())
		case *ast.CommClause:
			println("CommClause", n.Pos(), "~", n.End())
		case *ast.SelectStmt:
			println("SelectStmt", n.Pos(), "~", n.End())
		case *ast.ForStmt:
			println("ForStmt", n.Pos(), "~", n.End())
		case *ast.RangeStmt:
			println("RangeStmt", n.Pos(), "~", n.End())

			// Declarations
		case *ast.ImportSpec:
			println("ImportSpec", n.Name, n.Pos(), "~", n.End())
		case *ast.ValueSpec:
			println("ValueSpec", n.Pos(), "~", n.End())
		case *ast.TypeSpec:
			println("TypeSpec", n.Pos(), "~", n.End())
		case *ast.BadDecl:
			println("BadDecl", n.Pos(), "~", n.End())
		case *ast.GenDecl:
			println("GenDecl", n.Tok.String(), n.Pos(), "~", n.End())
		case *ast.FuncDecl:
			println("FuncDecl", n.Pos(), "~", n.End())
			// File
		case *ast.File:
			println("File", n.Name.Name, n.Pos(), "~", n.End())
		case *ast.Package:
			println("Package", n.Pos(), "~", n.End())

			// Interface
		case ast.Expr:
			println("Expr", n.Pos(), "~", n.End())
		case ast.Stmt:
			println("Stmt", n.Pos(), "~", n.End())
		case ast.Decl:
			println("Decl", n.Pos(), "~", n.End())
		case ast.Node:
			println("Node", n.Pos(), "~", n.End())
		case nil:
			println("<>")
		default:
			println("Default", n)
		}
		return true
	})

	return

	loader := SourceStructLoader{
		Dir: dir,
	}
	err = loader.Load()
	requireNoError(err, "initialize current structure loader")

	structure := loader.GetStruct()
	if *_debug {
		println()
		println("dir:", dir)
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

func setupLog() {
	log.SetFlags(0)
	log.SetPrefix(_commandName + ": ")
	flag.Usage = Usage
	flag.Parse()
}
