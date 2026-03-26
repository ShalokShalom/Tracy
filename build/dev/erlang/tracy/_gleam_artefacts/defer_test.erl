-module(defer_test).
-compile([no_auto_import, nowarn_unused_vars, nowarn_unused_function, nowarn_nomatch, inline]).
-define(FILEPATH, "test/defer_test.gleam").
-export([with_lock_runs_body_test/0, with_lock_returns_body_value_test/0, with_resource_ok_path_test/0, with_resource_open_failure_test/0, with_resource_body_error_test/0, with_audit_appends_entry_test/0, with_audit_preserves_existing_log_test/0]).

-file("test/defer_test.gleam", 4).
-spec with_lock_runs_body_test() -> nil.
with_lock_runs_body_test() ->
    _pipe = defer:with_lock(fun() -> nil end, fun() -> nil end, fun() -> 42 end),
    gleeunit@should:equal(_pipe, 42).

-file("test/defer_test.gleam", 9).
-spec with_lock_returns_body_value_test() -> nil.
with_lock_returns_body_value_test() ->
    _pipe = defer:with_lock(
        fun() -> nil end,
        fun() -> nil end,
        fun() -> <<"done"/utf8>> end
    ),
    gleeunit@should:equal(_pipe, <<"done"/utf8>>).

-file("test/defer_test.gleam", 14).
-spec with_resource_ok_path_test() -> nil.
with_resource_ok_path_test() ->
    _pipe = defer:with_resource(
        fun() -> {ok, <<"file_handle"/utf8>>} end,
        fun(_) -> nil end,
        fun(H) -> {ok, <<"read: "/utf8, H/binary>>} end
    ),
    gleeunit@should:equal(_pipe, {ok, <<"read: file_handle"/utf8>>}).

-file("test/defer_test.gleam", 21).
-spec with_resource_open_failure_test() -> nil.
with_resource_open_failure_test() ->
    _pipe = defer:with_resource(
        fun() -> {error, <<"no such file"/utf8>>} end,
        fun(_) -> nil end,
        fun(_) -> {ok, <<"should not run"/utf8>>} end
    ),
    gleeunit@should:equal(_pipe, {error, <<"no such file"/utf8>>}).

-file("test/defer_test.gleam", 28).
-spec with_resource_body_error_test() -> nil.
with_resource_body_error_test() ->
    _pipe = defer:with_resource(
        fun() -> {ok, 1} end,
        fun(_) -> nil end,
        fun(_) -> {error, <<"body failed"/utf8>>} end
    ),
    gleeunit@should:equal(_pipe, {error, <<"body failed"/utf8>>}).

-file("test/defer_test.gleam", 33).
-spec with_audit_appends_entry_test() -> nil.
with_audit_appends_entry_test() ->
    {Result, Log} = defer:with_audit(
        [],
        <<"action:create"/utf8>>,
        fun() -> 99 end
    ),
    _pipe = Result,
    gleeunit@should:equal(_pipe, 99),
    _pipe@1 = Log,
    gleeunit@should:equal(_pipe@1, [<<"action:create"/utf8>>]).

-file("test/defer_test.gleam", 39).
-spec with_audit_preserves_existing_log_test() -> nil.
with_audit_preserves_existing_log_test() ->
    {_, Log} = defer:with_audit(
        [<<"prev"/utf8>>],
        <<"action:update"/utf8>>,
        fun() -> nil end
    ),
    _pipe = Log,
    gleeunit@should:equal(_pipe, [<<"prev"/utf8>>, <<"action:update"/utf8>>]).
