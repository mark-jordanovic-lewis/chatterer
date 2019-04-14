package configurableNode

import (
  "errors"
  "fmt"
  "time"
)
// =================== \\
// Conductor functions \\
// =================== \\
// can use this for a map or slice of discussions:
//   discussions := []Conductor{ monologue, dialogue, ..., forum }
// think of it like an interface for function types
//
// func discussion()
// - map[string]chan Composure
// - channels are bi-directional
func dialogueConductor(self Composure, name string) (func ([]chan []string, chan bool) error) { // func dialogue(channels []chan Composure, kill chan bool) error {
  return func (channels []chan []string, kill chan bool) error {
    if len(channels) != 2 { return errors.New("Conductor Error: dialogue - channels must have two elements.") }
    for {
      select {
      case ingress := <- self.serverChannel:
        fmt.Printf("ingress json: %+v", ingress)
      case one := <- channels[0]:
        for _, in_string := range one {
          ping := make(chan []string)
          fmt.Println(in_string)
          // go self.discussions[name]
          reciept := <- ping
          fmt.Printf("recieved: %+v", reciept)
        }
        // parse string in relevant way
        // send to channel two
        fmt.Println("got:", one, "on channel 1")
        // channels[1] <- channels[0]
      case two := <- channels[1]:
        for _, in_string := range two {
          ping := make(chan []string)
          fmt.Println(in_string)
          // go self.discussions[name]
          reciept := <- ping
          fmt.Printf("recieved: %+v", reciept)
        }
        // parse string in relevant way
        // send to channel two
        fmt.Println("got:", two, "on channel 2")
        // channels[0] <- channels[1]
      case <- kill:
        return nil
      default:
        time.Sleep(time.Duration(50*time.Millisecond))
        fmt.Print(".")
      }
    }
  }
}
