package main

import (
	"fmt"
	"log"
	"path/filepath"
)

type pathFunc func(internal, filename string) string

type Package string

func (pkg Package) String() string {
	return string(pkg)
}

const (
	_usecase    Package = "usecase"
	_repository Package = "repository"
)

var (
	_usecasePathFn pathFunc = func(internal, filename string) string {
		return filepath.Join(internal, _usecase.String(), filename)
	}
	_repositoryPathFn pathFunc = func(internal, filename string) string {
		return filepath.Join(internal, _repository.String(), filename)
	}
)

type Generator struct {
	pkg      Package
	receiver string
	replace  bool
	imp      *Implementation
	filename string
	pathFunc pathFunc
}

func NewGenerator(pkg Package, receiver, filename string, replace bool, in *Interface, pf pathFunc) *Generator {
	return &Generator{
		pkg:      pkg,
		receiver: receiver,
		replace:  replace,
		imp:      NewImplementation(pkg, receiver, in),
		filename: filename,
		pathFunc: pf,
	}
}

func (g *Generator) DebugPrint() {
	println("\tpkg:", g.pkg)
	println("\treceiver:", g.receiver)
	println("\treplace:", g.replace)
	println("\timpName:", g.imp.Name)
	println("\timp:", g.imp.upperCase)
	for _, me := range g.imp.upperCase.Me {
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
	parser := NewFileUpdater(g.pkg.String(), g.pathFunc(internalDir, g.filename))
	file, err := parser.Parse()
	if err != nil {
		return fmt.Errorf("parse file (%s), err: %+v", parser.path, err)
	}

	/* find implementation and replace method */
	needInsertImplement := true
	for _, node := range file.Nodes {
		switch node.Type {
		case ntStruct:
			if g.imp.EqualName(node.Name) {
				needInsertImplement = false
				g.imp.SetCompareResultCharacterCase(node.Name)
			}
		case ntMethod:
			if !g.imp.EqualName(node.MethodReceiver) {
				continue
			}

			g.imp.SetCompareResultCharacterCase(node.MethodReceiver)
			/* find match method */
			if g.imp.Content().Me[node.Name] == nil {
				continue
			}

			if g.replace {
				node.Value = g.imp.Content().Me[node.Name].Method
			}

			delete(g.imp.Content().Me, node.Name)
		}
	}

	/* move to new file */

	if needInsertImplement {
		file.Nodes = append(file.Nodes, &FileNode{Value: g.imp.Content().Implementation})
	}

	for _, me := range g.imp.Content().Me {
		file.Nodes = append(file.Nodes, &FileNode{Value: me.Method})
	}

	/* save file */
	if err := parser.SaveFile(file); err != nil {
		return fmt.Errorf("save file (%s), err: %+v", parser.path, err)
	}

	if *_debug {
		log.Default().Printf("save file (%s) succeed", parser.path)
	}

	return nil
}
