package ir

import (
	"fmt"
	"golang.org/x/tools/go/ssa"
)

type Module struct {
	Name       string
	RecordTypes []RecordType
	Funcs      []Func
}

type RecordType struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name   string
	GoType ssa.Type
}

type Func struct {
	Name   string
	Target string      // param name "p"
	Update RecordUpdate
	Return string      // "Person"
}

type RecordUpdate struct {
	Field string // "Age"
	Value Expr   // "p.Age + 1"
}

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