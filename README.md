# together
Easily batch your OLTP code, enjoy the performance.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Background
-------------------------
One of the generally accepted way to increase throughput for data loading is batching.
Beside allowing us to amortize every step's cost(network RTT, parsing, planning, locking, etc),
we also gain capability to do advanced planning (like in [Calvin](http://cs.yale.edu/homes/thomson/publications/calvin-sigmod12.pdf)).
This pattern is much more prevalent in background ETL workflow, where lots of data is ready at once.
(and lots of good tools, such as [dbt](https://www.getdbt.com/) and [spark](https://spark.apache.org/)).

Human themselves usually can't spot the difference between 5ms and 10, 50, 100, or even 200ms.
So actually the option to batch interactive, OLTP-style requests to achieve more throughtput is interesting,
and most data store (database, queue, 3rd party API, etc) already have batching capability (such as `group commit`, `insert multiple`, `update join`, `select in`, etc).

Unfortunately, most business-logic code (READ: `almost all`) does not use this technique. There are some,
but only done per request basis, usually when the request has lots of data to insert/fetch at once.
It has 2 implications:

1. Number of transactions to handle is equal to the number of requests from users, which means it is hard
to handle sudden surge, something easy to do with batching.
2. While network indeed is getting faster, because of transaction and locking (for data integrity)
connection is held for relatively long time, causing latencies to add up from waiting, and that results
in low-throughput.

Actually, there is a prominent user of batching in OLTP scheme, that is, [GraphQL](graphql.org) with its [Dataloader](https://github.com/graphql/dataloader) pattern.
They can do it because each graphql request is basically a graph/tree request, meaning lots of data is ready to be queried at once. But it is still done on per request basis, which also means the previous 2 points still hold.

In my opinion, the main reason batching is not usually done on OLTP scheme is there are no good libraries to help business app developers gather data across requests easily.
It is a *tricky* problem, even for graphql library maintainer, which *only* do batching on per request basis, see [here](https://productionreadygraphql.com/blog/2020-05-21-graphql-to-sql/).
This library is an attempt to solve that, helping developers easily achieve high throughput plus backpressure ability to handle sudden surge.

Installation
-------------------------

```bash
go get -u github.com/aarondwi/together
```

Features
-------------------------

1. Small codebase (<1000 LoC).
2. Fast. On 2 years old laptop with Intel Core-i7 8550u, with batch worker simulating network call by sleeping for 2ms, `Cluster` reaching ~2 million invocation/s.
3. Easy API (just use `Submit` or equivalent call), and all params will be available to batch worker. You just need to return the call with same key as parameters.
4. Circumvent single lock contention using `Cluster` implementation.
5. Background waiting worker, so no goroutine creations on hot path (using tunable `WorkerPool`).
6. Waiting multiple results at once, to reduce latency (Using `Combiner` implementation).
7. Non-context and context variant available

Usages
-------------------------

See the [engine](https://github.com/aarondwi/together/blob/main/engine/engine_test.go), [cluster](https://github.com/aarondwi/together/blob/main/cluster/cluster_test.go), and [combiner](https://github.com/aarondwi/together/blob/main/combiner/combiner_test.go) test files directly for the most up-to-date example.

Notes
-------------------------

1. By `batching`, it does not mean a batch processor like [gobatch](https://github.com/MasterOfBinary/gobatch),
[spring batch](https://spring.io/projects/spring-batch), [dbt](https://www.getdbt.com/), [spark](https://spark.apache.org/), or anything like that. `Batching` here means aggregating multiple request into single request to backend, like how facebook manages its [memcache's flow](https://www.mimuw.edu.pl/~iwanicki/courses/ds/2016/presentations/08_Pawlowska.pdf).
2. This is designed to be used in business-level OLTP code, so it is not aiming to be *every-last-cpu-cycle* optimized
(in particular, this implementation use `interface{}`, until golang support generics).
If you have something that can be solved with this pattern,
but need a more optimized one, it is recommended to make something similar yourself.
3. The batching implementation waits on either number of message, or timeout (akin to [kafka](https://kafka.apache.org/)).
This is by design, because we either want to batch for throughput, or for saving (if you call 3rd party APIs
which has rate-limiter/pay-per-call, but allow multiple message in each call).
Other complex implementation have their own downsides, such as:
    * waiting only if more than specific number of connection, like in [PostgreSQL](https://postgresqlco.nf/doc/en/param/commit_siblings/), gonna make the API harder (and weird) to be incorporated into business-level code.
    * getting from the queue and working the batch as fast as possible, like in [Tarantool](https://dzone.com/articles/asynchronous-processing-with-in-memory-databases-o) or [Aurora](https://www.semanticscholar.org/paper/Amazon-Aurora%3A-On-Avoiding-Distributed-Consensus-Verbitski-Gupta/fa4a2b8ab110472c6d8b1b19baa81af21800468b), may results in better throughput/latency overall, but not as useful for the saving goal.
4. This library will never include `panic` handling, because IMO, it is a bad practice. `panic` should only be used
when keep going is dangerous for integrity, and the best solution is to just **crash**.
If you (or library you are using) still insist to use `panic`, please catch it and return error instead.
5. Even though this library lets you easily batch your requests, be cautious with record locking/isolation level.
One conflicting record may block/rollback entire batch.
