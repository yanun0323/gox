package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

type pathFunc func(internal, filename string) string

type Package string

func (k Package) String() string {
	return string(k)
}

const (
	_usecase    Package = "usecase"
	_repository Package = "repository"
	_entity     Package = "entity"
	_payload    Package = "payload"
)

var (
	_payloadPathFn pathFunc = func(internal, filename string) string {
		return filepath.Join(internal, "delivery", "http", _payload.String(), filename)
	}
	_entityPathFn pathFunc = func(internal, filename string) string {
		return filepath.Join(internal, "domain", _entity.String(), filename)
	}
	_usecasePathFn pathFunc = func(internal, filename string) string {
		return filepath.Join(internal, "domain", _usecase.String(), filename)
	}
	_repositoryPathFn pathFunc = func(internal, filename string) string {
		return filepath.Join(internal, "domain", _repository.String(), filename)
	}
)

type Generator struct {
	source     *Structure
	Payload    *Element
	Entity     *Element
	Repository *Element
	Usecase    *Element
}

func (g *Generator) ProvideSourceStructure(st Structure, filename string) {
	g.source = &st
	switch currentPackage() {
	case _payload:
		g.Payload.st = &st
		g.Payload.filename = filename
	case _entity:
		g.Entity.st = &st
		g.Entity.filename = filename
	case _repository:
		g.Repository.st = &st
		g.Repository.filename = filename
	case _usecase:
		g.Usecase.st = &st
		g.Usecase.filename = filename
	}
}

func (g *Generator) DebugPrint() {
	println()
	println("Generator:")
	println()
	println("Payload:")
	g.Payload.DebugPrint()
	println("Entity:")
	g.Entity.DebugPrint()
	println("Repository:")
	g.Repository.DebugPrint()
	println("Usecase:")
	g.Usecase.DebugPrint()
	println()
}

func (g *Generator) Gen() {
	g.Payload.Gen(g.source)
	g.Entity.Gen(g.source)
	g.Repository.Gen(g.source)
	g.Usecase.Gen(g.source)
}

func (g *Generator) Save(internalDir string, noComment bool) error {
	if err := g.Payload.Save(internalDir, noComment); err != nil {
		return fmt.Errorf("save payload, err: %+v", err)
	}

	if err := g.Entity.Save(internalDir, noComment); err != nil {
		return fmt.Errorf("save entity, err: %+v", err)
	}

	if err := g.Repository.Save(internalDir, noComment); err != nil {
		return fmt.Errorf("save repository, err: %+v", err)
	}

	if err := g.Usecase.Save(internalDir, noComment); err != nil {
		return fmt.Errorf("save usecase, err: %+v", err)
	}

	return nil
}

type Element struct {
	st *Structure
	me map[string]*Method
	fn map[string]*Function

	ElementParam
}

type Method struct {
	MethodName string
	Method     string
}

type Function struct {
	FunctionName string
	Function     string
}

type ElementParam struct {
	pkg                                                  Package
	filename                                             string
	replace                                              bool
	toPayload, toEntity, toRepository, toUseCase         bool
	fromPayload, fromEntity, fromRepository, fromUseCase bool
	unix                                                 bool
	timestamp                                            bool
	genJson                                              bool
	pathFunc                                             pathFunc
}

func NewElement(param ElementParam) *Element {
	if len(param.filename) == 0 && param.pkg != currentPackage() {
		return nil
	}

	return &Element{
		me:           make(map[string]*Method, 4),
		fn:           make(map[string]*Function, 4),
		ElementParam: param,
	}
}

func (elem *Element) DebugPrint() {
	if elem == nil {
		println("\tnil")
		println()
		return
	}
	println("\tstructName:", elem.st.StructName)
	println("\tstruct", elem.st.Struct)
	for _, me := range elem.me {
		println("\t\tmethodName:", me.MethodName)
		println("\t\tmethod:", me.Method)
	}
	println("\treplace:", elem.replace)
	println("\ttoPayload:", elem.toPayload)
	println("\ttoEntity:", elem.toEntity)
	println("\ttoRepository:", elem.toRepository)
	println("\ttoUseCase:", elem.toUseCase)
	println("\tunix:", elem.unix)
	println("\ttimestamp:", elem.timestamp)
	println("\tgenJson:", elem.genJson)
	println("\tfilename:", elem.filename)
	println("\tpkg:", elem.pkg)
	println()
}

