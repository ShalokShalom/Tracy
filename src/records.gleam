//// Phase 1 — Go Structs → Gleam Custom Types
////
//// Ein Go-Struct wird zu einem Custom Type mit einem einzigen Variant.
//// Pointer-Receiver-Mutation wird zu einer reinen Funktion mit Record-Update-Syntax.

// ---------------------------------------------------------------------------
// Typen
// ---------------------------------------------------------------------------

pub type Person {
  Person(name: String, age: Int)
}

pub type Address {
  Address(city: String)
}

pub type Employee {
  Employee(person: Person, address: Address, id: Int)
}

// ---------------------------------------------------------------------------
// Muster
// ---------------------------------------------------------------------------

/// Go: p := Person{Name: "Alice", Age: 30}
pub fn make_person(name: String, age: Int) -> Person {
  Person(name: name, age: age)
}

/// Go: p.Name
pub fn get_name(p: Person) -> String {
  p.name
}

/// Go: func (p *Person) Birthday() { p.Age++ }
/// Pointer-Receiver-Mutation → reines Record-Update.
pub fn birthday(p: Person) -> Person {
  Person(..p, age: p.age + 1)
}

/// Go: func (p *Person) Rename(n string) { p.Name = n }
pub fn rename(p: Person, new_name: String) -> Person {
  Person(..p, name: new_name)
}

/// Go: func (e *Employee) MoveTo(city string) { e.Address.City = city }
/// Verschachtelte Mutation → geschachtelte Record-Updates.
pub fn move_to(e: Employee, city: String) -> Employee {
  Employee(..e, address: Address(city: city))
  // no spread needed on Address
}
