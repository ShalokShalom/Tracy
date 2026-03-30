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

// --- Phase 2: nil → Option ---

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

// AnalyzeFuncPhase2 detects nil-related patterns in Go functions and produces OptionFunc IR.
// It identifies several patterns:
//   - Functions returning *T that return nil in some branches → Option(T)
//   - Functions taking *T params and nil-checking → case expression on Option
//   - Functions taking two *T and returning first non-nil → first_some pattern
//   - Functions mapping *T → *T via nil-check → option.map pattern
func AnalyzeFuncPhase2(fn *ssa.Function, pkgMembers map[string]ssa.Member) []ir.OptionFunc {
	if fn == nil || fn.Blocks == nil {
		return nil
	}

	sig := fn.Signature
	if sig == nil {
		return nil
	}

	// Skip methods (receivers) for now — Phase 2 focuses on free functions
	if sig.Recv() != nil {
		return nil
	}

	var optFuncs []ir.OptionFunc

	params := sig.Params()
	results := sig.Results()

	// Detect: function returns *T (nullable return)
	returnsNullable := results.Len() == 1 && isNullableType(results.At(0).Type())
	// Detect: function has nullable params
	hasNullableParam := false
	for i := 0; i < params.Len(); i++ {
		if isNullableType(params.At(i).Type()) {
			hasNullableParam = true
			break
		}
	}

	if !returnsNullable && !hasNullableParam {
		return nil
	}

	// Try to classify the function pattern
	if optFn := tryNilCheckReturn(fn, sig); optFn != nil {
		optFuncs = append(optFuncs, *optFn)
	} else if optFn := tryNilCoalesce(fn, sig); optFn != nil {
		optFuncs = append(optFuncs, *optFn)
	} else if optFn := tryNilMapFunc(fn, sig); optFn != nil {
		optFuncs = append(optFuncs, *optFn)
	} else if optFn := tryNilReturnFunc(fn, sig); optFn != nil {
		optFuncs = append(optFuncs, *optFn)
	}

	return optFuncs
}

// tryNilReturnFunc detects: func F(...) *T { if cond { return nil }; return &val }
// Maps to: fn f(...) -> Option(T) { case cond { True -> None; False -> Some(val) } }
func tryNilReturnFunc(fn *ssa.Function, sig *types.Signature) *ir.OptionFunc {
	results := sig.Results()
	if results.Len() != 1 {
		return nil
	}
	retType := results.At(0).Type()
	if !isNullableType(retType) {
		return nil
	}

	// Check if function has nil returns (return nil) and non-nil returns (return &x)
	hasNilReturn := false
	hasNonNilReturn := false

	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			ret, ok := instr.(*ssa.Return)
			if !ok || len(ret.Results) == 0 {
				continue
			}
			retVal := ret.Results[0]
			if isNilConst(retVal) {
				hasNilReturn = true
			} else {
				hasNonNilReturn = true
			}
		}
	}

	if !hasNilReturn || !hasNonNilReturn {
		return nil
	}

	// Check: does this function have nullable *pointer* params that are nil-checked?
	// If so, let tryNilMapFunc or tryNilCheckReturn handle it instead.
	// Slices/maps used as collections (not nil-checked) are OK.
	params := sig.Params()
	for i := 0; i < params.Len(); i++ {
		pt := params.At(i).Type()
		if _, isPtr := pt.(*types.Pointer); isPtr {
			if paramIsNilChecked(fn, fn.Params[i]) {
				return nil // let tryNilMapFunc or tryNilCheckReturn handle it
			}
		}
	}

	innerType := innerNullableType(retType)
	gleamInner := goTypeToGleamStr(innerType)

	// Build params
	var irParams []ir.FuncParam
	for i := 0; i < params.Len(); i++ {
		irParams = append(irParams, ir.FuncParam{
			Name: fn.Params[i].Name(),
			Typ:  params.At(i).Type(),
		})
	}

	// Try to extract condition and Some expression from SSA
	condition, someExpr := extractNilReturnBranches(fn)

	return &ir.OptionFunc{
		Name:       fn.Name(),
		Params:     irParams,
		ReturnType: fmt.Sprintf("Option(%s)", gleamInner),
		Pattern:    ir.NilReturnFunc,
		InnerType:  gleamInner,
		Condition:  condition,
		SomeExpr:   someExpr,
	}
}

