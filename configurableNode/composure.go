package configurableNode

import (
  // Heavy Lifting packages
  "net/http"

  // Std utilities
  "fmt"
)

// ================ \\
// Server Functions \\
// ================ \\

// Microservice Action One \\
// ----------------------- \\
type Action func(*Composure, string) (action func(chan []string))// maybe []byte would be better so JSON can be passed
type Conductor func(*Composure, string) (conductor func([]chan []string, chan bool))// type Conductor func([]chan Composure, chan bool)
type Composure struct{
  name string
  port int
  client *http.Client               // https://golang.org/pkg/net/http/#CookieJar
  serverChannel chan []interface{}       // have to do some clever stuff here to extract what type is sent
  actions map[string][]Action       // have to do some clever stuff here to extract what is sent
  conductors map[string]Conductor
  quit chan bool
  // friends map[string][]struct{
  //   port int
  //   suffix string
  // }
  // neighbours map[string]<-chan []string
}

func Compose(actions map[string][]Action, conductors map[string]Conductor, quit chan bool, name string, port int) Composure {
  return Composure {
    name: name,
    port: 8000,
    client: &http.Client {  },
    serverChannel: make(chan []interface{}),
    actions: actions,
    conductors: conductors,
    quit: make(chan bool),
  }
}
// run loop for Composure operations
func (self  Composure) run(channels map[string][]chan []string) {
  go self.listen()
  for title, conversation := range self.conductors {
    if channel, ok := channels[title]; ok { conversation(&self, title)(channel, self.quit) }
  }
  // need to know what to do about boot strapping the workers
  for {
    select {
    case <- self.quit:
      return
    case body := <- self.serverChannel:
      fmt.Printf("%s got %v", self.name, body)
    }
  }
}
