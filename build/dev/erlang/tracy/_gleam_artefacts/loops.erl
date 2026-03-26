-module(loops).
-compile([no_auto_import, nowarn_unused_vars, nowarn_unused_function, nowarn_nomatch, inline]).
-define(FILEPATH, "src/loops.gleam").
-export([transform/2, keep/2, reduce/3, sum/1, find_first/2, map_indexed/2, flat_transform/2]).

-file("src/loops.gleam", 3).
-spec transform(list(integer()), fun((integer()) -> integer())) -> list(integer()).
transform(Items, F) ->
    gleam@list:map(Items, F).

-file("src/loops.gleam", 7).
-spec keep(list(integer()), fun((integer()) -> boolean())) -> list(integer()).
keep(Items, Pred) ->
    gleam@list:filter(Items, Pred).

-file("src/loops.gleam", 11).
-spec reduce(
    list(integer()),
    integer(),
    fun((integer(), integer()) -> integer())
) -> integer().
reduce(Items, Init, F) ->
    gleam@list:fold(Items, Init, F).

-file("src/loops.gleam", 15).
-spec sum(list(integer())) -> integer().
sum(Items) ->
    gleam@list:fold(Items, 0, fun(Acc, X) -> Acc + X end).

-file("src/loops.gleam", 19).
-spec find_first(list(integer()), fun((integer()) -> boolean())) -> {ok,
        integer()} |
    {error, nil}.
find_first(Items, Pred) ->
    gleam@list:find(Items, Pred).

-file("src/loops.gleam", 24).
-spec map_indexed(list(AIE), fun((integer(), AIE) -> AIG)) -> list(AIG).
map_indexed(Items, F) ->
    gleam@list:index_map(Items, fun(Item, I) -> F(I, Item) end).

-file("src/loops.gleam", 28).
-spec flat_transform(list(list(integer())), fun((integer()) -> integer())) -> list(integer()).
flat_transform(Items, F) ->
    gleam@list:flat_map(Items, fun(Inner) -> gleam@list:map(Inner, F) end).
