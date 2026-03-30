# Extended Dining Hall

The extended dining hall problem adds a staging step before students sit down.

Each student:

1. gets food
2. may have to wait until it is socially safe to start dining
3. dines
4. may have to wait again before leaving
5. leaves

The synchronization rules are:

- a student may not invoke `dine` alone at an otherwise empty table
- a student may not invoke `leave` if that would leave exactly one diner alone
- another ready-to-eat student can release a student waiting to sit
- either a new arrival or the final diner can release a student waiting to leave
- the implementation must be reusable across many groups

## Layout

- `extended_dining_hall.go`: challenge implementation
- `extended_dining_hall_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable behavioral tests and benchmarks
- `solutions/`: reference implementations that also use the shared testsuite

## Testing

Run the challenge tests:

```sh
go test -tags challenge ./extended_dining_hall
```

Run the reference solutions:

```sh
go test ./extended_dining_hall/solutions
```

Benchmark the reference solutions:

```sh
go test ./extended_dining_hall/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./extended_dining_hall -run '^$' -bench .
```
