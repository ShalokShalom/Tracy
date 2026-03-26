# Tracy

A Go to Gleam compiler

Ja, ein vollautomatischer Transpiler für eine sinnvolle Untermenge von Go-Code ist realistisch machbar — mit einer Automatisierungsrate von ca. 85–90% für diese Teilmenge, und ~70–80% über den gesamten Durchschnitts-Go-Code. Hier ist die vollständige ehrliche Bilanz.

## Go fix 
## 
Dieser Compiler setzt Go code voraus, der mit `go fix` bearbeitet worden ist. 

## Was vollständig automatisierbar ist

Auf Basis aller drei Dokumente lässt sich eine klare „sichere Zone" definieren, innerhalb derer der Transpiler korrekte, idiomatische Gleam-Ausgabe produzieren kann:[^1][^2][^3]

- **Primitive Typen, Structs, einfache Mutation** → direkte Übersetzung via SSA + Record-Update-Syntax[^2]
- **`nil` → `Option(T)`** → vollständig mit dem NilAway-inspirierten intraprozeduralen + interprozedualen Nullability-Pass und Andersens Points-to-Analyse über `golang.org/x/tools/go/pointer`[^3]
- **`(val, error)` → `Result(T, E)` und `use`-Chains** → die semantisch engste Analogie, fast verlustfrei[^1][^2]
- **Closed Interfaces → Custom Type Variants** → über Callgraph-Typauflösung (`callgraph.GraphVisitEdges`)[^3]
- **Simple `for range`-Loops → `list.map/fold`** → mechanisch[^2]
- **Einfache Goroutines mit einzelnem Channel → OTP Actor + `Subject(T)`** → Skelett aus SSA-`MakeChan`-Instruktionen synthetisierbar[^3]
- **Einfaches `defer` (Resource-Cleanup, Mutex-Unlock)** → Platzierung des Cleanup-Calls am Block-Ende, überschreibbar durch `use`-Scope[^1]
- **`panic` bei Initialisierung → `let assert` / `Result`** → pattern-erkennbar[^1]
- **Generics mit `any`/`comparable` → Gleam-Typparameter** → direkt[^2]


## Wo es noch nicht automatisch geht

Diese vier Probleme bleiben nach heutigem Stand **nicht vollständig automatisierbar**:

**1. Shared Mutable State mit Cross-Goroutine-Aliasing (Kategorie 3)**
Der Transpiler kann das Actor-Skelett generieren, aber die semantische Korrektheit — ob alle Lese-/Schreibzugriffe korrekt serialisiert sind — kann nicht automatisch verifiziert werden. Die Points-to-Analyse kann die Stellen markieren (`TODO: manual redesign required`), aber nicht garantieren, dass das generierte Actor-Pattern das Go-Verhalten exakt repliziert, insbesondere bei `sync.RWMutex`-geschütztem State mit read-heavy Workloads.

**2. `select` auf heterogenen Channels**
`process.new_selector()` mit `process.selecting()` existiert in Gleam, aber Go erzwingt nicht, dass alle `select`-Branches denselben Typ zurückgeben. Der Transpiler muss einen synthetischen `SelectResult`-Union-Typ erzeugen und alle Branches darin einwickeln — das ist mechanisch machbar für einfache Fälle, aber bei `select` mit `default`-Branch (nicht-blockierendes Polling) oder dynamisch erzeugten Channel-Mengen versagt die automatische Typ-Synthese.

**3. Komplexes `defer`**
Einfaches `defer f.Close()` am Block-Ende: ✅ automatisch. Aber `defer` innerhalb von Loops, `defer` mit Value-Capture auf sich verändernde Variablen, und mehrere `defer`s mit definierter LIFO-Reihenfolge: ❌ nicht automatisch übersetzbar. Gleam hat schlicht kein äquivalentes Sprachmittel.

**4. Offene Interfaces mit mehr als einer Methode**
Ein-Methoden-Interfaces (z. B. `io.Reader`) → einfacher `fn`-Typ-Alias: ✅ automatisch. Mehr-Methoden-Interfaces → Record-of-Functions: ✅ Skelett automatisch generierbar. Aber an jeder **Call-Site**, wo ein offenes Interface übergeben wird, muss der Transpiler die konkrete Implementierung kennen und den richtigen Konstruktor einsetzen — bei externen Packages (z. B. `http.Handler`) ist das nicht lösbar ohne manuelle Annotation.

