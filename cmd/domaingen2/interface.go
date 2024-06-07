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
	Name                   string
	CompareResultUpperCase bool
	upperCase              *ImplementationContent
	lowerCase              *ImplementationContent
}

type ImplementationContent struct {
	Implementation string
	Me             map[string]*Method
}

func NewImplementation(pkg Package, receiver string, in *Interface) *Implementation {
	line := strings.Split(in.Interface, "\n")
	line = line[1 : len(line)-1]

	implementNameUpper := firstUpperCase(in.InterfaceName)
	implementNameLower := firstLowerCase(in.InterfaceName)
	receiverLower := firstLowerCase(receiver)

	meLower := make(map[string]*Method, len(line))
	meUpper := make(map[string]*Method, len(line))
	for _, l := range line {
		l = strings.TrimSpace(l)
		if len(l) == 0 || l[0] == '/' {
			continue
		}
		meName := strings.Split(l, "(")[0]
		meLower[meName] = &Method{
			MethodName: meName,
			Method:     fmt.Sprintf(_methodTemplate, receiverLower, implementNameLower, l),
		}
		meUpper[meName] = &Method{
			MethodName: meName,
			Method:     fmt.Sprintf(_methodTemplate, receiverLower, implementNameUpper, l),
		}
	}

	template := _usecaseTemplate
	if pkg == _repository {
		template = _repositoryTemplate
	}

	return &Implementation{
		Name: implementNameLower,
		lowerCase: &ImplementationContent{
			Implementation: fmt.Sprintf(template,
				implementNameUpper,
				implementNameLower,
				implementNameUpper, implementNameUpper, implementNameUpper,
				implementNameLower,
			),
			Me: meLower,
		},
		upperCase: &ImplementationContent{
			Implementation: fmt.Sprintf(template,
				implementNameUpper,
				implementNameUpper,
				implementNameUpper, implementNameUpper, implementNameUpper,
				implementNameUpper,
			),
			Me: meUpper,
		},
	}
}

func (im *Implementation) Content() *ImplementationContent {
	if im.CompareResultUpperCase {
		return im.upperCase
	}

	return im.lowerCase
}

func (im *Implementation) EqualName(text string) bool {
	return strings.EqualFold(text, im.Name)
}

func (im *Implementation) SetCompareResultCharacterCase(text string) {
	if !im.EqualName(text) {
		return
	}

	if text[0] >= 'A' && text[0] <= 'Z' {
		im.CompareResultUpperCase = true
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
