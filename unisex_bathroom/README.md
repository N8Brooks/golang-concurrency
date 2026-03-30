# Unisex Bathroom

The unisex bathroom problem models a shared bathroom with two constraints.

Men and women may use the bathroom, but they may not be inside at the same
time, and there may never be more than three people inside.

In this package, a unisex bathroom implementation must:

- allow multiple men to use the bathroom concurrently when no women are inside
- allow multiple women to use the bathroom concurrently when no men are inside
- prevent men and women from being in the bathroom at the same time
- prevent more than three people from being in the bathroom at once
- allow a waiting caller to stop waiting if its context is canceled
- remain reusable across multiple uses

## Layout

- `unisex_bathroom.go`: challenge implementation stub
- `unisex_bathroom_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable unisex-bathroom tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./unisex_bathroom
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./unisex_bathroom/solutions
```

Benchmark the hint solutions:

```sh
go test ./unisex_bathroom/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./unisex_bathroom -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