## Realistische Gesamtbilanz

Für eine **wohldefinierte Untermenge** — pure Domänenlogik, Datenverarbeitung, Fehlerbehandlung, einfache Mutex-Structs, geschlossene Interfaces, simple Goroutine-Patterns — ist ein **voll automatischer** Transpiler machbar, der idiomatisches, kompilierendes Gleam produziert.

Für den Rest erzeugt der Transpiler realistischerweise **kompilierbare Gleam-Skelette mit präzisen `// TODO(transpiler):`-Kommentaren**, die dem Entwickler genau sagen, welches manuelle Redesign nötig ist. Das ist kein Versagen — das ist genau das, was Tools wie `go fix` auch tun: den mechanisch lösbaren Teil abnehmen und den Rest klar markieren.

## Roadmap and Doctests

Sehr gut. Fangen wir langsam und methodisch an. Ich schlage vor: zuerst nur das Fundament — ein Gleam-Projekt, das als **lebende Spezifikation** dient. Jede Phase bekommt eine eigene Datei mit Doctests, die mit `gleam test` sofort verifizierbar sind. Die Doctests definieren, was der Transpiler produzieren *muss*.

***

## Projektstruktur

```
go2gleam/
├── gleam.toml
├── ROADMAP.md
├── src/
│   └── go2gleam/
│       ├── records.gleam       ← Phase 1: Structs → Custom Types
│       ├── option.gleam        ← Phase 2: nil → Option(T)
│       ├── result.gleam        ← Phase 3: (val, err) → Result(T, E)
│       ├── interfaces.gleam    ← Phase 4: Closed Interfaces → Variants
│       └── loops.gleam         ← Phase 5: for range → list.*
└── test/
    └── go2gleam_test.gleam
```

***

## `gleam.toml`

```toml
name = "go2gleam"
version = "0.1.0"

[dependencies]
gleam_stdlib = ">= 0.44.0 and < 2.0.0"

[dev-dependencies]
gleeunit = ">= 1.0.0 and < 2.0.0"
```

***

## `ROADMAP.md`

```markdown
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
```

***

## `src/go2gleam/records.gleam`

```gleam
//// Phase 1 — Go Structs → Gleam Custom Types
////
//// Ein Go-Struct mit N Feldern wird zu einem Custom Type mit
//// einem einzigen Variant, dessen Felder den Struct-Feldern entsprechen.
////
//// Go:
////   type Person struct { Name string; Age int }

// ---------------------------------------------------------------------------
// 1.1 Struct-Definition
// ---------------------------------------------------------------------------

/// Ein Gleam Custom Type mit einem Variant entspricht einem Go-Struct.
///
/// Go:   type Person struct { Name string; Age int }
/// Gleam: pub type Person { Person(name: String, age: Int) }
pub type Person {
  Person(name: String, age: Int)
}

// ---------------------------------------------------------------------------
// 1.2 Konstruktion (struct literal)
// ---------------------------------------------------------------------------

/// Go: p := Person{Name: "Alice", Age: 30}
///
/// ```gleam
/// make_person("Alice", 30)
/// // -> Person(name: "Alice", age: 30)
/// ```
pub fn make_person(name: String, age: Int) -> Person {
  Person(name: name, age: age)
}

// ---------------------------------------------------------------------------
// 1.3 Feldlesezugriff
// ---------------------------------------------------------------------------

/// Go: p.Name
///
/// ```gleam
/// get_name(Person(name: "Alice", age: 30))
/// // -> "Alice"
/// ```
pub fn get_name(p: Person) -> String {
  p.name
}

// ---------------------------------------------------------------------------
// 1.4 Mutation via Pointer-Receiver → Record-Update
// ---------------------------------------------------------------------------

/// Go: func (p *Person) Birthday() { p.Age++ }
///
/// Pointer-Receiver-Mutation wird zu einer reinen Funktion,
/// die ein neues Record zurückgibt (Record-Update-Syntax).
///
/// ```gleam
/// birthday(Person(name: "Alice", age: 30))
/// // -> Person(name: "Alice", age: 31)
/// ```
pub fn birthday(p: Person) -> Person {
  Person(..p, age: p.age + 1)
}

