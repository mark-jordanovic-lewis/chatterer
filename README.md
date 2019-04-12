# Chatterer Systems

A composer system is an internet IO handler. It's tasks are customisable via the thinkers it is holding in its []thinker slice and the conductors that facilitate communications between them (and other disparate packages).

- A composer is a http network node
  - it has a concurrently running system of internal communication and computation nodes
  - it's functionality is defined by the injected `worker` function types
    - `worker` takes a channel arg and 'returns' down the same channel
    - should have multiple duplicate workers for high traffic endpoint
  - i/o defined by `conductor` function types,
    -  routes data via channels between `worker` and other conductor functions

### Bootstrapping
Going to have to build some examples...

- User defines a set of Conductor and `worker` functions to model their solution
- functions are injected into Composure to form an async system of processes to handle
requests.
- it's basically a microservice

##### HTTP Req
- serveHTTP sends to server channel.
- at least one conductor listens on server channel
- data sent down channel to `worker` which computes result and

- should know about a bunch of clients to respond to
- should know about all the neighbours that it needs to talk to
- should be told how to handle inputs (whether from ) via it's thinker functions (written by the user-programmer)
- should have a set of conductor functions that direct communications between relevant channels (written by the user-programmer)
- should handle a set of http requests
  - ie/ be coherent with std REST APIs, externally
- Be able to consume external APIs via specialised composers
- The

## Components
#### Chatterer
- Chatterers should be flexible enough to instantiate other composers based in input
- Should handle http requests like a std Go Handler interface
- Can form the initial stage of some processing pipeline for a response

#### Thinker
#### Conductor
- at least one conductor should listen on the server channel

#### Init
- Initialise the environment
- Here is where you should pop all the composers that need initialising
#### There's a logger to deal with too


## Flow
#### Std Handler type
- HTTP req handled by composer
- ServeHTTP is as normal, possibly send logging data through conductor
- Response generated

####

## Todo
