# Throttle

This package models the LeetCode `Throttle` problem: invoke immediately, then
coalesce calls during the throttle window into one trailing call using the
latest arguments.

## Testing

Run the challenge tests:

```sh
go test -tags challenge ./throttle
```

Run the hint solutions:

```sh
go test ./throttle/solutions
```
