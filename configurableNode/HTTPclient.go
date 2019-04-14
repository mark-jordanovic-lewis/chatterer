package configurableNode

import (
  "fmt"
  "encoding/json"
  "bytes"
  "net/http"
)
// send because _everything_ is a POST req in my world - <JSON>is the way</JSON>
func (self Composure) send(url, suffix string, port int) (*http.Response, error) {
  json_body, err := buildBody(url, suffix, port)
  if err != nil { logErrorEvent("FAILED SEND", fmt.Sprintf("Error marshaling json_body in buildBody(): %s", err)) }
  addr := fmt.Sprintf("%s:%d/%s", url, port, suffix)
  resp, err := post(addr, json_body, self.client)
  if err != nil { return nil, logErrorEvent("FAILED SEND", fmt.Sprintf("Error posting to %s in post(): %s", addr, err)) }
  return resp, nil
}

// helper - build a json body []byte
func buildBody(url, suffix string, port int) ([]byte, error) {
  var body struct{
    url string    `json:"url"`
    port int      `json:"port"`
    suffix string `json:"suffix"`
  }
  body.url = url
  body.port = port
  body.suffix = suffix
  return json.Marshal(body)
}

// helper - just a wrapper around https request functionality
func post(url string, body []byte, client *http.Client) (*http.Response, error) {
  req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
  if err != nil { return nil, err }
  req.Header.Set("Content-Type", "application/json")
  resp, err := client.Do(req)
  if err != nil { return nil, err }
  return resp, nil
}
