package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/yanun0323/goast"
	"github.com/yanun0323/goast/kind"
	"github.com/yanun0323/goast/scope"
)

const _commandName = "modelgen"

var (
	_help        = flag.Bool("h", false, "show command help")
	_debug       = flag.Bool("v", false, "show debug information")
	_replace     = flag.Bool("replace", false, "replace all structure and method if there's already a same structure")
	_relative    = flag.Bool("relative", false, "target implementation structure name")
	_destination = flag.String("destination", "", "target file name to generate implementation")
	_package     = flag.String("package", "", "target implementation structure name")
	_name        = flag.String("name", "", "target implementation structure name")
	_function    = flag.String("function", "", "target implementation structure name")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "%s: 根據定義的介面 package, 生成程式碼到對應位置\n", _commandName)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-h\t\t顯示用法\n")
	fmt.Fprintf(os.Stderr, "\t-name\t\t目標結構的名稱\t\t\t-name=usecase\n")
	fmt.Fprintf(os.Stderr, "\t-destination\t\t目標檔案名稱\t\t\t-destination=../../usecase/member_usecase.go\n")
	fmt.Fprintf(os.Stderr, "\t-replace\t強制取代目標相同名稱的 Struct/Function/Method\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t範例:\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t//go:generate %s -destination=../../usecase/member.go -name=usecase -replace\n", _commandName)
	fmt.Fprintf(os.Stderr, "\n")
}

func main() {
	NoError(run())
}

func run() error {
	helper.setupLog()
	helper.debugPrint()

	if *_help {
		flag.Usage()
		return nil
	}

	helper.requireDestination()

	ast, goLine, pkg, err := parseAstFromGoGenerator()
	if err != nil {
		return err
	}

	targetScope, err := findTargetStruct(ast, goLine)
	if err != nil {
		return err
	}

	structName, err := findStructNameAndSetImplementName(targetScope)
	if err != nil {
		return err
	}

	var (
		importPkg           bool
		relativeScopes      map[string]goast.Scope
		relativeScopesNames []string
	)

	if *_relative {
		relativeScopes, relativeScopesNames = findRelativeScopes(ast, targetScope)
	} else {
		importPkg = addPackageNameInFrontOfParamType(targetScope, pkg)
	}

	println(structName, importPkg, relativeScopes, relativeScopesNames)

	return nil
}
func parseAstFromGoGenerator() (ast goast.Ast, goLine int, pkg string, err error) {
	_, file, err := helper.getDir()
	if err != nil {
		return nil, 0, "", fmt.Errorf("get directory, err: %w", err)
	}

	astObj, err := goast.ParseAst(file)
	if err != nil {
		return nil, 0, "", fmt.Errorf("parse ast, err: %w", err)
	}

	goLineNum, err := strconv.Atoi(os.Getenv("GOLINE"))
	if err != nil {
		return nil, 0, "", fmt.Errorf("parse GOLINE, err: %w", err)
	}

	pkgName := os.Getenv("GOPACKAGE")

	return astObj, goLineNum, pkgName, nil
}

func findTargetStruct(ast goast.Ast, goLine int) (goast.Scope, error) {
	var (
		lineMatched bool
		targetScope goast.Scope
	)

	ast.IterScope(func(s goast.Scope) bool {
		if lineMatched {
			if s.Line() != goLine {
				return false
			}

			goLine++

			switch s.Kind() {
			case scope.Type:
				if _, ok := s.GetStructName(); ok {
					targetScope = s
				}

				return false
			case scope.Comment:
				return true
			default:
				return false
			}
		} else {
			lineMatched = s.Line() == goLine
			if lineMatched {
				goLine++
			}
		}

		return true
	})

	if targetScope == nil {
		return nil, errors.New("target struct not found")
	}

	return targetScope, nil
}

func findStructNameAndSetImplementName(targetScope goast.Scope) (string, error) {
	structName, ok := targetScope.GetTypeName()
	if !ok || len(structName) == 0 {
		return "", errors.New("target struct name not found")
	}

	if len(*_name) == 0 {
		*_name = structName
	}

	return structName, nil
}

func addPackageNameInFrontOfParamType(targetScope goast.Scope, pkg string) (importPkg bool) {
	targetScope.Node().IterNext(func(n *goast.Node) bool {
		text := n.Text()
		if len(text) == 0 || n.Kind() != kind.ParamType || !helper.isFirstUpperCase(text, '*') {
			return true
		}

		importPkg = true
		n.SetText(helper.insertString(text, "*", pkg+"."))
		return true
	})

	return importPkg
}

func findRelativeScopes(ast goast.Ast, targetScope goast.Scope) (map[string]goast.Scope, []string) {
	var (
		astScopeLength = len(ast.Scope())

		isUnhandledStructScope = make(map[string]goast.Scope, astScopeLength)
		resultScopes           = make(map[string]goast.Scope, astScopeLength)
		resultScopeNames       = make([]string, 0, astScopeLength)

		findRelativeScope func(goast.Scope)
	)

	ast.IterScope(func(sc goast.Scope) bool {
		if sc.Kind() != scope.Type {
			return true
		}

		structName, ok := sc.GetStructName()
		if ok {
			isUnhandledStructScope[structName] = sc
		}

		return true
	})

	findRelativeScope = func(target goast.Scope) {
		target.Node().IterNext(func(n *goast.Node) bool {
			text := helper.tidyString(n.Text(), '*')
			if len(text) == 0 || n.Kind() != kind.ParamType || !helper.isFirstUpperCase(text) {
				return true
			}

			structScope, ok := isUnhandledStructScope[text]
			delete(isUnhandledStructScope, text)
			if ok {
				resultScopes[text] = structScope
				resultScopeNames = append(resultScopeNames, text)
				findRelativeScope(structScope)
			}

			return true
		})
	}

	findRelativeScope(targetScope)

	return resultScopes, resultScopeNames
}
