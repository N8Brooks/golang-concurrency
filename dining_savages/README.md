# Dining Savages

The dining savages problem has one cook and any number of savages sharing a
single communal pot.

The pot holds a fixed number of servings. When a savage wants to eat, it takes
one serving from the pot. If the pot is empty, a savage wakes the cook and
waits until the pot has been refilled.

In this package, a dining savages implementation must:

- prevent savages from taking a serving when the pot is empty
- let the cook refill the pot only when it is empty
- serialize access to the pot while servings are added or removed
- allow waiting goroutines to stop cleanly when their context is canceled

## Layout

- `dining_savages.go`: challenge implementation stub
- `dining_savages_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable dining savages tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./dining_savages
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./dining_savages/solutions
```

Benchmark the hint solutions:

```sh
go test ./dining_savages/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./dining_savages -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
