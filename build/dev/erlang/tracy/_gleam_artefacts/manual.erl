-module(manual).
-compile([no_auto_import, nowarn_unused_vars, nowarn_unused_function, nowarn_nomatch, inline]).
-define(FILEPATH, "src/manual.gleam").
-export([panic_stub/1, recover_stub/1, select_stub/0, complex_defer_stub/1]).

-file("src/manual.gleam", 1).
-spec panic_stub(binary()) -> {ok, any()} | {error, binary()}.
panic_stub(_) ->
    {error, <<"TODO(transpiler): panic — redesign as Result(T, E)"/utf8>>}.

-file("src/manual.gleam", 5).
-spec recover_stub(fun(() -> AMO)) -> {ok, AMO} | {error, binary()}.
recover_stub(_) ->
    {error, <<"TODO(transpiler): recover — redesign as Result(T, E)"/utf8>>}.

-file("src/manual.gleam", 9).
-spec select_stub() -> {ok, any()} | {error, binary()}.
select_stub() ->
    {error, <<"TODO(transpiler): select — redesign as OTP receive"/utf8>>}.

-file("src/manual.gleam", 13).
-spec complex_defer_stub(fun(() -> AMU)) -> {ok, AMU} | {error, binary()}.
complex_defer_stub(_) ->
    {error,
        <<"TODO(transpiler): complex defer — manual cleanup required"/utf8>>}.
