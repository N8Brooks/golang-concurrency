# Barrier

The barrier problem has a fixed number of participants.

Each participant has two phases:

1. a first-phase action
2. a second-phase action

The synchronization requirement is:

- each participant must run its first-phase action before its second-phase action
- no participant may run its second-phase action until every participant in the round has completed its first-phase action

In other words, all participants must wait at the barrier between their first
and second phases.

## Layout

- `barrier.go`: challenge implementation stub
- `barrier_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable barrier tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./barrier
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./barrier/solutions
```

This runs the reusable test suite against the provided example solutions.

Benchmark the hint solutions:

```sh
go test ./barrier/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./barrier -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
