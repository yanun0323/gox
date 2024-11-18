package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
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
	_noEmbed     = flag.Bool("noembed", false, "skip implementing embed interface functions")
	_destination = flag.String("destination", "", "target file name to generate implementation")
	_package     = flag.String("package", "", "target implementation structure name")
	_struct      = flag.Bool("struct", false, "generate struct")
	_name        = flag.String("name", "", "target implementation structure name")
	_constructor = flag.Bool("constructor", false, "generate constructor function")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "%s: generate an implementation from the interface \n", _commandName)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t-h\t\t\t\tshow usage\n")
	fmt.Fprintf(os.Stderr, "\t-struct\t\t\tgenerate struct\n")
	fmt.Fprintf(os.Stderr, "\t-name\t\t\t\timplemented struct name\t\t\t-name=usecase\n")
	fmt.Fprintf(os.Stderr, "\t-constructor\t\t\tgenerate constructor function\n")
	fmt.Fprintf(os.Stderr, "\t-package\t(require)\timplemented struct package name\n")
	fmt.Fprintf(os.Stderr, "\t-destination\t(require)\tgenerated file path\t\t\t-destination=../../usecase/member_usecase.go\n")
	fmt.Fprintf(os.Stderr, "\t-replace\t\t\tforce replace exist struct/func/method\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\texample:\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\t//go:generate %s -destination=../../usecase/member.go -name=usecase -replace -constructor\n", _commandName)
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

	// for _, ev := range []string{"GOARCH", "GOOS", "GOFILE", "GOLINE", "GOPACKAGE", "DOLLAR", "GOMODULE"} {
	// 	fmt.Println("\t", ev, "=", os.Getenv(ev))
	// }

	helper.requireTag()

	ast, goLine, pkg, curDir, err := parseAstFromGoGenerator()
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

	importPkg := false
	isSameFolder := isDestinationSameFolderToSource(curDir)
	if !isSameFolder {
		importPkg = addPackageNameInFrontOfParamType(targetScope, pkg)
	}

	methodNodes, methodNodesIndexTable := getInterfaceMethodNodes(ast, targetScope)

	desAst, destination, err := tryGetDestinationFile()
	if err != nil {
		return err
	}

	destinationFileNotFound := desAst == nil

	if destinationFileNotFound {
		return createNewDestinationFileAndSave(
			importPkg,
			isSameFolder,
			interfaceName,
			pkg,
			destination,
			methodNodes,
		)
	}

	return updateDestinationFileAndSave(
		desAst,
		isSameFolder,
		interfaceName,
		pkg,
		destination,
		methodNodes,
		methodNodesIndexTable,
	)
}

