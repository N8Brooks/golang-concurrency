# Cigarette Smokers

The cigarette smokers problem has one agent and three smokers.

The agent repeatedly places two different ingredients on the table:

- tobacco
- paper
- match

Each smoker has an infinite supply of exactly one ingredient. The smoker with
the complementary ingredient should pick up the two available ingredients, make
a cigarette, signal the agent, and then smoke.

In this package, a cigarette smokers implementation must:

- wake exactly the smoker that can complete the current set of ingredients
- avoid waking smokers that still cannot proceed
- allow the same smokers and agent to participate across multiple rounds
- stop cleanly when the context is canceled

This package models the interesting version of the problem: the agent interface
is fixed and the coordination logic belongs entirely to the smokers side.

## Layout

- `cigarette_smokers.go`: challenge implementation stub
- `cigarette_smokers_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable cigarette smokers tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./cigarette_smokers
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./cigarette_smokers/solutions
```

Benchmark the hint solutions:

```sh
go test ./cigarette_smokers/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./cigarette_smokers -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
