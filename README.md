# ![patan-logo](go-patan.png)-patan a small library for gathering metrics.

[![Build Status](https://travis-ci.org/toefel18/go-patan.svg?branch=master)](https://travis-ci.org/toefel18/go-patan) [![codecov.io](https://codecov.io/github/toefel18/go-patan/coverage.svg?branch=master "coverage")](https://codecov.io/github/toefel18/go-patan) [![Go Report Card](https://goreportcard.com/badge/github.com/toefel18/go-patan)](https://goreportcard.com/report/github.com/toefel18/go-patan)

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
    "github.com/toefel18/go-patan/metrics"
    "time"
    "fmt"
    "encoding/json"
)

func main() {
    stopwatch := metrics.StartStopwatch()
    time.Sleep(2 * time.Second)

    metrics.AddSample("mem.allocations", 167.0334)
    metrics.AddSample("mem.allocations", 111.9216)
    metrics.AddSample("mem.allocations", 133.4686)
    metrics.AddSample("mem.collects", 2)
    metrics.AddToCounter("active.sessions", 132)
    metrics.DecrementCounter("active.sessions")
    metrics.RecordElapsedTime("my.heavy.operation", stopwatch)   // A
    snapshot := metrics.Snapshot()                               // B

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

`metrics` is the default metrics instance, which is always and directly available when using patan. See code example, `metrics`
is the default available instance.
