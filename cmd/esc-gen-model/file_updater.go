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

type FileUpdater struct {
	path string
	pkg  string
}

func NewFileUpdater(pkg, path string) *FileUpdater {
	return &FileUpdater{
		path: path,
		pkg:  pkg,
	}
}

func (p *FileUpdater) Parse() (*File, error) {
	if len(p.path) == 0 {
		return nil, errors.New("empty path")
	}

	f, err := os.Open(p.path)
	if errors.Is(err, os.ErrNotExist) {
		file := &File{Nodes: []*FileNode{}}
		file.Append("package "+p.pkg, ntOther)
		return file, nil
	}

	if err != nil {
		return nil, errors.Errorf("open file (%s), err: %+v", p.path, err)
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.Errorf("read file (%s), err: %+v", p.path, err)
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
				file.Append(row, ntStruct)
				continue
			}
			file.Append(findClose(&i), ntStruct)
			continue
		}

		if strings.HasPrefix(row, "func (") {
			file.Append(findClose(&i), ntMethod)
			continue
		}

		if strings.HasPrefix(row, "func") {
			file.Append(findClose(&i), ntFunc)
			continue
		}

		file.Append(row, ntOther)
	}

	return file, nil
}

func (p *FileUpdater) SaveFile(file *File) error {
	if len(p.path) == 0 {
		return errors.New("empty path")
	}

	content := file.ToString()
	if len(content) == 0 {
		log.Default().Printf("skip empty content for %s", p.path)
		return nil
	}

	dir := filepath.Dir(p.path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return errors.Errorf("mkdir (%s) all, err: %+v", dir, err)
	}

	f, err := os.OpenFile(p.path, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return errors.Errorf("create file (%s), err: %+v", p.path, err)
	}
	defer f.Close()

	if err := f.Truncate(0); err != nil {
		return errors.Errorf("truncate file (%s), err: %+v", p.path, err)
	}

	formatted, err := format.Source([]byte(content))
	if err == nil {
		content = string(formatted)
	}

	if _, err := f.WriteString(content); err != nil {
		return errors.Errorf("write content into file (%s), err: %+v", p.path, err)
	}

	return nil
}

type File struct {
	Nodes []*FileNode
}

func (f *File) Append(val string, t NodeType) {
	name := ""
	methodReceiver := ""
	switch t {
	case ntStruct, ntFunc:
		spans := strings.Split(strings.TrimSpace(val), " ")
		if len(spans) >= 2 {
			name = strings.Split(spans[1], "(")[0]
		}
	case ntMethod:
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
	Type           NodeType /*
		0=other
		1=struct
		2=method
		3=function
	*/
}

type NodeType int

const (
	ntOther  NodeType = 0
	ntStruct NodeType = 1
	ntMethod NodeType = 2
	ntFunc   NodeType = 3
)
