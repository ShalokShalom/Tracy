# Tracy

### A transpiler from Go to Gleam — specifically focused on the Bubbletea ecosystem.

Go and Gleam are semantically very different languages; **Go** is imperative, mutable, and nil-based. 
**Gleam** is functional, immutable, and has no nil. 

Despite this, a surprisingly large portion of Go code can be automatically translated into Gleam.

This repository tries to achieve that.😎😎

The transpiler itself is written in Go, and the tests obviously in Gleam.

- `src/` - The **pattern library**: documented Gleam target patterns,
  one file per translation category. These files here act as the reality check for the compiler.

- `test/` - The **acceptance tests**: explicit `_test` functions using
  `should.equal`, one per pattern. `gleam test` is the acceptance criterion
  for each phase. 

- `cmd/` - The Go **command-line tool**, including analysis code for SSA, and compiler lives there.

## Roadmap

**Phase 1 — Primitive types & structs** (`records.gleam`)  
Go structs become custom types with a single variant.  
Pointer-receiver mutation becomes a pure function returning a record update.
<br/><br/>

**Phase 2 — nil → Option(T)** (`options.gleam`)  
Every nullable Go type (`*T`, interface, slice, map, chan, func) becomes `Option(T)`.  
Nil-checks become case expressions. Requires SSA + points-to analysis to identify which values can be nil.
<br/><br/>

**Phase 3 — (val, err) → Result(T, E)** (`results.gleam`)  
The semantically closest translation. Go's multi-return `(val, error)` maps directly to `Result(val, error)`.  
Error chains (`if err != nil {return err }`) become `use` expressions with `result.try`.
<br/><br/>

**Phase 4 — Closed interfaces → variants** (`interfaces.gleam`)  
Interfaces with a known, finite set of implementations become custom types
with one variant per implementation.  
Resolved via callgraph type analysis.
<br/><br/>

**Phase 5 — for range → list.*** (`loops.gleam`)  
Simple range loops over slices map cleanly to `list.map`, `list.filter`, and `list.fold`.
<br/><br/>

**Phase 6 — Simple defer** (`defer.gleam`)  
Resource cleanup and mutex unlock patterns are translated by placing the deferred call at the end of the block.  
Complex defer (loops, value capture) requires manual work.
<br/><br/>

**Phase 7 — Shared mutable state → Actor** (`actors.gleam`)  
Mutex-protected structs become OTP actors. The transpiler generates the actor skeleton with correct message types.  
Semantic correctness requires manual review. (Can eventually be made by Ada proofs.)
<br/><br/>

**Phase 8 — Goroutines + channels → OTP** (`concurrency.gleam`)  
Channel types and message shapes are inferred from SSA.  
The transpiler generates a compilable actor skeleton, with `select` on multiple channels and `context.Context` patterns asking for manual work.
<br/><br/>

**Phase 9 — select, complex defer, panic/recover** (`manual.gleam`)  
These constructs have no direct Gleam equivalent.  
The transpiler emits compilable stubs with precise comments, 
explaining what manual redesign is needed and which Gleam pattern to use.  

**Notes:**  
```
Phases 1 to 5 are fully automatic and cover roughly 65–70% of typical Go code, and phase 6 solves simple constructs fully automatically.

From here on out, it is currently questionable how to automate everything.

Complex cases of `defer`, and all the constructs of the phases 7 and 8, which produce compilable Gleam skeletons, and no working code.

Phase 9 marks everything that requires a human as a comment.  
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

