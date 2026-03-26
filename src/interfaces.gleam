import gleam/string

pub type Shape {
  Circle(radius: Float)
  Rectangle(width: Float, height: Float)
  Triangle(base: Float, height: Float)
}

const pi = 3.14159265358979

pub fn area(shape: Shape) -> Float {
  case shape {
    Circle(r) -> pi *. r *. r
    Rectangle(w, h) -> w *. h
    Triangle(b, h) -> 0.5 *. b *. h
  }
}

pub fn describe(shape: Shape) -> String {
  case shape {
    Circle(_) -> "circle"
    Rectangle(_, _) -> "rectangle"
    Triangle(_, _) -> "triangle"
  }
}

pub type Logger {
  ConsoleLogger(prefix: String)
  FileLogger(path: String, level: String)
}

pub fn log_level(logger: Logger) -> String {
  case logger {
    ConsoleLogger(_) -> "stdout"
    FileLogger(_, level) -> level
  }
}

pub type Validator {
  NonEmpty(value: String)
  MinLength(value: String, min: Int)
  MaxLength(value: String, max: Int)
}

pub fn validate(v: Validator) -> Result(String, String) {
  case v {
    NonEmpty(s) ->
      case s {
        "" -> Error("must not be empty")
        _ -> Ok(s)
      }
    MinLength(s, min) ->
      case string.length(s) >= min {
        True -> Ok(s)
        False -> Error("too short")
      }
    MaxLength(s, max) ->
      case string.length(s) <= max {
        True -> Ok(s)
        False -> Error("too long")
      }
  }
}
