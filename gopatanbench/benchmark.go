package main

import (
    "sync"
    "github.com/toefel18/go-patan/statistics/lockbased"
    "fmt"
    "time"
    "flag"
    "os"
    "log"
    "runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
    flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
    Benchmrk(10, 20000)
    Benchmrk(100, 20000)
    Benchmrk(1000, 20000)
    Benchmrk(10, 200000)
    Benchmrk(100, 200000)
}

func Benchmrk(threads int64, itemsPerThread int64) {
    millisStart := currentTimeMillis()
    wg := sync.WaitGroup{}
    subject := lockbased.NewFacade(lockbased.NewStore())
    for i := int64(0); i < threads; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            sw := subject.StartStopwatch()
            defer subject.RecordElapsedTime("goroutine.duration", sw)
            for i := int64(0); i < itemsPerThread; i++ {
                subject.IncrementCounter("concurrency.counter")
                subject.AddSample("concurrency.sample", i)
            }
        }()
    }
    wg.Wait()
    snapshot := subject.Snapshot()
    expectedItems := threads * itemsPerThread
    if snapshot.Counters()["concurrency.counter"] != expectedItems {
        panic(fmt.Sprint(expectedItems, "counters expected, but got", snapshot.Counters()["concurrency.counter"]))
    }
    if snapshot.Durations()["goroutine.duration"].SampleCount() != threads {
        panic(fmt.Sprint("There should be",threads, "durations registered but got", snapshot.Durations()["goroutine.duration"].SampleCount()))
    }
    if snapshot.Samples()["concurrency.sample"].SampleCount() != threads * itemsPerThread{
        panic(fmt.Sprint(expectedItems, "samples expected but got", snapshot.Samples()["concurrency.sample"].SampleCount()))
    }
    millisEnd := currentTimeMillis()
    fmt.Println(threads, "threads with", itemsPerThread, "items took", (millisEnd - millisStart))
}

func currentTimeMillis() int64 {
    return time.Now().UnixNano() / time.Millisecond.Nanoseconds()
}
