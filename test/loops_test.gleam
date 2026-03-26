import gleeunit/should
import loops.{
  find_first, flat_transform, keep, map_indexed, reduce, sum, transform,
}

pub fn transform_doubles_test() {
  transform([1, 2, 3], fn(x) { x * 2 }) |> should.equal([2, 4, 6])
}

pub fn transform_empty_test() {
  transform([], fn(x) { x * 2 }) |> should.equal([])
}

pub fn transform_increments_test() {
  transform([10, 20], fn(x) { x + 1 }) |> should.equal([11, 21])
}

pub fn keep_even_test() {
  keep([1, 2, 3, 4, 5, 6], fn(x) { x % 2 == 0 }) |> should.equal([2, 4, 6])
}

pub fn keep_all_filtered_out_test() {
  keep([1, 3, 5], fn(x) { x % 2 == 0 }) |> should.equal([])
}

pub fn keep_all_pass_test() {
  keep([2, 4, 6], fn(x) { x % 2 == 0 }) |> should.equal([2, 4, 6])
}

pub fn reduce_sum_test() {
  reduce([1, 2, 3, 4], 0, fn(acc, x) { acc + x }) |> should.equal(10)
}

pub fn reduce_product_test() {
  reduce([1, 2, 3, 4], 1, fn(acc, x) { acc * x }) |> should.equal(24)
}

pub fn reduce_empty_returns_init_test() {
  reduce([], 42, fn(acc, x) { acc + x }) |> should.equal(42)
}

pub fn sum_nonempty_test() {
  sum([1, 2, 3, 4, 5]) |> should.equal(15)
}

pub fn sum_empty_test() {
  sum([]) |> should.equal(0)
}

pub fn find_first_found_test() {
  find_first([1, 2, 3, 4], fn(x) { x > 2 }) |> should.equal(Ok(3))
}

pub fn find_first_not_found_test() {
  find_first([1, 2, 3], fn(x) { x > 10 }) |> should.equal(Error(Nil))
}

pub fn map_indexed_test() {
  map_indexed(["a", "b", "c"], fn(i, v) { #(i, v) })
  |> should.equal([#(0, "a"), #(1, "b"), #(2, "c")])
}

pub fn flat_transform_test() {
  flat_transform([[1, 2], [3, 4]], fn(x) { x * 10 })
  |> should.equal([10, 20, 30, 40])
}

pub fn flat_transform_empty_outer_test() {
  flat_transform([], fn(x) { x }) |> should.equal([])
}
