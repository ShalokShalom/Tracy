// analyzer.go – Go SSA → ir.Module for Phase 1 (records)

package analyze

import (
	"fmt"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/ssa"

	"codeberg.org/shalokshalom/Tracy/internal/ir"
)

func AnalyzeFunc(fn *ssa.Function, pkgMembers map[string]ssa.Member) (*ir.Module, error) {
	if fn == nil || fn.Blocks == nil {
		return nil, nil
	}

	m := &ir.Module{Name: fn.Pkg.Pkg.Name()}

	// PHASE 1: Scan package members for structs → RecordTypes
	for name, member := range pkgMembers {
		if typeVal, ok := member.(*ssa.Type); ok {
			typ := typeVal.Type()

			if strct, ok := typ.Underlying().(*types.Struct); ok {
				rt := ir.RecordType{Name: name}
				for i := 0; i < strct.NumFields(); i++ {
					fld := strct.Field(i)
					rt.Fields = append(rt.Fields, ir.Field{
						Name:   fld.Name(),
						GoType: fld.Type(),
					})
				}
				m.RecordTypes = append(m.RecordTypes, rt)
				fmt.Printf("Found RecordType: %s\n", name)
			}
		}
	}

	// PHASE 1: Analyze function for record-update patterns.
	// Handles both pointer-receiver methods and value-receiver functions
	// that return the same struct type.

	sig := fn.Signature
	if sig == nil {
		return m, nil
	}

	// Case 1: Pointer-receiver method (e.g., func (p *Point) Move(...))
	if sig.Recv() != nil {
		if ptr, ok := sig.Recv().Type().(*types.Pointer); ok {
			if _, isStruct := ptr.Elem().Underlying().(*types.Struct); isStruct {
				typeName := shortTypeName(ptr.Elem().String())
				rt := findRecordType(m.RecordTypes, typeName)
				if rt == nil {
					return m, nil
				}

				updates := analyzeFieldUpdates(fn, rt)
				if len(updates) > 0 {
					m.Funcs = append(m.Funcs, ir.Func{
						Name:    fn.Name(),
						Target:  "self",
						Updates: updates,
						Return:  rt.Name,
					})
				}
			}
		}
		return m, nil
	}

	// Case 2: Value-receiver function (e.g., func Move(p Point, dx, dy float64) Point)
	// Detect functions that take a struct as first param and return the same struct type.
	results := sig.Results()
	params := sig.Params()
	if results.Len() == 1 && params.Len() >= 1 {
		retType := results.At(0).Type()
		firstParamType := params.At(0).Type()

		retName := shortTypeName(retType.String())
		paramName := shortTypeName(firstParamType.String())

		if retName == paramName {
			rt := findRecordType(m.RecordTypes, retName)
			if rt != nil {
				updates := analyzeSSABody(fn, rt)
				var extraParams []ir.FuncParam
				for i := 1; i < params.Len(); i++ {
					extraParams = append(extraParams, ir.FuncParam{
						Name: fn.Params[i].Name(),
						Typ:  params.At(i).Type(),
					})
				}

				selfName := fn.Params[0].Name()
				m.Funcs = append(m.Funcs, ir.Func{
					Name:    fn.Name(),
					Target:  selfName,
					Params:  extraParams,
					Updates: updates,
					Return:  rt.Name,
				})
			}
		}
	}

	return m, nil
}

// analyzeFieldUpdates walks SSA blocks for FieldAddr + Store patterns
// used in pointer-receiver methods.
func analyzeFieldUpdates(fn *ssa.Function, rt *ir.RecordType) []ir.RecordUpdate {
	var updates []ir.RecordUpdate
	seen := map[string]bool{}
	aliases := buildParamAliases(fn)

	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			store, ok := instr.(*ssa.Store)
			if !ok {
				continue
			}
			fa, ok := store.Addr.(*ssa.FieldAddr)
			if !ok {
				continue
			}
			strct, ok := fa.X.Type().(*types.Pointer)
			if !ok {
				continue
			}
			st, ok := strct.Elem().Underlying().(*types.Struct)
			if !ok {
				continue
			}
			if fa.Field < 0 || fa.Field >= st.NumFields() {
				continue
			}
			fieldName := st.Field(fa.Field).Name()
			if seen[fieldName] {
				continue
			}
			seen[fieldName] = true

			updates = append(updates, ir.RecordUpdate{
				Self:  ir.Var{Name: "self"},
				Field: fieldName,
				Value: ssaValueToExprWith(store.Val, aliases),
			})
		}
	}
	return updates
}

// analyzeSSABody inspects SSA for value-receiver functions that construct
// a new struct return value. It looks for MakeStruct or individual field
// assignments to build IR updates.
func analyzeSSABody(fn *ssa.Function, rt *ir.RecordType) []ir.RecordUpdate {
	var updates []ir.RecordUpdate

	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			// Look for Return instructions to find struct construction
			ret, ok := instr.(*ssa.Return)
			if !ok || len(ret.Results) == 0 {
				continue
			}
			// Trace back the returned value
			retVal := ret.Results[0]
			updates = traceStructConstruction(retVal, fn, rt)
			if len(updates) > 0 {
				return updates
			}
		}
	}

	// Fallback: walk all blocks looking for BinOp patterns that look like field updates
	updates = analyzeFieldUpdates(fn, rt)
	return updates
}

