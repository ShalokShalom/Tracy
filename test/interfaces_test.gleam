import gleeunit/should
import interfaces.{
  Circle, ConsoleLogger, FileLogger, MaxLength, MinLength, NonEmpty, Rectangle,
  Triangle, area, describe, log_level, validate,
}

pub fn circle_area_test() {
  area(Circle(1.0)) |> should.equal(3.14159265358979)
}

pub fn rectangle_area_test() {
  area(Rectangle(4.0, 5.0)) |> should.equal(20.0)
}

pub fn triangle_area_test() {
  area(Triangle(6.0, 4.0)) |> should.equal(12.0)
}

pub fn describe_circle_test() {
  describe(Circle(1.0)) |> should.equal("circle")
}

pub fn describe_rectangle_test() {
  describe(Rectangle(1.0, 2.0)) |> should.equal("rectangle")
}

pub fn describe_triangle_test() {
  describe(Triangle(3.0, 4.0)) |> should.equal("triangle")
}

pub fn console_logger_level_test() {
  log_level(ConsoleLogger("INFO")) |> should.equal("stdout")
}

pub fn file_logger_level_test() {
  log_level(FileLogger("/var/log/app.log", "warn")) |> should.equal("warn")
}

pub fn validate_non_empty_ok_test() {
  validate(NonEmpty("hello")) |> should.equal(Ok("hello"))
}

pub fn validate_non_empty_error_test() {
  validate(NonEmpty("")) |> should.equal(Error("must not be empty"))
}

pub fn validate_min_length_ok_test() {
  validate(MinLength("hello", 3)) |> should.equal(Ok("hello"))
}

pub fn validate_min_length_exact_test() {
  validate(MinLength("abc", 3)) |> should.equal(Ok("abc"))
}

pub fn validate_min_length_too_short_test() {
  validate(MinLength("ab", 3)) |> should.equal(Error("too short"))
}

pub fn validate_max_length_ok_test() {
  validate(MaxLength("hi", 10)) |> should.equal(Ok("hi"))
}

pub fn validate_max_length_too_long_test() {
  validate(MaxLength("toolongstring", 5)) |> should.equal(Error("too long"))
}
