-module(actors).
-compile([no_auto_import, nowarn_unused_vars, nowarn_unused_function, nowarn_nomatch, inline]).
-define(FILEPATH, "src/actors.gleam").
-export([handle_message/2]).
-export_type([counter_msg/0, counter_state/0]).

-type counter_msg() :: increment |
    decrement |
    reset |
    {get_count, fun((integer()) -> nil)}.

-type counter_state() :: {counter_state, integer()}.

-file("src/actors.gleam", 12).
-spec handle_message(counter_state(), counter_msg()) -> counter_state().
handle_message(State, Msg) ->
    case Msg of
        increment ->
            {counter_state, erlang:element(2, State) + 1};

        decrement ->
            {counter_state, erlang:element(2, State) - 1};

        reset ->
            {counter_state, 0};

        {get_count, Reply} ->
            Reply(erlang:element(2, State)),
            State
    end.
