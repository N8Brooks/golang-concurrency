# Debounce

This package models the debounce helper: repeated calls within the debounce
window collapse into a single delayed invocation using the latest arguments.

## Layout

- `debounce.go`: challenge stub
- `debounce_test.go`: shared tests for the challenge
- `testsuite/`: reusable tests and benchmarks
- `solutions/`: hint implementations

## Testing

Run the challenge tests:

```sh
go test -tags challenge ./debounce
```

Run the hint solutions:

```sh
go test ./debounce/solutions
```
