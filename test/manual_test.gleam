import gleeunit/should
import manual.{complex_defer_stub, panic_stub, recover_stub, select_stub}

pub fn panic_stub_compiles_test() {
  panic_stub("something went wrong")
  |> should.equal(Error("TODO(transpiler): panic — redesign as Result(T, E)"))
}

pub fn recover_stub_compiles_test() {
  recover_stub(fn() { 42 })
  |> should.equal(Error("TODO(transpiler): recover — redesign as Result(T, E)"))
}

pub fn select_stub_compiles_test() {
  select_stub()
  |> should.equal(Error("TODO(transpiler): select — redesign as OTP receive"))
}

pub fn complex_defer_stub_compiles_test() {
  complex_defer_stub(fn() { "result" })
  |> should.equal(Error(
    "TODO(transpiler): complex defer — manual cleanup required",
  ))
}
