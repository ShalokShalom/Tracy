# go2gleam — Transpiler Roadmap

Der Transpiler besteht aus zwei unabhängigen Teilen:

1. **Analysis-Tool** (`cmd/`, Go):
   Liest Go-Quellcode → SSA-Form → Nullability/Points-to-Analyse → Gleam-AST → Quellcode

2. **Pattern-Library** (`src/go2gleam/`, Gleam):
   Dokumentierte, getestete Zielmuster. Das ist die Ground-Truth für den Code-Generator.
   `gleam test` ist das Acceptance-Criterion für jede Phase.

---

## Phasen

| # | Name                          | Datei                 | Automatisierbar | Status |
|---|-------------------------------|-----------------------|-----------------|--------|
| 0 | go fix (Vorverarbeitung)      | –                     | ✅ vollständig   | –      |
| 1 | Primitive Typen & Structs     | records.gleam         | ✅ vollständig   | 🔲     |
| 2 | nil → Option(T)               | option.gleam          | ✅ mit Analyse   | 🔲     |
| 3 | (val, err) → Result(T, E)     | result.gleam          | ✅ vollständig   | 🔲     |
| 4 | Closed Interfaces → Variants  | interfaces.gleam      | ✅ via Callgraph | 🔲     |
| 5 | for range → list.*            | loops.gleam           | ✅ vollständig   | 🔲     |
| 6 | defer (einfach)               | defer.gleam           | ⚠️ einfache Fälle| 🔲     |
| 7 | Shared Mutable State → Actor  | actors.gleam          | ⚠️ Skelett       | 🔲     |
| 8 | Goroutines + Channels → OTP   | concurrency.gleam     | ⚠️ Skelett       | 🔲     |
| 9 | select, defer komplex, panic  | manual.gleam          | ✗ TODO-Marker   | 🔲     |

---

## Acceptance-Criterion

Jede Phase gilt als abgeschlossen, wenn:
- `gleam test` für die zugehörige Datei ohne Fehler durchläuft
- Der Go-seitige Code-Generator denselben Output erzeugt wie die Muster in der Datei
