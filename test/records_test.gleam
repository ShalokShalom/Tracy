import gleeunit/should
import records.{
  type Address, type Employee, type Person,
  Address, Employee, Person,
  birthday, get_name, make_person, move_to, rename,
}

pub fn make_person_test() {
  make_person("Alice", 30)
  |> should.equal(Person(name: "Alice", age: 30))
}

pub fn get_name_test() {
  Person(name: "Alice", age: 30)
  |> get_name()
  |> should.equal("Alice")
}

pub fn birthday_test() {
  Person(name: "Alice", age: 30)
  |> birthday()
  |> should.equal(Person(name: "Alice", age: 31))
}

pub fn birthday_preserves_name_test() {
  Person(name: "Alice", age: 30)
  |> birthday()
  |> get_name()
  |> should.equal("Alice")
}

pub fn rename_test() {
  Person(name: "Alice", age: 30)
  |> rename("Bob")
  |> should.equal(Person(name: "Bob", age: 30))
}

pub fn rename_preserves_age_test() {
  Person(name: "Alice", age: 30)
  |> rename("Bob")
  |> fn(p) { p.age }
  |> should.equal(30)
}

pub fn move_to_test() {
  let emp =
    Employee(
      person: Person(name: "Alice", age: 30),
      address: Address(city: "Vienna"),
      id: 1,
    )
  emp
  |> move_to("Graz")
  |> fn(e) { e.address.city }
  |> should.equal("Graz")
}

pub fn move_to_preserves_person_test() {
  let emp =
    Employee(
      person: Person(name: "Alice", age: 30),
      address: Address(city: "Vienna"),
      id: 1,
    )
  emp
  |> move_to("Graz")
  |> fn(e) { e.person.name }
  |> should.equal("Alice")
}

pub fn person_fields_test() {
  let p: Person = Person(name: "Alice", age: 30)
  p.age
  |> should.equal(30)
}

pub fn address_fields_test() {
  let a: Address = Address(city: "Vienna")
  a.city
  |> should.equal("Vienna")
}

pub fn employee_fields_test() {
  let e: Employee =
  Employee(person: Person(name: "Alice", age: 30), address: Address(city: "Vienna"), id: 42)
  e.id
  |> should.equal(42)
}
