# Baboon Crossing

The baboon crossing problem has baboons swinging across a canyon on a single
rope.

Baboons traveling in opposite directions must never be on the rope at the same
time, the rope can hold at most five baboons, and a steady stream in one
direction must not starve baboons waiting to go the other way.

In this package, an implementation must:

- allow baboons to cross only with others going the same direction
- enforce a maximum of five baboons on the rope
- prevent one direction from starving the other
- let waiting baboons stop waiting if their context is canceled
- be reusable across multiple crossing waves

## Layout

- `baboon_crossing.go`: challenge implementation stub
- `baboon_crossing_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./baboon_crossing
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./baboon_crossing/solutions
```

Benchmark the hint solutions:

```sh
go test ./baboon_crossing/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./baboon_crossing -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