func parseAstFromGoGenerator() (ast goast.Ast, goLine int, pkg string, curDir string, err error) {
	dir, file, err := helper.getDir()
	if err != nil {
		return nil, 0, "", "", fmt.Errorf("get directory, err: %w", err)
	}

	astObj, err := goast.ParseAst(file)
	if err != nil {
		return nil, 0, "", "", fmt.Errorf("parse ast, err: %w", err)
	}

	goLineNum, err := strconv.Atoi(os.Getenv("GOLINE"))
	if err != nil {
		return nil, 0, "", "", fmt.Errorf("parse GOLINE, err: %w", err)
	}

	pkgName := os.Getenv("GOPACKAGE")

	return astObj, goLineNum, pkgName, dir, nil
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

func isDestinationSameFolderToSource(curDir string) bool {
	targetFile := *_destination
	if !filepath.IsAbs(targetFile) {
		targetFile, _ = filepath.Abs(targetFile)
	}
	targetDir := filepath.Dir(targetFile)
	return curDir == targetDir
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

func getInterfaceMethodNodes(ast goast.Ast, targetScope goast.Scope) ([]*goast.Node, map[string]int) {
	funcNodes := []*goast.Node{}
	funcNodesIndexTable := map[string]int{}
	embedInterfaceNames := []string{}
	targetScopeName, _ := targetScope.GetInterfaceName()

	targetScope.Node().IterNext(func(n *goast.Node) bool {
		switch n.Kind() {
		case kind.FuncName:
			funcNodesIndexTable[n.Text()] = len(funcNodes)
			funcNodes = append(funcNodes, n)
			_ = n.RemovePrev()
		case kind.TypeName:
			if !*_noEmbed {
				if len(targetScopeName) != 0 && targetScopeName != n.Text() {
					embedInterfaceNames = append(embedInterfaceNames, n.Text())
				}
			}
		}
		return true
	})

	for _, name := range embedInterfaceNames {
		s, ok := findInterfaceScope(ast, name)
		if !ok {
			continue
		}

		fns, fnit := getInterfaceMethodNodes(ast, s)
		idxOffset := len(funcNodes)
		funcNodes = append(funcNodes, fns...)
		for k, v := range fnit {
			funcNodesIndexTable[k] = v + idxOffset
		}
	}

	for _, fnNode := range funcNodes {
		fnNode.IterNext(func(n *goast.Node) bool {
			if n.Kind() == kind.NewLine {
				_ = n.RemovePrev()
				return false
			}

			return true
		})
	}

	return funcNodes, funcNodesIndexTable
}

func findInterfaceScope(ast goast.Ast, name string) (goast.Scope, bool) {
	var result goast.Scope
	ast.IterScope(func(s goast.Scope) bool {
		in, ok := s.GetInterfaceName()
		if ok && in == name {

			result = s
			return false
		}

		return true
	})
	return result, result != nil
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

func createNewDestinationFileAndSave(moduleName, isSameFolder bool, interfaceName, pkg, destination string, methodNodes []*goast.Node) error {

	text := fmt.Sprintf("%s\n%s\n%s\n%s\n",
		genPackageString(),
		genImportString(moduleName),
		genImplementationString(),
		genConstructorString(interfaceName, pkg, isSameFolder),
	)

	scs, err := goast.ParseScope(0, []byte(text))
	if err != nil {
		return fmt.Errorf("parse scope for creating struct, err: %w", err)
	}

	for _, fnNode := range methodNodes {
		fnNode = addMethodImplementationPrefixSuffix(fnNode, "")
		scs = append(scs, goast.NewScope(fnNode.Line(), scope.Func, fnNode))
	}

	newAst, err := goast.NewAst(scs...)
	if err != nil {
		return fmt.Errorf("new ast, err: %w", err)
	}

	{
		// // XXX: Remove me
		// bd := strings.Builder{}
		// for _, sc := range newAst.Scope() {
		// 	sc.Node().IterNext(func(n *goast.Node) bool {
		// 		bd.WriteString(n.Text())
		// 		return true
		// 	})
		// }

		// result, err := imports.Process(destination, []byte(bd.String()), nil)
		// if err != nil {
		// 	return fmt.Errorf("auto imports, err: %w", err)
		// }

		// println("imported:\n", string(result))
	}

	if err := newAst.Save(destination, true); err != nil {
		return fmt.Errorf("save new ast, err: %w", err)
	}

	return nil
}

func genPackageString() string {
	return fmt.Sprintf("package %s\n", *_package)
}

func genImportString(moduleName bool) string {
	if !moduleName {
		return ""
	}

	alias, importPath, err := helper.getSourceImportString()
	if err != nil {
		return ""
	}

	if len(alias) != 0 {
		return fmt.Sprintf("import (\n\t%s \"%s\"\n)\n", alias, importPath)

	}

	return fmt.Sprintf("import (\n\t\"%s\"\n)\n", importPath)
}

func genImplementationString() string {
	if !*_struct {
		return ""
	}

	if *_replace {
		return fmt.Sprintf("type %s struct {\n\t// Replace by %s\n\t// TODO: Implement me\n}\n", *_name, _commandName)
	} else {
		return fmt.Sprintf("type %s struct {\n\t// TODO: Implement me\n}\n", *_name)
	}
}

func genConstructorString(interfaceName, pkg string, isSameFolder bool) string {
	if !*_constructor {
		return ""
	}

	returnType := pkg + "." + interfaceName
	if isSameFolder {
		returnType = interfaceName
	}

	if *_replace {
		return fmt.Sprintf("func %s() (%s, error) {\n\t// Replace by %s\n\t// TODO: Implement me\n\treturn &%s{}, nil\n}\n", constructFuncName(interfaceName), returnType, _commandName, *_name)
	} else {
		return fmt.Sprintf("func %s() (%s, error) {\n\t// TODO: Implement me\n\treturn &%s{}, nil\n}\n", constructFuncName(interfaceName), returnType, *_name)
	}
}

// add the 'func(x *X)' to the start of func node
// and the '{}' to the end of func node
func addMethodImplementationPrefixSuffix(methodNode *goast.Node, receiverName string) *goast.Node {
	tail := methodNode.Last()
	for {
		switch tail.Kind() {
		case kind.NewLine, kind.Space, kind.Tab, kind.CurlyBracketRight:
			tail = tail.Prev()
			continue
		}
		break
	}

	tail.ReplaceNext(goast.NewNodes(tail.Line(), "{"))
	tail = tail.Last()

	if *_replace {
		tail.ReplaceNext(goast.NewNodes(tail.Line(), "\n", "\t", fmt.Sprintf("// Replace by %s", _commandName)))
		tail = tail.Last()
	}

	tail.ReplaceNext(goast.NewNodes(tail.Line(), "\n", "\t", "// TODO: Implement me"))
	tail = tail.Last()

	// if rn, ok := generateReturnValue(methodNode); ok && rn != nil {
	// 	tail.ReplaceNext(goast.NewNodes(tail.Line(), "\n", "\t", "return", " "))
	// 	tail = tail.Last()

	// 	tail.ReplaceNext(rn)
	// 	tail = tail.Last()
	// }

	tail.ReplaceNext(goast.NewNodes(tail.Line(), "\n", "}", "\n", "\n", "\n"))
	tail = tail.Last()

	if len(receiverName) == 0 {
		receiverName = string(helper.firstLowerCase(*_name)[0])
		lowercaseName := strings.ToLower(*_name)
		if strings.Contains(lowercaseName, "usecase") {
			receiverName = "use"
		} else if strings.Contains(lowercaseName, "repo") {
			receiverName = "repo"
		}
	}

	head := goast.NewNodes(methodNode.Line(), "\n", "func", "(", receiverName, " ", "*"+*_name, ")", " ")
	head.Last().ReplaceNext(methodNode)

	return head
}

func generateReturnValue(methodNode *goast.Node) (*goast.Node, bool) {
	funcParenthesisCount := 0
	returnValueHead := methodNode.IterNext(func(n *goast.Node) bool {
		return n.Kind() != kind.ParenthesisLeft
	}).IterNext(func(n *goast.Node) bool {
		switch n.Kind() {
		case kind.ParenthesisLeft:
			funcParenthesisCount++
		case kind.ParenthesisRight:
			funcParenthesisCount--
		}

		return funcParenthesisCount != 0
	}).Next()

	hasReturnValue := false
	returnValueHead.IterNext(func(n *goast.Node) bool {
		switch n.Kind() {
		case kind.Space:
			return true
		case kind.NewLine:
			return false
		default:
			hasReturnValue = true
			return false
		}
	})

	if !hasReturnValue {
		return nil, false
	}

	parenthesisReturnValue := false
	returnValueHead.IterNext(func(n *goast.Node) bool {
		switch n.Kind() {
		case kind.ParenthesisLeft:
			parenthesisReturnValue = true
			return false
		case kind.Space, kind.Comment:
			return true
		default:
			return false
		}
	})

	result := &goast.Node{}
	addResult := func(n *goast.Node) {
		result.ReplaceNext(n)
	}

	cleanResult := func() *goast.Node {
		result = result.First().Next()
		result.RemovePrev()
		return result
	}

	if !parenthesisReturnValue {
		returnValueHead.IterNext(func(n *goast.Node) bool {
			switch n.Kind() {
			case kind.NewLine:
				return false
			case kind.CurlyBracketLeft:
				return true
			default:
				addResult(n.Copy())
				return true
			}
		})

		return cleanResult(), true
	}

	ns := helper.extractParenthesisParameters(returnValueHead)
	if len(ns) == 0 {
		return nil, false
	}

	for _, n := range ns {
		n.IterNext(func(n *goast.Node) bool {
			switch n.Kind() {
			case kind.Comment:
			default:
				addResult(n.Copy())
			}

			return true
		})

		addResult(goast.NewNode(0, ",", kind.Comma))
	}

	return cleanResult(), true
}

func constructFuncName(interfaceName string) string {
	return fmt.Sprintf("New%s", interfaceName)
}

func updateDestinationFileAndSave(desAst goast.Ast, isSameFolder bool, interfaceName, pkg string, destination string, methodNodes []*goast.Node, methodNodesIndexTable map[string]int) error {
	// find implementation is exist or not
	var (
		isPackageExist     bool
		isStructExist      bool
		isConstructorExist bool
		scopes             []goast.Scope

		newFuncName = constructFuncName(interfaceName)
	)

	existReceiverName := ""

	if *_replace {
		/* keep other code */
		desAst.IterScope(func(sc goast.Scope) bool {
			if sc.Kind() == scope.Package {
				isPackageExist = true
			}

			/* keep struct */
			name, ok := sc.GetStructName()
			if ok && strings.EqualFold(name, *_name) {
				isStructExist = true
				// return true
			}

			/* keep construct */
			fnName, ok := sc.GetFuncName()
			if ok && strings.EqualFold(fnName, newFuncName) {
				isConstructorExist = true
				// return true
			}

			/* drop method */
			receiverName, receiverType, _, ok := findScopeMethod(sc)
			if ok && helper.EqualFold(receiverType, *_name, '*') {
				if len(receiverName) != 0 {
					existReceiverName = receiverName
				}

				return true
			}

			scopes = append(scopes, sc)

			return true
		})
	} else {
		/* find isStructExist, isConstructorExist and if methods exist */
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

			receiverName, receiverType, methodName, ok := findScopeMethod(sc)
			if !ok {
				return true
			}

			if !helper.EqualFold(receiverType, *_name, '*') {
				return true
			}

			if len(receiverName) != 0 {
				existReceiverName = receiverName
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

	if !isConstructorExist && *_constructor {
		scs, err := goast.ParseScope(0, []byte(genConstructorString(interfaceName, pkg, isSameFolder)))
		if err != nil {
			return fmt.Errorf("parse scope for constructor, err: %w", err)
		}

		scopes = append(scopes, scs...)
	}

	for _, fnNode := range methodNodes {
		if fnNode == nil {
			continue
		}
		fnNode = addMethodImplementationPrefixSuffix(fnNode, existReceiverName)
		scopes = append(scopes, goast.NewScope(0, scope.Func, fnNode))
	}

	resultAst := desAst.SetScope(scopes)

	return resultAst.Save(destination, true)
}

func findScopeMethod(sc goast.Scope) (receiverName, receiverType, methodName string, ok bool) {
	if sc.Kind() != scope.Func {
		return "", "", "", false
	}

	var (
		rName   string
		rnFound bool
		rType   string
		rtFound bool
		mName   string
		mFound  bool
	)
	sc.Node().IterNext(func(n *goast.Node) bool {
		switch n.Kind() {
		case kind.MethodReceiverName:
			rName = n.Text()
			rnFound = true
		case kind.MethodReceiverType:
			rType = n.Text()
			rtFound = true
		case kind.FuncName:
			mName = n.Text()
			mFound = true
		}

		return !mFound
	})

	return rName, rType, mName, rnFound && rtFound && mFound
}
