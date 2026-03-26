import gleeunit/should
import results.{divide, pipeline, unwrap_or_zero, with_context}

pub fn divide_ok_test() {
  divide(10.0, 2.0)
  |> should.equal(Ok(5.0))
}

pub fn divide_error_test() {
  divide(10.0, 0.0)
  |> should.equal(Error("division by zero"))
}

pub fn pipeline_ok_test() {
  pipeline(2.0)
  |> should.equal(Ok(3.0))
}

pub fn pipeline_error_test() {
  pipeline(0.0)
  |> should.equal(Error("division by zero"))
}

pub fn with_context_ok_test() {
  divide(10.0, 2.0)
  |> with_context("calculation")
  |> should.equal(Ok(5.0))
}

pub fn with_context_error_test() {
  divide(10.0, 0.0)
  |> with_context("calculation")
  |> should.equal(Error("calculation: division by zero"))
}

pub fn unwrap_or_zero_ok_test() {
  divide(10.0, 2.0)
  |> unwrap_or_zero()
  |> should.equal(5.0)
}

pub fn unwrap_or_zero_error_test() {
  divide(10.0, 0.0)
  |> unwrap_or_zero()
  |> should.equal(0.0)
}
