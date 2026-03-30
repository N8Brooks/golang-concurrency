# FIFO Barbershop

The FIFO barbershop problem is the sleeping barber problem with one additional
constraint: customers must be served in the order they enter the shop.

This package uses the textbook's `n` convention: `NewFIFOBarbershop(n)` sets
the maximum number of customers allowed inside the shop at once, including the
customer currently in the barber chair. That means the waiting room has
`n-1` chairs.

Each customer does one of two things:

1. enter the shop and eventually get a haircut
2. balk immediately if the shop is full

The barber repeatedly:

1. sleeps until a customer is available
2. invites the longest-waiting customer to the chair
3. cuts that customer's hair before starting the next haircut

In this package, a FIFO barbershop implementation must:

- satisfy all of the ordinary barbershop constraints
- serve customers in the same order they pass the turnstile
- still allow a sleeping barber or waiting customer to stop waiting if its context is canceled

The provided hint solution follows the textbook structure: each customer gets
its own semaphore, and the barber wakes customers by removing them from a FIFO
queue.

## Layout

- `fifo_barbershop.go`: challenge implementation stub
- `fifo_barbershop_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable FIFO-barbershop tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./fifo_barbershop
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./fifo_barbershop/solutions
```

Benchmark the hint solutions:

```sh
go test ./fifo_barbershop/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./fifo_barbershop -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
