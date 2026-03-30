# No-Starve Unisex Bathroom

The no-starve unisex bathroom problem adds a fairness constraint to the basic
unisex bathroom problem.

Men and women may use the bathroom, but they may not be inside at the same
time, there may never be more than three people inside, and a long stream of
one sex should not starve the other.

In this package, a no-starve unisex bathroom implementation must:

- allow multiple men to use the bathroom concurrently when no women are inside
- allow multiple women to use the bathroom concurrently when no men are inside
- prevent men and women from being in the bathroom at the same time
- prevent more than three people from being in the bathroom at once
- prevent later arrivals from bypassing a queued entrant of the opposite sex
- allow a waiting caller to stop waiting if its context is canceled
- remain reusable across multiple uses

## Layout

- `no_starve_unisex_bathroom.go`: challenge implementation stub
- `no_starve_unisex_bathroom_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable no-starve unisex-bathroom tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./no_starve_unisex_bathroom
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./no_starve_unisex_bathroom/solutions
```

Benchmark the hint solutions:

```sh
go test ./no_starve_unisex_bathroom/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./no_starve_unisex_bathroom -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
