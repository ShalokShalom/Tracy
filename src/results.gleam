//// Phase 3 — Go (val, error) → Gleam Result(T, E)
////
//// Go's Multi-Return (val, error) entspricht Result(val, error) fast 1:1.
//// "if err != nil { return err }"-Ketten werden zu use-Chains mit result.try.

import gleam/result

/// Go:
///   func Divide(a, b float64) (float64, error) {
///     if b == 0 { return 0, errors.New("division by zero") }
///     return a / b, nil
///   }
pub fn divide(a: Float, b: Float) -> Result(Float, String) {
  case b {
    0.0 -> Error("division by zero")
    _ -> Ok(a /. b)
  }
}

/// Go:
///   a, err := stepOne(); if err != nil { return err }
///   b, err := stepTwo(a); if err != nil { return err }
///   return stepThree(b), nil
///
/// use flacht die Result-Kette auf — semantisch identisch zu Go's error chain.
pub fn pipeline(input: Float) -> Result(Float, String) {
  use a <- result.try(divide(10.0, input))
  use b <- result.try(divide(10.0, a))
  Ok(b +. 1.0)
}

/// Go: fmt.Errorf("context: %w", err)
pub fn with_context(
  r: Result(Float, String),
  ctx: String,
) -> Result(Float, String) {
  result.map_error(r, fn(e) { ctx <> ": " <> e })
}

/// Go: val, _ := divide(10, 2)
/// Ignorierter Fehler → result.unwrap mit explizitem Default.
/// Der Transpiler emittiert einen Kommentar wenn _ den Fehler ignoriert.
pub fn unwrap_or_zero(r: Result(Float, String)) -> Float {
  result.unwrap(r, 0.0)
}
