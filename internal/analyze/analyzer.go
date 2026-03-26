package analyze

import (
	"fmt"

	"codeberg.org/shalokshalom/Tracy/internal/ir"
	"golang.org/x/tools/go/ssa"
	"go/types"
)

func AnalyzeFunc(fn *ssa.Function, pkgMembers map[string]ssa.Member) (*ir.Module, error) {
	if fn == nil || fn.Blocks == nil {
		return nil, nil
	}

	m := &ir.Module{Name: fn.Pkg.Pkg.Name()}

	// Scan package members for structs -> RecordTypes
	for name, member := range pkgMembers {
		if typeVal, ok := member.(*ssa.Type); ok {
			typ := typeVal.Type() // func() types.Type for checks
			if strct, ok := typ.Underlying().(*types.Struct); ok {
				rt := ir.RecordType{Name: name}
				for i := 0; i < strct.NumFields(); i++ {
					fld := strct.Field(i)
					rt.Fields = append(rt.Fields, ir.Field{
						Name:   fld.Name(),
						GoType: *typeVal, // ssa.Type (pointer to value)
					})
				}
				m.RecordTypes = append(m.RecordTypes, rt)
				fmt.Printf(" Found RecordType: %s\n", name)
			}
		}
	}

	var recordUpdate *ir.RecordUpdate

	// Look for record update patterns
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			if store, ok := instr.(*ssa.Store); ok {
				if fieldAddr, ok := store.Addr.(*ssa.FieldAddr); ok {
					targetValue := fieldAddr.X
					if targetValue == nil {
						continue
					}
					targetName := targetValue.Name()
					if targetName == "" {
						continue
					}

					fieldIndex := fieldAddr.Field
					fieldName := getFieldName(targetValue.Type(), fieldIndex)
					storeValStr := analyzeValue(store.Val)

					fmt.Printf(" Found record update: %s.%s = %s\n", targetName, fieldName, storeValStr)

					recordUpdate = &ir.RecordUpdate{
						Field: fieldName,
						Value: parseExpr(storeValStr),
					}
					break
				}
			}
			if recordUpdate != nil {
				break
			}
		}
		if recordUpdate != nil {
			break
		}
	}

	if recordUpdate != nil {
		paramName := "p"
		if len(fn.Params) > 0 {
			param := fn.Params[0]
			if param.Name() != "" {
				paramName = param.Name()
			}
		}

		returnType := "Record"
		sig := fn.Signature
		if sig != nil {
			results := sig.Results()
			if results != nil && results.Len() > 0 {
				retVar := results.At(0)
				retType := retVar.Type()
				if named, ok := retType.Underlying().(*types.Named); ok {
					returnType = named.Obj().Name()
				} else {
					returnType = retType.String()
				}
			}
		}

		m.Funcs = append(m.Funcs, ir.Func{
			Name:   fn.Name(),
			Target: paramName,
			Update: *recordUpdate,
			Return: returnType,
		})
	}

	return m, nil
}

func analyzeValue(val ssa.Value) string {
	if val == nil {
		return "0"
	}
	switch v := val.(type) {
	case *ssa.Const:
		return fmt.Sprintf("%v", v.Value)
	case *ssa.Parameter, *ssa.FreeVar:
		return v.Name()
	case *ssa.BinOp:
		return fmt.Sprintf("(%s %s %s)", analyzeValue(v.X), v.Op.String(), analyzeValue(v.Y))
	default:
		if namer, ok := val.(interface{ Name() string }); ok {
			return namer.Name()
		}
		return "unknown"
	}
}

func parseExpr(s string) ir.Expr {
	return ir.Var{Name: s}
}

func getFieldName(typ types.Type, fieldIndex int) string {
	if ptr, ok := typ.Underlying().(*types.Pointer); ok {
		typ = ptr.Elem().Underlying()
	}
	if strct, ok := typ.(*types.Struct); ok {
		if fieldIndex < strct.NumFields() {
			return strct.Field(fieldIndex).Name()
		}
	}
	return fmt.Sprintf("field%d", fieldIndex)
}