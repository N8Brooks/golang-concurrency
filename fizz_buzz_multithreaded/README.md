# Fizz Buzz Multithreaded

The multithreaded FizzBuzz problem coordinates four worker threads over the
numbers `1..n`.

One worker prints `fizz`, one prints `buzz`, one prints `fizzbuzz`, and one
prints plain numbers. Together they must produce the exact ordered FizzBuzz
sequence.

In this package, an implementation must:

- print each position exactly once
- route each position to the correct worker
- preserve the exact overall order from `1` through `n`
- allow the non-number workers to start first and wait for their turns

## Layout

- `fizz_buzz_multithreaded.go`: challenge implementation stub
- `fizz_buzz_multithreaded_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable behavioral tests and benchmarks
- `solutions/`: reference implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./fizz_buzz_multithreaded
```

That is expected to fail until the stub is implemented.

Run the reference solutions:

```sh
go test ./fizz_buzz_multithreaded/solutions
```

Benchmark the reference solutions:

```sh
go test ./fizz_buzz_multithreaded/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./fizz_buzz_multithreaded -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
