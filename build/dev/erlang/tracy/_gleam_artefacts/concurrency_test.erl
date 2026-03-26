-module(concurrency_test).
-compile([no_auto_import, nowarn_unused_vars, nowarn_unused_function, nowarn_nomatch, inline]).
-define(FILEPATH, "test/concurrency_test.gleam").
-export([worker_processes_job_test/0, worker_counts_jobs_test/0, worker_tracks_last_job_test/0, worker_stop_preserves_state_test/0]).

-file("test/concurrency_test.gleam", 4).
-spec worker_processes_job_test() -> nil.
worker_processes_job_test() ->
    _pipe = {worker_state, 0, <<""/utf8>>},
    _pipe@1 = concurrency:handle_worker(_pipe, {job, <<"task-1"/utf8>>}),
    gleeunit@should:equal(_pipe@1, {worker_state, 1, <<"task-1"/utf8>>}).

-file("test/concurrency_test.gleam", 10).
-spec worker_counts_jobs_test() -> nil.
worker_counts_jobs_test() ->
    S = begin
        _pipe = {worker_state, 0, <<""/utf8>>},
        _pipe@1 = concurrency:handle_worker(_pipe, {job, <<"a"/utf8>>}),
        _pipe@2 = concurrency:handle_worker(_pipe@1, {job, <<"b"/utf8>>}),
        concurrency:handle_worker(_pipe@2, {job, <<"c"/utf8>>})
    end,
    _pipe@3 = erlang:element(2, S),
    gleeunit@should:equal(_pipe@3, 3).

-file("test/concurrency_test.gleam", 19).
-spec worker_tracks_last_job_test() -> nil.
worker_tracks_last_job_test() ->
    S = begin
        _pipe = {worker_state, 0, <<""/utf8>>},
        _pipe@1 = concurrency:handle_worker(_pipe, {job, <<"first"/utf8>>}),
        concurrency:handle_worker(_pipe@1, {job, <<"second"/utf8>>})
    end,
    _pipe@2 = erlang:element(3, S),
    gleeunit@should:equal(_pipe@2, <<"second"/utf8>>).

-file("test/concurrency_test.gleam", 27).
-spec worker_stop_preserves_state_test() -> nil.
worker_stop_preserves_state_test() ->
    S = {worker_state, 2, <<"x"/utf8>>},
    _pipe = concurrency:handle_worker(S, stop),
    gleeunit@should:equal(_pipe, S).
