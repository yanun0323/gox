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
	switch field {
	case "CreateTime", "CreatedTime", "UpdateTime", "UpdatedTime", "StartTime", "StartedTime", "EndTime", "EndedTime":
		return true
	default:
		return false
	}
}

const (
	_methodTemplate = `
	func (%s *%s) %s() *%s { /* generate by ` + _commandName + ` */
		return &%s{
			%s
		}
	}
`
)

func (st *Structure) GenMethod(pkg, methodName string) *Method {
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
			receiver, st.StructName, methodName, pkg+"."+st.StructName,
			pkg+"."+st.StructName,
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
