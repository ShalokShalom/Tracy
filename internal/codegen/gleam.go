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

	// Build a map of name → *ir.RecordType for codegen.
	recordTypes := map[string]*ir.RecordType{}
	for i := range mod.RecordTypes {
		rt := &mod.RecordTypes[i]
		recordTypes[strings.ToLower(rt.Name)] = rt
	}

	// Emit record types as Gleam variant types.
	for _, rt := range mod.RecordTypes {
		name := initialLower(rt.Name)
		fmt.Fprintf(f, "pub type %s {\n", name)
		fmt.Fprintf(f, "  %s(", name)
		for i, field := range rt.Fields {
			sep := ""
			if i > 0 {
				sep = ", "
			}
			gleamType := goTypeToGleam(field.GoType)
			fmt.Fprintf(f, "%s%s: %s", sep, initialLower(field.Name), gleamType)
		}
		fmt.Fprintf(f, ")\n}\n\n")
	}

	// Emit functions
	for _, fn := range mod.Funcs {
		if fn.Update == nil {
			continue
		}

		// Use self‑style receiver, lowercase name.
		self := fn.Target
		retName := strings.ToLower(fn.Return)

		// Find the record type this function returns.
		rt := recordTypes[retName]
		if rt == nil {
			// Can’t emit a proper constructor; skip or fallback.
			fmt.Fprintf(f, "pub fn %s(%s: any) -> any {\n", lowercaseFirst(fn.Name), self)
			fmt.Fprintf(f, "  todo(\"cannot emit record constructor for unknown type: %s\")\n", fn.Return)
			fmt.Fprint(f, "}\n\n")
			continue
		}

		// Emit function header.
		fmt.Fprintf(f, "pub fn %s(%s: %s) -> %s {\n",
			lowercaseFirst(fn.Name), self, retName, retName)

		// Emit constructor body:
		// Person(name: p.name, age: p.age + 1)
		var fields []string
		for _, field := range rt.Fields {
			fname := initialLower(field.Name)
			if field.Name == fn.Update.Field {
				// Updated field: p.field + 1
				fields = append(fields, fmt.Sprintf("%s = %s.%s + 1", fname, self, fname))
			} else {
				// Unchanged field.
				fields = append(fields, fmt.Sprintf("%s = %s.%s", fname, self, fname))
			}
		}

		fmt.Fprintf(f, "  %s(%s)\n", retName, strings.Join(fields, ", "))
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

func initialLower(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = rune(strings.ToLower(string(runes[0]))[0])
	return string(runes)
}

func lowercaseFirst(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}