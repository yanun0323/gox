package main

import (
	"fmt"
	"strings"
)

type Interface struct {
	InterfaceName string
	Interface     string
}

type Implementation struct {
	ImplementationName string
	Implementation     string
	Me                 map[string]*Method
}

func NewImplementation(pkg Package, receiver string, in *Interface) *Implementation {
	line := strings.Split(in.Interface, "\n")
	line = line[1 : len(line)-1]

	domain := in.InterfaceName
	implementName := firstLowerCase(domain)

	me := make(map[string]*Method, len(line))
	for _, l := range line {
		l = strings.TrimSpace(l)
		if len(l) == 0 || l[0] == '/' {
			continue
		}
		meName := strings.Split(l, "(")[0]
		me[meName] = &Method{
			MethodName: meName,
			Method:     fmt.Sprintf(_methodTemplate, receiver, implementName, l),
		}
	}

	template := _usecaseTemplate
	if pkg == _repository {
		template = _repositoryTemplate
	}

	return &Implementation{
		ImplementationName: implementName,
		Implementation: fmt.Sprintf(template,
			domain,
			implementName,
			domain, domain, domain,
			implementName,
		),
		Me: me,
	}
}

func (in *Interface) GetStructType() string {
	spans := strings.Split(strings.TrimSpace(strings.Split(in.Interface, "\n")[0]), " ")
	if len(spans) <= 2 {
		return ""
	}
	return strings.TrimSpace(spans[2])
}

type Method struct {
	MethodName string
	Method     string
}

const (
	_methodTemplate = `
	func (%s *%s) %s {
		// TODO: implement me
	}
`

	_usecaseTemplate = `
	type %sParam struct {
		dig.In
	
		Config       *cfg.Config[configs.Config] ` + "`" + `name:"config"` + "`" + `
		Store        *redis.ClusterClient ` + "`" + `name:"redis"` + "`" + `
	}
	
	type %s struct {
		config       *cfg.Config[configs.Config]
		store        *redis.ClusterClient
	}
	
	func New%s(param %sParam) usecase.%s {
		return &%s{
			config:       param.Config,
			store:        param.Store,
		}
	}
`

	_repositoryTemplate = `

	type %sParam struct {
		dig.In
	
		DB *gorm.DB ` + "`" + `name:"dbM"` + "`" + `
	}
	
	type %s struct {
		db *gorm.DB
	}
	
	func New%s(param %sParam) repository.%s {
		return &%s{
			db: param.DB,
		}
	}
	`
)
