package main

import (
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type FileParser struct {
	Dir                                           string
	Entity, UseCase, Repository                   string
	EntityStruct, UseCaseStruct, RepositoryStruct Structure
	entityFile, useCaseFile, repositoryFile       *File
	internalDir                                   string
}

func (fp *FileParser) ParseFile() error {
	internalDir, err := fp.findInternalDir()
	if err != nil {
		return errors.Wrap(err, "find internal dir")
	}
	fp.internalDir = internalDir

	if len(fp.Entity) != 0 {
		f, err := fp.parsingFile(filepath.Join(internalDir, "domain", "entity", fp.Entity), "entity")
		if err != nil {
			return errors.Wrap(err, "parsing entity file")
		}
		fp.entityFile = f
	}

	if len(fp.UseCase) != 0 {
		f, err := fp.parsingFile(filepath.Join(internalDir, "domain", "usecase", fp.UseCase), "usecase")
		if err != nil {
			return errors.Wrap(err, "parsing use case file")
		}
		fp.useCaseFile = f
	}

	if len(fp.Repository) != 0 {
		f, err := fp.parsingFile(filepath.Join(internalDir, "domain", "repository", fp.Repository), "repository")
		if err != nil {
			return errors.Wrap(err, "parsing repository file")
		}
		fp.repositoryFile = f
	}

	return nil
}

func (fp *FileParser) InsertStruct() {
	if fp.entityFile != nil {
		fp.insertStruct(fp.entityFile, fp.EntityStruct)
	}

	if fp.useCaseFile != nil {
		fp.insertStruct(fp.useCaseFile, fp.UseCaseStruct)
	}

	if fp.repositoryFile != nil {
		fp.insertStruct(fp.repositoryFile, fp.RepositoryStruct)
	}
}

func (*FileParser) insertStruct(f *File, st Structure) {
	replacedStruct := false
	replacedMethod := false
	for _, node := range f.Nodes {
		switch node.Type {
		case 1:
			if node.Name == st.StructName {
				node.Value = st.Struct
				replacedStruct = true
			}
		case 2:
			println()
			println("compare node")
			println("node:", fmt.Sprintf("%+v", *node))
			println("structure:", fmt.Sprintf("%+v", st))
			println(node.MethodReceiver == st.StructName)
			println(node.Name == st.MethodName)
			if node.MethodReceiver == st.StructName &&
				node.Name == st.MethodName {
				node.Value = st.Method
				replacedMethod = true
			}
		}
	}
	if !replacedStruct {
		f.Append(st.Struct, 1)
	}

	if !replacedMethod {
		f.Append(st.Method, 2)
	}
}

func (fp *FileParser) SaveFile() error {
	if len(fp.internalDir) == 0 {
		return errors.New("empty internal dir")
	}

	if len(fp.Entity) != 0 {
		err := fp.saveFile(fp.entityFile, filepath.Join(fp.internalDir, "domain", "entity", fp.Entity))
		if err != nil {
			return errors.Errorf("save entity file, err: %+v", err)
		}
	}

	if len(fp.UseCase) != 0 {
		err := fp.saveFile(fp.useCaseFile, filepath.Join(fp.internalDir, "domain", "usecase", fp.UseCase))
		if err != nil {
			return errors.Errorf("save use case file, err: %+v", err)
		}
	}

	if len(fp.Repository) != 0 {
		err := fp.saveFile(fp.repositoryFile, filepath.Join(fp.internalDir, "domain", "repository", fp.Repository))
		if err != nil {
			return errors.Errorf("save repository file, err: %+v", err)
		}
	}

	return nil
}

func (fp *FileParser) findInternalDir() (string, error) {

	sep := string(filepath.Separator)
	spans := strings.Split(fp.Dir, sep+"internal"+sep)
	if len(spans) == 1 {
		return "", errors.New("missing internal folder in working path")
	}
	prefix := strings.Join(spans[:len(spans)-1], sep)

	internal := filepath.Join(prefix, "internal")
	if *_debug {
		println("internal folder path:", internal)
	}

	return internal, nil
}

func (*FileParser) parsingFile(path string, pkg string) (*File, error) {
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		file := &File{Nodes: []*FileNode{}}
		file.Append("package "+pkg, 0)
		return file, nil
	}

	if err != nil {
		return nil, errors.Errorf("open file (%s), err: %+v", path, err)
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.Errorf("read file (%s), err: %+v", path, err)
	}

	rows := strings.Split(string(buf), "\n")
	file := &File{}

	findClose := func(i *int) string {
		bracesPrefix := 0
		findFirst := false
		cache := make([]string, 0, 10)
		for bracesPrefix != 0 || !findFirst {
			for _, char := range rows[*i] {
				switch char {
				case '{':
					findFirst = true
					bracesPrefix++
				case '}':
					bracesPrefix--
				}
			}
			cache = append(cache, rows[*i])
			*i++
		}
		*i--
		return strings.Join(cache, "\n")
	}

	for i := 0; i < len(rows); i++ {
		row := rows[i]
		trimmed := strings.TrimSpace(row)
		if strings.HasPrefix(trimmed, "type") {
			if !strings.HasSuffix(trimmed, "{") {
				file.Append(row, 1)
				continue
			}
			file.Append(findClose(&i), 1)
			continue
		}

		if strings.HasPrefix(row, "func (") {
			file.Append(findClose(&i), 2)
			continue
		}

		if strings.HasPrefix(row, "func") {
			file.Append(findClose(&i), 3)
			continue
		}

		file.Append(row, 0)
	}

	return file, nil
}

