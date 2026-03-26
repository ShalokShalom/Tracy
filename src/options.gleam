//// Phase 2 — Go nil → Gleam Option(T)
////
//// Jeder nullable Go-Typ (*T, interface, []T, map, chan, func) wird zu Option(T).
//// nil-Checks (if x != nil) werden zu case-Expressions.

import gleam/option.{type Option, None, Some}

pub type User {
  User(name: String, age: Int)
}

pub type Order {
  Order(user: Option(User), amount: Int)
}

/// Go: func FindUser(id int) *User { if id == 1 { return &User{...} }; return nil }
pub fn find_user(id: Int) -> Option(User) {
  case id {
    1 -> Some(User(name: "Alice", age: 30))
    _ -> None
  }
}

/// Go: if user != nil { fmt.Println(user.Name) }
/// SSA-BinOp-Knoten wird zu case-Expression.
pub fn use_user_name(user: Option(User)) -> String {
  case user {
    Some(u) -> u.name
    None -> "unknown"
  }
}

/// Go: order.User != nil && order.User.Name != ""
/// Verkettete nil-Zugriffe → option.map
pub fn order_user_name(order: Order) -> Option(String) {
  option.map(order.user, fn(u) { u.name })
}

/// Go: if user == nil { user = &User{Name: "Guest"} }
/// nil-Coalescing → option.unwrap
pub fn with_default(user: Option(User)) -> User {
  option.unwrap(user, User(name: "Guest", age: 0))
}