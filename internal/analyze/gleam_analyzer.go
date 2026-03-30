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

// ParseGleamFile extracts types/fns from a Gleam source file.
func ParseGleamFile(path string) (map[string]GleamMember, error) {
	members := make(map[string]GleamMember)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentStruct string
	var braceDepth int

	reStruct := regexp.MustCompile(`^pub\s+type\s+(\w+)\s*\{`)
	reFn := regexp.MustCompile(`^pub\s+fn\s+(\w+)\s*\(`)
	reField := regexp.MustCompile(`^\s*(\w+)\s*:\s*(.+?)\s*,?\s*$`)

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Type definition start
		if currentStruct == "" {
			if match := reStruct.FindStringSubmatch(trimmed); match != nil {
				currentStruct = match[1]
				braceDepth = 1
				members[currentStruct] = GleamMember{
					Name: currentStruct,
					Type: &GleamType{Name: currentStruct},
				}
				continue
			}
		}

		// Inside a type definition — track brace depth
		if currentStruct != "" {
			for _, ch := range trimmed {
				if ch == '{' {
					braceDepth++
				} else if ch == '}' {
					braceDepth--
				}
			}

			if braceDepth <= 0 {
				currentStruct = ""
				braceDepth = 0
				continue
			}

			// Parse fields
			if match := reField.FindStringSubmatch(trimmed); match != nil {
				if mem, ok := members[currentStruct]; ok && mem.Type != nil {
					mem.Type.Fields = append(mem.Type.Fields, GleamField{
						Name: match[1],
						Type: strings.TrimSpace(match[2]),
					})
					members[currentStruct] = mem
				}
			}
			continue
		}

		// Function definition
		if match := reFn.FindStringSubmatch(trimmed); match != nil {
			fnName := match[1]
			members[fnName] = GleamMember{Name: fnName}
		}
	}

	return members, scanner.Err()
}

// AnalyzeGleam builds IR directly from a Gleam source file.
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