func (*FileParser) saveFile(file *File, path string) error {
	content := file.ToString()
	if len(content) == 0 {
		log.Default().Printf("skip empty content for %s", path)
		return nil
	}

	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return errors.Errorf("mkdir (%s) all, err: %+v", dir, err)
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return errors.Errorf("create file (%s), err: %+v", path, err)
	}
	defer f.Close()

	if err := f.Truncate(0); err != nil {
		return errors.Errorf("truncate file (%s), err: %+v", path, err)
	}

	formatted, err := format.Source([]byte(content))
	if err == nil {
		content = string(formatted)
	}

	if _, err := f.WriteString(content); err != nil {
		return errors.Errorf("write content into file (%s), err: %+v", path, err)
	}

	return nil
}

type File struct {
	Nodes []*FileNode
}

func (f *File) Append(val string, t int) {
	name := ""
	methodReceiver := ""
	switch t {
	case 1, 3:
		spans := strings.Split(strings.TrimSpace(val), " ")
		if len(spans) >= 2 {
			name = strings.Split(spans[1], "(")[0]
		}
	case 2:
		trimmed := strings.TrimSpace(val)
		bracketSpans := strings.Split(trimmed, "(")
		if len(bracketSpans) <= 2 {
			break
		}
		bracketSpaceSpans := strings.Split(bracketSpans[1], " ")
		name = bracketSpaceSpans[len(bracketSpaceSpans)-1]
		methodReceiver = f.findMethodReceiver(trimmed)
	}

	fn := FileNode{
		Value:          val,
		Name:           name,
		MethodReceiver: methodReceiver,
		Type:           t,
	}
	if *_debug {
		println("file node:", fmt.Sprintf("%+v", fn))
	}

	f.Nodes = append(f.Nodes, &fn)
}

func (*File) findMethodReceiver(row string) string {
	a := strings.IndexByte(row, '(')
	if a == -1 {
		requireNoError(errors.New("missing '(' in method"))
	}
	b := strings.IndexByte(row, ')')
	if b == -1 {
		requireNoError(errors.New("missing ')' in method"))
	}

	receivers := strings.Split(row[a+1:b], " ")
	receiver := receivers[len(receivers)-1]
	if receiver[0] == '*' {
		return receiver[1:]
	}
	return receiver
}

func (f *File) ToString() string {
	cache := make([]string, 0, len(f.Nodes))
	for _, n := range f.Nodes {
		cache = append(cache, n.Value)
	}
	return strings.Join(cache, "\n")
}

type FileNode struct {
	Value          string
	Name           string
	MethodReceiver string
	Type           int /*
		0=other
		1=struct
		2=method
		3=function
	*/
}
