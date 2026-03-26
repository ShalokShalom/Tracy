import gleam/option.{type Option, None, Some}
import gleeunit/should
import options.{find_int, first_some, flat_map_opt, greet, map_opt}

pub fn find_existing_element_test() {
  find_int([1, 2, 3, 4], 3) |> should.equal(Some(3))
}

pub fn find_first_element_test() {
  find_int([7, 8, 9], 7) |> should.equal(Some(7))
}

pub fn find_last_element_test() {
  find_int([7, 8, 9], 9) |> should.equal(Some(9))
}

pub fn find_missing_element_test() {
  find_int([1, 2, 3], 99) |> should.equal(None)
}

pub fn find_empty_list_test() {
  find_int([], 1) |> should.equal(None)
}

pub fn greet_with_name_test() {
  greet(Some("Bob")) |> should.equal("Hello, Bob")
}

pub fn greet_without_name_test() {
  greet(None) |> should.equal("Hello, stranger")
}

pub fn first_some_prefers_first_test() {
  first_some(Some("a"), Some("b")) |> should.equal(Some("a"))
}

pub fn first_some_falls_back_to_second_test() {
  first_some(None, Some("b")) |> should.equal(Some("b"))
}

pub fn first_some_both_none_test() {
  let expected: Option(String) = None
  first_some(None, None) |> should.equal(expected)
}

pub fn map_opt_some_test() {
  map_opt(Some(5), fn(x) { x * 2 }) |> should.equal(Some(10))
}

pub fn map_opt_none_test() {
  map_opt(None, fn(x) { x * 2 }) |> should.equal(None)
}

pub fn flat_map_opt_some_to_some_test() {
  flat_map_opt(Some(4), fn(x) { Some(x * 10) }) |> should.equal(Some(40))
}

pub fn flat_map_opt_some_to_none_test() {
  flat_map_opt(Some(0), fn(_) { None }) |> should.equal(None)
}

pub fn flat_map_opt_none_short_circuits_test() {
  flat_map_opt(None, fn(x) { Some(x * 10) }) |> should.equal(None)
}
