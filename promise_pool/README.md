# Promise Pool

This package models bounded-concurrency execution of asynchronous work. The
pool runs at most `n` promises at a time and resolves when all work has
finished.

## Layout

- `promise_pool.go`: challenge stub
- `promise_pool_test.go`: shared tests for the challenge
- `testsuite/`: reusable tests and benchmarks
- `solutions/`: hint implementations

## Testing

Run the challenge tests:

```sh
go test -tags challenge ./promise_pool
```

Run the hint solutions:

```sh
go test ./promise_pool/solutions
```
