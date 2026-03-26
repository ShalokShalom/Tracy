import gleam/list
import gleam/option.{type Option, None, Some}

pub fn find_int(items: List(Int), target: Int) -> Option(Int) {
  list.find(items, fn(x) { x == target })
  |> option.from_result
}

pub fn greet(name: Option(String)) -> String {
  case name {
    Some(n) -> "Hello, " <> n
    None -> "Hello, stranger"
  }
}

pub fn first_some(a: Option(a), b: Option(a)) -> Option(a) {
  case a {
    Some(_) -> a
    None -> b
  }
}

pub fn map_opt(opt: Option(Int), f: fn(Int) -> Int) -> Option(Int) {
  option.map(opt, f)
}

pub fn flat_map_opt(opt: Option(Int), f: fn(Int) -> Option(Int)) -> Option(Int) {
  option.then(opt, f)
}
