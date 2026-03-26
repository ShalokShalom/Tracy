-module(tracy).
-compile([no_auto_import, nowarn_unused_vars, nowarn_unused_function, nowarn_nomatch, inline]).
-define(FILEPATH, "src/tracy.gleam").
-export([start/0, increment/1, decrement/1, reset/1, get_count/1, main/0]).

-file("src/tracy.gleam", 11).
-spec start() -> {ok, gleam@erlang@process:subject(actors:counter_msg())} |
    {error, gleam@otp@actor:start_error()}.
start() ->
    Start_result = begin
        _pipe = gleam@otp@actor:new({counter_state, 0}),
        _pipe@1 = gleam@otp@actor:on_message(
            _pipe,
            fun(State, Msg) ->
                gleam@otp@actor:continue(actors:handle_message(State, Msg))
            end
        ),
        gleam@otp@actor:start(_pipe@1)
    end,
    case Start_result of
        {ok, {started, _, Subject}} ->
            {ok, Subject};

        {error, E} ->
            {error, E}
    end.

-file("src/tracy.gleam", 26).
-spec increment(gleam@erlang@process:subject(actors:counter_msg())) -> nil.
increment(Counter) ->
    gleam@otp@actor:send(Counter, increment).

-file("src/tracy.gleam", 30).
-spec decrement(gleam@erlang@process:subject(actors:counter_msg())) -> nil.
decrement(Counter) ->
    gleam@otp@actor:send(Counter, decrement).

-file("src/tracy.gleam", 34).
-spec reset(gleam@erlang@process:subject(actors:counter_msg())) -> nil.
reset(Counter) ->
    gleam@otp@actor:send(Counter, reset).

-file("src/tracy.gleam", 38).
-spec get_count(gleam@erlang@process:subject(actors:counter_msg())) -> integer().
get_count(Counter) ->
    gleam@otp@actor:call(
        Counter,
        1000,
        fun(Reply_subject) ->
            {get_count,
                fun(N) -> gleam@erlang@process:send(Reply_subject, N) end}
        end
    ).

-file("src/tracy.gleam", 44).
-spec main() -> nil.
main() ->
    Counter@1 = case start() of
        {ok, Counter} -> Counter;
        _assert_fail ->
            erlang:error(#{gleam_error => let_assert,
                        message => <<"Pattern match failed, no pattern matched the value."/utf8>>,
                        file => <<?FILEPATH/utf8>>,
                        module => <<"tracy"/utf8>>,
                        function => <<"main"/utf8>>,
                        line => 45,
                        value => _assert_fail,
                        start => 1040,
                        'end' => 1072,
                        pattern_start => 1051,
                        pattern_end => 1062})
    end,
    increment(Counter@1),
    increment(Counter@1),
    decrement(Counter@1),
    Count = get_count(Counter@1),
    case Count of
        1 -> nil;
        _assert_fail@1 ->
            erlang:error(#{gleam_error => let_assert,
                        message => <<"Pattern match failed, no pattern matched the value."/utf8>>,
                        file => <<?FILEPATH/utf8>>,
                        module => <<"tracy"/utf8>>,
                        function => <<"main"/utf8>>,
                        line => 50,
                        value => _assert_fail@1,
                        start => 1171,
                        'end' => 1191,
                        pattern_start => 1182,
                        pattern_end => 1183})
    end,
    reset(Counter@1).
