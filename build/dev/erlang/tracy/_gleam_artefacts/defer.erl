-module(defer).
-compile([no_auto_import, nowarn_unused_vars, nowarn_unused_function, nowarn_nomatch, inline]).
-define(FILEPATH, "src/defer.gleam").
-export([with_lock/3, with_resource/3, with_audit/3]).

-file("src/defer.gleam", 3).
-spec with_lock(fun(() -> any()), fun(() -> nil), fun(() -> ACC)) -> ACC.
with_lock(Acquire, Cleanup, Body) ->
    _ = Acquire(),
    Result = Body(),
    _ = Cleanup(),
    Result.

-file("src/defer.gleam", 10).
-spec with_resource(
    fun(() -> {ok, ACD} | {error, ACE}),
    fun((ACD) -> nil),
    fun((ACD) -> {ok, ACH} | {error, ACE})
) -> {ok, ACH} | {error, ACE}.
with_resource(Open, Close, Use_) ->
    case Open() of
        {error, E} ->
            {error, E};

        {ok, Resource} ->
            Result = Use_(Resource),
            _ = Close(Resource),
            Result
    end.

-file("src/defer.gleam", 25).
-spec with_audit(list(binary()), binary(), fun(() -> ACN)) -> {ACN,
    list(binary())}.
with_audit(Log, Entry, Body) ->
    Result = Body(),
    {Result, lists:append(Log, [Entry])}.