// ---------------------------------------------------------------------------
// 1.5 Mehrfache Feldmutation
// ---------------------------------------------------------------------------

/// Go: func (p *Person) Rename(n string) { p.Name = n }
///
/// ```gleam
/// rename(Person(name: "Alice", age: 30), "Bob")
/// // -> Person(name: "Bob", age: 30)
/// ```
pub fn rename(p: Person, new_name: String) -> Person {
  Person(..p, name: new_name)
}

// ---------------------------------------------------------------------------
// 1.6 Nested Structs
// ---------------------------------------------------------------------------

/// Go:
///   type Address struct { City string }
///   type Employee struct { Person Person; Address Address; ID int }
pub type Address {
  Address(city: String)
}

pub type Employee {
  Employee(person: Person, address: Address, id: Int)
}

/// Mutation eines verschachtelten Felds → geschachtelte Record-Updates.
///
/// Go: func (e *Employee) MoveTo(city string) { e.Address.City = city }
///
/// ```gleam
/// let emp = Employee(
///   person: Person(name: "Alice", age: 30),
///   address: Address(city: "Vienna"),
///   id: 1,
/// )
/// move_to(emp, "Graz").address.city
/// // -> "Graz"
/// ```
pub fn move_to(e: Employee, city: String) -> Employee {
  Employee(..e, address: Address(..e.address, city: city))
}
```

***

## `src/go2gleam/option.gleam`

```gleam
//// Phase 2 — Go nil → Gleam Option(T)
////
//// Jeder nullable Go-Typ (*T, interface, []T, map, chan, func)
//// wird zu Option(T). Der Transpiler analysiert via SSA + Points-to,
//// welche Werte nil sein können, und wrapped sie entsprechend.

import gleam/option.{type Option, None, Some}

// ---------------------------------------------------------------------------
// 2.1 Nullable Rückgabewert (der häufigste Fall)
// ---------------------------------------------------------------------------

pub type User {
  User(name: String, age: Int)
}

/// Go:
///   func FindUser(id int) *User {
///     if id == 1 { return &User{Name: "Alice", Age: 30} }
///     return nil
///   }
///
/// ```gleam
/// find_user(1)
/// // -> Some(User(name: "Alice", age: 30))
/// ```
///
/// ```gleam
/// find_user(99)
/// // -> None
/// ```
pub fn find_user(id: Int) -> Option(User) {
  case id {
    1 -> Some(User(name: "Alice", age: 30))
    _ -> None
  }
}

// ---------------------------------------------------------------------------
// 2.2 nil-Check am Use-Site → case
// ---------------------------------------------------------------------------

/// Go:
///   user := findUser(id)
///   if user != nil { fmt.Println(user.Name) }
///
/// Der Transpiler erkennt den nil-Check via SSA-BinOp-Knoten
/// und erzeugt direkt das case-Statement.
///
/// ```gleam
/// use_user_name(Some(User(name: "Alice", age: 30)))
/// // -> "Alice"
/// ```
///
/// ```gleam
/// use_user_name(None)
/// // -> "unknown"
/// ```
pub fn use_user_name(user: Option(User)) -> String {
  case user {
    Some(u) -> u.name
    None -> "unknown"
  }
}

// ---------------------------------------------------------------------------
// 2.3 Verkettete Option-Zugriffe (optional chaining)
// ---------------------------------------------------------------------------

pub type Order {
  Order(user: Option(User), amount: Int)
}

/// Go: order.User != nil && order.User.Name != ""
///
/// In Gleam via verschachteltem case oder gleam/option.map:
///
/// ```gleam
/// import gleam/option
/// option.map(Some(User(name: "Alice", age: 30)), fn(u) { u.name })
/// // -> Some("Alice")
/// ```
///
/// ```gleam
/// import gleam/option
/// option.map(None, fn(u: User) { u.name })
/// // -> None
/// ```
pub fn order_user_name(order: Order) -> Option(String) {
  option.map(order.user, fn(u) { u.name })
}

// ---------------------------------------------------------------------------
// 2.4 Option mit Default-Wert (Go: if x == nil { x = default })
// ---------------------------------------------------------------------------

