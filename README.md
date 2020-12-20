# l4-loadbalancer
 A simple L4 loadbalancer written in Go

 ## Features

 - Layer 4 load balancing
 - Passive healthchecks
 - Roundrobin selection
 - Scalable connection handling due to go routines
 - Race free due to use of mutexes and atomic operations

## TODO

- Add leastconn algorithm
- Add service discovery
- Add ability to define multiple listeners
- Read from YAML config file

