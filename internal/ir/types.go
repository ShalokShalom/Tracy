// types.go – IR for Go → Gleam phases 1–3 (records, options, results)

package ir

import (
	"fmt"
	"go/types"
	"strings"
)

// Module is the top unit that maps roughly to one Gleam module.

type Module struct {
	Name          string
	RecordTypes   []RecordType
	OptionTypes   []OptionType
	ResultTypes   []ResultType
	Funcs         []Func
	OptionMatches []OptionMatch
}

// Record IR

type RecordType struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name   string
	GoType types.Type
}

// Expression IR

type Expr interface {
	String() string
}

type Var struct {
	Name string
}

func (v Var) String() string { return v.Name }

type BinOp struct {
	Op   string
	Left Expr
	Right Expr
}

func (b BinOp) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left, b.Op, b.Right)
}

type LitInt int64

func (l LitInt) String() string { return fmt.Sprintf("%d", int64(l)) }

type RecordValue struct {
	Type   *RecordType
	Fields map[string]Expr
}

func (r RecordValue) String() string {
	var parts []string
	for name, expr := range r.Fields {
		parts = append(parts, fmt.Sprintf("%s: %s", name, expr))
	}
	return fmt.Sprintf("{ %s }", strings.Join(parts, ", "))
}

type RecordField struct {
	Record Expr
	Field  string
}

func (r RecordField) String() string {
	return fmt.Sprintf("%s.%s", r.Record, r.Field)
}

type RecordUpdate struct {
	Self  Expr
	Field string
	Value Expr
}

func (u RecordUpdate) String() string {
	return fmt.Sprintf("..%s { %s = %s }", u.Self, u.Field, u.Value)
}

// Option / Result IR

type OptionType struct {
	Name  string
	Inner types.Type
}

type OptionSome struct {
	Inner Expr
	Type  *OptionType
}

func (o OptionSome) String() string {
	return fmt.Sprintf("Some(%s)", o.Inner)
}

type OptionNone struct {
	Type *OptionType
}

func (o OptionNone) String() string {
	return "None"
}

type OptionMatch struct {
	Value       Expr
	SomePattern string
	SomeBody    Expr
	OkBody      Expr
	NoneBody    Expr
}

func (m OptionMatch) String() string {
	return fmt.Sprintf("case %s { Some(%s) -> %s; None -> %s }",
		m.Value, m.SomePattern, m.OkBody, m.NoneBody)
}

type ResultType struct {
	Name    string
	OkType  types.Type
	ErrType types.Type
}

type ResultOk struct {
	Value Expr
	Type  *ResultType
}

func (r ResultOk) String() string {
	return fmt.Sprintf("Ok(%s)", r.Value)
}

type ResultError struct {
	Error Expr
	Type  *ResultType
}

func (r ResultError) String() string {
	return fmt.Sprintf("Error(%s)", r.Error)
}

type UseResult struct {
	Value      Expr
	OkVar      string
	OkBody     Expr
	ErrorVar   string
	ErrorBody  Expr
}

func (u UseResult) String() string {
	return fmt.Sprintf("use Ok(%s) -> %s; Error(%s) -> %s",
		u.OkVar, u.OkBody, u.ErrorVar, u.ErrorBody)
}

// Functions

type Func struct {
	Name   string
	Target string
	Update *RecordUpdate
	Return string
}

type ResultFunc struct {
	Name         string
	Params       []FuncParam
	OkExpr       Expr
	ErrExpr      Expr
	HasErrorChain bool
}

type FuncParam struct {
	Name string
	Typ  types.Type
}