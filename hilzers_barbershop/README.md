# Hilzer's Barbershop

Hilzer's barbershop problem extends the classic sleeping-barber setup.

The shop has three barbers, a sofa with four seats, standing room for more
customers, a hard cap of twenty customers in the shop, and a single cash
register. Customers must move through the shop in order, the standing queue
and sofa queue are FIFO, up to three haircuts may happen at once, and payment
must be serialized.

In this package, an implementation must:

- reject customers when the shop is already at capacity
- move customers from standing room to the sofa in FIFO order
- move customers from the sofa to barber chairs in FIFO order
- allow up to three concurrent haircuts
- ensure payment is accepted one customer at a time
- ensure a customer exits only after payment is accepted

## Layout

- `hilzers_barbershop.go`: challenge implementation stub
- `hilzers_barbershop_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./hilzers_barbershop
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./hilzers_barbershop/solutions
```

Benchmark the hint solutions:

```sh
go test ./hilzers_barbershop/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./hilzers_barbershop -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
