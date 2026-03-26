-module(loops_test).
-compile([no_auto_import, nowarn_unused_vars, nowarn_unused_function, nowarn_nomatch, inline]).
-define(FILEPATH, "test/loops_test.gleam").
-export([transform_doubles_test/0, transform_empty_test/0, transform_increments_test/0, keep_even_test/0, keep_all_filtered_out_test/0, keep_all_pass_test/0, reduce_sum_test/0, reduce_product_test/0, reduce_empty_returns_init_test/0, sum_nonempty_test/0, sum_empty_test/0, find_first_found_test/0, find_first_not_found_test/0, map_indexed_test/0, flat_transform_test/0, flat_transform_empty_outer_test/0]).

-file("test/loops_test.gleam", 6).
-spec transform_doubles_test() -> nil.
transform_doubles_test() ->
    _pipe = loops:transform([1, 2, 3], fun(X) -> X * 2 end),
    gleeunit@should:equal(_pipe, [2, 4, 6]).

-file("test/loops_test.gleam", 10).
-spec transform_empty_test() -> nil.
transform_empty_test() ->
    _pipe = loops:transform([], fun(X) -> X * 2 end),
    gleeunit@should:equal(_pipe, []).

-file("test/loops_test.gleam", 14).
-spec transform_increments_test() -> nil.
transform_increments_test() ->
    _pipe = loops:transform([10, 20], fun(X) -> X + 1 end),
    gleeunit@should:equal(_pipe, [11, 21]).

-file("test/loops_test.gleam", 18).
-spec keep_even_test() -> nil.
keep_even_test() ->
    _pipe = loops:keep([1, 2, 3, 4, 5, 6], fun(X) -> (X rem 2) =:= 0 end),
    gleeunit@should:equal(_pipe, [2, 4, 6]).

-file("test/loops_test.gleam", 22).
-spec keep_all_filtered_out_test() -> nil.
keep_all_filtered_out_test() ->
    _pipe = loops:keep([1, 3, 5], fun(X) -> (X rem 2) =:= 0 end),
    gleeunit@should:equal(_pipe, []).

-file("test/loops_test.gleam", 26).
-spec keep_all_pass_test() -> nil.
keep_all_pass_test() ->
    _pipe = loops:keep([2, 4, 6], fun(X) -> (X rem 2) =:= 0 end),
    gleeunit@should:equal(_pipe, [2, 4, 6]).

-file("test/loops_test.gleam", 30).
-spec reduce_sum_test() -> nil.
reduce_sum_test() ->
    _pipe = loops:reduce([1, 2, 3, 4], 0, fun(Acc, X) -> Acc + X end),
    gleeunit@should:equal(_pipe, 10).

-file("test/loops_test.gleam", 34).
-spec reduce_product_test() -> nil.
reduce_product_test() ->
    _pipe = loops:reduce([1, 2, 3, 4], 1, fun(Acc, X) -> Acc * X end),
    gleeunit@should:equal(_pipe, 24).

-file("test/loops_test.gleam", 38).
-spec reduce_empty_returns_init_test() -> nil.
reduce_empty_returns_init_test() ->
    _pipe = loops:reduce([], 42, fun(Acc, X) -> Acc + X end),
    gleeunit@should:equal(_pipe, 42).

-file("test/loops_test.gleam", 42).
-spec sum_nonempty_test() -> nil.
sum_nonempty_test() ->
    _pipe = loops:sum([1, 2, 3, 4, 5]),
    gleeunit@should:equal(_pipe, 15).

-file("test/loops_test.gleam", 46).
-spec sum_empty_test() -> nil.
sum_empty_test() ->
    _pipe = loops:sum([]),
    gleeunit@should:equal(_pipe, 0).

-file("test/loops_test.gleam", 50).
-spec find_first_found_test() -> nil.
find_first_found_test() ->
    _pipe = loops:find_first([1, 2, 3, 4], fun(X) -> X > 2 end),
    gleeunit@should:equal(_pipe, {ok, 3}).

-file("test/loops_test.gleam", 54).
-spec find_first_not_found_test() -> nil.
find_first_not_found_test() ->
    _pipe = loops:find_first([1, 2, 3], fun(X) -> X > 10 end),
    gleeunit@should:equal(_pipe, {error, nil}).

-file("test/loops_test.gleam", 58).
-spec map_indexed_test() -> nil.
map_indexed_test() ->
    _pipe = loops:map_indexed(
        [<<"a"/utf8>>, <<"b"/utf8>>, <<"c"/utf8>>],
        fun(I, V) -> {I, V} end
    ),
    gleeunit@should:equal(
        _pipe,
        [{0, <<"a"/utf8>>}, {1, <<"b"/utf8>>}, {2, <<"c"/utf8>>}]
    ).

-file("test/loops_test.gleam", 63).
-spec flat_transform_test() -> nil.
flat_transform_test() ->
    _pipe = loops:flat_transform([[1, 2], [3, 4]], fun(X) -> X * 10 end),
    gleeunit@should:equal(_pipe, [10, 20, 30, 40]).

-file("test/loops_test.gleam", 68).
-spec flat_transform_empty_outer_test() -> nil.
flat_transform_empty_outer_test() ->
    _pipe = loops:flat_transform([], fun(X) -> X end),
    gleeunit@should:equal(_pipe, []).
