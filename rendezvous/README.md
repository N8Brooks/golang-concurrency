# Rendezvous

The rendezvous problem has two participants, `A` and `B`.

Each participant has two phases:

1. a first-phase action (`a1` / `b1`)
2. a second-phase action (`a2` / `b2`)

The synchronization requirement is:

- `a1` must happen before `a2`
- `b1` must happen before `b2`
- neither `a2` nor `b2` may happen until both `a1` and `b1` have completed

In other words, `A` and `B` must rendezvous between their first and second
phases.

## Layout

- `rendezvous.go`: challenge implementation stub
- `rendezvous_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable rendezvous tests
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./rendezvous
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./rendezvous/solutions
```

This runs the reusable test suite against the provided example solutions.

Benchmark the hint solutions:

```sh
go test ./rendezvous/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./rendezvous -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
