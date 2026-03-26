-module(records).
-compile([no_auto_import, nowarn_unused_vars, nowarn_unused_function, nowarn_nomatch, inline]).
-define(FILEPATH, "src/records.gleam").
-export([move/3, birthday/1, with_debug/1, with_port/2]).
-export_type([point/0, person/0, config/0]).

-type point() :: {point, float(), float()}.

-type person() :: {person, binary(), integer()}.

-type config() :: {config, binary(), integer(), boolean()}.

-file("src/records.gleam", 5).
-spec move(point(), float(), float()) -> point().
move(P, Dx, Dy) ->
    {point, erlang:element(2, P) + Dx, erlang:element(3, P) + Dy}.

-file("src/records.gleam", 13).
-spec birthday(person()) -> person().
birthday(P) ->
    {person, erlang:element(2, P), erlang:element(3, P) + 1}.

-file("src/records.gleam", 21).
-spec with_debug(config()) -> config().
with_debug(C) ->
    {config, erlang:element(2, C), erlang:element(3, C), true}.

-file("src/records.gleam", 25).
-spec with_port(config(), integer()) -> config().
with_port(C, Port) ->
    {config, erlang:element(2, C), Port, erlang:element(4, C)}.
