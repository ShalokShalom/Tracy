-module(interfaces).
-compile([no_auto_import, nowarn_unused_vars, nowarn_unused_function, nowarn_nomatch, inline]).
-define(FILEPATH, "src/interfaces.gleam").
-export([describe/1, log_level/1, validate/1, area/1]).
-export_type([shape/0, logger/0, validator/0]).

-type shape() :: {circle, float()} |
    {rectangle, float(), float()} |
    {triangle, float(), float()}.

-type logger() :: {console_logger, binary()} | {file_logger, binary(), binary()}.

-type validator() :: {non_empty, binary()} |
    {min_length, binary(), integer()} |
    {max_length, binary(), integer()}.

-file("src/interfaces.gleam", 19).
-spec describe(shape()) -> binary().
describe(Shape) ->
    case Shape of
        {circle, _} ->
            <<"circle"/utf8>>;

        {rectangle, _, _} ->
            <<"rectangle"/utf8>>;

        {triangle, _, _} ->
            <<"triangle"/utf8>>
    end.

-file("src/interfaces.gleam", 32).
-spec log_level(logger()) -> binary().
log_level(Logger) ->
    case Logger of
        {console_logger, _} ->
            <<"stdout"/utf8>>;

        {file_logger, _, Level} ->
            Level
    end.

-file("src/interfaces.gleam", 45).
-spec validate(validator()) -> {ok, binary()} | {error, binary()}.
validate(V) ->
    case V of
        {non_empty, S} ->
            case S of
                <<""/utf8>> ->
                    {error, <<"must not be empty"/utf8>>};

                _ ->
                    {ok, S}
            end;

        {min_length, S@1, Min} ->
            case string:length(S@1) >= Min of
                true ->
                    {ok, S@1};

                false ->
                    {error, <<"too short"/utf8>>}
            end;

        {max_length, S@2, Max} ->
            case string:length(S@2) =< Max of
                true ->
                    {ok, S@2};

                false ->
                    {error, <<"too long"/utf8>>}
            end
    end.

-file("src/interfaces.gleam", 11).
-spec area(shape()) -> float().
area(Shape) ->
    case Shape of
        {circle, R} ->
            (3.14159265358979 * R) * R;

        {rectangle, W, H} ->
            W * H;

        {triangle, B, H@1} ->
            (0.5 * B) * H@1
    end.
