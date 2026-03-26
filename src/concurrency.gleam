pub type WorkerMsg {
  Job(payload: String)
  Stop
}

pub type WorkerState {
  WorkerState(processed: Int, last_result: String)
}

pub fn handle_worker(state: WorkerState, msg: WorkerMsg) -> WorkerState {
  case msg {
    Job(payload) ->
      WorkerState(processed: state.processed + 1, last_result: payload)
    Stop -> state
  }
}
