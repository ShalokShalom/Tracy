import gleeunit/should
import records.{Config, Person, Point, birthday, move, with_debug, with_port}

pub fn point_creation_test() {
  let p = Point(x: 1.0, y: 2.0)
  p.x |> should.equal(1.0)
  p.y |> should.equal(2.0)
}

pub fn move_returns_new_record_test() {
  let p = Point(x: 0.0, y: 0.0)
  move(p, 3.0, 4.0) |> should.equal(Point(x: 3.0, y: 4.0))
}

pub fn move_does_not_mutate_original_test() {
  let p = Point(x: 1.0, y: 1.0)
  let _ = move(p, 5.0, 5.0)
  p |> should.equal(Point(x: 1.0, y: 1.0))
}

pub fn move_negative_delta_test() {
  let p = Point(x: 10.0, y: 10.0)
  move(p, -3.0, -4.0) |> should.equal(Point(x: 7.0, y: 6.0))
}

pub fn birthday_increments_age_test() {
  let alice = Person(name: "Alice", age: 30)
  birthday(alice) |> should.equal(Person(name: "Alice", age: 31))
}

pub fn birthday_preserves_name_test() {
  let bob = Person(name: "Bob", age: 25)
  birthday(bob).name |> should.equal("Bob")
}

pub fn birthday_chained_test() {
  Person(name: "Eve", age: 20)
  |> birthday
  |> birthday
  |> birthday
  |> should.equal(Person(name: "Eve", age: 23))
}

pub fn config_with_debug_test() {
  let c = Config(host: "localhost", port: 8080, debug: False)
  with_debug(c)
  |> should.equal(Config(host: "localhost", port: 8080, debug: True))
}

pub fn config_with_port_test() {
  let c = Config(host: "localhost", port: 8080, debug: False)
  with_port(c, 9090)
  |> should.equal(Config(host: "localhost", port: 9090, debug: False))
}

pub fn config_update_preserves_other_fields_test() {
  let c = Config(host: "example.com", port: 443, debug: True)
  let updated = with_port(c, 80)
  updated.host |> should.equal("example.com")
  updated.debug |> should.equal(True)
}
