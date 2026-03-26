import actors.{
  CounterState, Decrement, GetCount, Increment, Reset, handle_message,
}
import gleeunit/should

pub fn increment_increases_count_test() {
  CounterState(count: 0)
  |> handle_message(Increment)
  |> should.equal(CounterState(count: 1))
}

pub fn decrement_decreases_count_test() {
  CounterState(count: 5)
  |> handle_message(Decrement)
  |> should.equal(CounterState(count: 4))
}

pub fn reset_zeroes_count_test() {
  CounterState(count: 99)
  |> handle_message(Reset)
  |> should.equal(CounterState(count: 0))
}

pub fn get_count_does_not_change_state_test() {
  CounterState(count: 7)
  |> handle_message(GetCount(fn(_) { Nil }))
  |> should.equal(CounterState(count: 7))
}

pub fn sequential_messages_test() {
  CounterState(count: 0)
  |> handle_message(Increment)
  |> handle_message(Increment)
  |> handle_message(Increment)
  |> handle_message(Decrement)
  |> should.equal(CounterState(count: 2))
}

pub fn reset_after_increments_test() {
  CounterState(count: 0)
  |> handle_message(Increment)
  |> handle_message(Increment)
  |> handle_message(Reset)
  |> should.equal(CounterState(count: 0))
}
