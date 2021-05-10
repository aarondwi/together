# together
Easily batch your API logic, enjoy the performance.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Background
-------------------------
One of the generally accepted way to increase throughput for data loading is batching.
This increase in throughput is gotten from amortized cost of doing every step (network RTT, parsing, planning, locking, etc).
This pattern is much more prevalent in background ETL workflow, where lots of data ready at once.
(and lots of good tools, such as [dbt](https://www.getdbt.com/) and [spark](https://spark.apache.org/)).

Human themselves usually can't spot the difference between 5ms and 10-20ms.
So actually the option to batch interactive, OLTP-style request to achieve more throughtput is interesting,
and most data store (database, queue, 3rd party API, etc) already have batching capability (such as `group commit`, `insert multiple`, `update join`, `select in`, etc).
Unfortunately, most business-logic code (READ: `almost all`) rarely uses this technique, and instead
do everything on per request basis. While network getting faster,
multiple calls also add latency up, and during that time connection + lock is held
(especially because the need for transaction), resulting in low-throughput.

This library is an attempt to help developer easily incorporate batching logic into business-level OLTP code.

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

1. By **batching**, it does not mean batch processor like [gobatch](https://github.com/MasterOfBinary/gobatch),
[spring batch](https://spring.io/projects/spring-batch), [dbt](https://www.getdbt.com/), [spark](https://spark.apache.org/), or anything like that. **Batching** here means aggregating multiple request into single request to backend, akin to deduplciation.
2. This is designed to be used in business-level OLTP code, so it is not aiming to be *every-last-cpu-cycle* optimized
(in particular, this implementation use `interface{}`, until golang support generics).
If you have something that can be solved with this pattern,
but need a more optimized one, it is recommended to code something similar yourself.
3. The batching implementation waits on either number of message, or timeout. This is by design,
because we either want to batch for throughput, or for saving (if you call 3rd party APIs which has rate-limiter,
but allow multiple message in each call). Other more complex implementation have their own downsides, such as:
    * waiting only if more than specific number of connection, like in [PostgreSQL](https://postgresqlco.nf/doc/en/param/commit_siblings/), gonna make the API harder (and weird) to be incorporated into business-level code.
    * getting from the queue and working the batch as fast as possible, like in [Tarantool](https://dzone.com/articles/asynchronous-processing-with-in-memory-databases-o), may results in better performance/latency overall, but not as useful for the saving goal.
4. This library will never include `panic` handling, because IMO, it is a bad practice. `panic` should only be used
when keep going is dangerous for integrity, and the best solution is to just **crash**.
If you want to use `panic`, please catch it and return error instead.

TODO
-------------------------

1. Add `joiner` implementation (waiting for multiple call at once)
2. CI with github actions