// traceStructConstruction follows the SSA value chain to extract field expressions
// from struct literal construction.
func traceStructConstruction(val ssa.Value, fn *ssa.Function, rt *ir.RecordType) []ir.RecordUpdate {
	var updates []ir.RecordUpdate

	// If the return value is directly an Alloc + series of FieldAddr stores,
	// extract those.
	switch v := val.(type) {
	case *ssa.Alloc:
		// Find all stores into this alloc's fields
		for _, ref := range *v.Referrers() {
			fa, ok := ref.(*ssa.FieldAddr)
			if !ok {
				continue
			}
			strct, ok := fa.X.Type().(*types.Pointer)
			if !ok {
				continue
			}
			st, ok := strct.Elem().Underlying().(*types.Struct)
			if !ok {
				continue
			}
			if fa.Field < 0 || fa.Field >= st.NumFields() {
				continue
			}
			fieldName := st.Field(fa.Field).Name()
			aliases := buildParamAliases(fn)
			for _, faRef := range *fa.Referrers() {
				store, ok := faRef.(*ssa.Store)
				if !ok {
					continue
				}
				selfName := "self"
				if len(fn.Params) > 0 {
					selfName = fn.Params[0].Name()
				}
				updates = append(updates, ir.RecordUpdate{
					Self:  ir.Var{Name: selfName},
					Field: fieldName,
					Value: ssaValueToExprWith(store.Val, aliases),
				})
			}
		}
	case *ssa.UnOp:
		// Dereference of an alloc
		if v.Op == token.MUL {
			return traceStructConstruction(v.X, fn, rt)
		}
	}

	return updates
}

// buildParamAliases maps SSA Alloc values that are local copies of parameters
// back to the parameter name.  In SSA, value-receiver params are copied into
// a local alloc (e.g. t0 = local Point (p)); this map lets us resolve t0 → p.
func buildParamAliases(fn *ssa.Function) map[ssa.Value]string {
	aliases := map[ssa.Value]string{}
	if fn == nil || len(fn.Blocks) == 0 {
		return aliases
	}
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			store, ok := instr.(*ssa.Store)
			if !ok {
				continue
			}
			alloc, ok := store.Addr.(*ssa.Alloc)
			if !ok {
				continue
			}
			param, ok := store.Val.(*ssa.Parameter)
			if !ok {
				continue
			}
			aliases[alloc] = param.Name()
		}
	}
	return aliases
}

// ssaValueToExpr recursively converts an SSA value into an IR expression.
// aliases maps SSA alloc values to parameter names so that local copies
// of parameters are rendered with the correct name.
func ssaValueToExpr(v ssa.Value) ir.Expr {
	return ssaValueToExprWith(v, nil)
}

func ssaValueToExprWith(v ssa.Value, aliases map[ssa.Value]string) ir.Expr {
	switch val := v.(type) {
	case *ssa.BinOp:
		isFloat := isFloatType(val.X.Type()) || isFloatType(val.Y.Type())
		return ir.BinOp{
			Op:      val.Op.String(),
			Left:    ssaValueToExprWith(val.X, aliases),
			Right:   ssaValueToExprWith(val.Y, aliases),
			IsFloat: isFloat,
		}
	case *ssa.UnOp:
		if val.Op == token.MUL {
			return ssaValueToExprWith(val.X, aliases)
		}
		return ir.Var{Name: val.Name()}
	case *ssa.Const:
		if val.Value != nil {
			return ir.Var{Name: val.Value.ExactString()}
		}
		return ir.Var{Name: "0"}
	case *ssa.Parameter:
		return ir.Var{Name: val.Name()}
	case *ssa.Alloc:
		if aliases != nil {
			if name, ok := aliases[val]; ok {
				return ir.Var{Name: name}
			}
		}
		return ir.Var{Name: val.Name()}
	case *ssa.FieldAddr:
		strct, ok := val.X.Type().(*types.Pointer)
		if ok {
			if st, ok := strct.Elem().Underlying().(*types.Struct); ok {
				if val.Field >= 0 && val.Field < st.NumFields() {
					return ir.RecordField{
						Record: ssaValueToExprWith(val.X, aliases),
						Field:  st.Field(val.Field).Name(),
					}
				}
			}
		}
		return ir.Var{Name: val.Name()}
	case *ssa.Field:
		if st, ok := val.X.Type().Underlying().(*types.Struct); ok {
			if val.Field >= 0 && val.Field < st.NumFields() {
				return ir.RecordField{
					Record: ssaValueToExprWith(val.X, aliases),
					Field:  st.Field(val.Field).Name(),
				}
			}
		}
		return ir.Var{Name: val.Name()}
	default:
		return ir.Var{Name: v.Name()}
	}
}

func findRecordType(recs []ir.RecordType, typeName string) *ir.RecordType {
	for i := range recs {
		if recs[i].Name == typeName || shortTypeName(recs[i].Name) == typeName {
			return &recs[i]
		}
	}
	return nil
}

// shortTypeName extracts the short type name from a fully-qualified Go type string.
// e.g., "codeberg.org/shalokshalom/Tracy/fixtures/records.Point" → "Point"
func shortTypeName(fullName string) string {
	if idx := strings.LastIndex(fullName, "."); idx >= 0 {
		return fullName[idx+1:]
	}
	return fullName
}

// isFloatType checks if a Go type is a floating-point type.
func isFloatType(t types.Type) bool {
	if t == nil {
		return false
	}
	if basic, ok := t.Underlying().(*types.Basic); ok {
		return basic.Kind() == types.Float32 || basic.Kind() == types.Float64
	}
	return false
}

// --- helpers for future phases (Phase 2: nil → Option) ---
// These are not yet wired into the analysis pipeline.
// They will be used when Phase 2 adds nullable-type detection.

func isNullableType(t types.Type) bool {
	switch t.(type) {
	case *types.Pointer, *types.Interface, *types.Slice, *types.Map, *types.Chan, *types.Signature:
		return true
	default:
		return false
	}
}

func innerNullableType(t types.Type) types.Type {
	if ptr, ok := t.(*types.Pointer); ok {
		return ptr.Elem()
	}
	return t
}
