package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"

    "codeberg.org/shalokshalom/Tracy/internal/analyze"
    "codeberg.org/shalokshalom/Tracy/internal/codegen"
    "golang.org/x/tools/go/packages"
    "golang.org/x/tools/go/ssa"
    "golang.org/x/tools/go/ssa/ssautil"
)

func main() {
    if len(os.Args) != 2 {
        log.Fatalf("Usage: %s <file/package-path>", os.Args[0])
    }
    targetPath := os.Args[1]

    base := filepath.Base(targetPath)

    if strings.HasSuffix(base, ".go") || filepath.Ext(targetPath) == "" {
        analyzeGo(targetPath)
    } else if strings.HasSuffix(base, ".gleam") {
        analyzeGleam(targetPath)
    } else {
        log.Fatalf("Unsupported file type: %s", base)
    }
}

func analyzeGo(targetPath string) {
    // 1. Load, parse, type-check
    cfg := &packages.Config{
        Mode: packages.NeedName |
            packages.NeedFiles |
            packages.NeedCompiledGoFiles |
            packages.NeedTypes |
            packages.NeedTypesInfo |
            packages.NeedSyntax,
    }
    initial, err := packages.Load(cfg, targetPath)
    if err != nil {
        log.Fatalf("Failed to load packages: %v", err)
    }
    if packages.PrintErrors(initial) > 0 {
        log.Fatalf("Packages contain errors")
    }

    // 2. Build SSA
    prog, ssaPkgs := ssautil.Packages(initial, 0)
    prog.Build()

    // 3. Analyze
    for _, p := range ssaPkgs {
        if p == nil {
            continue
        }
        fmt.Printf("Analyzing Package: %s\n", p.Pkg.Name())

        for name, member := range p.Members {
            if typeMember, ok := member.(*ssa.Type); ok {
                fmt.Printf(" Found Struct: %s\n", name)
                fmt.Printf("  Type: %s\n", typeMember.Type().Underlying().String())
            }
        }

        for name, member := range p.Members {
            if fn, ok := member.(*ssa.Function); ok {
                fmt.Printf(" Analyzing Function: %s\n", name)
                irMod, _ := analyze.AnalyzeFunc(fn, p.Members)
                if irMod != nil {
                    fmt.Printf("  IR: %+v\n", irMod)
                    
                    // Generate Gleam
                    gleamPath := fmt.Sprintf("src/generated_%s.gleam", p.Pkg.Name())
                    if err := codegen.EmitGleam(irMod, gleamPath); err == nil {
                        fmt.Printf("  Generated: %s\n", gleamPath)
                    }
                }
            }
        }
    }
}

func analyzeGleam(path string) {
    fmt.Printf("Analyzing Gleam: %s\n", path)
    irMod, err := analyze.AnalyzeGleam(path)
    if err != nil {
        log.Fatalf("Gleam analysis failed: %v", err)
    }
    if irMod != nil {
        fmt.Printf("  IR: %+v\n", irMod)
    }
}