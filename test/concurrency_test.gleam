import concurrency.{Job, Stop, WorkerState, handle_worker}
import gleeunit/should

pub fn worker_processes_job_test() {
  WorkerState(processed: 0, last_result: "")
  |> handle_worker(Job("task-1"))
  |> should.equal(WorkerState(processed: 1, last_result: "task-1"))
}

pub fn worker_counts_jobs_test() {
  let s =
    WorkerState(processed: 0, last_result: "")
    |> handle_worker(Job("a"))
    |> handle_worker(Job("b"))
    |> handle_worker(Job("c"))
  s.processed |> should.equal(3)
}

pub fn worker_tracks_last_job_test() {
  let s =
    WorkerState(processed: 0, last_result: "")
    |> handle_worker(Job("first"))
    |> handle_worker(Job("second"))
  s.last_result |> should.equal("second")
}

pub fn worker_stop_preserves_state_test() {
  let s = WorkerState(processed: 2, last_result: "x")
  handle_worker(s, Stop) |> should.equal(s)
}
