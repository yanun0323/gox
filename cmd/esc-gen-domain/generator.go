package main

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type pathFunc func(internal, filename string) string

var (
	_usecase    = "usecase"
	_repository = "repository"

	_usecasePathFn pathFunc = func(internal, filename string) string {
		return filepath.Join(internal, _usecase, filename)
	}
	_repositoryPathFn pathFunc = func(internal, filename string) string {
		return filepath.Join(internal, _repository, filename)
	}
)

type Generator struct {
	pkg      string
	receiver string
	replace  bool
	imp      *Implementation
	filename string
	pathFunc pathFunc
}

func NewGenerator(pkg, receiver, filename string, replace bool, inter *Interface, pf pathFunc) *Generator {
	return &Generator{
		pkg:      pkg,
		receiver: receiver,
		replace:  replace,
		imp:      NewImplementation(pkg, receiver, inter),
		filename: filename,
		pathFunc: pf,
	}
}

func (g *Generator) DebugPrint() {
	println("\tpkg:", g.pkg)
	println("\treceiver:", g.receiver)
	println("\treplace:", g.replace)
	println("\timpName:", g.imp.ImplementationName)
	println("\timp:", g.imp.Implementation)
	for _, me := range g.imp.Me {
		println("\t\tmeName:", me.MethodName)
		println("\t\tme:", me.Method)
	}
	println("\tfilename:", g.filename)
	println("\tpathFunc:", g.pathFunc)
}

func (g *Generator) Save(internalDir string) error {
	if g == nil {
		return nil
	}

	/* load file */
	parser := NewFileUpdater(g.pkg, g.pathFunc(internalDir, g.filename))
	file, err := parser.Parse()
	if err != nil {
		return errors.Errorf("parse file (%s), err: %+v", parser.path, err)
	}

	/* find implementation and replace method */
	needInsertImplement := true
	for _, node := range file.Nodes {
		switch node.Type {
		case ntStruct:
			if strings.EqualFold(node.Name, g.imp.ImplementationName) {
				needInsertImplement = false
			}
		case ntMethod:
			if node.MethodReceiver != g.imp.ImplementationName {
				continue
			}

			/* find match method */
			if g.imp.Me[node.Name] == nil {
				continue
			}

			if g.replace {
				node.Value = g.imp.Me[node.Name].Method
			}

			delete(g.imp.Me, node.Name)
		}
	}

	/* move to new file */

	if needInsertImplement {
		file.Nodes = append(file.Nodes, &FileNode{Value: g.imp.Implementation})
	}

	for _, me := range g.imp.Me {
		file.Nodes = append(file.Nodes, &FileNode{Value: me.Method})
	}

	/* save file */
	if err := parser.SaveFile(file); err != nil {
		return errors.Errorf("save file (%s), err: %+v", parser.path, err)
	}

	log.Default().Printf("save file (%s) succeed", parser.path)

	return nil
}