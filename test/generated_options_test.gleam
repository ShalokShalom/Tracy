import gleam/option.{type Option, None, Some}
import gleeunit/should
import generated_options

// --- greet_user: nil-check pattern (GreetUser) ---

pub fn greet_user_with_some_test() {
  let user = Some(generated_options.User(name: "Alice"))
  generated_options.greet_user(user)
  |> should.equal("Hello, Alice")
}

pub fn greet_user_with_none_test() {
  generated_options.greet_user(None)
  |> should.equal("Hello, stranger")
}

pub fn greet_user_preserves_name_test() {
  let user = Some(generated_options.User(name: "Bob"))
  generated_options.greet_user(user)
  |> should.equal("Hello, Bob")
}

// --- first_user: coalesce pattern (FirstUser) ---

pub fn first_user_both_some_test() {
  let a = Some(generated_options.User(name: "Alice"))
  let b = Some(generated_options.User(name: "Bob"))
  generated_options.first_user(a, b)
  |> should.equal(Some(generated_options.User(name: "Alice")))
}

pub fn first_user_first_none_test() {
  let b = Some(generated_options.User(name: "Bob"))
  generated_options.first_user(None, b)
  |> should.equal(Some(generated_options.User(name: "Bob")))
}

pub fn first_user_both_none_test() {
  let expected: Option(generated_options.User) = None
  generated_options.first_user(None, None)
  |> should.equal(expected)
}

pub fn first_user_second_none_test() {
  let a = Some(generated_options.User(name: "Alice"))
  generated_options.first_user(a, None)
  |> should.equal(Some(generated_options.User(name: "Alice")))
}

// --- increment_maybe_age: map pattern (IncrementMaybeAge) ---

pub fn increment_maybe_age_some_test() {
  generated_options.increment_maybe_age(Some(25))
  |> should.equal(Some(26))
}

pub fn increment_maybe_age_none_test() {
  generated_options.increment_maybe_age(None)
  |> should.equal(None)
}

pub fn increment_maybe_age_zero_test() {
  generated_options.increment_maybe_age(Some(0))
  |> should.equal(Some(1))
}

pub fn increment_maybe_age_negative_test() {
  generated_options.increment_maybe_age(Some(-1))
  |> should.equal(Some(0))
}

// --- get_user: nil-return pattern (GetUser) ---
// Note: get_user's Some branch uses todo (struct construction from SSA
// not fully resolved), so we can only test the None branch.

pub fn get_user_invalid_id_returns_none_test() {
  generated_options.get_user(0)
  |> should.equal(None)
}

pub fn get_user_negative_id_returns_none_test() {
  generated_options.get_user(-5)
  |> should.equal(None)
}

// --- User type ---

pub fn user_type_construction_test() {
  let u = generated_options.User(name: "Test")
  u.name |> should.equal("Test")
}
