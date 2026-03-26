package analyze

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"codeberg.org/shalokshalom/Tracy/internal/ir"
)

type GleamMember struct {
	Name string
	Type *GleamType
}

type GleamType struct {
	Name   string
	Fields []GleamField
}

type GleamField struct {
	Name string
	Type string
}

// ParseGleamFile extracts types/fns
func ParseGleamFile(path string) (map[string]GleamMember, error) {
	members := make(map[string]GleamMember)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentStruct string
	reStruct := regexp.MustCompile(`pub\s+type\s+(\w+)\s*{`)
	reFn := regexp.MustCompile(`pub\s+fn\s+(\w+)\s*\(`)
	reField := regexp.MustCompile(`^\s*(\w+)\s*:\s*\w+`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Struct
		if match := reStruct.FindStringSubmatch(line); match != nil {
			currentStruct = match[1]
			members[currentStruct] = GleamMember{Name: currentStruct, Type: &GleamType{Name: currentStruct}}
			continue
		}
		// End struct
		if currentStruct != "" && strings.Contains(line, "}") {
			currentStruct = ""
			continue
		}
		// Fields
		if currentStruct != "" {
			if match := reField.FindStringSubmatch(line); match != nil {
				if mem, ok := members[currentStruct]; ok && mem.Type != nil {
					mem.Type.Fields = append(mem.Type.Fields, GleamField{Name: match[1]})
				}
			}
		}

		// Fn
		if match := reFn.FindStringSubmatch(line); match != nil {
			fnName := match[1]
			members[fnName] = GleamMember{Name: fnName}
		}
	}

	return members, scanner.Err()
}

// AnalyzeGleam builds IR directly
func AnalyzeGleam(path string) (*ir.Module, error) {
	members, err := ParseGleamFile(path)
	if err != nil {
		return nil, err
	}

	m := &ir.Module{Name: "gleam_module"}

	// RecordTypes from Gleam types
	for _, mem := range members {
		if mem.Type != nil {
			rt := ir.RecordType{Name: mem.Name}
			for _, f := range mem.Type.Fields {
				rt.Fields = append(rt.Fields, ir.Field{Name: f.Name})
			}
			m.RecordTypes = append(m.RecordTypes, rt)
			fmt.Printf(" Found Gleam RecordType: %s\n", mem.Name)
		}
	}

	// TODO: Fn analysis + updates

	return m, nil
}