func (elem *Element) Gen(source *Structure) {
	if elem == nil {
		return
	}

	if elem.st == nil {
		elem.st = NewStructureFrom(elem.pkg, source, elem.unix, elem.timestamp, elem.genJson)
	}

	if elem.toPayload {
		m := elem.st.GenMethod(_payload, "ToPayload")
		elem.me[m.MethodName] = m
	}

	if elem.toEntity {
		m := elem.st.GenMethod(_entity, "ToEntity")
		elem.me[m.MethodName] = m
	}

	if elem.toRepository {
		m := elem.st.GenMethod(_repository, "ToRepository")
		elem.me[m.MethodName] = m
	}

	if elem.toUseCase {
		m := elem.st.GenMethod(_usecase, "ToUseCase")
		elem.me[m.MethodName] = m
	}

	if elem.fromPayload {
		f := elem.st.GenFunction(_payload, "FromPayload")
		elem.fn[f.FunctionName] = f
	}

	if elem.fromEntity {
		f := elem.st.GenFunction(_entity, "FromEntity")
		elem.fn[f.FunctionName] = f
	}

	if elem.fromRepository {
		f := elem.st.GenFunction(_repository, "FromRepository")
		elem.fn[f.FunctionName] = f
	}

	if elem.fromUseCase {
		f := elem.st.GenFunction(_usecase, "FromUseCase")
		elem.fn[f.FunctionName] = f
	}
}

const (
	_commentKeyword  = "esc-model-gen"
	_commentTemplate = `
	// Code generated by by esc-model-gen.
	// Source: %s%s`
)

func (elem *Element) Save(internalDir string, noComment bool) error {
	if elem == nil {
		return nil
	}

	isSourceStruct := currentPackage() == elem.pkg

	/* load file */
	parser := NewFileUpdater(elem.pkg.String(), elem.pathFunc(internalDir, elem.filename))
	file, err := parser.Parse()
	if err != nil {
		return fmt.Errorf("parse file (%s), err: %+v", parser.path, err)
	}
	comment := ""

	if !noComment {
		replaceCaution := ""
		if elem.replace {
			replaceCaution = `
			//
			// REPLACE CAUTION!!`
		}

		source := currentPackage().String() + "." + elem.st.StructName
		comment = fmt.Sprintf(_commentTemplate, source, replaceCaution)
	}

	needInsertStruct := true
	/* replace struct and method */
	for _, node := range file.Nodes {
		switch node.Type {
		case ntStruct:
			if node.Name != elem.st.StructName {
				continue
			}

			needInsertStruct = false

			if strings.Contains(node.Comment, _commentKeyword) && !isSourceStruct {
				node.Comment = comment
			}

			if elem.replace && !isSourceStruct {
				node.Comment = comment
				node.Value = elem.st.Struct
			}
		case ntMethod:
			if node.MethodReceiver != elem.st.StructName {
				continue
			}

			/* find match method */
			if elem.me[node.Name] == nil {
				continue
			}

			if strings.Contains(node.Comment, _commentKeyword) {
				node.Comment = comment
			}

			if elem.replace {
				node.Comment = comment
				node.Value = elem.me[node.Name].Method
			}

			delete(elem.me, node.Name)
		case ntFunc:
			/* find match function */
			if elem.fn[node.Name] == nil {
				continue
			}

			if strings.Contains(node.Comment, _commentKeyword) {
				node.Comment = comment
			}

			if elem.replace {
				node.Comment = comment
				node.Value = elem.fn[node.Name].Function
			}

			delete(elem.fn, node.Name)
		}
	}

	structMaxLine := 5
	newFile := &File{
		Nodes: make([]*FileNode, 0, len(file.Nodes)+structMaxLine),
	}

	/* move to new file */
	for _, node := range file.Nodes {
		switch node.Type {
		case ntStruct:
			if node.Name != elem.st.StructName {
				newFile.Nodes = append(newFile.Nodes, node)
				continue
			}

			newFile.Nodes = append(newFile.Nodes, node)

			for _, fn := range elem.fn {
				newFile.Nodes = append(newFile.Nodes, &FileNode{
					Value:   fn.Function,
					Comment: comment,
				})
			}

			for _, me := range elem.me {
				newFile.Nodes = append(newFile.Nodes, &FileNode{
					Value:   me.Method,
					Comment: comment,
				})
			}
		default:
			newFile.Nodes = append(newFile.Nodes, node)
		}

	}

	if needInsertStruct {
		newFile.Nodes = append(newFile.Nodes, &FileNode{
			Value:   elem.st.Struct,
			Comment: comment,
		})

		for _, fn := range elem.fn {
			newFile.Nodes = append(newFile.Nodes, &FileNode{
				Value:   fn.Function,
				Comment: comment,
			})
		}

		for _, me := range elem.me {
			newFile.Nodes = append(newFile.Nodes, &FileNode{
				Value:   me.Method,
				Comment: comment,
			})
		}
	}

	/* save file */
	if err := parser.SaveFile(newFile); err != nil {
		return fmt.Errorf("save file (%s), err: %+v", parser.path, err)
	}

	if *_debug {
		log.Default().Printf("save file (%s) succeed", parser.path)
	}

	return nil
}