// tryNilCheckReturn detects: func F(p *T) R { if p == nil { return X } return Y(p) }
// Maps to: fn f(p: Option(T)) -> R { case p { Some(v) -> Y(v); None -> X } }
func tryNilCheckReturn(fn *ssa.Function, sig *types.Signature) *ir.OptionFunc {
	params := sig.Params()
	results := sig.Results()

	if results.Len() != 1 {
		return nil
	}

	// Find a nullable param that is nil-checked
	nilCheckedParam := -1
	for i := 0; i < params.Len(); i++ {
		if isNullableType(params.At(i).Type()) {
			// Check if this param is nil-compared in the SSA
			if paramIsNilChecked(fn, fn.Params[i]) {
				nilCheckedParam = i
				break
			}
		}
	}

	if nilCheckedParam < 0 {
		return nil
	}

	// Must NOT also be a coalesce or map pattern
	// (coalesce: all params nullable and return nullable; map: 1 nullable param, nullable return)
	retType := results.At(0).Type()
	if isNullableType(retType) {
		// Could be coalesce or map — don't handle here
		return nil
	}

	// This is a nil-check-and-return pattern like GreetUser
	paramType := params.At(nilCheckedParam).Type()
	innerType := innerNullableType(paramType)
	gleamInner := goTypeToGleamStr(innerType)
	paramName := fn.Params[nilCheckedParam].Name()

	// Try to extract Some/None bodies from SSA
	someBody, noneBody := extractNilCheckBranches(fn, fn.Params[nilCheckedParam])

	var irParams []ir.FuncParam
	for i := 0; i < params.Len(); i++ {
		irParams = append(irParams, ir.FuncParam{
			Name: fn.Params[i].Name(),
			Typ:  params.At(i).Type(),
		})
	}

	gleamReturn := goTypeToGleamStr(retType)

	return &ir.OptionFunc{
		Name:       fn.Name(),
		Params:     irParams,
		ReturnType: gleamReturn,
		Pattern:    ir.NilCheckReturn,
		ParamName:  paramName,
		SomeBody:   someBody,
		NoneBody:   noneBody,
		InnerType:  gleamInner,
	}
}

// tryNilCoalesce detects: func F(a, b *T) *T { if a != nil { return a } return b }
// Maps to: fn f(a: Option(T), b: Option(T)) -> Option(T) { case a { Some(_) -> a; None -> b } }
func tryNilCoalesce(fn *ssa.Function, sig *types.Signature) *ir.OptionFunc {
	params := sig.Params()
	results := sig.Results()

	if params.Len() != 2 || results.Len() != 1 {
		return nil
	}

	// Both params nullable, return nullable, same type
	if !isNullableType(params.At(0).Type()) || !isNullableType(params.At(1).Type()) {
		return nil
	}
	if !isNullableType(results.At(0).Type()) {
		return nil
	}

	innerType := innerNullableType(results.At(0).Type())
	gleamInner := goTypeToGleamStr(innerType)

	var irParams []ir.FuncParam
	for i := 0; i < params.Len(); i++ {
		irParams = append(irParams, ir.FuncParam{
			Name: fn.Params[i].Name(),
			Typ:  params.At(i).Type(),
		})
	}

	return &ir.OptionFunc{
		Name:       fn.Name(),
		Params:     irParams,
		ReturnType: fmt.Sprintf("Option(%s)", gleamInner),
		Pattern:    ir.NilCoalesce,
		InnerType:  gleamInner,
	}
}

