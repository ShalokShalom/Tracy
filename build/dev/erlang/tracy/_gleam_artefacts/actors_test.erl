-module(actors_test).
-compile([no_auto_import, nowarn_unused_vars, nowarn_unused_function, nowarn_nomatch, inline]).
-define(FILEPATH, "test/actors_test.gleam").
-export([increment_increases_count_test/0, decrement_decreases_count_test/0, reset_zeroes_count_test/0, get_count_does_not_change_state_test/0, sequential_messages_test/0, reset_after_increments_test/0]).

-file("test/actors_test.gleam", 6).
-spec increment_increases_count_test() -> nil.
increment_increases_count_test() ->
    _pipe = {counter_state, 0},
    _pipe@1 = actors:handle_message(_pipe, increment),
    gleeunit@should:equal(_pipe@1, {counter_state, 1}).

-file("test/actors_test.gleam", 12).
-spec decrement_decreases_count_test() -> nil.
decrement_decreases_count_test() ->
    _pipe = {counter_state, 5},
    _pipe@1 = actors:handle_message(_pipe, decrement),
    gleeunit@should:equal(_pipe@1, {counter_state, 4}).

-file("test/actors_test.gleam", 18).
-spec reset_zeroes_count_test() -> nil.
reset_zeroes_count_test() ->
    _pipe = {counter_state, 99},
    _pipe@1 = actors:handle_message(_pipe, reset),
    gleeunit@should:equal(_pipe@1, {counter_state, 0}).

-file("test/actors_test.gleam", 24).
-spec get_count_does_not_change_state_test() -> nil.
get_count_does_not_change_state_test() ->
    _pipe = {counter_state, 7},
    _pipe@1 = actors:handle_message(_pipe, {get_count, fun(_) -> nil end}),
    gleeunit@should:equal(_pipe@1, {counter_state, 7}).

-file("test/actors_test.gleam", 30).
-spec sequential_messages_test() -> nil.
sequential_messages_test() ->
    _pipe = {counter_state, 0},
    _pipe@1 = actors:handle_message(_pipe, increment),
    _pipe@2 = actors:handle_message(_pipe@1, increment),
    _pipe@3 = actors:handle_message(_pipe@2, increment),
    _pipe@4 = actors:handle_message(_pipe@3, decrement),
    gleeunit@should:equal(_pipe@4, {counter_state, 2}).

-file("test/actors_test.gleam", 39).
-spec reset_after_increments_test() -> nil.
reset_after_increments_test() ->
    _pipe = {counter_state, 0},
    _pipe@1 = actors:handle_message(_pipe, increment),
    _pipe@2 = actors:handle_message(_pipe@1, increment),
    _pipe@3 = actors:handle_message(_pipe@2, reset),
    gleeunit@should:equal(_pipe@3, {counter_state, 0}).
