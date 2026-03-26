pub fn panic_stub(_msg: String) -> Result(a, String) {
  Error("TODO(transpiler): panic — redesign as Result(T, E)")
}

pub fn recover_stub(_body: fn() -> a) -> Result(a, String) {
  Error("TODO(transpiler): recover — redesign as Result(T, E)")
}

pub fn select_stub() -> Result(a, String) {
  Error("TODO(transpiler): select — redesign as OTP receive")
}

pub fn complex_defer_stub(_body: fn() -> a) -> Result(a, String) {
  Error("TODO(transpiler): complex defer — manual cleanup required")
}
