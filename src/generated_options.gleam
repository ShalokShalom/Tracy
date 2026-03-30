// Generated from Go — Phase 2: nil → Option(T)

import gleam/option.{type Option, None, Some}

pub type User {
  User(name: String)
}

pub fn greet_user(user: Option(User)) -> String {
  case user {
    Some(value) -> "Hello, " <> value.name
    None -> "Hello, stranger"
  }
}

pub fn first_user(a: Option(User), b: Option(User)) -> Option(User) {
  case a {
    Some(_) -> a
    None -> b
  }
}

pub fn increment_maybe_age(age: Option(Int)) -> Option(Int) {
  option.map(age, fn(value) { value + 1 })
}

pub fn find_int(_items: List(Int), _target: Int) -> Option(Int) {
  // Go: returns nil or &value (condition involves loop/complex SSA)
  todo
}

pub fn get_user(id: Int) -> Option(User) {
  case id <= 0 {
    True -> None
    False -> todo // Some(User(...))
  }
}

