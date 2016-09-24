# ![patan-logo](go-patan.png) GO-patan, a small library for gathering statistics.

Golang port of the java sampling library [patan](https://github.com/toefel18/patan) that provides: 
  - counters; keeping track of how many times *something* has taken place
  - sampling; collecting samples and describing their distribution
  - durations; measuring the duration of a task as a special case of sampling
  
The library provides an API and comes with a default implementation safe to be used in a multi-threaded environment.
