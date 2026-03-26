-module(results_test).
-compile([no_auto_import, nowarn_unused_vars, nowarn_unused_function, nowarn_nomatch, inline]).
-define(FILEPATH, "test/results_test.gleam").
-export([parse_age_valid_test/0, parse_age_zero_test/0, parse_age_boundary_max_test/0, parse_age_out_of_range_test/0, parse_age_negative_test/0, parse_age_invalid_string_test/0, parse_positive_valid_test/0, parse_positive_zero_is_invalid_test/0, parse_positive_negative_is_invalid_test/0, parse_positive_bad_string_test/0, double_parse_valid_test/0, double_parse_propagates_parse_error_test/0, double_parse_propagates_validation_error_test/0, parse_and_clamp_within_range_test/0, parse_and_clamp_at_boundary_test/0, parse_and_clamp_exceeds_max_test/0, lookup_user_alice_test/0, lookup_user_bob_test/0, lookup_user_not_found_test/0, greet_user_found_test/0, greet_user_not_found_propagates_test/0]).

-file("test/results_test.gleam", 7).
-spec parse_age_valid_test() -> nil.
parse_age_valid_test() ->
    _pipe = results:parse_age(<<"25"/utf8>>),
    gleeunit@should:equal(_pipe, {ok, 25}).

-file("test/results_test.gleam", 11).
-spec parse_age_zero_test() -> nil.
parse_age_zero_test() ->
    _pipe = results:parse_age(<<"0"/utf8>>),
    gleeunit@should:equal(_pipe, {ok, 0}).

-file("test/results_test.gleam", 15).
-spec parse_age_boundary_max_test() -> nil.
parse_age_boundary_max_test() ->
    _pipe = results:parse_age(<<"150"/utf8>>),
    gleeunit@should:equal(_pipe, {ok, 150}).

-file("test/results_test.gleam", 19).
-spec parse_age_out_of_range_test() -> nil.
parse_age_out_of_range_test() ->
    _pipe = results:parse_age(<<"200"/utf8>>),
    gleeunit@should:equal(
        _pipe,
        {error, {validation_error, <<"age out of range: 200"/utf8>>}}
    ).

-file("test/results_test.gleam", 24).
-spec parse_age_negative_test() -> nil.
parse_age_negative_test() ->
    _pipe = results:parse_age(<<"-1"/utf8>>),
    gleeunit@should:equal(
        _pipe,
        {error, {validation_error, <<"age out of range: -1"/utf8>>}}
    ).

-file("test/results_test.gleam", 29).
-spec parse_age_invalid_string_test() -> nil.
parse_age_invalid_string_test() ->
    _pipe = results:parse_age(<<"abc"/utf8>>),
    gleeunit@should:equal(
        _pipe,
        {error, {parse_error, <<"invalid integer: abc"/utf8>>}}
    ).

-file("test/results_test.gleam", 34).
-spec parse_positive_valid_test() -> nil.
parse_positive_valid_test() ->
    _pipe = results:parse_positive(<<"42"/utf8>>),
    gleeunit@should:equal(_pipe, {ok, 42}).

-file("test/results_test.gleam", 38).
-spec parse_positive_zero_is_invalid_test() -> nil.
parse_positive_zero_is_invalid_test() ->
    _pipe = results:parse_positive(<<"0"/utf8>>),
    gleeunit@should:equal(
        _pipe,
        {error, {validation_error, <<"must be positive"/utf8>>}}
    ).

-file("test/results_test.gleam", 43).
-spec parse_positive_negative_is_invalid_test() -> nil.
parse_positive_negative_is_invalid_test() ->
    _pipe = results:parse_positive(<<"-5"/utf8>>),
    gleeunit@should:equal(
        _pipe,
        {error, {validation_error, <<"must be positive"/utf8>>}}
    ).

-file("test/results_test.gleam", 48).
-spec parse_positive_bad_string_test() -> nil.
parse_positive_bad_string_test() ->
    _pipe = results:parse_positive(<<"x"/utf8>>),
    gleeunit@should:equal(
        _pipe,
        {error, {parse_error, <<"not a number: x"/utf8>>}}
    ).

-file("test/results_test.gleam", 52).
-spec double_parse_valid_test() -> nil.
double_parse_valid_test() ->
    _pipe = results:double_parse(<<"7"/utf8>>),
    gleeunit@should:equal(_pipe, {ok, 14}).

-file("test/results_test.gleam", 56).
-spec double_parse_propagates_parse_error_test() -> nil.
double_parse_propagates_parse_error_test() ->
    _pipe = results:double_parse(<<"bad"/utf8>>),
    gleeunit@should:equal(
        _pipe,
        {error, {parse_error, <<"not a number: bad"/utf8>>}}
    ).

-file("test/results_test.gleam", 60).
-spec double_parse_propagates_validation_error_test() -> nil.
double_parse_propagates_validation_error_test() ->
    _pipe = results:double_parse(<<"0"/utf8>>),
    gleeunit@should:equal(
        _pipe,
        {error, {validation_error, <<"must be positive"/utf8>>}}
    ).

-file("test/results_test.gleam", 64).
-spec parse_and_clamp_within_range_test() -> nil.
parse_and_clamp_within_range_test() ->
    _pipe = results:parse_and_clamp(<<"5"/utf8>>, 10),
    gleeunit@should:equal(_pipe, {ok, 5}).

-file("test/results_test.gleam", 68).
-spec parse_and_clamp_at_boundary_test() -> nil.
parse_and_clamp_at_boundary_test() ->
    _pipe = results:parse_and_clamp(<<"10"/utf8>>, 10),
    gleeunit@should:equal(_pipe, {ok, 10}).

-file("test/results_test.gleam", 72).
-spec parse_and_clamp_exceeds_max_test() -> nil.
parse_and_clamp_exceeds_max_test() ->
    _pipe = results:parse_and_clamp(<<"15"/utf8>>, 10),
    gleeunit@should:equal(
        _pipe,
        {error, {validation_error, <<"exceeds max: 10"/utf8>>}}
    ).

-file("test/results_test.gleam", 77).
-spec lookup_user_alice_test() -> nil.
lookup_user_alice_test() ->
    _pipe = results:lookup_user(1),
    gleeunit@should:equal(_pipe, {ok, <<"Alice"/utf8>>}).

-file("test/results_test.gleam", 81).
-spec lookup_user_bob_test() -> nil.
lookup_user_bob_test() ->
    _pipe = results:lookup_user(2),
    gleeunit@should:equal(_pipe, {ok, <<"Bob"/utf8>>}).

-file("test/results_test.gleam", 85).
-spec lookup_user_not_found_test() -> nil.
lookup_user_not_found_test() ->
    _pipe = results:lookup_user(99),
    gleeunit@should:equal(
        _pipe,
        {error, {not_found, <<"user not found: 99"/utf8>>}}
    ).

-file("test/results_test.gleam", 90).
-spec greet_user_found_test() -> nil.
greet_user_found_test() ->
    _pipe = results:greet_user(1),
    gleeunit@should:equal(_pipe, {ok, <<"Hello, Alice!"/utf8>>}).

-file("test/results_test.gleam", 94).
-spec greet_user_not_found_propagates_test() -> nil.
greet_user_not_found_propagates_test() ->
    _pipe = results:greet_user(99),
    gleeunit@should:equal(
        _pipe,
        {error, {not_found, <<"user not found: 99"/utf8>>}}
    ).
