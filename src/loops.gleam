import gleam/list

pub fn transform(items: List(Int), f: fn(Int) -> Int) -> List(Int) {
  list.map(items, f)
}

pub fn keep(items: List(Int), pred: fn(Int) -> Bool) -> List(Int) {
  list.filter(items, pred)
}

pub fn reduce(items: List(Int), init: Int, f: fn(Int, Int) -> Int) -> Int {
  list.fold(items, init, f)
}

pub fn sum(items: List(Int)) -> Int {
  list.fold(items, 0, fn(acc, x) { acc + x })
}

pub fn find_first(items: List(Int), pred: fn(Int) -> Bool) -> Result(Int, Nil) {
  list.find(items, pred)
}

// stdlib's index_map passes fn(element, index), so we swap to expose fn(index, element)
pub fn map_indexed(items: List(a), f: fn(Int, a) -> b) -> List(b) {
  list.index_map(items, fn(item, i) { f(i, item) })
}

pub fn flat_transform(items: List(List(Int)), f: fn(Int) -> Int) -> List(Int) {
  list.flat_map(items, fn(inner) { list.map(inner, f) })
}
