# Multiplex

The multiplex problem allows up to `n` participants to execute a critical
section at the same time.

Each participant has one phase:

1. wait for an available slot
2. run its critical section
3. release the slot for another participant

The synchronization requirement is:

- at most `n` critical sections may be active at once
- additional participants must wait until a slot becomes available
- a waiting participant should return early if its context is canceled

In other words, a multiplex is a generalized mutex with a positive concurrency
limit instead of a limit of one.

## Layout

- `multiplex.go`: challenge implementation stub
- `multiplex_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable multiplex tests
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./multiplex
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./multiplex/solutions
```

This runs the reusable test suite against the provided example solutions.

Benchmark the hint solutions:

```sh
go test ./multiplex/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./multiplex -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
