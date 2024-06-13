package main

import (
	"flag"
	"fmt"
	"os"
)

const _commandName = "modelgen"

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

	return nil
}
