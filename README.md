# Implementation of a MM1 queue in Go
An MM1 queue is a model with a single server, arrivals determined by a Poisson process and job service time with exponential distribution.

## Structure of the implementation
The core of the implementation consists in a cycle repeated as many times as the customers we want to process. the cycle consists in:
- determining if the next event will be an arrive or a departure
- invoke the right function (arrive/departure)
- update statistical counters for later calculations

There are also an initialization function and a report function with self-explanatory goals.

## Arrive / Departure functions

After the detection of the next event type the program:

### In case of an arrive

- schedules next arrive 
- if the server is already BUSY it adds a customere in queue with the right arrival time
- if the server is empty the server  is set to busy and we schedule the next departure time
  
### In case of a departure

- if the number of customers in queue is 0 the server is set to NOT BUSY, we schedule the next departure to an large time to ensure an arrive as next event
- if there are customers in queue we decrement the number of customers and schedule the next departure time
- the customer served counter is incremented

## statistical counters
At every iteration of the cycle we keep count of some statistics for future use.

- the area underlying the number of people at any time in the queue. (the sum of the rectangles determined by the number of people in queue in every slice of time)
- the area underlying the server utilization in time.
  



