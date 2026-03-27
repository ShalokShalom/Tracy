// cmd/tracy/main.go – Go → IR → Gleam end‑to‑end

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"codeberg.org/shalokshalom/Tracy/internal/analyze"
	"codeberg.org/shalokshalom/Tracy/internal/codegen"
	"codeberg.org/shalokshalom/Tracy/internal/ir"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <file-or-package-path>", os.Args[0])
	}
	targetPath := os.Args[1]

	base := filepath.Base(targetPath)
	ext := filepath.Ext(targetPath)

	switch {
	case ext == ".go" || (ext == "" && !strings.HasSuffix(base, ".gleam")):
		analyzeGo(targetPath)
	case ext == ".gleam":
		analyzeGleam(targetPath)
	default:
		log.Fatalf("Unsupported file type: %s", base)
	}
}

// analyzeGo analyzes Go code and emits Gleam modules.

func analyzeGo(targetPath string) {
	// 1. Load, parse, type-check
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedSyntax,
		Dir: ".",
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

	// 3. Analyze packages
	for _, p := range ssaPkgs {
		if p == nil {
			continue
		}
		fmt.Printf("Analyzing Package: %s\n", p.Pkg.Path())

		// Build pkgMembers for structs (Phase 1).
		pkgMembers := map[string]ssa.Member{}
		for name, member := range p.Members {
			if typeMember, ok := member.(*ssa.Type); ok {
				fmt.Printf(" Struct: %s\n", name)
				fmt.Printf("  Type: %s\n", typeMember.Type().Underlying())
				pkgMembers[name] = typeMember
			}
		}

		// Build one IR module per package.
		var mod *ir.Module // ← corrected here
		var firstError error

		for name, member := range p.Members {
			if fn, ok := member.(*ssa.Function); ok {
				fmt.Printf(" Function: %s\n", name)

				irMod, err := analyze.AnalyzeFunc(fn, pkgMembers)
				if err != nil {
					if firstError == nil {
						firstError = err
					}
					log.Printf("Skipped %s.%s: %v", p.Pkg.Name(), name, err)
					continue
				}

				if irMod == nil {
					continue
				}

				if mod == nil {
					mod = irMod
				} else {
					mod.RecordTypes = append(mod.RecordTypes, irMod.RecordTypes...)
					mod.Funcs = append(mod.Funcs, irMod.Funcs...)
				}
			}
		}

		if mod == nil {
			if firstError != nil {
				log.Printf("No usable IR for package %s: %v", p.Pkg.Path(), firstError)
			}
			continue
		}

		// 4. Emit Gleam
		pkgName := p.Pkg.Name()
		gleamPath := fmt.Sprintf("src/generated_%s.gleam", pkgName)

		if err := codegen.EmitGleam(mod, gleamPath); err != nil {
			log.Printf("Codegen for %s failed: %v", pkgName, err)
			continue
		}

		fmt.Printf(" Generated: %s\n", gleamPath)
	}
}

// analyzeGleam is a stub for future Gleam‑side analysis.

func analyzeGleam(path string) {
	fmt.Printf("Analyzing Gleam: %s\n", path)
	// irMod, err := analyze.AnalyzeGleam(path)
	// if err != nil {
	// 	log.Fatalf("Gleam analysis failed: %v", err)
	// }
	// if irMod != nil {
	// 	fmt.Printf("  IR: %+v\n", irMod)
	// }
}