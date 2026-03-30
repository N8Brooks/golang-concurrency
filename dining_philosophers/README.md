# Dining Philosophers

The dining philosophers problem models five philosophers sitting around a
table with five forks. Philosopher `i` needs fork `i` on the right and fork
`(i+1)%5` on the left before eating.

Each philosopher performs one iteration of:

1. think
2. wait until both adjacent forks are available
3. eat while holding both forks
4. put both forks down

In this package, a dining philosophers implementation must:

- ensure neighboring philosophers never eat at the same time
- avoid deadlock when all five philosophers become hungry together
- allow non-neighboring philosophers to eat concurrently
- allow a waiting philosopher to stop waiting if its context is canceled
- be reusable across many dining rounds

The provided hint solutions mirror the textbook's deadlock-free approaches:

- `footman.go`: limits the table to four philosophers at a time
- `asymmetric.go`: makes one philosopher pick up forks in the opposite order

Tanenbaum's well-known state-machine solution is intentionally not included as
a hint here because it is deadlock-free but can still starve a philosopher.

## Layout

- `dining_philosophers.go`: challenge implementation stub
- `dining_philosophers_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable dining-philosophers tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./dining_philosophers
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./dining_philosophers/solutions
```

Benchmark the hint solutions:

```sh
go test ./dining_philosophers/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./dining_philosophers -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
