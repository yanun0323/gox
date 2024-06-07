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
	_replace     = flag.Bool("replace", false, "replace all structure and method if there's already a same structure")
	_debug       = flag.Bool("v", false, "show debug information")
	_help        = flag.Bool("h", false, "show command help")
	_destination = flag.String("destination", "", "target file name to generate implementation")
	_name        = flag.String("name", "", "target implementation structure name")
	_package     = flag.String("package", "", "target implementation structure name")
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
	environmentPrint()
	debugPrint()

	if *_help {
		flag.Usage()
		return
	}

	requireDestination()
	err := run()
	NoError(err)
}

func run() error {
	_, file, err := getDir()
	if err != nil {
		return fmt.Errorf("get directory, err: %w", err)
	}

	ast, err := goast.ParseAst(file)
	if err != nil {
		return fmt.Errorf("parse ast, err: %w", err)
	}

	goLine, err := strconv.Atoi(os.Getenv("GOLINE"))
	if err != nil {
		return fmt.Errorf("parse GOLINE, err: %w", err)
	}

	pkg := os.Getenv("GOPACKAGE")

	// find target interface
	var (
		lineMatched bool
		targetScope goast.Scope
	)

	ast.IterScope(func(s goast.Scope) bool {
		lineMatched = lineMatched || s.Line() == goLine
		if lineMatched && s.Kind() == scope.Type {
			targetScope = s
			return false
		}
		return true
	})

	if targetScope == nil {
		return errors.New("target interface not found")
	}

	// find interface name and set implement name
	interfaceName, ok := targetScope.GetTypeName()
	if !ok || len(interfaceName) == 0 {
		return errors.New("target interface name not found")
	}

	if len(*_name) == 0 {
		*_name = firstLowerCase(interfaceName)
	}

	// add package name in front of paramType
	importPkg := false

	targetScope.Node().IterNext(func(n *goast.Node) bool {
		text := n.Text()
		if len(text) == 0 || n.Kind() != kind.ParamType || !isFirstUpperCase(text) {
			return true
		}

		importPkg = true

		if text[0] == '*' {
			n.SetText("*" + pkg + "." + string(text[1:]))
			return true
		}

		n.SetText(pkg + "." + text)
		return true
	})

	// get func node table
	funcIndexTable := map[string]int{}
	funcSlice := []*goast.Node{}
	targetScope.Node().IterNext(func(n *goast.Node) bool {
		if n.Kind() == kind.FuncName {
			funcIndexTable[n.Text()] = len(funcSlice)
			funcSlice = append(funcSlice, n)
			_ = n.RemovePrev()
		}
		return true
	})

	// add the 'func(x *X)' to the start of func node
	// and the '{}' to the end of func node
	for i, n := range funcSlice {
		println()
		n.DebugPrint()
		println()
		tail := n.IterNext(func(nn *goast.Node) bool { return nn.Next() != nil })
		for {
			switch tail.Kind() {
			case kind.NewLine, kind.Space, kind.Tab, kind.CurlyBracketRight:
				tail = tail.Prev()
				continue
			}
			break
		}
		tail.ReplaceNext(goast.NewNodes(tail.Line(), "{", "\n", "\t", "// TODO: Implement me", "\n", "}", "\n", "\n"))

		head := goast.NewNodes(n.Line(), "func", "(", string(firstLowerCase(*_name)[0]), " ", "*"+*_name, ")", " ")
		headTail := head.IterNext(func(n *goast.Node) bool { return n.Next() != nil })
		headTail.ReplaceNext(n)
		funcSlice[i] = head
	}

	// try get destination file
	destination := *_destination
	if !strings.HasSuffix(destination, ".go") {
		destination = destination + ".go"
	}

	desAst, err := goast.ParseAst(destination)
	if err != nil && !errors.Is(err, goast.ErrNotExist) {
		return fmt.Errorf("parse destination ast, err: %w", err)
	}

	importText := ""
	if importPkg {
		s, err := getSourceImportString()
		if err == nil {
			importText = fmt.Sprintf("import (\n\t\"%s\"\n)", s)
		}
	}

	if desAst == nil {
		// create new ast

		// create struct
		text := fmt.Sprintf(`package %s

%s

type %s struct {
	// TODO: Implement me
}

func New%s() %s.%s {
	// TODO: Implement me
	return &%s{}
}

`, *_package, importText, *_name, interfaceName, pkg, interfaceName, *_name)

		scs, err := goast.ParseScope(0, []byte(text))
		if err != nil {
			return fmt.Errorf("parse scope for creating struct, err: %w", err)
		}

		for _, fnNode := range funcSlice {
			scs = append(scs, goast.NewScope(fnNode.Line(), scope.Func, fnNode))
		}

		newAst, err := goast.NewAst(scs...)
		if err != nil {
			return fmt.Errorf("new ast, err: %w", err)
		}

		println(newAst.String())

		if err := newAst.Save(destination); err != nil {
			return fmt.Errorf("save new ast, err: %w", err)
		}

		return nil
	}

	// insert func into ast
	return nil
}
