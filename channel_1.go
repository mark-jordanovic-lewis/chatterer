package main

import (
  // Heavy Lifting packages
  "encoding/json"
  "bytes"
  "net/http"
  "io"
  // Error and Log handling
  "errors"
  "log"
  "os"
  "os/exec"
  "path/filepath"
  "runtime"
  // Std utilities
  "fmt"
  "strconv"
  "math/rand"
  "time"
)

// =========================== \\
// Globals (here be monsters!) \\
// =========================== \\
var (
    logger *log.Logger
)

// ============= \\
// Program setup \\
// ============= \\
func init() {
  initialiseLogger()
}


// ================ \\
// Server Functions \\
// ================ \\

// Microservice Action One \\
// ----------------------- \\
type Worker func(*Composure) (start func(chan []string))// maybe []byte would be better so JSON can be passed
type Conductor func(*Composure) (start func([]chan []string, chan bool))// type Conductor func([]chan Composure, chan bool)
type Composure struct{
  name string
  port int
  client *http.Client
  serverChannel chan interface{}   // have to do some clever stuff here to extract what type is sent
  discussions map[string]Conductor
  actions map[string][]Worker       // have to do some clever stuff here to extract what is sent
  quit chan bool
  friends map[string][]struct{
    port int
    suffix string
  }
  neighbours map[string]<-chan []string
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
    // send to server_channel here
    io.WriteString(res, fmt.Sprintf("extracted_params: %+v", acceptable_params))
  }
}

// open a listening channel`
func (self Composure) listen() {
  http.ListenAndServe(fmt.Sprintf(":%d", self.port), self.mux())
}

// run loop for Composure operations
func (self  Composure) run(channels map[string][]chan []string) {
  go self.listen()
  for title, conversation := range self.discussions {
    if channel, ok := channels[title]; ok { conversation(&self)(channel, self.quit) }
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

// router
func (self Composure) mux() *http.ServeMux {
  mux := http.NewServeMux()
  mux.Handle(fmt.Sprintf("/%s/", self.name), self)
  return mux
}

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
func dialogue(channels []chan []string, kill chan bool) error { // func dialogue(channels []chan Composure, kill chan bool) error {
  if len(channels) != 2 { return errors.New("Conductor Error: dialogue - channels must have two elements.") }
  for {
    select {
    case one := <- channels[0]:
      // parse string in relevant way
      // send to channel two
      fmt.Println("got:", one, "on channel 1")
      // channels[1] <- channels[0]
    case two := <- channels[1]:
      // parse string in relevant way
      // send to channel one
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

// ==== \\
// MAIN \\
// ==== \\
func main() {
  defer func() {
    if err := recover(); err != nil {
      fmt.Println("The panic was recovered:\n", err)
    }
  }()
  rand.Seed(time.Now().UnixNano())

  chan_one := make(chan []string)
  chan_two := make(chan []string)
  chan_one <- []string{"beep"}
  go func() { <-chan_one }()
  chan_two <- []string{"beep"}
  go func() { <-chan_two }()
  chan_one = make(chan []string)
}

// ================= \\
// Logging Functions \\
// ================= \\
// initialise the logger
// - make a:
//     channel log
//     error log
//     ingress log
func initialiseLogger() {
  // make and set location of log file
  date := time.Now().Format("2006-01-02")
  var logpath = getCurrentExecDir() + "/log/" + date + "_server.log"
  os.MkdirAll(getCurrentExecDir()+"/log/", os.ModePerm)
  // open file for create and append
  var file, err = os.OpenFile(logpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
  if err != nil { panic(fmt.Sprintf("Could not open logfile `%s`: %s", logpath, err)) }
  // make new logger io
  log_writer := io.MultiWriter(os.Stdout, file)
  logger = log.New(log_writer, " ", log.LstdFlags)
  logger.Println("LogFile : " + logpath)
}

// log to file
func logErrorEvent(logLevel string, message string) error {
  _, file, no, ok := runtime.Caller(1)
  var logLineData string
  if ok {
    file = shortenFilePath(file)
    logLineData = fmt.Sprintf("%s ==== [ ERROR : %s ] %s : %d - %s", timestamp(), logLevel, file,  no, message)
  } else {
    logLineData = fmt.Sprintf("%s ==== [ ERROR : %s ] %s (Caller not found) : %d - %s", timestamp(), logLevel, file, no, message)
  }
  logger.Println(logLineData)
  return errors.New(message)
}

func logIngressEvent(logLevel string, message string) {
  _, file, no, ok := runtime.Caller(1)
  var logLineData string
  if ok {
    file = shortenFilePath(file)
    logLineData = fmt.Sprintf("%s - %s %s : %d - %s", timestamp(), logLevel, file, no, message)
  } else {
    logLineData = fmt.Sprintf("%s - %s %s (Caller not found) : %d - %s", timestamp(), logLevel, file, no, message)
  }
  logger.Println(logLineData)
}

func logChannelEvent(logLevel string, message string) {
  // direction to find channel easier --> | <-- | <-->
  _, file, no, ok := runtime.Caller(1)
  var logLineData string
  if ok {
    file = shortenFilePath(file)
    logLineData = fmt.Sprintf("%s - %s %s : %d <---------> %s", timestamp(), logLevel, file, no, message)
  } else {
    logLineData = fmt.Sprintf("%s - %s %s (Caller not found) : %d - %s", timestamp(), logLevel, file, no, message)
  }
  logger.Println(logLineData)
}

func timestamp() string {
  var minute string
  if min := time.Now().Minute(); min < 0 {
    minute = ":0"+strconv.Itoa(min)
  } else {
    minute = strconv.Itoa(min)
  }
  return strconv.Itoa(time.Now().Hour()) + minute
}

// shortens file path a/b/c/d.go -> d.go
func shortenFilePath(file string) string {
  short := file
  for i := len(file) - 1; i > 0; i-- {
    if file[i] == '/' {
        short = file[i+1:]
        break
    }
  }
  file = short
  return file
}

func getCurrentExecDir() string {
  path, err := exec.LookPath(os.Args[0])
  if err != nil { panic(fmt.Sprintf("exec.LookPath(%s), err: %s\n", os.Args[0], err)) }
  absPath, err := filepath.Abs(path)
  if err != nil { panic(fmt.Sprintf("filepath.Abs(%s), err: %s\n", path, err)) }
  return filepath.Dir(absPath)
}
