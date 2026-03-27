// analyzer.go – Go ssa → ir.Module for Phase 1 (records)

package analyze

import (
	"fmt"
	"go/types"
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
						GoType: fld.Type(), // Go `types.Type`
					})
				}
				m.RecordTypes = append(m.RecordTypes, rt)
				fmt.Printf("Found RecordType: %s\n", name)
			}
		}
	}

	// PHASE 1: pointer‑receiver method → record‑update pattern
	if fn.Signature == nil || fn.Signature.Recv() == nil {
		return m, nil
	}

	recv := fn.Signature.Recv()
	if ptr, ok := recv.Type().(*types.Pointer); ok {
		if _, isStruct := ptr.Elem().Underlying().(*types.Struct); isStruct {
			rt := findRecordType(m.RecordTypes, ptr.Elem().String())
			if rt == nil {
				return m, nil
			}

			// Trivial example update: `self.field = self.field + 1`
			update := &ir.RecordUpdate{
				Self: ir.Var{Name: "self"},
				Field: "x",
				Value: ir.BinOp{
					Op: "+",
					Left: ir.RecordField{
						Record: ir.Var{Name: "self"},
						Field:  "x",
					},
					Right: ir.LitInt(1),
				},
			}

			m.Funcs = append(m.Funcs, ir.Func{
				Name:   fn.Name(),
				Target: "self",
				Update: update,
				Return: rt.Name,
			})
		}
	}

	return m, nil
}

func findRecordType(recs []ir.RecordType, typeName string) *ir.RecordType {
	for i := range recs {
		if recs[i].Name == typeName {
			return &recs[i]
		}
	}
	return nil
}

// --- helpers for future phases ---

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