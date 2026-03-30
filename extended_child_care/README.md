# Extended Child Care

The extended child care problem uses the same safety rule as the basic child
care problem:

- there must always be at least one adult present for every three children

The difference is in how an adult waiting to leave is treated.

In the basic semaphore solution, a departing adult may effectively reserve
three child permits while waiting, which can block children even when their
entry would still be legal. In this extended version, an adult that is waiting
to leave still counts as present until it actually exits, so children should not
wait unnecessarily.

That means:

- adults may always enter the center
- a child may only enter when doing so would keep `children <= 3 * adults`
- a child may always leave
- an adult may only leave when doing so would still keep `children <= 3 * adults`
  after that adult has actually exited
- a waiting adult must not prevent additional children from entering when
  those children would still satisfy the ratio

In this package, an extended child-care implementation must:

- allow adults to enter immediately
- block children until enough adult capacity is available
- block an adult that is trying to leave when leaving would violate the ratio
- allow children to continue entering while an adult is waiting to leave, as
  long as the current number of adults still makes their entry legal
- allow a waiting child to stop waiting if its context is canceled
- remain reusable across multiple waves of adults and children

## Layout

- `extended_child_care.go`: challenge implementation stub
- `extended_child_care_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable extended-child-care tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./extended_child_care
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./extended_child_care/solutions
```

Benchmark the hint solutions:

```sh
go test ./extended_child_care/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./extended_child_care -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
