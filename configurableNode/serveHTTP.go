package configurableNode

import (
  "io"
  "encoding/json"
  "net/http"
  "fmt"
)

// open a listening channel`
func (self Composure) listen() {
  http.ListenAndServe(fmt.Sprintf(":%d", self.port), self.mux())
}

// make Composure an implementor of Handler interface
func (self Composure) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  decoder := json.NewDecoder(req.Body)
  var acceptable_params struct{
    url string
    port int
    suffix string
  }
  if err := decoder.Decode(&acceptable_params); err != nil {
    logErrorEvent("INGRESS DECODE ERROR",fmt.Sprintf("Error decoding json: %s", err))
    io.WriteString(res, fmt.Sprintf("%s", err))
  } else {
    self.serverChannel <- []interface{}{ acceptable_params }  // this blocks, buffer the server channel.ß
    io.WriteString(res, fmt.Sprintf("extracted_params: %+v", acceptable_params))
  }
}

// router
func (self Composure) mux() *http.ServeMux {
  mux := http.NewServeMux()
  mux.Handle(fmt.Sprintf("/%s/", self.name), self)
  return mux
}
