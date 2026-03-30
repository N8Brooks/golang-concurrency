# Child Care

The child care problem models a center with one safety rule:

- there must always be at least one adult present for every three children

That means:

- adults may always enter the center
- a child may only enter when doing so would keep `children <= 3 * adults`
- a child may always leave
- an adult may only leave when doing so would still keep `children <= 3 * adults`

In this package, a child-care implementation must:

- allow adults to enter immediately
- block children until enough adult capacity is available
- block an adult that is trying to leave when leaving would violate the ratio
- avoid the classic deadlock where two adults split the three available child
  permits between them while trying to leave at the same time
- allow a waiting child to stop waiting if its context is canceled
- remain reusable across multiple waves of adults and children

## Layout

- `child_care.go`: challenge implementation stub
- `child_care_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable child-care tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./child_care
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./child_care/solutions
```

Benchmark the hint solutions:

```sh
go test ./child_care/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./child_care -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
