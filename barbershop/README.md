# Barbershop

The barbershop problem models a single barber, a barber chair, and a waiting
room with limited capacity.

This package uses the textbook's `n` convention: `NewBarbershop(n)` sets the
maximum number of customers allowed inside the shop at once, including the
customer currently in the barber chair. That means the waiting room has
`n-1` chairs.

Each customer does one of two things:

1. enter the shop and eventually get a haircut
2. balk immediately if the shop is full

The barber repeatedly:

1. sleeps until a customer is available
2. invites exactly one customer to the chair
3. cuts that customer's hair before starting the next haircut

In this package, a barbershop implementation must:

- let a customer wait if the shop is not yet full
- call `balk` instead of waiting when the shop is full
- let a sleeping barber wait until a customer arrives
- ensure each `cutHair` runs concurrently with exactly one `getHairCut`
- ensure the barber does not start the next haircut before the current one finishes
- be reusable across many customers and haircuts
- allow a sleeping barber or waiting customer to stop waiting if its context is canceled

The provided hint solution follows the textbook structure: a scoreboard for
occupancy plus two rendezvous-style handshakes.

## Layout

- `barbershop.go`: challenge implementation stub
- `barbershop_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable barbershop tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./barbershop
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./barbershop/solutions
```

Benchmark the hint solutions:

```sh
go test ./barbershop/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./barbershop -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
