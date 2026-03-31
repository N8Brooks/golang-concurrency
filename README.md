# Concurrency Problems

This repository collects concurrency exercises in a shared template format:

- challenge implementations are opt-in with `-tags challenge`
- reusable behavioral tests live in each package's `testsuite/`
- working hint implementations live in each package's `solutions/`

## The Little Book of Semaphores

> Many of the synchronization exercises in this repo are adapted from Allen B.
> Downey's _The Little Book of Semaphores_. Chapters 1 and 2 are background;
> the problem index below starts where the book's problem-oriented table of
> contents starts and follows the same chapter structure.

### 3. Basic synchronization patterns

- 3.3 Rendezvous: [`rendezvous`](./rendezvous)
- 3.4 Mutex: [`mutex`](./mutex)
- 3.5 Multiplex: [`multiplex`](./multiplex)
- 3.6 Barrier: [`barrier`](./barrier)
- 3.7 Reusable barrier: [`barrier`](./barrier)
- 3.8 Queue: [`queue`](./queue)
- 3.8 Exclusive queue: [`exclusive_queue`](./exclusive_queue)

### 4. Classical synchronization problems

- 4.1 Producer-consumer problem: [`producer_consumer`](./producer_consumer)
- 4.1 Finite buffer producer-consumer: [`finite_buffer_producer_consumer`](./finite_buffer_producer_consumer)
- 4.2 Readers-writers problem: [`readers_writers`](./readers_writers)
- 4.2 No-starve readers-writers: [`no_starve_readers_writers`](./no_starve_readers_writers)
- 4.2 Writer-priority readers-writers: [`writers_priority_readers_writers`](./writers_priority_readers_writers)
- 4.3 No-starve mutex: [`no_starve_mutex`](./no_starve_mutex)
- 4.4 Dining philosophers: [`dining_philosophers`](./dining_philosophers)
- 4.5 Cigarette smokers problem: [`cigarette_smokers`](./cigarette_smokers)
- 4.5 Generalized Smokers Problem: [`generalized_smokers`](./generalized_smokers)

### 5. Less classical synchronization problems

- 5.1 The dining savages problem: [`dining_savages`](./dining_savages)
- 5.2 The barbershop problem: [`barbershop`](./barbershop)
- 5.3 The FIFO barbershop: [`fifo_barbershop`](./fifo_barbershop)
- 5.4 Hilzer's Barbershop problem: [`hilzers_barbershop`](./hilzers_barbershop)
- 5.5 The Santa Claus problem: [`santa_claus`](./santa_claus)
- 5.6 Building H2O: [`building_h2o`](./building_h2o)
- 5.7 River crossing problem: [`river_crossing`](./river_crossing)
- 5.8 The roller coaster problem: [`roller_coaster`](./roller_coaster)
- 5.8 Multi-car Roller Coaster problem: [`multi_car_roller_coaster`](./multi_car_roller_coaster)

### 6. Not-so-classical problems

- 6.1 The search-insert-delete problem: [`search_insert_delete`](./search_insert_delete)
- 6.2 The unisex bathroom problem: [`unisex_bathroom`](./unisex_bathroom)
- 6.2 No-starve unisex bathroom problem: [`no_starve_unisex_bathroom`](./no_starve_unisex_bathroom)
- 6.3 Baboon crossing problem: [`baboon_crossing`](./baboon_crossing)
- 6.4 The Modus Hall Problem: [`modus_hall`](./modus_hall)

### 7. Not remotely classical problems

- 7.1 The sushi bar problem: [`sushi_bar`](./sushi_bar)
- 7.2 The child care problem: [`child_care`](./child_care)
- 7.2 Extended child care problem: [`extended_child_care`](./extended_child_care)
- 7.3 The room party problem: [`room_party`](./room_party)
- 7.4 The Senate Bus problem: [`senate_bus`](./senate_bus)
- 7.5 The Faneuil Hall problem: [`faneuil_hall`](./faneuil_hall)
- 7.5 Extended Faneuil Hall problem: [`extended_faneuil_hall`](./extended_faneuil_hall)
- 7.6 Dining Hall problem: [`dining_hall`](./dining_hall)
- 7.6 Extended Dining Hall problem: [`extended_dining_hall`](./extended_dining_hall)

## LeetCode Concurrency Problems

The following LeetCode concurrency problems are also represented in this repo:

- [`bounded_blocking_queue`](./bounded_blocking_queue)
- [`fizz_buzz_multithreaded`](./fizz_buzz_multithreaded)
- [`print_foobar_alternately`](./print_foobar_alternately)
- [`print_in_order`](./print_in_order)
- [`print_zero_even_odd`](./print_zero_even_odd)
- [`traffic_light`](./traffic_light)
- [`web_crawler`](./web_crawler)

## LeetCode Async Problems

The following LeetCode-derived async/concurrency problems now follow the same
template architecture:

- [`promise_all`](./promise_all)
- [`add_two_promises`](./add_two_promises)
- [`promise_time_limit`](./promise_time_limit)
- [`promise_pool`](./promise_pool)
- [`sleep`](./sleep)
- [`debounce`](./debounce)
- [`throttle`](./throttle)
- [`cache_with_time_limit`](./cache_with_time_limit)
- [`event_emitter`](./event_emitter)
- [`timeout_cancellation`](./timeout_cancellation)
- [`interval_cancellation`](./interval_cancellation)
- [`query_batching`](./query_batching)
- [`promisify`](./promisify)
- [`custom_interval`](./custom_interval)
- [`delay_all`](./delay_all)
- [`promise_all_settled`](./promise_all_settled)
