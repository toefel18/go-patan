# ![patan-logo](go-patan.png)-patan a small library for gathering statistics.

[![Build Status](https://travis-ci.org/toefel18/go-patan.svg?branch=master)](https://travis-ci.org/toefel18/go-patan) [![codecov.io](https://codecov.io/github/toefel18/go-patan/coverage.svg?branch=master "coverage")](https://codecov.io/github/toefel18/go-patan)

Golang port of the java sampling library [patan](https://github.com/toefel18/patan) that provides: 
  - counters; keeping track of how many times *something* has taken place
  - sampling; collecting samples and describing their distribution
  - durations; measuring the duration of a task as a special case of sampling
  
The API is not identical to the java version, some methods are named differently. 

When serializing a snapshot to JSON, it should produce the same output as the java version.

Start with:
```
go get github.com/toefel18/go-patan
```
  
Usage:
```go
package main

import (
    "github.com/toefel18/go-patan/statistics"
    "time"
    "fmt"
    "encoding/json"
)

func main() {
    stopwatch := statistics.StartStopwatch()
    time.Sleep(2 * time.Second)

    statistics.AddSample("mem.allocations", 167.0334)
    statistics.AddSample("mem.allocations", 111.9216)
    statistics.AddSample("mem.allocations", 133.4686)
    statistics.AddSample("mem.collects", 2)
    statistics.AddToCounter("active.sessions", 132)
    statistics.DecrementCounter("active.sessions")
    statistics.RecordElapsedTime("my.heavy.operation", stopwatch)   // A
    snapshot := statistics.Snapshot()                               // B

    duration, exists := snapshot.Durations()["my.heavy.operation"]
    if exists {
        fmt.Print("duration has sample count of ", duration.SampleCount()) // samplecount = 1
    } else {
        fmt.Print("duration does not exist")
    }

    // JSONizing a snapshot that can be published through REST
    if json, err := json.MarshalIndent(snapshot, "", "  "); err == nil {
        fmt.Println(string(json))
    }
}
```
The output will be:
```json
{
  "timestampStarted": 1480792554683,
  "timestampTaken": 1480792556683,
  "durations": {
    "my.heavy.operation": {
      "sampleCount": 1,
      "minimum": 2000.129942,
      "maximum": 2000.129942,
      "mean": 2000.129942,
      "stdDeviation": 0
    }
  },
  "counters": {
    "active.sessions": 131
  },
  "samples": {
    "mem.allocations": {
      "sampleCount": 3,
      "minimum": 111.9216,
      "maximum": 167.0334,
      "mean": 137.47453333333334,
      "stdDeviation": 27.77342707001304
    },
    "mem.collects": {
      "sampleCount": 1,
      "minimum": 2,
      "maximum": 2,
      "mean": 2,
      "stdDeviation": 0
    }
  }
}
```

The default statistics instance, which is always and directly available when using patan (see code example, `statistics` 
is the default available instance) uses a lock-based implementation. For more info on implementations, read on.

## API implementations
There are two implementations, one based on channels and another based on locks. Micro benchmarks 
were in favor of the lock-based implementation and it also shows direct consistency whereas the 
channel-based implementation is eventual consistent. 

#### Lock-based  implementation
The implementation is direct consistent.

#### Channel-based  implementation
This implementation is eventual consitent and the invocation order does not necessarily reflect the 
processing order! Users should not depend on that. Consider the example code above when it would
use a channel-based implementation:

It is possible (and even likely) that snapshot does not have my.heavy.operation yet, meaning that
B is executed earlier than A! This differs from the java version of Patan and is a consequence
of the non-blocking setup with channels. This is OK because patan is meant to give insight in
the distribution of data over a longer period of time, not for individual measurements.
git 