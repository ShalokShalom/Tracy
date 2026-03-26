package analyze

import (
	"fmt"

	"codeberg.org/shalokshalom/Tracy/internal/ir"
	"golang.org/x/tools/go/ssa"
	"go/types"
)

func AnalyzeFunc(fn *ssa.Function) *ir.Module {
	if fn == nil || fn.Blocks == nil {
		return nil
	}

	m := &ir.Module{Name: fn.Pkg.Pkg.Name()}

	var recordUpdate *ir.RecordUpdate

	// Look for record update patterns
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			if store, ok := instr.(*ssa.Store); ok {
				if fieldAddr, ok := store.Addr.(*ssa.FieldAddr); ok {
					targetValue := fieldAddr.X // Already ssa.Value
					if targetValue == nil {
						continue
					}
					targetName := targetValue.Name() // Direct on Value
					if targetName == "" {
						continue
					}

					// Found: Store to {target}.{field}
					fieldIndex := fieldAddr.Field

					// Get actual field name from struct type
					fieldName := getFieldName(targetValue.Type(), fieldIndex)

					fmt.Printf(" Found record update: %s.%s\n", targetName, fieldName)

					recordUpdate = &ir.RecordUpdate{
						Field: fieldName,
						Value: ir.Var{Name: fmt.Sprintf("%s.%s + 1", targetName, fieldName)},
					}
					break // Break instr loop after first match
				}
			}
			if recordUpdate != nil {
				break // Break instr loop
			}
		}
		if recordUpdate != nil {
			break // Break block loop
		}
	}

	if recordUpdate != nil {
		// Get first param name safely
		paramName := "p"
		if len(fn.Params) > 0 {
			param := fn.Params[0]
			if param.Name() != "" {
				paramName = param.Name()
			}
		}

		// Get return type name safely
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

	return m
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