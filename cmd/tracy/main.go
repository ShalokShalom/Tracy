package main

import (
	"fmt"
	"log"
	"os"

	"codeberg.org/shalokshalom/Tracy/internal/analyze"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <package-path>", os.Args[0])
	}
	targetPath := os.Args[1]

	// 1. Load, parse, and type-check
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
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

	// 3. Analyze each package
	for _, p := range ssaPkgs {
		if p == nil {
			continue
		}

		fmt.Printf("Analyzing Package: %s\n", p.Pkg.Name())
		
		// Print structs
		for name, member := range p.Members {
			if typeMember, ok := member.(*ssa.Type); ok {
				fmt.Printf("  Found Struct: %s\n", name)
				fmt.Printf("    Type: %s\n", typeMember.Type().Underlying().String())
			}
		}

		// Analyze functions for record update patterns
		for name, member := range p.Members {
			if fn, ok := member.(*ssa.Function); ok {
				fmt.Printf("\n  Analyzing Function: %s\n", name)
				
				irMod := analyze.AnalyzeFunc(fn)
				if irMod != nil {
					fmt.Printf("    IR: %+v\n", irMod)
				}
			}
		}
	}
}