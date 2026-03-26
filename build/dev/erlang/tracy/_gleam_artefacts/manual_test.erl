-module(manual_test).
-compile([no_auto_import, nowarn_unused_vars, nowarn_unused_function, nowarn_nomatch, inline]).
-define(FILEPATH, "test/manual_test.gleam").
-export([panic_stub_compiles_test/0, recover_stub_compiles_test/0, select_stub_compiles_test/0, complex_defer_stub_compiles_test/0]).

-file("test/manual_test.gleam", 4).
-spec panic_stub_compiles_test() -> nil.
panic_stub_compiles_test() ->
    _pipe = manual:panic_stub(<<"something went wrong"/utf8>>),
    gleeunit@should:equal(
        _pipe,
        {error, <<"TODO(transpiler): panic — redesign as Result(T, E)"/utf8>>}
    ).

-file("test/manual_test.gleam", 9).
-spec recover_stub_compiles_test() -> nil.
recover_stub_compiles_test() ->
    _pipe = manual:recover_stub(fun() -> 42 end),
    gleeunit@should:equal(
        _pipe,
        {error, <<"TODO(transpiler): recover — redesign as Result(T, E)"/utf8>>}
    ).

-file("test/manual_test.gleam", 14).
-spec select_stub_compiles_test() -> nil.
select_stub_compiles_test() ->
    _pipe = manual:select_stub(),
    gleeunit@should:equal(
        _pipe,
        {error, <<"TODO(transpiler): select — redesign as OTP receive"/utf8>>}
    ).

-file("test/manual_test.gleam", 19).
-spec complex_defer_stub_compiles_test() -> nil.
complex_defer_stub_compiles_test() ->
    _pipe = manual:complex_defer_stub(fun() -> <<"result"/utf8>> end),
    gleeunit@should:equal(
        _pipe,
        {error,
            <<"TODO(transpiler): complex defer — manual cleanup required"/utf8>>}
    ).