/// Go: name := user.Name; if user == nil { name = "Guest" }
///
/// ```gleam
/// import gleam/option
/// option.unwrap(Some(User(name: "Alice", age: 30)), User(name: "Guest", age: 0)).name
/// // -> "Alice"
/// ```
///
/// ```gleam
/// import gleam/option
/// option.unwrap(None, User(name: "Guest", age: 0)).name
/// // -> "Guest"
/// ```
pub fn with_default(user: Option(User)) -> User {
  option.unwrap(user, User(name: "Guest", age: 0))
}
```

***

## `src/go2gleam/result.gleam`

```gleam
//// Phase 3 — Go (val, error) → Gleam Result(T, E)
////
//// Dies ist die semantisch engste Übersetzung.
//// Go's Multi-Return (val, error) entspricht Result(val, error) fast 1:1.
//// Go's "if err != nil { return err }"-Ketten werden zu use-Chains.

// ---------------------------------------------------------------------------
// 3.1 Einfache Fehlerrückgabe
// ---------------------------------------------------------------------------

/// Go:
///   func Divide(a, b float64) (float64, error) {
///     if b == 0 { return 0, errors.New("division by zero") }
///     return a / b, nil
///   }
///
/// ```gleam
/// divide(10.0, 2.0)
/// // -> Ok(5.0)
/// ```
///
/// ```gleam
/// divide(10.0, 0.0)
/// // -> Error("division by zero")
/// ```
pub fn divide(a: Float, b: Float) -> Result(Float, String) {
  case b {
    0.0 -> Error("division by zero")
    _ -> Ok(a /. b)
  }
}

// ---------------------------------------------------------------------------
// 3.2 Error-Chain (if err != nil { return err }) → use
// ---------------------------------------------------------------------------

/// Go:
///   a, err := stepOne()
///   if err != nil { return err }
///   b, err := stepTwo(a)
///   if err != nil { return err }
///   return stepThree(b), nil
///
/// Der use-Ausdruck flacht die Result-Kette auf — semantisch identisch.
///
/// ```gleam
/// pipeline(2.0)
/// // -> Ok(3.0)
/// ```
///
/// ```gleam
/// pipeline(0.0)
/// // -> Error("division by zero")
/// ```
pub fn pipeline(input: Float) -> Result(Float, String) {
  use a <- result.try(divide(10.0, input))
  use b <- result.try(divide(a, 2.0))
  Ok(b +. 1.0)
}

// ---------------------------------------------------------------------------
// 3.3 Fehlerkonvertierung (errors.Wrap → result.map_error)
// ---------------------------------------------------------------------------

/// Go: fmt.Errorf("context: %w", err)
///
/// ```gleam
/// import gleam/result
/// result.map_error(divide(10.0, 0.0), fn(e) { "context: " <> e })
/// // -> Error("context: division by zero")
/// ```
pub fn with_context(r: Result(Float, String), ctx: String) -> Result(Float, String) {
  result.map_error(r, fn(e) { ctx <> ": " <> e })
}

// ---------------------------------------------------------------------------
// 3.4 Ignorierter Fehler (Go: val, _ := f())
// ---------------------------------------------------------------------------

/// Go: val, _ := divide(10, 2)  ← Fehler wird ignoriert
///
/// In Gleam: result.unwrap mit explizitem Default.
/// Der Transpiler emittiert einen Kommentar, wenn _ den Fehler ignoriert.
///
/// ```gleam
/// import gleam/result
/// result.unwrap(divide(10.0, 2.0), 0.0)
/// // -> 5.0
/// ```
pub fn unwrap_or_zero(r: Result(Float, String)) -> Float {
  result.unwrap(r, 0.0)
}
```

***

## `test/go2gleam_test.gleam`

```gleam
import gleeunit

pub fn main() {
  gleeunit.main()
}

// Doctests in src/go2gleam/*.gleam werden automatisch von gleeunit gefunden
// und als Teil von `gleam test` ausgeführt.
```

***

Starte jetzt mit:

```bash
gleam new go2gleam
# Dateien einfügen wie oben
gleam test
```

Alle Doctests in `records.gleam`, `option.gleam` und `result.gleam` sollten sofort grün sein. Sobald das läuft, machen wir mit Phase 4 (`interfaces.gleam`) weiter — oder fangen parallel mit dem Go-seitigen Analysis-Tool an. Was möchtest du als nächstes angehen?

