# Add Two Promises

This package models the LeetCode-style problem of waiting for two integer
promises and returning a promise for their sum.

## Layout

- `add_two_promises.go`: challenge stub
- `add_two_promises_test.go`: shared tests for the challenge
- `testsuite/`: reusable tests and benchmarks
- `solutions/`: hint implementations

## Testing

Run the challenge tests:

```sh
go test -tags challenge ./add_two_promises
```

Run the hint solutions:

```sh
go test ./add_two_promises/solutions
```
