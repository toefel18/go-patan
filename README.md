# ![patan-logo](go-patan.png)-patan a small library for gathering statistics.

Golang port of the java sampling library [patan](https://github.com/toefel18/patan) that provides: 
  - counters; keeping track of how many times *something* has taken place
  - sampling; collecting samples and describing their distribution
  - durations; measuring the duration of a task as a special case of sampling
  
The API is not identical to the java version, some methods are named differently. 

When serializing a snapshot to JSON, it will differ from the java version of patan because
occurrences are renamed to counters.  
