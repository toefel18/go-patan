# ![patan-logo](go-patan.png)-patan a small library for gathering statistics.

Golang port of the java sampling library [patan](https://github.com/toefel18/patan) that provides: 
  - counters; keeping track of how many times *something* has taken place
  - sampling; collecting samples and describing their distribution
  - durations; measuring the duration of a task as a special case of sampling
  
The API is not identical to the java version, some methods are named differently. 

When serializing a snapshot to JSON, it will differ from the java version of patan because
occurrences are renamed to counters.  


**USAGE NOTE**: Invocation order does not necessarily reflect the processing order! Users should
 not depend on that. Consider the following example:

```go
package main
import (
  "time"
  "fmt"
  "github.com/toefel18/go-patan/statistics"
)
func main() {
    stopwatch := statistics.StartStopwatch()
    time.Sleep(2 * time.Second)                 
                                 
    statistics.RecordElapsedTime("my.heavy.operation", stopwatch)   // A
    snapshot := statistics.Snapshot()                               // B
    
    duration, exists := snapshot.Durations()["my.heavy.operation"]
    if exists {
        fmt.Print("duration has sample count of ", duration.SampleCount()) // samplecount = 1
    } else {
        fmt.Print("duration does not exist")
    }
}
```

It is possible (and even likely) that snapshot doesn't have my.heavy.operation yet, meaning that
B is executed earlier than A! This differs from the java version of Patan and is a consequence
of the non-blocking setup with channels. This is OK because patan is meant to give insight in
the distribution of data over a longer period of time, not for individual measurements.