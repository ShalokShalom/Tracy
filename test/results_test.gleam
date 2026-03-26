import gleeunit/should
import results.{
  NotFound, ParseError, ValidationError, double_parse, greet_user, lookup_user,
  parse_age, parse_and_clamp, parse_positive,
}

pub fn parse_age_valid_test() {
  parse_age("25") |> should.equal(Ok(25))
}

pub fn parse_age_zero_test() {
  parse_age("0") |> should.equal(Ok(0))
}

pub fn parse_age_boundary_max_test() {
  parse_age("150") |> should.equal(Ok(150))
}

pub fn parse_age_out_of_range_test() {
  parse_age("200")
  |> should.equal(Error(ValidationError("age out of range: 200")))
}

pub fn parse_age_negative_test() {
  parse_age("-1")
  |> should.equal(Error(ValidationError("age out of range: -1")))
}

pub fn parse_age_invalid_string_test() {
  parse_age("abc")
  |> should.equal(Error(ParseError("invalid integer: abc")))
}

pub fn parse_positive_valid_test() {
  parse_positive("42") |> should.equal(Ok(42))
}

pub fn parse_positive_zero_is_invalid_test() {
  parse_positive("0")
  |> should.equal(Error(ValidationError("must be positive")))
}

pub fn parse_positive_negative_is_invalid_test() {
  parse_positive("-5")
  |> should.equal(Error(ValidationError("must be positive")))
}

pub fn parse_positive_bad_string_test() {
  parse_positive("x") |> should.equal(Error(ParseError("not a number: x")))
}

pub fn double_parse_valid_test() {
  double_parse("7") |> should.equal(Ok(14))
}

pub fn double_parse_propagates_parse_error_test() {
  double_parse("bad") |> should.equal(Error(ParseError("not a number: bad")))
}

pub fn double_parse_propagates_validation_error_test() {
  double_parse("0") |> should.equal(Error(ValidationError("must be positive")))
}

pub fn parse_and_clamp_within_range_test() {
  parse_and_clamp("5", 10) |> should.equal(Ok(5))
}

pub fn parse_and_clamp_at_boundary_test() {
  parse_and_clamp("10", 10) |> should.equal(Ok(10))
}

pub fn parse_and_clamp_exceeds_max_test() {
  parse_and_clamp("15", 10)
  |> should.equal(Error(ValidationError("exceeds max: 10")))
}

pub fn lookup_user_alice_test() {
  lookup_user(1) |> should.equal(Ok("Alice"))
}

pub fn lookup_user_bob_test() {
  lookup_user(2) |> should.equal(Ok("Bob"))
}

pub fn lookup_user_not_found_test() {
  lookup_user(99)
  |> should.equal(Error(NotFound("user not found: 99")))
}

pub fn greet_user_found_test() {
  greet_user(1) |> should.equal(Ok("Hello, Alice!"))
}

pub fn greet_user_not_found_propagates_test() {
  greet_user(99)
  |> should.equal(Error(NotFound("user not found: 99")))
}
