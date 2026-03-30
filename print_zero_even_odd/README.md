# Print Zero Even Odd

The print zero even odd problem has three workers:

- `Zero`, which prints `0`
- `Odd`, which prints odd numbers in increasing order
- `Even`, which prints even numbers in increasing order

Together they must produce the sequence:

- `010203...0n`

In this package, an implementation must:

- print exactly one `0` before each number from `1` through `n`
- let `Odd` print only odd numbers in increasing order
- let `Even` print only even numbers in increasing order
- produce the full sequence in the correct order for any valid `n`

## Layout

- `zero_even_odd.go`: challenge implementation stub
- `zero_even_odd_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./print_zero_even_odd
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./print_zero_even_odd/solutions
```

Benchmark the hint solutions:

```sh
go test ./print_zero_even_odd/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./print_zero_even_odd -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
