package configurableNode

import (
  //"errors"
  "testing"
  "time"
  "fmt"
)

func TestCompose(t *testing.T) {
  fmt.Println()
  fmt.Println("TestCompose")
  fmt.Println("================")
  actions := make(map[string][]Action)
  actions["name"] = []Action { mockAction() }
  conductors := make(map[string]Conductor)
  conductors["name"] = mockConductor()
  quit := make(chan bool)
  composure := Compose(actions, conductors, quit, "name", 8080)
  fmt.Println("The endpoint system is instantiated")
  fmt.Println("-----------------------------------")
  fmt.Println("actions map generated")
  if len(composure.actions) != 1 || len(composure.actions["name"]) != 1 {
    t.Errorf("Did not add mock action")
  }
  fmt.Println("conductors map generated")
  if len(composure.conductors) != 1 {
    t.Errorf("Did not add mock action")
  }
  fmt.Println("IO channels setup")
  fmt.Println("-----------------")
  fmt.Println("server channel permits passing anything (it's up to the user to deal with input)")
  // This is going to be a big test - let's leave this till I make some actions and composors
  fmt.Println("http client pings t'interwebs")
  _, err := composure.client.Get("http://www.yahoo.com")
  if err != nil { t.Errorf("Could not ping t'interwebs") }
  fmt.Println("quit recieves and sends bool")
  go func(){
    for {
      select {
      case <- composure.quit: return
      case <- time.Tick(1 * time.Second): t.Errorf("Timed out waitng for quit send")
      }
    }
  }()
  composure.quit <- true
}


// mocks \\
func mockAction() Action {
  action := func(self *Composure, name string) (func(chan []string)) {
    return func(_io chan []string) {}
  }
  return action
}

func mockConductor() Conductor {
  conductor := func(self *Composure, name string) (func([]chan []string, chan bool)) {
    return func(_ios []chan []string, kill chan bool) {}
  }
  return conductor
}
