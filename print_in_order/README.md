# Print In Order

The print-in-order problem has three operations:

1. `first`
2. `second`
3. `third`

Callers may invoke the corresponding methods in any order and from different
goroutines, but the callbacks must still run in the fixed order
`first -> second -> third`.

In this package, a print-in-order implementation must:

- allow `First`, `Second`, and `Third` to be invoked in any order
- ensure `second` does not run before `first`
- ensure `third` does not run before `second`
- be reusable across many rounds
- allow a waiting caller to stop waiting if its context is canceled

## Layout

- `print_in_order.go`: challenge implementation stub
- `print_in_order_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable print-in-order tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./print_in_order
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./print_in_order/solutions
```

Benchmark the hint solutions:

```sh
go test ./print_in_order/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./print_in_order -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
