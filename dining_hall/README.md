# Dining Hall

The dining hall problem models students eating at a shared table.

Each student:

1. starts dining
2. eventually becomes ready to leave once dining finishes
3. leaves the table

The synchronization requirement is that no student may ever be left dining
alone because everyone else left first. The only problematic state is one
student still dining while one other student is ready to leave.

In this package, a dining hall implementation must:

- allow students to dine concurrently
- delay a ready-to-leave student when their departure would leave exactly one student dining alone
- allow a new arrival or the final dining student to release that waiting student
- be reusable across many arrivals and departures

## Layout

- `dining_hall.go`: challenge implementation
- `dining_hall_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable dining hall tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

Run the challenge tests:

```sh
go test -tags challenge ./dining_hall
```

Run the hint solutions:

```sh
go test ./dining_hall/solutions
```

Benchmark the hint solutions:

```sh
go test ./dining_hall/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./dining_hall -run '^$' -bench .
```