// tryNilMapFunc detects: func F(p *T) *T { if p == nil { return nil } v := *p; ... return &result }
// Maps to: fn f(p: Option(T)) -> Option(T) { option.map(p, fn(v) { ... }) }
func tryNilMapFunc(fn *ssa.Function, sig *types.Signature) *ir.OptionFunc {
	params := sig.Params()
	results := sig.Results()

	if results.Len() != 1 || params.Len() != 1 {
		return nil
	}

	if !isNullableType(params.At(0).Type()) || !isNullableType(results.At(0).Type()) {
		return nil
	}

	// Verify it nil-checks the param and returns nil in that case
	if !paramIsNilChecked(fn, fn.Params[0]) {
		return nil
	}

	innerType := innerNullableType(results.At(0).Type())
	gleamInner := goTypeToGleamStr(innerType)
	paramName := fn.Params[0].Name()

	var irParams []ir.FuncParam
	irParams = append(irParams, ir.FuncParam{
		Name: paramName,
		Typ:  params.At(0).Type(),
	})

	return &ir.OptionFunc{
		Name:       fn.Name(),
		Params:     irParams,
		ReturnType: fmt.Sprintf("Option(%s)", gleamInner),
		Pattern:    ir.NilMapFunc,
		InnerType:  gleamInner,
	}
}

// paramIsNilChecked returns true if the SSA shows a nil-comparison on the given parameter.
func paramIsNilChecked(fn *ssa.Function, param *ssa.Parameter) bool {
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			binOp, ok := instr.(*ssa.BinOp)
			if !ok {
				continue
			}
			if binOp.Op != token.EQL && binOp.Op != token.NEQ {
				continue
			}
			// Check if one operand is the param and the other is nil
			if (binOp.X == param && isNilConst(binOp.Y)) ||
				(binOp.Y == param && isNilConst(binOp.X)) {
				return true
			}
		}
	}
	return false
}

// isNilConst checks if an SSA value is the nil constant.
func isNilConst(v ssa.Value) bool {
	c, ok := v.(*ssa.Const)
	if !ok {
		return false
	}
	return c.Value == nil && c.IsNil()
}

// extractNilCheckBranches tries to extract the Some/None return expressions
// from a nil-check pattern like: if p == nil { return X } return Y
func extractNilCheckBranches(fn *ssa.Function, param *ssa.Parameter) (someBody ir.Expr, noneBody ir.Expr) {
	// Find the If instruction that tests param == nil
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			ifInstr, ok := instr.(*ssa.If)
			if !ok {
				continue
			}
			// The condition should be a BinOp comparing param to nil
			cond, ok := ifInstr.Cond.(*ssa.BinOp)
			if !ok {
				continue
			}
			isParamNilCheck := (cond.X == param && isNilConst(cond.Y)) ||
				(cond.Y == param && isNilConst(cond.X))
			if !isParamNilCheck {
				continue
			}

			// block.Succs[0] = true branch, block.Succs[1] = false branch
			if len(block.Succs) != 2 {
				continue
			}

			trueBranch := block.Succs[0]
			falseBranch := block.Succs[1]

			trueRet := findReturnExpr(trueBranch, fn)
			falseRet := findReturnExpr(falseBranch, fn)

			if cond.Op == token.EQL {
				// if param == nil: true → None body, false → Some body
				noneBody = trueRet
				someBody = falseRet
			} else {
				// if param != nil: true → Some body, false → None body
				someBody = trueRet
				noneBody = falseRet
			}
			return
		}
	}
	return nil, nil
}

// findReturnExpr finds a return value expression in a block (or its successors).
func findReturnExpr(block *ssa.BasicBlock, fn *ssa.Function) ir.Expr {
	aliases := buildParamAliases(fn)
	for _, instr := range block.Instrs {
		ret, ok := instr.(*ssa.Return)
		if !ok || len(ret.Results) == 0 {
			continue
		}
		return ssaValueToExprWith(ret.Results[0], aliases)
	}
	// Check successors (one level deep)
	for _, succ := range block.Succs {
		for _, instr := range succ.Instrs {
			ret, ok := instr.(*ssa.Return)
			if !ok || len(ret.Results) == 0 {
				continue
			}
			return ssaValueToExprWith(ret.Results[0], aliases)
		}
	}
	return ir.Var{Name: "todo"}
}

