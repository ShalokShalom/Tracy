pub type CounterMsg {
  Increment
  Decrement
  Reset
  GetCount(reply_with: fn(Int) -> Nil)
}

pub type CounterState {
  CounterState(count: Int)
}

pub fn handle_message(state: CounterState, msg: CounterMsg) -> CounterState {
  case msg {
    Increment -> CounterState(count: state.count + 1)
    Decrement -> CounterState(count: state.count - 1)
    Reset -> CounterState(count: 0)
    GetCount(reply) -> {
      reply(state.count)
      state
    }
  }
}
