-module(options_test).
-compile([no_auto_import, nowarn_unused_vars, nowarn_unused_function, nowarn_nomatch, inline]).
-define(FILEPATH, "test/options_test.gleam").
-export([find_existing_element_test/0, find_first_element_test/0, find_last_element_test/0, find_missing_element_test/0, find_empty_list_test/0, greet_with_name_test/0, greet_without_name_test/0, first_some_prefers_first_test/0, first_some_falls_back_to_second_test/0, first_some_both_none_test/0, map_opt_some_test/0, map_opt_none_test/0, flat_map_opt_some_to_some_test/0, flat_map_opt_some_to_none_test/0, flat_map_opt_none_short_circuits_test/0]).

-file("test/options_test.gleam", 5).
-spec find_existing_element_test() -> nil.
find_existing_element_test() ->
    _pipe = options:find_int([1, 2, 3, 4], 3),
    gleeunit@should:equal(_pipe, {some, 3}).

-file("test/options_test.gleam", 9).
-spec find_first_element_test() -> nil.
find_first_element_test() ->
    _pipe = options:find_int([7, 8, 9], 7),
    gleeunit@should:equal(_pipe, {some, 7}).

-file("test/options_test.gleam", 13).
-spec find_last_element_test() -> nil.
find_last_element_test() ->
    _pipe = options:find_int([7, 8, 9], 9),
    gleeunit@should:equal(_pipe, {some, 9}).

-file("test/options_test.gleam", 17).
-spec find_missing_element_test() -> nil.
find_missing_element_test() ->
    _pipe = options:find_int([1, 2, 3], 99),
    gleeunit@should:equal(_pipe, none).

-file("test/options_test.gleam", 21).
-spec find_empty_list_test() -> nil.
find_empty_list_test() ->
    _pipe = options:find_int([], 1),
    gleeunit@should:equal(_pipe, none).

-file("test/options_test.gleam", 25).
-spec greet_with_name_test() -> nil.
greet_with_name_test() ->
    _pipe = options:greet({some, <<"Bob"/utf8>>}),
    gleeunit@should:equal(_pipe, <<"Hello, Bob"/utf8>>).

-file("test/options_test.gleam", 29).
-spec greet_without_name_test() -> nil.
greet_without_name_test() ->
    _pipe = options:greet(none),
    gleeunit@should:equal(_pipe, <<"Hello, stranger"/utf8>>).

-file("test/options_test.gleam", 33).
-spec first_some_prefers_first_test() -> nil.
first_some_prefers_first_test() ->
    _pipe = options:first_some({some, <<"a"/utf8>>}, {some, <<"b"/utf8>>}),
    gleeunit@should:equal(_pipe, {some, <<"a"/utf8>>}).

-file("test/options_test.gleam", 37).
-spec first_some_falls_back_to_second_test() -> nil.
first_some_falls_back_to_second_test() ->
    _pipe = options:first_some(none, {some, <<"b"/utf8>>}),
    gleeunit@should:equal(_pipe, {some, <<"b"/utf8>>}).

-file("test/options_test.gleam", 41).
-spec first_some_both_none_test() -> nil.
first_some_both_none_test() ->
    Expected = none,
    _pipe = options:first_some(none, none),
    gleeunit@should:equal(_pipe, Expected).

-file("test/options_test.gleam", 46).
-spec map_opt_some_test() -> nil.
map_opt_some_test() ->
    _pipe = options:map_opt({some, 5}, fun(X) -> X * 2 end),
    gleeunit@should:equal(_pipe, {some, 10}).

-file("test/options_test.gleam", 50).
-spec map_opt_none_test() -> nil.
map_opt_none_test() ->
    _pipe = options:map_opt(none, fun(X) -> X * 2 end),
    gleeunit@should:equal(_pipe, none).

-file("test/options_test.gleam", 54).
-spec flat_map_opt_some_to_some_test() -> nil.
flat_map_opt_some_to_some_test() ->
    _pipe = options:flat_map_opt({some, 4}, fun(X) -> {some, X * 10} end),
    gleeunit@should:equal(_pipe, {some, 40}).

-file("test/options_test.gleam", 58).
-spec flat_map_opt_some_to_none_test() -> nil.
flat_map_opt_some_to_none_test() ->
    _pipe = options:flat_map_opt({some, 0}, fun(_) -> none end),
    gleeunit@should:equal(_pipe, none).

-file("test/options_test.gleam", 62).
-spec flat_map_opt_none_short_circuits_test() -> nil.
flat_map_opt_none_short_circuits_test() ->
    _pipe = options:flat_map_opt(none, fun(X) -> {some, X * 10} end),
    gleeunit@should:equal(_pipe, none).
