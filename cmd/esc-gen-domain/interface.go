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

func NewImplementation(pkg, receiver string, inter *Interface) *Implementation {
	line := strings.Split(inter.Interface, "\n")
	line = line[1 : len(line)-1]

	domain := inter.InterfaceName
	implement := firstLowerCase(domain)

	me := make(map[string]*Method, len(line))
	for _, l := range line {
		l = strings.TrimSpace(l)
		if len(l) == 0 || l[0] == '/' {
			continue
		}
		meName := strings.Split(l, "(")[0]
		me[meName] = &Method{
			MethodName: meName,
			Method:     fmt.Sprintf(_methodTemplate, receiver, implement, l),
		}
	}

	return &Implementation{
		ImplementationName: implement,
		Implementation: fmt.Sprintf(_implementationTemplate,
			domain,
			implement,
			domain, domain, pkg, domain,
			implement,
		),
		Me: me,
	}
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

	_implementationTemplate = `
	type %sParam struct {
		dig.In
	
		Config       *cfg.Config[configs.Config] ` + "`" + `name:"config"` + "`" + `
		Store        *redis.ClusterClient ` + "`" + `name:"redis"` + "`" + `
	}
	
	type %s struct {
		config       *cfg.Config[configs.Config]
		store        *redis.ClusterClient
	}
	
	func New%s(param %sParam) %s.%s {
		return &%s{
			config:       param.Config,
			store:        param.Store,
		}
	}
`
)
