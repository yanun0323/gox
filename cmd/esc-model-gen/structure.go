package main

import (
	"fmt"
	"strings"
)

type Structure struct {
	StructName string
	Struct     string
}

func NewStructureFrom(st *Structure, unix, timestamp, keepTag bool) *Structure {
	lines := strings.Split(st.Struct, "\n")
	rows := lines[1 : len(lines)-1]
	newLines := make([]string, 0, len(lines)+1)

	newLines = append(newLines, lines[0])

	for _, row := range rows {
		spans := strings.Split(strings.TrimSpace(row), " ")
		if len(spans) <= 1 {
			newLines = append(newLines, row)
			continue
		}

		if isTimeField(spans[0]) {
			if unix {
				for i := range spans {
					if spans[i] == "string" {
						spans[i] = "int64"
						break
					}
				}
			}

			if timestamp {
				for i := range spans {
					if spans[i] == "int64" {
						spans[i] = "string"
						break
					}
				}
			}
		}

		if keepTag {
			newLines = append(newLines, strings.Join(spans, " "))
			continue
		}

		/* clean tag */
		newSpans := make([]string, 0, len(spans))
		tag := false
		for _, span := range spans {
			if len(span) != 0 && span[0] == '`' {
				tag = true
			}

			if !tag {
				newSpans = append(newSpans, span)
			}

			if len(span) != 0 && span[len(span)-1] == '`' {
				tag = false
			}
		}
		newLines = append(newLines, strings.Join(newSpans, " "))
	}

	newLines = append(newLines, lines[len(lines)-1])

	return &Structure{
		StructName: st.StructName,
		Struct:     strings.Join(newLines, "\n"),
	}
}

func isTimeField(field string) bool {
	return strings.HasSuffix(field, "Time") || strings.HasSuffix(field, "time")
}

const (
	_methodTemplate = `
	func (%s *%s) %s() *%s {
		return &%s{
			%s
		}
	}
`
	_functionTemplate = `
	func %s(%s *%s) *%s {
		return &%s{
			%s
		}
	}
`
)

func (st *Structure) GetStructType() string {
	spans := strings.Split(strings.TrimSpace(strings.Split(st.Struct, "\n")[0]), " ")
	if len(spans) <= 2 {
		return ""
	}
	return strings.TrimSpace(spans[2])
}

func (st *Structure) GenMethod(pkg Package, methodName string) *Method {
	receiver := "elem"
	fields := st.getFields()
	setters := make([]string, 0, len(fields))
	for _, field := range fields {
		set := field + ":" + receiver + "." + field + ","
		setters = append(setters, set)
	}
	return &Method{
		MethodName: methodName,
		Method: fmt.Sprintf(_methodTemplate,
			receiver, st.StructName, methodName, pkg.String()+"."+st.StructName,
			pkg.String()+"."+st.StructName,
			strings.Join(setters, "\n"),
		),
	}
}

func (st *Structure) GenFunction(pkg Package, functionNameSuffix string) *Function {
	receiver := "elem"
	fields := st.getFields()
	setters := make([]string, 0, len(fields))
	for _, field := range fields {
		set := field + ":" + receiver + "." + field + ","
		setters = append(setters, set)
	}

	functionName := fmt.Sprintf("New%s%s", st.StructName, functionNameSuffix)
	return &Function{
		FunctionName: functionName,
		Function: fmt.Sprintf(_functionTemplate,
			functionName, receiver, pkg.String()+"."+st.StructName, st.StructName,
			st.StructName,
			strings.Join(setters, "\n"),
		),
	}
}

func (st *Structure) getFields() []string {
	rows := strings.Split(st.Struct, "\n")
	rows = rows[1 : len(rows)-1]
	fields := make([]string, 0, len(rows))

	for _, row := range rows {
		/* get field  */
		field := strings.Split(strings.TrimSpace(row), " ")[0]
		if len(field) == 0 {
			continue
		}

		/* check embed struct */
		spans := strings.Split(field, ".")
		field = spans[len(spans)-1]
		fields = append(fields, field)
	}
	return fields
}
