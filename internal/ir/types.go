package ir

import (
    "fmt"
    "golang.org/x/tools/go/ssa"
)

type Module struct {
    Name        string
    RecordTypes []RecordType
    Funcs       []Func
    ResultTypes []ResultType  // ✅ Moved inside Module
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
    Target string
    Update RecordUpdate
    Return string
}

type RecordUpdate struct {
    Field string
    Value Expr
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

// Phase 3 Result types
type ResultType struct {
    Name   string
    OkType ssa.Type
    ErrType ssa.Type
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
    Typ  ssa.Type
}