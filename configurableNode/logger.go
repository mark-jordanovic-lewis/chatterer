package configurableNode

import (
  "log"
  "os"
  "time"
  "io"
  "fmt"
  "runtime"
  "errors"
  "strconv"
  "path/filepath"
  "os/exec"
)
// =========================== \\
// Globals (here be monsters!) \\
// =========================== \\
var Logger *log.Logger

// ============= \\
// Program setup \\
// ============= \\
func init() {
  initialiseLogger()
}

// ================= \\
// Logging Functions \\
// ================= \\
// initialise the Logger
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
  // make new Logger io
  log_writer := io.MultiWriter(os.Stdout, file)
  Logger = log.New(log_writer, " ", log.LstdFlags)
  Logger.Println("LogFile : " + logpath)
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
  Logger.Println(logLineData)
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
  Logger.Println(logLineData)
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
  Logger.Println(logLineData)
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
