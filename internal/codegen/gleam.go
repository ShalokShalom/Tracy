package codegen

import (
	"fmt"
	"go/types"
	"os"
	"strings"
	"unicode"

	"codeberg.org/shalokshalom/Tracy/internal/ir"
)

func EmitGleam(mod *ir.Module, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprint(f, "// Generated from Go records\n\n")

	// Deduplicate record types by name.
	seen := map[string]bool{}
	var dedupedRecords []ir.RecordType
	for _, rt := range mod.RecordTypes {
		if !seen[rt.Name] {
			seen[rt.Name] = true
			dedupedRecords = append(dedupedRecords, rt)
		}
	}
	mod.RecordTypes = dedupedRecords

	// Build a map of name → *ir.RecordType for codegen.
	recordTypes := map[string]*ir.RecordType{}
	for i := range mod.RecordTypes {
		rt := &mod.RecordTypes[i]
		recordTypes[rt.Name] = rt
	}

	// Emit record types as Gleam custom types (PascalCase required).
	for _, rt := range mod.RecordTypes {
		name := ensurePascalCase(rt.Name)
		fmt.Fprintf(f, "pub type %s {\n", name)
		fmt.Fprintf(f, "  %s(", name)
		for i, field := range rt.Fields {
			sep := ""
			if i > 0 {
				sep = ", "
			}
			gleamType := goTypeToGleam(field.GoType)
			fmt.Fprintf(f, "%s%s: %s", sep, lowercaseFirst(field.Name), gleamType)
		}
		fmt.Fprintf(f, ")\n}\n\n")
	}

	// Emit functions
	for _, fn := range mod.Funcs {
		if len(fn.Updates) == 0 {
			continue
		}

		self := fn.Target
		rt := recordTypes[fn.Return]
		if rt == nil {
			fmt.Fprintf(f, "pub fn %s(%s) {\n", lowercaseFirst(fn.Name), self)
			fmt.Fprintf(f, "  todo(\"cannot emit record constructor for unknown type: %s\")\n", fn.Return)
			fmt.Fprint(f, "}\n\n")
			continue
		}

		retName := ensurePascalCase(fn.Return)

		// Build function signature with extra params.
		var paramParts []string
		paramParts = append(paramParts, fmt.Sprintf("%s: %s", self, retName))
		for _, p := range fn.Params {
			paramParts = append(paramParts, fmt.Sprintf("%s: %s", p.Name, goTypeToGleam(p.Typ)))
		}

		fmt.Fprintf(f, "pub fn %s(%s) -> %s {\n",
			lowercaseFirst(fn.Name), strings.Join(paramParts, ", "), retName)

		// Determine which fields are updated.
		updatedFields := map[string]ir.Expr{}
		for _, u := range fn.Updates {
			updatedFields[u.Field] = u.Value
		}

		// If all fields are updated, emit a full constructor.
		// Otherwise, use record update syntax: Type(..self, field: expr)
		if len(updatedFields) == len(rt.Fields) {
			var fields []string
			for _, field := range rt.Fields {
				fname := lowercaseFirst(field.Name)
				if expr, ok := updatedFields[field.Name]; ok {
					fields = append(fields, fmt.Sprintf("%s: %s", fname, emitExpr(expr)))
				} else {
					fields = append(fields, fmt.Sprintf("%s: %s.%s", fname, self, fname))
				}
			}
			fmt.Fprintf(f, "  %s(%s)\n", retName, strings.Join(fields, ", "))
		} else {
			var updateParts []string
			for _, u := range fn.Updates {
				fname := lowercaseFirst(u.Field)
				updateParts = append(updateParts, fmt.Sprintf("%s: %s", fname, emitExpr(u.Value)))
			}
			fmt.Fprintf(f, "  %s(..%s, %s)\n", retName, self, strings.Join(updateParts, ", "))
		}

		fmt.Fprint(f, "}\n\n")
	}

	return nil
}

// emitExpr recursively renders an IR expression as Gleam code.
func emitExpr(e ir.Expr) string {
	switch v := e.(type) {
	case ir.Var:
		return v.Name
	case ir.LitInt:
		return fmt.Sprintf("%d", int64(v))
	case ir.BinOp:
		op := gleamOp(v.Op, v.IsFloat)
		return fmt.Sprintf("%s %s %s", emitExpr(v.Left), op, emitExpr(v.Right))
	case ir.RecordField:
		return fmt.Sprintf("%s.%s", emitExpr(v.Record), lowercaseFirst(v.Field))
	case ir.RecordValue:
		var parts []string
		for name, expr := range v.Fields {
			parts = append(parts, fmt.Sprintf("%s: %s", lowercaseFirst(name), emitExpr(expr)))
		}
		if v.Type != nil {
			return fmt.Sprintf("%s(%s)", ensurePascalCase(v.Type.Name), strings.Join(parts, ", "))
		}
		return fmt.Sprintf("{ %s }", strings.Join(parts, ", "))
	case ir.RecordUpdate:
		return fmt.Sprintf("..%s { %s: %s }", emitExpr(v.Self), lowercaseFirst(v.Field), emitExpr(v.Value))
	default:
		return e.String()
	}
}

// goTypeToGleam maps Go types to Gleam types using exact type inspection.
func goTypeToGleam(t types.Type) string {
	if t == nil {
		return "String"
	}

	// Use go/types Basic kind for exact matching.
	if basic, ok := t.(*types.Basic); ok {
		switch basic.Kind() {
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
			return "Int"
		case types.Float32, types.Float64:
			return "Float"
		case types.String:
			return "String"
		case types.Bool:
			return "Bool"
		default:
			return "String"
		}
	}

	// Named types — check underlying
	if named, ok := t.(*types.Named); ok {
		return goTypeToGleam(named.Underlying())
	}

	return "String"
}

// ensurePascalCase ensures the first letter is uppercase (Gleam requires PascalCase
// for type names and constructors).
func ensurePascalCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// gleamOp translates a Go arithmetic operator to the Gleam equivalent.
// Gleam uses +. -. *. /. for float operations.
func gleamOp(op string, isFloat bool) string {
	if !isFloat {
		return op
	}
	switch op {
	case "+":
		return "+."
	case "-":
		return "-."
	case "*":
		return "*."
	case "/":
		return "/."
	default:
		return op
	}
}

func lowercaseFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}
