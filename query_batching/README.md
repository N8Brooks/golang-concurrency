# Query Batching

This package models the LeetCode `Query Batching` problem: execute the first
query immediately, then batch subsequent keys that arrive within the throttle
window.

## Testing

Run the challenge tests:

```sh
go test -tags challenge ./query_batching
```

Run the hint solutions:

```sh
go test ./query_batching/solutions
```
