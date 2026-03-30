# Tracy

### A transpiler from Go to Gleam

Go and Gleam are semantically very different languages - on the surface.  

**Go** is imperative, mutable, and nil-based, and **Gleam** is functional, immutable, and has no nil. *Right?*

**Simultaneously, both languages also have a lot in common:**  
Simplicity, value-based error handling, enforced formatting, fast compilation, and a small standard library. 

So, a surprisingly large portion of Go code can be automatically translated into Gleam.
And this is what we are trying to do here. ;)

The transpiler itself is written in Go, and the tests in Gleam, obviously. 

We are trying to get compatibility with the BubbleTea ecosystem:  

Both because we love it, and secondly, since we believe that the Elm architecture is a great pattern for Gleam.  
It will also enrich the Erlang ecosystem, and especially so our new implementation of the [Erlang VM (BEAM) in Ada.](https://codeberg.org/ShalokShalom/Linda)

- `src/` - The **pattern library**:  
   This one documents Gleam target patterns, one file per translation category.  
   The compiler conforms to these files, which serve as the reality check for the project.  
   Gleam standards force this directory to be called `src`, counterintuitively. 

- `test/` - The **acceptance tests**:  
  Consists of explicit `_test` functions using `should.equal`, again one per pattern.  
  Run the tests with `gleam test`, which is the acceptance criterion for each phase. 

- `cmd/` - The Go **command-line tool**:  
  The actual compiler, including analysis code for SSA and similar code, lives there.  
  I am still working on a proper command-line interface. Also thinking about a Bubble Tea TUI for it.

## Roadmap

**Phase 1 — Primitive types & structs** (`records.gleam`) — *in progress*

Go structs become custom types with a single variant.
Pointer-receiver mutation becomes a pure function returning a record update.
*Status: The analyzer extracts struct definitions and detects record-update patterns via SSA. Codegen emits valid Gleam types and functions. Not all Go patterns are covered yet.*

**Phase 2 — nil → Option(T)** (`options.gleam`) — *planned*

Every nullable Go type (`*T`, interface, slice, map, chan, func) becomes `Option(T)`.
Nil-checks become case expressions. Requires SSA + points-to analysis to identify which values can be nil.
*Status: IR types and nullable-type helpers exist but are not yet wired into the analysis pipeline.*

**Phase 3 — (val, err) → Result(T, E)** (`results.gleam`)  

The semantically closest translation. Go's multi-return `(val, error)` maps directly to `Result(val, error)`.  
Error chains (`if err != nil {return err }`) become `use` expressions with `result.try`.

**Phase 4 — Closed interfaces → variants** (`interfaces.gleam`)  

Interfaces with a known, finite set of implementations become custom types
with one variant per implementation.  
Resolved via callgraph type analysis.

**Phase 5 — for range → list.*** (`loops.gleam`)  

Simple range loops over slices map cleanly to `list.map`, `list.filter`, and `list.fold`.

**Phase 6 — Simple defer** (`defer.gleam`)  

Resource cleanup and mutex unlock patterns are translated by placing the deferred call at the end of the block.  
Complex defer (loops, value capture) requires manual work.

**Phase 7 — Shared mutable state → Actor** (`actors.gleam`)  

Mutex-protected structs become OTP actors. The transpiler generates the actor skeleton with correct message types.  
Semantic correctness requires manual review. (Can eventually be made by Ada proofs.)

**Phase 8 — Goroutines + channels → OTP** (`concurrency.gleam`)  

Channel types and message shapes are inferred from SSA.  
The transpiler generates a compilable actor skeleton, with `select` on multiple channels and `context.Context` patterns asking for manual work.

**Phase 9 — select, complex defer, panic/recover** (`manual.gleam`)  

These constructs have no direct Gleam equivalent.  
The transpiler emits compilable stubs with precise comments, 
explaining what manual redesign is needed and which Gleam pattern to use.  

**Notes:**  
```
Phases 1 to 5 are automatic and cover roughly 65–70% of typical Go code. Phase 6 solves simple constructs fully automatically.

From here on out, it is currently questionable how to automate everything.

Complex cases of `defer`, and all the constructs of the phases 7 and 8 produce compilable Gleam skeletons, and no working code.

Phase 9 only marks everything that requires a human as a comment.  
```

## Quickstart

```bash
gleam new Tracy
# copy files from src/ and test/
gleam test
```

## Dependencies

```toml
[dependencies]
gleam_stdlib = ">= 0.34.0"
gleam_otp = ">= 1.0.0"
gleam_erlang = ">= 0.25.0" 

[dev-dependencies]
gleeunit = ">= 1.9.0 and < 2.0.0"
```

## Why Gleam?

Among all BEAM languages, Gleam is the strongest migration target for Go:

- Both languages are **statically typed** — no type guessing required
- Gleam's `Result(T, E)` mirrors Go's `(val, error)` almost exactly
- Gleam's OTP actors are a clean equivalent for Go's goroutines and mutexes
- Exhaustive pattern matching makes Go's panic-based assertions unnecessary

## Larger Project

This project is part of an effort to modernize and strengthen the Erlang ecosystem.
Among the other projects is a Rust to Ada transpiler and a complete rewrite of Erlang's VM in Ada. 

