package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/yanun0323/goast"
	"github.com/yanun0323/goast/kind"
	"github.com/yanun0323/goast/scope"
)

const _commandName = "domaingen"

var (
	_help        = flag.Bool("h", false, "show command help")
	_debug       = flag.Bool("v", false, "show debug information")
	_replace     = flag.Bool("replace", false, "replace all structure and method if there's already a same structure")
	_destination = flag.String("destination", "", "target file name to generate implementation")
	_package     = flag.String("package", "", "target implementation structure name")
	_name        = flag.String("name", "", "target implementation structure name")
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

	targetScope, err := findTargetInterface(ast, goLine)
	if err != nil {
		return err
	}

	interfaceName, err := findInterfaceNameAndSetImplementName(targetScope)
	if err != nil {
		return err
	}

	importPkg := addPackageNameInFrontOfParamType(targetScope, pkg)

	methodNodes, methodNodesIndexTable := getInterfaceMethodNodes(targetScope)

	addMethodImplementationPrefixSuffix(methodNodes)

	desAst, destination, err := tryGetDestinationFile()
	if err != nil {
		return err
	}

	destinationFileNotFound := desAst == nil

	if destinationFileNotFound {
		return createNewDestinationFileAndSave(
			importPkg,
			interfaceName,
			pkg,
			destination,
			methodNodes,
		)
	}

	return updateDestinationFileAndSave(
		desAst,
		interfaceName,
		pkg,
		destination,
		methodNodes,
		methodNodesIndexTable,
	)
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

func findTargetInterface(ast goast.Ast, goLine int) (goast.Scope, error) {
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
				if _, ok := s.GetInterfaceName(); ok {
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
		return nil, errors.New("target interface not found")
	}

	return targetScope, nil
}

func findInterfaceNameAndSetImplementName(targetScope goast.Scope) (string, error) {
	interfaceName, ok := targetScope.GetTypeName()
	if !ok || len(interfaceName) == 0 {
		return "", errors.New("target interface name not found")
	}

	if len(*_name) == 0 {
		*_name = helper.firstLowerCase(interfaceName)
	}

	return interfaceName, nil
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

func getInterfaceMethodNodes(targetScope goast.Scope) ([]*goast.Node, map[string]int) {
	funcNodes := []*goast.Node{}
	funcNodesIndexTable := map[string]int{}

	targetScope.Node().IterNext(func(n *goast.Node) bool {
		if n.Kind() == kind.FuncName {
			funcNodesIndexTable[n.Text()] = len(funcNodes)
			funcNodes = append(funcNodes, n)
			_ = n.RemovePrev()
		}
		return true
	})

	return funcNodes, funcNodesIndexTable
}

// add the 'func(x *X)' to the start of func node
// and the '{}' to the end of func node
func addMethodImplementationPrefixSuffix(methodNodes []*goast.Node) {
	for i, n := range methodNodes {
		tail := n.IterNext(func(nn *goast.Node) bool { return nn.Next() != nil })
		for {
			switch tail.Kind() {
			case kind.NewLine, kind.Space, kind.Tab, kind.CurlyBracketRight:
				tail = tail.Prev()
				continue
			}
			break
		}

		if *_replace {
			tail.ReplaceNext(goast.NewNodes(tail.Line(), "{", "\n", "\t", fmt.Sprintf("// Replace by %s", _commandName), "\n", "// TODO: Implement me", "\n", "}", "\n", "\n", "\n"))
		} else {
			tail.ReplaceNext(goast.NewNodes(tail.Line(), "{", "\n", "\t", "// TODO: Implement me", "\n", "}", "\n", "\n", "\n"))
		}

		head := goast.NewNodes(n.Line(), "\n", "func", "(", string(helper.firstLowerCase(*_name)[0]), " ", "*"+*_name, ")", " ")
		headTail := head.IterNext(func(n *goast.Node) bool { return n.Next() != nil })
		headTail.ReplaceNext(n)
		methodNodes[i] = head
	}
}

func tryGetDestinationFile() (goast.Ast, string, error) {
	destination := *_destination
	if !strings.HasSuffix(destination, ".go") {
		destination = destination + ".go"
	}

	desAst, err := goast.ParseAst(destination)
	if err != nil && !errors.Is(err, goast.ErrNotExist) {
		return nil, "", fmt.Errorf("parse destination ast, err: %w", err)
	}

	return desAst, destination, nil
}

func createNewDestinationFileAndSave(importPkg bool, interfaceName, pkg, destination string, methodNodes []*goast.Node) error {
	text := fmt.Sprintf("%s\n%s\n%s\n%s\n",
		genPackageString(),
		genImportString(importPkg),
		genImplementationString(),
		genConstructorString(interfaceName, pkg),
	)

	scs, err := goast.ParseScope(0, []byte(text))
	if err != nil {
		return fmt.Errorf("parse scope for creating struct, err: %w", err)
	}

	for _, fnNode := range methodNodes {
		scs = append(scs, goast.NewScope(fnNode.Line(), scope.Func, fnNode))
	}

	newAst, err := goast.NewAst(scs...)
	if err != nil {
		return fmt.Errorf("new ast, err: %w", err)
	}

	if err := newAst.Save(destination); err != nil {
		return fmt.Errorf("save new ast, err: %w", err)
	}

	return nil
}

func genPackageString() string {
	return fmt.Sprintf("package %s\n", *_package)
}

func genImportString(importPkg bool) string {
	importText := ""
	if importPkg {
		s, err := helper.getSourceImportString()
		if err == nil {
			importText = fmt.Sprintf("import (\n\t\"%s\"\n)\n", s)
		}
	}

	return importText
}

func genImplementationString() string {
	if *_replace {
		return fmt.Sprintf("type %s struct {\n\t// Replace by %s\n\t// TODO: Implement me\n}\n", *_name, _commandName)
	} else {
		return fmt.Sprintf("type %s struct {\n\t// TODO: Implement me\n}\n", *_name)
	}
}

func genConstructorString(interfaceName, pkg string) string {
	if *_replace {
		return fmt.Sprintf("func %s() %s.%s {\n\t// Replace by %s\n\t// TODO: Implement me\n\treturn &%s{}\n}\n", constructFuncName(interfaceName), pkg, interfaceName, _commandName, *_name)
	} else {
		return fmt.Sprintf("func %s() %s.%s {\n\t// TODO: Implement me\n\treturn &%s{}\n}\n", constructFuncName(interfaceName), pkg, interfaceName, *_name)
	}
}

func constructFuncName(interfaceName string) string {
	return fmt.Sprintf("New%s", interfaceName)
}

func updateDestinationFileAndSave(desAst goast.Ast, interfaceName, pkg string, destination string, methodNodes []*goast.Node, methodNodesIndexTable map[string]int) error {
	// find implementation is exist or not
	var (
		isPackageExist     bool
		isStructExist      bool
		isConstructorExist bool
		scopes             []goast.Scope

		newFuncName = constructFuncName(interfaceName)
	)

	if *_replace {
		desAst.IterScope(func(sc goast.Scope) bool {
			if sc.Kind() == scope.Package {
				isPackageExist = true
			}

			name, ok := sc.GetStructName()
			if ok && strings.EqualFold(name, *_name) {
				return true
			}

			fnName, ok := sc.GetFuncName()
			if ok && strings.EqualFold(fnName, newFuncName) {
				return true
			}

			receiverType, _, ok := findScopeMethod(sc)
			if ok && helper.EqualFold(receiverType, *_name, '*') {
				return true
			}

			scopes = append(scopes, sc)

			return true
		})
	} else {
		// find isStructExist, isConstructorExist and if methods exist
		desAst.IterScope(func(sc goast.Scope) bool {
			scopes = append(scopes, sc)

			if sc.Kind() == scope.Package {
				isPackageExist = true
			}

			name, ok := sc.GetStructName()
			if ok && strings.EqualFold(name, *_name) {
				isStructExist = true
			}

			fnName, ok := sc.GetFuncName()
			if ok && strings.EqualFold(fnName, newFuncName) {
				isConstructorExist = true
			}

			receiverType, methodName, ok := findScopeMethod(sc)
			if !ok {
				return true
			}

			if !helper.EqualFold(receiverType, *_name, '*') {
				return true
			}

			i := methodNodesIndexTable[methodName]
			if i < len(methodNodes) {
				methodNodes[i] = nil
			}

			return true
		})
	}

	if !isPackageExist {
		scs, err := goast.ParseScope(0, []byte(genPackageString()))
		if err != nil {
			return fmt.Errorf("parse scope for package, err: %w", err)
		}

		scopes = append(scs, scopes...)
	}

	if !isStructExist {
		scs, err := goast.ParseScope(0, []byte(genImplementationString()))
		if err != nil {
			return fmt.Errorf("parse scope for struct, err: %w", err)
		}

		scopes = append(scopes, scs...)
	}

	if !isConstructorExist {
		scs, err := goast.ParseScope(0, []byte(genConstructorString(interfaceName, pkg)))
		if err != nil {
			return fmt.Errorf("parse scope for constructor, err: %w", err)
		}

		scopes = append(scopes, scs...)
	}

	for _, fnNode := range methodNodes {
		if fnNode == nil {
			continue
		}

		scopes = append(scopes, goast.NewScope(0, scope.Func, fnNode))
	}

	resultAst := desAst.SetScope(scopes)

	return resultAst.Save(destination)
}

func findScopeMethod(sc goast.Scope) (receiverType string, methodName string, ok bool) {
	if sc.Kind() != scope.Func {
		return "", "", false
	}

	var (
		rName  string
		rFound bool
		mName  string
		mFound bool
	)
	sc.Node().IterNext(func(n *goast.Node) bool {
		switch n.Kind() {
		case kind.MethodReceiverType:
			rName = n.Text()
			rFound = true
		case kind.FuncName:
			mName = n.Text()
			mFound = true
		}

		return !mFound
	})

	return rName, mName, rFound && mFound
}
