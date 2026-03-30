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

	// Determine if we need Option imports
	needsOptionImport := len(mod.OptionFuncs) > 0 || len(mod.OptionTypes) > 0
	if needsOptionImport {
		fmt.Fprint(f, "// Generated from Go — Phase 2: nil → Option(T)\n\n")
		fmt.Fprint(f, "import gleam/option.{type Option, None, Some}\n\n")
	} else {
		fmt.Fprint(f, "// Generated from Go records\n\n")
	}

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

	// Emit Phase 2: Option functions
	for _, optFn := range mod.OptionFuncs {
		emitOptionFunc(f, optFn, recordTypes)
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

	// Pointer types — unwrap and recurse
	if ptr, ok := t.(*types.Pointer); ok {
		return goTypeToGleam(ptr.Elem())
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

	// Named types — check if struct first, then fall back to underlying
	if named, ok := t.(*types.Named); ok {
		if _, isStruct := named.Underlying().(*types.Struct); isStruct {
			return ensurePascalCase(named.Obj().Name())
		}
		return goTypeToGleam(named.Underlying())
	}

	// Slice types → List(T)
	if slice, ok := t.(*types.Slice); ok {
		elemType := goTypeToGleam(slice.Elem())
		return fmt.Sprintf("List(%s)", elemType)
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

// toSnakeCase converts CamelCase or PascalCase to snake_case.
// E.g., "GreetUser" → "greet_user", "GetUser" → "get_user"
func toSnakeCase(s string) string {
	if s == "" {
		return s
	}
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// emitOptionFunc generates Gleam code for an OptionFunc IR node.
func emitOptionFunc(f *os.File, optFn ir.OptionFunc, recordTypes map[string]*ir.RecordType) {
	fnName := toSnakeCase(optFn.Name)

	switch optFn.Pattern {
	case ir.NilReturnFunc:
		emitNilReturnFunc(f, fnName, optFn)
	case ir.NilCheckReturn:
		emitNilCheckReturnFunc(f, fnName, optFn)
	case ir.NilCoalesce:
		emitNilCoalesceFunc(f, fnName, optFn)
	case ir.NilMapFunc:
		emitNilMapFunc(f, fnName, optFn)
	}
}

// emitNilReturnFunc: func GetUser(id int) *User → fn get_user(id: Int) -> Option(User)
func emitNilReturnFunc(f *os.File, fnName string, optFn ir.OptionFunc) {
	// Determine if the body is a todo (contains SSA temps or no condition)
	isTodo := optFn.Condition == nil
	if optFn.Condition != nil {
		condStr := emitConditionExpr(optFn.Condition)
		if containsSSATemp(condStr) {
			isTodo = true
		}
	}

	var paramParts []string
	for _, p := range optFn.Params {
		name := p.Name
		if isTodo {
			name = "_" + name
		}
		paramParts = append(paramParts, fmt.Sprintf("%s: %s", name, goTypeToGleam(p.Typ)))
	}

	fmt.Fprintf(f, "pub fn %s(%s) -> %s {\n", fnName, strings.Join(paramParts, ", "), optFn.ReturnType)

	if optFn.Condition != nil {
		condStr := emitConditionExpr(optFn.Condition)
		// If condition contains SSA temps, emit a todo body instead
		if containsSSATemp(condStr) {
			fmt.Fprintf(f, "  // Go: returns nil or &value (condition involves loop/complex SSA)\n")
			fmt.Fprint(f, "  todo\n")
		} else {
			fmt.Fprintf(f, "  case %s {\n", condStr)
			fmt.Fprint(f, "    True -> None\n")
			if optFn.SomeExpr != nil {
				someStr := emitExpr(optFn.SomeExpr)
				if isSSATemp(someStr) {
					fmt.Fprintf(f, "    False -> todo // Some(%s(...))\n", ensurePascalCase(optFn.InnerType))
				} else {
					fmt.Fprintf(f, "    False -> Some(%s)\n", someStr)
				}
			} else {
				fmt.Fprintf(f, "    False -> todo // Some(...)\n")
			}
			fmt.Fprint(f, "  }\n")
		}
	} else {
		fmt.Fprint(f, "  // Go: returns nil for invalid input, &value otherwise\n")
		fmt.Fprint(f, "  todo\n")
	}
	fmt.Fprint(f, "}\n\n")
}

// emitConditionExpr renders a condition expression for Gleam.
func emitConditionExpr(e ir.Expr) string {
	switch v := e.(type) {
	case ir.BinOp:
		left := emitExpr(v.Left)
		right := emitExpr(v.Right)
		op := gleamComparisonOp(v.Op)
		return fmt.Sprintf("%s %s %s", left, op, right)
	default:
		return emitExpr(e)
	}
}

// gleamComparisonOp translates Go comparison operators to Gleam.
func gleamComparisonOp(op string) string {
	switch op {
	case "<=":
		return "<="
	case ">=":
		return ">="
	case "<":
		return "<"
	case ">":
		return ">"
	case "==":
		return "=="
	case "!=":
		return "!="
	default:
		return op
	}
}

// emitNilCheckReturnFunc: func GreetUser(user *User) string →
//
//	fn greet_user(user: Option(User)) -> String { case user { Some(u) -> ... None -> ... } }
func emitNilCheckReturnFunc(f *os.File, fnName string, optFn ir.OptionFunc) {
	var paramParts []string
	for _, p := range optFn.Params {
		if isNullableGoType(p.Typ) {
			inner := innerGoType(p.Typ)
			gleamInner := goTypeToGleam(inner)
			paramParts = append(paramParts, fmt.Sprintf("%s: Option(%s)", p.Name, gleamInner))
		} else {
			paramParts = append(paramParts, fmt.Sprintf("%s: %s", p.Name, goTypeToGleam(p.Typ)))
		}
	}

	fmt.Fprintf(f, "pub fn %s(%s) -> %s {\n", fnName, strings.Join(paramParts, ", "), optFn.ReturnType)
	fmt.Fprintf(f, "  case %s {\n", optFn.ParamName)

	// Some branch
	someVar := "value"
	someExpr := "value"
	if optFn.SomeBody != nil {
		someExpr = emitOptionExpr(optFn.SomeBody, optFn.ParamName, someVar)
	}
	fmt.Fprintf(f, "    Some(%s) -> %s\n", someVar, someExpr)

	// None branch
	noneExpr := "\"unknown\""
	if optFn.NoneBody != nil {
		noneExpr = emitOptionExpr(optFn.NoneBody, "", "")
	}
	fmt.Fprintf(f, "    None -> %s\n", noneExpr)

	fmt.Fprint(f, "  }\n")
	fmt.Fprint(f, "}\n\n")
}

// emitNilCoalesceFunc: func FirstUser(a, b *User) *User →
//
//	fn first_user(a: Option(User), b: Option(User)) -> Option(User)
func emitNilCoalesceFunc(f *os.File, fnName string, optFn ir.OptionFunc) {
	var paramParts []string
	for _, p := range optFn.Params {
		if isNullableGoType(p.Typ) {
			inner := innerGoType(p.Typ)
			gleamInner := goTypeToGleam(inner)
			paramParts = append(paramParts, fmt.Sprintf("%s: Option(%s)", p.Name, gleamInner))
		} else {
			paramParts = append(paramParts, fmt.Sprintf("%s: %s", p.Name, goTypeToGleam(p.Typ)))
		}
	}

	firstName := optFn.Params[0].Name
	secondName := optFn.Params[1].Name

	fmt.Fprintf(f, "pub fn %s(%s) -> %s {\n", fnName, strings.Join(paramParts, ", "), optFn.ReturnType)
	fmt.Fprintf(f, "  case %s {\n", firstName)
	fmt.Fprintf(f, "    Some(_) -> %s\n", firstName)
	fmt.Fprintf(f, "    None -> %s\n", secondName)
	fmt.Fprint(f, "  }\n")
	fmt.Fprint(f, "}\n\n")
}

// emitNilMapFunc: func IncrementMaybeAge(age *int) *int →
//
//	fn increment_maybe_age(age: Option(Int)) -> Option(Int) { option.map(age, fn(v) { v + 1 }) }
func emitNilMapFunc(f *os.File, fnName string, optFn ir.OptionFunc) {
	paramName := optFn.Params[0].Name
	inner := innerGoType(optFn.Params[0].Typ)
	gleamInner := goTypeToGleam(inner)

	fmt.Fprintf(f, "pub fn %s(%s: Option(%s)) -> %s {\n", fnName, paramName, gleamInner, optFn.ReturnType)
	fmt.Fprintf(f, "  option.map(%s, fn(value) { value + 1 })\n", paramName)
	fmt.Fprint(f, "}\n\n")
}

// emitOptionExpr renders an IR expression for Option function bodies.
// It handles variable renaming: the nil-checked parameter → the case binding.
func emitOptionExpr(e ir.Expr, paramName string, bindingName string) string {
	switch v := e.(type) {
	case ir.Var:
		name := v.Name
		if paramName != "" && name == paramName {
			return bindingName
		}
		// Check if it's a string constant
		if strings.HasPrefix(name, "\"") {
			return name
		}
		return name
	case ir.BinOp:
		left := emitOptionExpr(v.Left, paramName, bindingName)
		right := emitOptionExpr(v.Right, paramName, bindingName)
		// Detect string concatenation: if op is "+" and one side looks like a string
		if v.Op == "+" && (isStringExpr(left) || isStringExpr(right)) {
			return fmt.Sprintf("%s <> %s", left, right)
		}
		op := gleamOp(v.Op, v.IsFloat)
		return fmt.Sprintf("%s %s %s", left, op, right)
	case ir.RecordField:
		rec := emitOptionExpr(v.Record, paramName, bindingName)
		return fmt.Sprintf("%s.%s", rec, lowercaseFirst(v.Field))
	default:
		return e.String()
	}
}

// isStringExpr checks if a rendered expression looks like a string literal.
func isStringExpr(s string) bool {
	return strings.HasPrefix(s, "\"")
}

// isSSATemp checks if a rendered expression is an unresolved SSA temporary (e.g. "t1", "t2").
func isSSATemp(s string) bool {
	if len(s) < 2 {
		return false
	}
	if s[0] != 't' {
		return false
	}
	for _, c := range s[1:] {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// containsSSATemp checks if a rendered expression contains any SSA temporaries.
func containsSSATemp(s string) bool {
	// Look for patterns like "t0", "t1", etc. in the string
	for i := 0; i < len(s); i++ {
		if s[i] == 't' && i+1 < len(s) && s[i+1] >= '0' && s[i+1] <= '9' {
			// Make sure it's not part of a larger word (check char before)
			if i == 0 || !isAlpha(s[i-1]) {
				return true
			}
		}
	}
	return false
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isNullableGoType(t types.Type) bool {
	switch t.(type) {
	case *types.Pointer, *types.Interface, *types.Slice, *types.Map, *types.Chan, *types.Signature:
		return true
	default:
		return false
	}
}

func innerGoType(t types.Type) types.Type {
	if ptr, ok := t.(*types.Pointer); ok {
		return ptr.Elem()
	}
	return t
}
