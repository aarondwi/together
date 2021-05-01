# together
Easily batch your API logic, enjoy the performance.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Background
-------------------------
Databases that guarantees durability (D in ACID) are usually considered to be slow.
This is because they need to call `fsync`, which is really slow, latency wise (around 1-2ms).
Calling it once every new data will not have nice throughput.
That is why most databases (RDBMSes, etcd, mongo, etc) will batch the `fsync` call,
so the latency is amortized among many calls, increasing throughput.
It is also shown by a paper called [bp-wrapper](http://ranger.uta.edu/~sjiang/pubs/papers/ding-09-BP-Wrapper.pdf)
shows that by batching (and prefetching, but this one is not our concern here) can
reduce any contention to an acceptable value.

Even though most database (store?) already uses batching, but most business-logic code
(basically, almost everything backend-related) rarely uses this technique, which in turn
causing databases to do every step (network RTT, parsing, planning, etc) multiple times
before the business code can return response. While network getting faster,
multiple calls also add latency up, and during that time connection is held, resulting in low-throughput.

This library is an attempt to easily incorporate batching logic into business-level code.

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

1. This is designed to be used in business-level code, so it is not aiming to be used
for lower-level code that needs more optimized version (in particular, this implementation use `interface{}`).
If you wanted to use this, it is preferable to code something similar yourself, specific for your use case, based on this implementation.
2. The batching implementation waits on either number of message, or timeout. This is by design,
because we either want to batch for throughput, or for saving (if you call 3rd party APIs which has rate-limiter,
but allow multiple message in each call). Other more complex implementation have their own downsides:
    * waiting only if more than specific number of connection, like in [PostgreSQL](https://postgresqlco.nf/doc/en/param/commit_siblings/), gonna make the API harder (and weird) to be incorporated into business-level code.
    * getting from the queue and working the batch as fast as possible, like in [Tarantool](https://dzone.com/articles/asynchronous-processing-with-in-memory-databases-o), may results in better performance/throughtput overall, but not as useful for the saving goal.
3. This library will never include `panic` handling, because IMO, it is a bad practice. `panic` should only be used
when keep going is dangerous for integrity, and the best solution is to just *crash*.

Alternative
-------------------------

1. [gobatch](https://github.com/MasterOfBinary/gobatch)
2. [async/batch](https://github.com/grab/async/blob/master/batch.go)

TODO
-------------------------

1. Add `cluster` implementation (multi-engine scalability)
2. Add `joiner` implementation (waiting for multiple call at once)
3. CI with github actions
