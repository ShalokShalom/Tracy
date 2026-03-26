import gleam/int
import gleam/result

pub type AppError {
  ParseError(String)
  ValidationError(String)
  NotFound(String)
}

pub fn parse_age(s: String) -> Result(Int, AppError) {
  case int.parse(s) {
    Ok(n) if n >= 0 && n <= 150 -> Ok(n)
    Ok(n) -> Error(ValidationError("age out of range: " <> int.to_string(n)))
    Error(_) -> Error(ParseError("invalid integer: " <> s))
  }
}

pub fn parse_positive(s: String) -> Result(Int, AppError) {
  case int.parse(s) {
    Ok(n) if n > 0 -> Ok(n)
    Ok(_) -> Error(ValidationError("must be positive"))
    Error(_) -> Error(ParseError("not a number: " <> s))
  }
}

pub fn double_parse(s: String) -> Result(Int, AppError) {
  use n <- result.try(parse_positive(s))
  Ok(n * 2)
}

pub fn parse_and_clamp(s: String, max: Int) -> Result(Int, AppError) {
  use n <- result.try(parse_positive(s))
  use _ <- result.try(case n <= max {
    True -> Ok(Nil)
    False -> Error(ValidationError("exceeds max: " <> int.to_string(max)))
  })
  Ok(n)
}

pub fn lookup_user(id: Int) -> Result(String, AppError) {
  case id {
    1 -> Ok("Alice")
    2 -> Ok("Bob")
    _ -> Error(NotFound("user not found: " <> int.to_string(id)))
  }
}

pub fn greet_user(id: Int) -> Result(String, AppError) {
  use name <- result.try(lookup_user(id))
  Ok("Hello, " <> name <> "!")
}
