import actors.{
  type CounterMsg, CounterState, Decrement, GetCount, Increment, Reset,
  handle_message,
}
import gleam/erlang/process
import gleam/otp/actor

pub type Counter =
  process.Subject(CounterMsg)

pub fn start() -> Result(Counter, actor.StartError) {
  let start_result =
    actor.new(CounterState(count: 0))
    |> actor.on_message(fn(state, msg) {
      actor.continue(handle_message(state, msg))
    })
    |> actor.start

  // Pattern match the Started record to extract the Subject cleanly
  case start_result {
    Ok(actor.Started(pid: _, data: subject)) -> Ok(subject)
    Error(e) -> Error(e)
  }
}

pub fn increment(counter: Counter) -> Nil {
  actor.send(counter, Increment)
}

pub fn decrement(counter: Counter) -> Nil {
  actor.send(counter, Decrement)
}

pub fn reset(counter: Counter) -> Nil {
  actor.send(counter, Reset)
}

pub fn get_count(counter: Counter) -> Int {
  actor.call(counter, 1000, fn(reply_subject) {
    GetCount(reply_with: fn(n) { process.send(reply_subject, n) })
  })
}

pub fn main() {
  let assert Ok(counter) = start()
  increment(counter)
  increment(counter)
  decrement(counter)
  let count = get_count(counter)
  let assert 1 = count
  reset(counter)
}
