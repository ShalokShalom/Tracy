import defer.{with_audit, with_lock, with_resource}
import gleeunit/should

pub fn with_lock_runs_body_test() {
  with_lock(fn() { Nil }, fn() { Nil }, fn() { 42 })
  |> should.equal(42)
}

pub fn with_lock_returns_body_value_test() {
  with_lock(fn() { Nil }, fn() { Nil }, fn() { "done" })
  |> should.equal("done")
}

pub fn with_resource_ok_path_test() {
  with_resource(fn() { Ok("file_handle") }, fn(_) { Nil }, fn(h) {
    Ok("read: " <> h)
  })
  |> should.equal(Ok("read: file_handle"))
}

pub fn with_resource_open_failure_test() {
  with_resource(fn() { Error("no such file") }, fn(_) { Nil }, fn(_) {
    Ok("should not run")
  })
  |> should.equal(Error("no such file"))
}

pub fn with_resource_body_error_test() {
  with_resource(fn() { Ok(1) }, fn(_) { Nil }, fn(_) { Error("body failed") })
  |> should.equal(Error("body failed"))
}

pub fn with_audit_appends_entry_test() {
  let #(result, log) = with_audit([], "action:create", fn() { 99 })
  result |> should.equal(99)
  log |> should.equal(["action:create"])
}

pub fn with_audit_preserves_existing_log_test() {
  let #(_, log) = with_audit(["prev"], "action:update", fn() { Nil })
  log |> should.equal(["prev", "action:update"])
}
