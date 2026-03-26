-module(options).
-compile([no_auto_import, nowarn_unused_vars, nowarn_unused_function, nowarn_nomatch, inline]).
-define(FILEPATH, "src/options.gleam").
-export([find_int/2, greet/1, first_some/2, map_opt/2, flat_map_opt/2]).

-file("src/options.gleam", 4).
-spec find_int(list(integer()), integer()) -> gleam@option:option(integer()).
find_int(Items, Target) ->
    _pipe = gleam@list:find(Items, fun(X) -> X =:= Target end),
    gleam@option:from_result(_pipe).

-file("src/options.gleam", 9).
-spec greet(gleam@option:option(binary())) -> binary().
greet(Name) ->
    case Name of
        {some, N} ->
            <<"Hello, "/utf8, N/binary>>;

        none ->
            <<"Hello, stranger"/utf8>>
    end.

-file("src/options.gleam", 16).
-spec first_some(gleam@option:option(AOE), gleam@option:option(AOE)) -> gleam@option:option(AOE).
first_some(A, B) ->
    case A of
        {some, _} ->
            A;

        none ->
            B
    end.

-file("src/options.gleam", 23).
-spec map_opt(gleam@option:option(integer()), fun((integer()) -> integer())) -> gleam@option:option(integer()).
map_opt(Opt, F) ->
    gleam@option:map(Opt, F).

-file("src/options.gleam", 27).
-spec flat_map_opt(
    gleam@option:option(integer()),
    fun((integer()) -> gleam@option:option(integer()))
) -> gleam@option:option(integer()).
flat_map_opt(Opt, F) ->
    gleam@option:then(Opt, F).
