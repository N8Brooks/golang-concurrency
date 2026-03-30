# Room Party

The room party problem models students entering and leaving a room and a Dean
of Students who may either search the room or break up a party.

This package uses the threshold from the textbook: the Dean may enter to break
up the party once the room reaches `50` students. The Dean may also enter when
the room is empty to conduct a search.

Each student does one visit:

1. wait until the Dean is not in the room
2. enter and party
3. leave the room

The Dean does one visit:

1. if the room is empty, enter and search
2. if the room has between `1` and `49` students, wait
3. if the room has `50` or more students, enter and break up the party
4. once in the room, wait until all students have left before leaving

In this package, a room party implementation must:

- allow any number of students to be in the room concurrently
- allow the Dean to enter only when the room is empty or has at least `50` students
- block new students from entering while the Dean is in the room
- still allow students already in the room to leave while the Dean is present
- keep the Dean in the room until all students have left
- be reusable across many visits
- allow a waiting student or a waiting Dean to stop waiting if its context is canceled

The provided hint solution follows the same state-machine idea as the textbook
solution, but adapts it to the repository's callback-based API.

## Layout

- `room_party.go`: challenge implementation stub
- `room_party_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable room-party tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./room_party
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./room_party/solutions
```

Benchmark the hint solutions:

```sh
go test ./room_party/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./room_party -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
