package generator

import (
	"regexp"
	"strings"
)

type Method struct {
	Name    string
	Arg     string
	ArgType string
	Ret     string
}

type Service struct {
	Package string
	Name    string
	Methods []Method
}

type Struct struct {
	Name   string
	Fields map[string]string
}

type DSL struct {
	Package string
	Service Service
	Structs []Struct
}

func ParseDSL(dsl string) DSL {
	// Yeah, like... I could use a proper language parser approach with a lexer, AST and all of that. But... have a
	// bunch of regexes instead
	rePackage := regexp.MustCompile(`^package\s+(\w+)$`)
	reService := regexp.MustCompile(`^service\s+(\w+)\s+{$`)
	reMethod := regexp.MustCompile(`^\s+(\w+)\((\w+):\s*(\w+)\):\s*(\w+)`)
	reStruct := regexp.MustCompile(`^struct\s+(\w+)\s+{$`)
	reField := regexp.MustCompile(`^\s+(\w+):\s*(\w+)`)
	reEnd := regexp.MustCompile(`^}$`)

	lines := strings.Split(dsl, "\n")
	var parsed DSL

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if matches := rePackage.FindStringSubmatch(line); matches != nil {
			parsed.Package = matches[1]
			continue
		}
		if matches := reService.FindStringSubmatch(line); matches != nil {
			service := Service{Name: matches[1]}
			for {
				i++
				if i == len(lines) {
					break
				}
				line := lines[i]
				if reEnd.MatchString(line) {
					break
				}
				if matches := reMethod.FindStringSubmatch(line); matches != nil {
					method := Method{Name: matches[1], Arg: matches[2], ArgType: matches[3], Ret: matches[4]}
					service.Methods = append(service.Methods, method)
				}
			}
			parsed.Service = service
			continue
		}

		if matches := reStruct.FindStringSubmatch(line); matches != nil {
			strct := Struct{Name: matches[1], Fields: make(map[string]string)}
			for {
				i++
				if i == len(lines) {
					break
				}
				line := lines[i]
				if reEnd.MatchString(line) {
					break
				}
				if matches := reField.FindStringSubmatch(line); matches != nil {
					fieldName := matches[1]
					fieldType := matches[2]
					strct.Fields[fieldName] = fieldType
				}
			}
			parsed.Structs = append(parsed.Structs, strct)
			continue
		}
	}

	return parsed
}
