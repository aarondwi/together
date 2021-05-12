# together
Easily batch your OLTP code, enjoy the performance.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Background
-------------------------
One of the generally accepted way to increase throughput for data loading is batching.
This increase in throughput is gotten from amortized cost of doing every step (network RTT, parsing, planning, locking, etc).
This pattern is much more prevalent in background ETL workflow, where lots of data ready at once.
(and lots of good tools, such as [dbt](https://www.getdbt.com/) and [spark](https://spark.apache.org/)).

Human themselves usually can't spot the difference between 5ms and 10-20ms (or even, 100ms).
So actually the option to batch interactive, OLTP-style requests to achieve more throughtput is interesting,
and most data store (database, queue, 3rd party API, etc) already have batching capability (such as `group commit`, `insert multiple`, `update join`, `select in`, etc).

Unfortunately, most business-logic code (READ: `almost all`) does not use this technique, and instead
do everything on per request basis. It has 2 implications:
    1. Number of transactions to handle is equal to the number of requests from users, which means it is hard
    to handle sudden surge, something easy to do with batching.
    2. While network indeed is getting faster, because of transaction and locking (important for data integrity)
    connection is held for relatively long time, causing latencies to add up from waiting, and that results
    in low-throughput.

This library is an attempt to help developers easily incorporate batching logic into business-level OLTP code,
easily achieving high throughput plus backpressure to handle sudden surge.

Installation
-------------------------

```bash
go get -u github.com/aarondwi/together
```

Usages
-------------------------

See the [tests](https://github.com/aarondwi/together/blob/main/engine_test.go) directly for the most up-to-date example.

Notes
-------------------------

1. By `batching`, it does not mean batch processor like [gobatch](https://github.com/MasterOfBinary/gobatch),
[spring batch](https://spring.io/projects/spring-batch), [dbt](https://www.getdbt.com/), [spark](https://spark.apache.org/), or anything like that. `Batching` here means aggregating multiple request into single request to backend, akin to deduplication.
2. This is designed to be used in business-level OLTP code, so it is not aiming to be *every-last-cpu-cycle* optimized
(in particular, this implementation use `interface{}`, until golang support generics).
If you have something that can be solved with this pattern,
but need a more optimized one, it is recommended to make something similar yourself.
3. The batching implementation waits on either number of message, or timeout (akin to [kafka](https://kafka.apache.org/)).
This is by design,
because we either want to batch for throughput, or for saving (if you call 3rd party APIs
which has rate-limiter/pay-per-call, but allow multiple message in each call).
Other complex implementation have their own downsides, such as:
    * waiting only if more than specific number of connection, like in [PostgreSQL](https://postgresqlco.nf/doc/en/param/commit_siblings/), gonna make the API harder (and weird) to be incorporated into business-level code.
    * getting from the queue and working the batch as fast as possible, like in [Tarantool](https://dzone.com/articles/asynchronous-processing-with-in-memory-databases-o) or [Aurora](https://www.semanticscholar.org/paper/Amazon-Aurora%3A-On-Avoiding-Distributed-Consensus-Verbitski-Gupta/fa4a2b8ab110472c6d8b1b19baa81af21800468b), may results in better throughput/latency overall, but not as useful for the saving goal.
4. This library will never include `panic` handling, because IMO, it is a bad practice. `panic` should only be used
when keep going is dangerous for integrity, and the best solution is to just **crash**.
If you (or library you are using) still insist to use `panic`, please catch it and return error instead.
5. Even though this library lets you easily get high throughput, be cautious with locking inside database.
One blocked row may block entire batch.

TODO
-------------------------

1. Add `joiner` implementation (waiting for multiple call at once)
2. CI with github actions
