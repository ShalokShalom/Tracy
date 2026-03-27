// package codegen

package codegen

import (
	"fmt"
	"os"
	"strings"

	"codeberg.org/shalokshalom/Tracy/internal/ir"
	"go/types"
)

func EmitGleam(mod *ir.Module, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprint(f, "// Generated from Go records\n\n")

	// Emit record types (deduplicate)
	seen := make(map[string]bool)
	for _, rt := range mod.RecordTypes {
		if seen[rt.Name] {
			continue
		}
		seen[rt.Name] = true

		fmt.Fprintf(f, "pub type %s {\n", rt.Name)
		for _, field := range rt.Fields {
			gleamType := goTypeToGleam(field.GoType)
			fmt.Fprintf(f, "  %s %s\n", field.Name, gleamType)
		}
		fmt.Fprint(f, "}\n\n")
	}

	// Emit functions
	for _, fn := range mod.Funcs {
		if fn.Update == nil {
			continue
		}
		fmt.Fprintf(f, "pub fn %s(%s: %s) -> %s {\n",
			fn.Name, fn.Target, fn.Return, fn.Return)
		fmt.Fprintf(f, "  %s\n", fn.Update)
		fmt.Fprint(f, "}\n\n")
	}

	return nil
}

func goTypeToGleam(t types.Type) string {
	typeName := t.String()
	typeName = strings.ToLower(typeName)

	if strings.Contains(typeName, "int") {
		return "Int"
	}
	if strings.Contains(typeName, "float") {
		return "Float"
	}
	if strings.Contains(typeName, "string") {
		return "String"
	}
	if strings.Contains(typeName, "bool") {
		return "Bool"
	}

	return "String"
}