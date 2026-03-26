-module(concurrency).
-compile([no_auto_import, nowarn_unused_vars, nowarn_unused_function, nowarn_nomatch, inline]).
-define(FILEPATH, "src/concurrency.gleam").
-export([handle_worker/2]).
-export_type([worker_msg/0, worker_state/0]).

-type worker_msg() :: {job, binary()} | stop.

-type worker_state() :: {worker_state, integer(), binary()}.

-file("src/concurrency.gleam", 10).
-spec handle_worker(worker_state(), worker_msg()) -> worker_state().
handle_worker(State, Msg) ->
    case Msg of
        {job, Payload} ->
            {worker_state, erlang:element(2, State) + 1, Payload};

        stop ->
            State
    end.
