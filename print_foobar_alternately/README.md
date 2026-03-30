# Print FooBar Alternately

The print FooBar alternately problem has two threads that must print `foo` and
`bar` in alternating order.

Given `n`, one thread should call `printFoo` exactly `n` times and the other
should call `printBar` exactly `n` times, producing `foobar` repeated `n`
times.

In this package, a FooBar implementation must:

- alternate `foo` and `bar` correctly
- print exactly `n` copies of each token
- block `bar` until the corresponding `foo` has run

## Layout

- `print_foobar_alternately.go`: challenge implementation stub
- `print_foobar_alternately_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable FooBar tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./print_foobar_alternately
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./print_foobar_alternately/solutions
```

Benchmark the hint solutions:

```sh
go test ./print_foobar_alternately/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./print_foobar_alternately -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
