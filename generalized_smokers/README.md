# Generalized Smokers

The generalized smokers problem removes the agent handshake from the classic
smokers problem.

Instead of waiting for a smoker after placing ingredients, the agent can keep
adding more tobacco, paper, and matches. That means ingredients can accumulate
on the table, and the solution has to track counts rather than simple
presence/absence flags.

In this package, an implementation must:

- consume ingredients as they arrive over time
- wake the smoker whose own ingredient completes a pair on the table
- handle accumulated supplies correctly when multiple ingredients pile up
- let the long-running goroutines stop cleanly when their context is canceled

## Layout

- `generalized_smokers.go`: challenge implementation stub
- `generalized_smokers_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./generalized_smokers
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./generalized_smokers/solutions
```

Benchmark the hint solutions:

```sh
go test ./generalized_smokers/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./generalized_smokers -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
