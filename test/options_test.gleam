import gleam/option.{None, Some}
import gleeunit/should
import options.{
  type Order, type User, Order, User, find_user, order_user_name, use_user_name,
  with_default,
}

pub fn find_user_found_test() {
  find_user(1)
  |> should.equal(Some(User(name: "Alice", age: 30)))
}

pub fn find_user_not_found_test() {
  find_user(99)
  |> should.equal(None)
}

pub fn use_user_name_some_test() {
  Some(User(name: "Alice", age: 30))
  |> use_user_name()
  |> should.equal("Alice")
}

pub fn use_user_name_none_test() {
  None
  |> use_user_name()
  |> should.equal("unknown")
}

pub fn order_user_name_some_test() {
  Order(user: Some(User(name: "Alice", age: 30)), amount: 100)
  |> order_user_name()
  |> should.equal(Some("Alice"))
}

pub fn order_user_name_none_test() {
  Order(user: None, amount: 100)
  |> order_user_name()
  |> should.equal(None)
}

pub fn with_default_some_test() {
  Some(User(name: "Alice", age: 30))
  |> with_default()
  |> should.equal(User(name: "Alice", age: 30))
}

pub fn with_default_none_test() {
  None
  |> with_default()
  |> should.equal(User(name: "Guest", age: 0))
}

pub fn user_fields_test() {
  let u: User = User(name: "Bob", age: 25)
  u.age
  |> should.equal(25)
}

pub fn order_fields_test() {
  let o: Order = Order(user: None, amount: 50)
  o.amount
  |> should.equal(50)
}
