# Tracy

A transpiler from Go to Gleam — built in phases, verified by tests.

## Idea

Go and Gleam are semantically very different languages.

Go is imperative, mutable, and nil-based. 
Gleam is functional, immutable, and has no nil. 

Despite this, a surprisingly large portion of Go code can
be translated automatically to Gleam.

This repository tries to achieve that. 

The transpiler itself is written in Go, and the tests obviously in Gleam.

- `src/` — The **pattern library**: documented Gleam target patterns,
  one file per translation category. These files are the ground truth for the
  code generator. Whatever is written here is what the transpiler must produce.

- `test/` — The **acceptance tests**: explicit `_test` functions using
  `should.equal`, one per pattern. `gleam test` is the acceptance criterion
  for each phase. Green tests mean the phase is complete.

   `cmd/` The Go-side analysis code, for SSA, points-to, code generator, etc lives there.

## Phases

**Phase 1 — Primitive types & structs** (`records.gleam`) — fully automatic  
Go structs become custom types with a single variant. Pointer-receiver
mutation becomes a pure function returning a record update.

**Phase 2 — nil → Option(T)** (`options.gleam`) — automatic with analysis  
Every nullable Go type (`*T`, interface, slice, map, chan, func) becomes
`Option(T)`. nil-checks become case expressions. Requires SSA + points-to
analysis to identify which values can be nil.

**Phase 3 — (val, err) → Result(T, E)** (`results.gleam`) — fully automatic  
The semantically closest translation. Go's multi-return `(val, error)`
maps directly to `Result(val, error)`. Error chains (`if err != nil {
return err }`) become `use` expressions with `result.try`.

**Phase 4 — Closed interfaces → variants** (`interfaces.gleam`) — automatic via callgraph  
Interfaces with a known, finite set of implementations become custom types
with one variant per implementation. Resolved via callgraph type analysis.

**Phase 5 — for range → list.*** (`loops.gleam`) — fully automatic  
Simple range loops over slices map cleanly to `list.map`, `list.filter`,
and `list.fold`.

**Phase 6 — Simple defer** (`defer.gleam`) — automatic for simple cases  
Resource cleanup and mutex unlock patterns are translated by placing the
deferred call at the end of the block. Complex defer (loops, value capture)
requires manual work.

**Phase 7 — Shared mutable state → Actor** (`actors.gleam`) — skeleton automatic  
Mutex-protected structs become OTP actors. The transpiler generates the
actor skeleton with correct message types; semantic correctness requires
manual review.

**Phase 8 — Goroutines + channels → OTP** (`concurrency.gleam`) — skeleton automatic  
Channel types and message shapes are inferred from SSA. The transpiler
generates a compilable actor skeleton. `select` on multiple channels and
`context.Context` patterns require manual redesign.

**Phase 9 — select, complex defer, panic/recover** (`manual.gleam`) — TODO markers only  
These constructs have no direct Gleam equivalent. The transpiler emits
compilable stubs with precise `// TODO(transpiler):` comments explaining
what manual redesign is needed and which Gleam pattern to use.

Phases 1–5 are fully automatic and cover roughly 65–70% of typical Go code.
Phases 6–8 produce compilable Gleam skeletons. Phase 9 marks everything
that requires a human.

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

