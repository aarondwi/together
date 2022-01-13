# together

Batch your call. Easily backpressure. Enjoy the performance.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Background

OLTP code has distinctive pattern from others, which is large number of small, simple requests.
Each of these requests are generally very cheap, but needs lots of overhead (usually networking cost to DB or other upstream APIs).
Calling these upstreams for every single request means paying everything multiple times (from contention, network latency, etc).
This means CPU and time are used more on the expensive stuff, instead of the more important business logic stuff.

Async patterns emerged to solve efficiency issue from older, thread-based pattern. But while they indeed saves on resources, they don't directly translate to better utilization. Biggest reasons for this are traditional protocol (the most famous one being HTTP and/or Postgres/MySQL protocol) are blocking protocols, so the asynchronous-ism just happen inside the instance, while all the expensive network-related stuff still done synchronously. Async runtimes also don't promote enough back-pressure and/or cancellation needed (see [here](https://lucumr.pocoo.org/2020/1/1/async-pressure/) for good contexts, and [its HN thread](https://news.ycombinator.com/item?id=21927427)), so actually they just make bottlenecks on DB/other upstreams happen faster than before.

Human themselves usually can't spot the difference between 5ms and 10, 50, 100, or even 200ms. So the option to trade bit of latencies for throughput is worth it.
Most data store (database, queue, 3rd party API, etc) already have batching capability (`group commit`, `insert multiple`, `update join`, `select in`, and their equivalents). This allows stateless tier to do some work (batching) for the upstream, basically changes I/O-heavy to CPU-heavy stuff, which is far easier to optimize.

Unfortunately, most business-logic code (READ: `almost all`) does not use this technique. There are some,
but only done per request basis, usually when the request has lots of data to insert/fetch at once.
It has 2 implications:

1. Number of transactions to handle is equal to the number of requests from users, which means it is hard
to handle sudden surge, something easy to do with batching.
2. While network hardware indeed is getting faster (albeit slowly), because of transaction and locking (for data integrity)
connection is held for relatively long time, causing latencies to add up from waiting, context-switches, etc. And that results in lower throughput than what is possible.

There is a prominent user of batching in OLTP scheme, that is, [GraphQL](graphql.org) with its [Dataloader](https://github.com/graphql/dataloader) pattern.
They can do it because each graphql request is basically a graph/tree request, meaning lots of data is ready to be queried at once. But it is still done on per request basis, which also means the previous 2 points still hold.

This library is an attempt to make it easier to combine separate individual requests, helping developers easily achieve high throughput plus effectively free backpressure ability to handle sudden surge.

## Installation

```bash
go get -u github.com/aarondwi/together
```

## Features

1. Small and clear codebase (~1000 LoC), excluding tests.
2. General enough to be used for any batching case, and can easily be abstracted higher.
3. Fast. On my test laptop, when callers submit 1 by 1, with workers simulating relatively fast network call by sleeping for 1-3ms, reaching ~1-2 million work/s. When callers submit a batch of 256, it reached ~6 million/s
4. Easy promise-like API (just use `Submit` or equivalent call), and all params will be available to batch worker. You just need to return the call with same key as the given parameters.
5. Circumvent single lock contention using `Cluster` implementation.
6. Optional background worker, so no goroutine creations on hot path (using tunable `WorkerPool`).
7. Waiting multiple results at once, to reduce latency (Using `Combiner` implementation).
8. Non-context and context variant available. (for timeout-based, cancellations, hedge-requests, etc)
9. Separating submitting and waiting results, to allow fire-and-forget cases.
10. Submit `Many` idiom, to directly put bunch of params with single lock. Useful especially for upstream services.

## Usages

To use this library, see the [engine](https://github.com/aarondwi/together/blob/main/engine/engine_test.go), [cluster](https://github.com/aarondwi/together/blob/main/cluster/cluster_test.go), and [combiner](https://github.com/aarondwi/together/blob/main/combiner/combiner_test.go) test files directly for the most up-to-date example.

For example how to write typical business logic as batch, please see [here](https://github.com/aarondwi/batch-logic-example), and those you should know and be careful of when designing batch logics, see [here](https://aarondwi.github.io/TogetherNotes)

## Notes

1. This is **NOT** a batch processor like [spring batch](https://spring.io/projects/spring-batch), [dbt](https://www.getdbt.com/), [spark](https://spark.apache.org/), or anything like that. `This library does combining/deduplicating/scatter-gather multiple request into (preferably) single request to backend, like how Facebook manages its [memcache's flow](https://www.mimuw.edu.pl/~iwanicki/courses/ds/2016/presentations/08_Pawlowska.pdf) or Quora with their [asynq](https://github.com/quora/asynq).
2. This is designed to be used in high level, business code, so it is not aiming to be *every-last-cpu-cycle* optimized (in particular, this implementation use `interface{}`, which is yet another pointer + allocation, until golang support generics).
If you have something that can be solved with this pattern, but need a more optimized one, it is recommended to make something similar yourself.
3. The batching implementation waits on either number of message (soft and hard limit) and timeout (akin to [kafka](https://kafka.apache.org/)'s `batch.size` and `linger.ms`), combined with [Smart Batching](http://mechanical-sympathy.blogspot.com/2011/10/smart-batching.html), to ensure we get the largest batch allowed.
This is by design, because we either want to batch for throughput or for saving (if you call 3rd party APIs
which has rate-limiter/pay-per-call, but allow multiple message in each call), but some API still have a hard limit on how many messages a single call can contain.
Other implementations have their own downsides, such as:
    * waiting only if more than specific number of connection, like in [PostgreSQL](https://postgresqlco.nf/doc/en/param/commit_siblings/), gonna make the API harder (and weird) to be incorporated into business-level code.
    * Naive smart batching, like in [Tarantool](https://dzone.com/articles/asynchronous-processing-with-in-memory-databases-o) or [Aurora](https://www.semanticscholar.org/paper/Amazon-Aurora%3A-On-Avoiding-Distributed-Consensus-Verbitski-Gupta/fa4a2b8ab110472c6d8b1b19baa81af21800468b), may results in better throughput and/or latency and/or utilization overall, but not as useful for the saving goal.
4. This library will never include `panic` handling, because IMO, it is a bad practice. `panic` should only be used when keep going is dangerous for integrity, and the best solution is to just **crash**.
If you (or a library you are using) still insist to use `panic`, please `recover` it and return error instead.
5. A batch here is **NOT**, and should never be considered as, a unit of `atomicity` (A in ACID). The goal is to let lots of data to be processed at once, increasing efficiency, but each content should be able to succeed or fail individually.
6. For now, there are no plans to support dynamic, adaptive setup (a la [Netflix adaptive concurrency limit](https://netflixtechblog.medium.com/performance-under-load-3e6fa9a60581)).
Besides cause this library gonna need more tuning (number of worker, batch size, waiting size, how to handle savings, etc) which makes it really really complex, together's Batch Buffering already absorbs most of the contention from requests, and upstream services easily become CPU bottlenecked. Adaptivity just gonna make CPU not operating at maximum available capacity.
Instead, allow the batch to be a bit bigger. This will already allow you to serve most requests even on heavy surge instead of a downtime
7. Even though the focus is for OLTP code, batch/ETL code can also use this when enriching result from multiple API
8. This library is much more general purpose than [golang's singleflight](https://pkg.go.dev/golang.org/x/sync/singleflight). You can implement singleflight on top of this library (albeit with bit different wait semantic). And while singleflight force you to share the resulting object (means don't mutate), this library does not (although you can)

## Setup Recommendations

1. For business logic setup, set normal large batch (~64-128 is good). For lots key-value access from a single requests, 256, 512 or more is good. The large number is more useful for backpressure purpose, and while bigger batch may have higher latency overall, it is better to still serve all with higher latency (of course, while still under acceptable human perception, e.g. 1-2s) than just straight downtime.
2. Not so much worker per `engine` instance (2-8 should be enough)
3. Waiting number to be at most the same as typical duration of a batch (if a full batch needs ~20ms, 10-20ms batch waiting time is good, getting good enough balance between latency, throughput, contention reduction via buffering, and call savings)
4. Separate `engine` and `cluster` instance for each needs (For example, placing order and getting item details are very different requests, with very different complexity and duration of requests). But, use same `Workerpool` instance (with quite large number of goroutines, e.g. >5-10K) to amortize all the waiting goroutines

For example, if an `engine` instance has a batch size of 128, 4 workers, full batch work time ~10ms (this is a rather slow one for batched key-value access, but make sense for more complex business logic), and assuming those batches keep filled cause of spike traffic, this engine instance can do `4 workers * 128/batch * (1 second / 10 ms)` = 51200 rps, which already far surpass most businesses' needs.

## Notes for benchmarks

We use 1 message per `Submit()` for the normal usage to mimic the outermost services, which need to combine many small messages. The `SubmitMany()` benchmarks use a batch of 256 to mimic upstream services/databases, which can receive batches from outer services, instead of one by one.

## Nice to have

1. Support for generic, once golang supports it
2. Workerpool to have rate-limited, max new goroutine per second and global max, so not fire-and-forget goroutines only, but amortized to a number of works
3. Cancellations for batch(?)
4. Actually implementing adaptivity:

    * Batch size hard limit is fixed. This is a hard requirement
    * Number of worker from 1. Most of the time is enough, try climbing slowly. Latency is the metric
    * Wait time can increase slowly, to cope with overload. But at this state, should worker only 1? -> reduce load
    * shared tracking context for multiple engines? separate?

5. Key-value based batch, as key-value use case is the most probable to be in need of high throughput
6. `Combiner` possibly return internal channel directly, letting user manage `resolutionHelper` in a pipelined way
7. Move Cluster and Combiner to config object (?)
8. Add fuzzy test
