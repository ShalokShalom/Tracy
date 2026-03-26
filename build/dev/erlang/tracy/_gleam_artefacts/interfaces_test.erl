-module(interfaces_test).
-compile([no_auto_import, nowarn_unused_vars, nowarn_unused_function, nowarn_nomatch, inline]).
-define(FILEPATH, "test/interfaces_test.gleam").
-export([circle_area_test/0, rectangle_area_test/0, triangle_area_test/0, describe_circle_test/0, describe_rectangle_test/0, describe_triangle_test/0, console_logger_level_test/0, file_logger_level_test/0, validate_non_empty_ok_test/0, validate_non_empty_error_test/0, validate_min_length_ok_test/0, validate_min_length_exact_test/0, validate_min_length_too_short_test/0, validate_max_length_ok_test/0, validate_max_length_too_long_test/0]).

-file("test/interfaces_test.gleam", 7).
-spec circle_area_test() -> nil.
circle_area_test() ->
    _pipe = interfaces:area({circle, 1.0}),
    gleeunit@should:equal(_pipe, 3.14159265358979).

-file("test/interfaces_test.gleam", 11).
-spec rectangle_area_test() -> nil.
rectangle_area_test() ->
    _pipe = interfaces:area({rectangle, 4.0, 5.0}),
    gleeunit@should:equal(_pipe, 20.0).

-file("test/interfaces_test.gleam", 15).
-spec triangle_area_test() -> nil.
triangle_area_test() ->
    _pipe = interfaces:area({triangle, 6.0, 4.0}),
    gleeunit@should:equal(_pipe, 12.0).

-file("test/interfaces_test.gleam", 19).
-spec describe_circle_test() -> nil.
describe_circle_test() ->
    _pipe = interfaces:describe({circle, 1.0}),
    gleeunit@should:equal(_pipe, <<"circle"/utf8>>).

-file("test/interfaces_test.gleam", 23).
-spec describe_rectangle_test() -> nil.
describe_rectangle_test() ->
    _pipe = interfaces:describe({rectangle, 1.0, 2.0}),
    gleeunit@should:equal(_pipe, <<"rectangle"/utf8>>).

-file("test/interfaces_test.gleam", 27).
-spec describe_triangle_test() -> nil.
describe_triangle_test() ->
    _pipe = interfaces:describe({triangle, 3.0, 4.0}),
    gleeunit@should:equal(_pipe, <<"triangle"/utf8>>).

-file("test/interfaces_test.gleam", 31).
-spec console_logger_level_test() -> nil.
console_logger_level_test() ->
    _pipe = interfaces:log_level({console_logger, <<"INFO"/utf8>>}),
    gleeunit@should:equal(_pipe, <<"stdout"/utf8>>).

-file("test/interfaces_test.gleam", 35).
-spec file_logger_level_test() -> nil.
file_logger_level_test() ->
    _pipe = interfaces:log_level(
        {file_logger, <<"/var/log/app.log"/utf8>>, <<"warn"/utf8>>}
    ),
    gleeunit@should:equal(_pipe, <<"warn"/utf8>>).

-file("test/interfaces_test.gleam", 39).
-spec validate_non_empty_ok_test() -> nil.
validate_non_empty_ok_test() ->
    _pipe = interfaces:validate({non_empty, <<"hello"/utf8>>}),
    gleeunit@should:equal(_pipe, {ok, <<"hello"/utf8>>}).

-file("test/interfaces_test.gleam", 43).
-spec validate_non_empty_error_test() -> nil.
validate_non_empty_error_test() ->
    _pipe = interfaces:validate({non_empty, <<""/utf8>>}),
    gleeunit@should:equal(_pipe, {error, <<"must not be empty"/utf8>>}).

-file("test/interfaces_test.gleam", 47).
-spec validate_min_length_ok_test() -> nil.
validate_min_length_ok_test() ->
    _pipe = interfaces:validate({min_length, <<"hello"/utf8>>, 3}),
    gleeunit@should:equal(_pipe, {ok, <<"hello"/utf8>>}).

-file("test/interfaces_test.gleam", 51).
-spec validate_min_length_exact_test() -> nil.
validate_min_length_exact_test() ->
    _pipe = interfaces:validate({min_length, <<"abc"/utf8>>, 3}),
    gleeunit@should:equal(_pipe, {ok, <<"abc"/utf8>>}).

-file("test/interfaces_test.gleam", 55).
-spec validate_min_length_too_short_test() -> nil.
validate_min_length_too_short_test() ->
    _pipe = interfaces:validate({min_length, <<"ab"/utf8>>, 3}),
    gleeunit@should:equal(_pipe, {error, <<"too short"/utf8>>}).

-file("test/interfaces_test.gleam", 59).
-spec validate_max_length_ok_test() -> nil.
validate_max_length_ok_test() ->
    _pipe = interfaces:validate({max_length, <<"hi"/utf8>>, 10}),
    gleeunit@should:equal(_pipe, {ok, <<"hi"/utf8>>}).

-file("test/interfaces_test.gleam", 63).
-spec validate_max_length_too_long_test() -> nil.
validate_max_length_too_long_test() ->
    _pipe = interfaces:validate({max_length, <<"toolongstring"/utf8>>, 5}),
    gleeunit@should:equal(_pipe, {error, <<"too long"/utf8>>}).
