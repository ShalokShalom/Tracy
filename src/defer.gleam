import gleam/list

pub fn with_lock(acquire: fn() -> a, cleanup: fn() -> Nil, body: fn() -> b) -> b {
  let _ = acquire()
  let result = body()
  let _ = cleanup()
  result
}

pub fn with_resource(
  open: fn() -> Result(a, e),
  close: fn(a) -> Nil,
  use_: fn(a) -> Result(b, e),
) -> Result(b, e) {
  case open() {
    Error(e) -> Error(e)
    Ok(resource) -> {
      let result = use_(resource)
      let _ = close(resource)
      result
    }
  }
}

pub fn with_audit(
  log: List(String),
  entry: String,
  body: fn() -> a,
) -> #(a, List(String)) {
  let result = body()
  #(result, list.append(log, [entry]))
}
