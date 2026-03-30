# Sushi Bar

The sushi bar problem has a fixed number of seats.

If a customer arrives while a seat is available, they may sit down immediately.
But if a customer arrives when the bar is full, they must wait until the entire
current party leaves before sitting down.

In this package, a sushi bar implementation must:

- let customers sit immediately while seats remain available
- block later arrivals once the bar becomes full
- keep blocked arrivals waiting until the entire current party leaves
- ensure at most `n` customers are dining at once
- be reusable across many parties
- allow a waiting customer to stop waiting if its context is canceled

The provided semaphore-based implementation follows the textbook scoreboard
approach.

## Layout

- `sushi_bar.go`: challenge implementation
- `sushi_bar_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable sushi bar tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

Run the challenge tests:

```sh
go test -tags challenge ./sushi_bar
```

Run the hint solutions:

```sh
go test ./sushi_bar/solutions
```

Benchmark the hint solutions:

```sh
go test ./sushi_bar/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./sushi_bar -run '^$' -bench .
```