// extractNilReturnBranches finds the if-condition and non-nil return expression
// in a NilReturnFunc pattern.
func extractNilReturnBranches(fn *ssa.Function) (condition ir.Expr, someExpr ir.Expr) {
	aliases := buildParamAliases(fn)

	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			ifInstr, ok := instr.(*ssa.If)
			if !ok {
				continue
			}

			// The condition should be a comparison
			cond, ok := ifInstr.Cond.(*ssa.BinOp)
			if !ok {
				continue
			}

			if len(block.Succs) != 2 {
				continue
			}

			condExpr := ssaValueToExprWith(cond, aliases)
			condition = condExpr

			trueBranch := block.Succs[0]
			falseBranch := block.Succs[1]

			// Find which branch returns nil and which returns a value
			trueRet := findReturnValue(trueBranch, fn)
			falseRet := findReturnValue(falseBranch, fn)

			if trueRet != nil && isNilConst(trueRet) {
				// True branch returns nil → false branch has the Some value
				someExpr = findNonNilReturnExpr(falseBranch, fn, aliases)
			} else if falseRet != nil && isNilConst(falseRet) {
				// False branch returns nil → true branch has the Some value
				someExpr = findNonNilReturnExpr(trueBranch, fn, aliases)
			}

			if condition != nil {
				return
			}
		}
	}
	return nil, nil
}

// findReturnValue finds the raw SSA return value in a block.
func findReturnValue(block *ssa.BasicBlock, fn *ssa.Function) ssa.Value {
	for _, instr := range block.Instrs {
		ret, ok := instr.(*ssa.Return)
		if !ok || len(ret.Results) == 0 {
			continue
		}
		return ret.Results[0]
	}
	// Check successors
	for _, succ := range block.Succs {
		for _, instr := range succ.Instrs {
			ret, ok := instr.(*ssa.Return)
			if !ok || len(ret.Results) == 0 {
				continue
			}
			return ret.Results[0]
		}
	}
	return nil
}

// findNonNilReturnExpr finds and converts the non-nil return expression to IR.
func findNonNilReturnExpr(block *ssa.BasicBlock, fn *ssa.Function, aliases map[ssa.Value]string) ir.Expr {
	for _, instr := range block.Instrs {
		ret, ok := instr.(*ssa.Return)
		if !ok || len(ret.Results) == 0 {
			continue
		}
		return ssaValueToExprWith(ret.Results[0], aliases)
	}
	for _, succ := range block.Succs {
		for _, instr := range succ.Instrs {
			ret, ok := instr.(*ssa.Return)
			if !ok || len(ret.Results) == 0 {
				continue
			}
			return ssaValueToExprWith(ret.Results[0], aliases)
		}
	}
	return ir.Var{Name: "todo"}
}

// goTypeToGleamStr converts a Go type to a Gleam type string.
func goTypeToGleamStr(t types.Type) string {
	if t == nil {
		return "String"
	}

	// Handle pointer types → unwrap
	if ptr, ok := t.(*types.Pointer); ok {
		return goTypeToGleamStr(ptr.Elem())
	}

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

	if named, ok := t.(*types.Named); ok {
		// For named struct types, return the short type name
		name := named.Obj().Name()
		if _, isStruct := named.Underlying().(*types.Struct); isStruct {
			return name
		}
		return goTypeToGleamStr(named.Underlying())
	}

	if slice, ok := t.(*types.Slice); ok {
		elemType := goTypeToGleamStr(slice.Elem())
		return fmt.Sprintf("List(%s)", elemType)
	}

	return "String"
}
