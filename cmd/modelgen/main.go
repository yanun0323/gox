package main

import "flag"

const _commandName = "modelgen"

var (
	_replace = flag.Bool("replace", false, "replace all structure and method if there's already a same structure")
	_debug   = flag.Bool("v", false, "show debug information")
	_help    = flag.Bool("h", false, "show command help")
)

func main() {
	println("Hello, world")
}
