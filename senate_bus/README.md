# Senate Bus

The Senate bus problem models riders waiting at a bus stop for buses with a
fixed capacity.

When a bus arrives, all currently waiting riders may board, up to the bus's
capacity. Riders who arrive while boarding is in progress must wait for the
next bus. When all riders selected for the current bus have boarded, the bus
may depart.

In this package, a Senate bus implementation must:

- let waiting riders board when a bus arrives
- allow a bus to depart immediately if no riders are waiting
- limit boarding to at most the bus capacity
- prevent riders who arrive during boarding from joining the current bus
- allow a waiting rider to stop waiting if its context is canceled
- remain reusable across multiple buses

## Layout

- `senate_bus.go`: challenge implementation stub
- `senate_bus_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable Senate bus tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./senate_bus
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./senate_bus/solutions
```

Benchmark the hint solutions:

```sh
go test ./senate_bus/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./senate_bus -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
