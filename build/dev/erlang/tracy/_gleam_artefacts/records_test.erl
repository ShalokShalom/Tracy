-module(records_test).
-compile([no_auto_import, nowarn_unused_vars, nowarn_unused_function, nowarn_nomatch, inline]).
-define(FILEPATH, "test/records_test.gleam").
-export([point_creation_test/0, move_returns_new_record_test/0, move_does_not_mutate_original_test/0, move_negative_delta_test/0, birthday_increments_age_test/0, birthday_preserves_name_test/0, birthday_chained_test/0, config_with_debug_test/0, config_with_port_test/0, config_update_preserves_other_fields_test/0]).

-file("test/records_test.gleam", 4).
-spec point_creation_test() -> nil.
point_creation_test() ->
    P = {point, 1.0, 2.0},
    _pipe = erlang:element(2, P),
    gleeunit@should:equal(_pipe, 1.0),
    _pipe@1 = erlang:element(3, P),
    gleeunit@should:equal(_pipe@1, 2.0).

-file("test/records_test.gleam", 10).
-spec move_returns_new_record_test() -> nil.
move_returns_new_record_test() ->
    P = {point, +0.0, +0.0},
    _pipe = records:move(P, 3.0, 4.0),
    gleeunit@should:equal(_pipe, {point, 3.0, 4.0}).

-file("test/records_test.gleam", 15).
-spec move_does_not_mutate_original_test() -> nil.
move_does_not_mutate_original_test() ->
    P = {point, 1.0, 1.0},
    _ = records:move(P, 5.0, 5.0),
    _pipe = P,
    gleeunit@should:equal(_pipe, {point, 1.0, 1.0}).

-file("test/records_test.gleam", 21).
-spec move_negative_delta_test() -> nil.
move_negative_delta_test() ->
    P = {point, 10.0, 10.0},
    _pipe = records:move(P, -3.0, -4.0),
    gleeunit@should:equal(_pipe, {point, 7.0, 6.0}).

-file("test/records_test.gleam", 26).
-spec birthday_increments_age_test() -> nil.
birthday_increments_age_test() ->
    Alice = {person, <<"Alice"/utf8>>, 30},
    _pipe = records:birthday(Alice),
    gleeunit@should:equal(_pipe, {person, <<"Alice"/utf8>>, 31}).

-file("test/records_test.gleam", 31).
-spec birthday_preserves_name_test() -> nil.
birthday_preserves_name_test() ->
    Bob = {person, <<"Bob"/utf8>>, 25},
    _pipe = erlang:element(2, records:birthday(Bob)),
    gleeunit@should:equal(_pipe, <<"Bob"/utf8>>).

-file("test/records_test.gleam", 36).
-spec birthday_chained_test() -> nil.
birthday_chained_test() ->
    _pipe = {person, <<"Eve"/utf8>>, 20},
    _pipe@1 = records:birthday(_pipe),
    _pipe@2 = records:birthday(_pipe@1),
    _pipe@3 = records:birthday(_pipe@2),
    gleeunit@should:equal(_pipe@3, {person, <<"Eve"/utf8>>, 23}).

-file("test/records_test.gleam", 44).
-spec config_with_debug_test() -> nil.
config_with_debug_test() ->
    C = {config, <<"localhost"/utf8>>, 8080, false},
    _pipe = records:with_debug(C),
    gleeunit@should:equal(_pipe, {config, <<"localhost"/utf8>>, 8080, true}).

-file("test/records_test.gleam", 50).
-spec config_with_port_test() -> nil.
config_with_port_test() ->
    C = {config, <<"localhost"/utf8>>, 8080, false},
    _pipe = records:with_port(C, 9090),
    gleeunit@should:equal(_pipe, {config, <<"localhost"/utf8>>, 9090, false}).

-file("test/records_test.gleam", 56).
-spec config_update_preserves_other_fields_test() -> nil.
config_update_preserves_other_fields_test() ->
    C = {config, <<"example.com"/utf8>>, 443, true},
    Updated = records:with_port(C, 80),
    _pipe = erlang:element(2, Updated),
    gleeunit@should:equal(_pipe, <<"example.com"/utf8>>),
    _pipe@1 = erlang:element(4, Updated),
    gleeunit@should:equal(_pipe@1, true).
