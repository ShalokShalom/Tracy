-module(results).
-compile([no_auto_import, nowarn_unused_vars, nowarn_unused_function, nowarn_nomatch, inline]).
-define(FILEPATH, "src/results.gleam").
-export([parse_age/1, parse_positive/1, double_parse/1, parse_and_clamp/2, lookup_user/1, greet_user/1]).
-export_type([app_error/0]).

-type app_error() :: {parse_error, binary()} |
    {validation_error, binary()} |
    {not_found, binary()}.

-file("src/results.gleam", 10).
-spec parse_age(binary()) -> {ok, integer()} | {error, app_error()}.
parse_age(S) ->
    case gleam_stdlib:parse_int(S) of
        {ok, N} when (N >= 0) andalso (N =< 150) ->
            {ok, N};

        {ok, N@1} ->
            {error,
                {validation_error,
                    <<"age out of range: "/utf8,
                        (erlang:integer_to_binary(N@1))/binary>>}};

        {error, _} ->
            {error, {parse_error, <<"invalid integer: "/utf8, S/binary>>}}
    end.

-file("src/results.gleam", 18).
-spec parse_positive(binary()) -> {ok, integer()} | {error, app_error()}.
parse_positive(S) ->
    case gleam_stdlib:parse_int(S) of
        {ok, N} when N > 0 ->
            {ok, N};

        {ok, _} ->
            {error, {validation_error, <<"must be positive"/utf8>>}};

        {error, _} ->
            {error, {parse_error, <<"not a number: "/utf8, S/binary>>}}
    end.

-file("src/results.gleam", 26).
-spec double_parse(binary()) -> {ok, integer()} | {error, app_error()}.
double_parse(S) ->
    gleam@result:'try'(parse_positive(S), fun(N) -> {ok, N * 2} end).

-file("src/results.gleam", 31).
-spec parse_and_clamp(binary(), integer()) -> {ok, integer()} |
    {error, app_error()}.
parse_and_clamp(S, Max) ->
    gleam@result:'try'(
        parse_positive(S),
        fun(N) -> gleam@result:'try'(case N =< Max of
                    true ->
                        {ok, nil};

                    false ->
                        {error,
                            {validation_error,
                                <<"exceeds max: "/utf8,
                                    (erlang:integer_to_binary(Max))/binary>>}}
                end, fun(_) -> {ok, N} end) end
    ).

-file("src/results.gleam", 40).
-spec lookup_user(integer()) -> {ok, binary()} | {error, app_error()}.
lookup_user(Id) ->
    case Id of
        1 ->
            {ok, <<"Alice"/utf8>>};

        2 ->
            {ok, <<"Bob"/utf8>>};

        _ ->
            {error,
                {not_found,
                    <<"user not found: "/utf8,
                        (erlang:integer_to_binary(Id))/binary>>}}
    end.

-file("src/results.gleam", 48).
-spec greet_user(integer()) -> {ok, binary()} | {error, app_error()}.
greet_user(Id) ->
    gleam@result:'try'(
        lookup_user(Id),
        fun(Name) ->
            {ok, <<<<"Hello, "/utf8, Name/binary>>/binary, "!"/utf8>>}
        end
    ).
