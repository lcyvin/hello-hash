package main

import (
  "os"
  "fmt"
  "io"
  "time"
  "crypto/md5"
  "encoding/json"
  "net/http"
)

const (
  DEFAULT_PORT string = "3031"
)

var (
  port string
  installPath string
)

type Response struct {
  Timestamp int64  `json:"timestamp"`
  Checksum []byte `json:"checksum"`
}

func NewResponse() (*Response, error) {
  r := &Response{
    Timestamp: time.Now().Unix(),
  }

  f, err := os.Open(installPath)
  if err != nil {
    return nil, err
  }

  defer f.Close()
  hash := md5.New()
  _, err = io.Copy(hash, f)
  if err != nil {
    return nil, err
  }

  r.Checksum = hash.Sum(nil)
  return r, nil
}

func ResponseHandler(w http.ResponseWriter, r *http.Request) {
  resp, err := NewResponse()
  if err != nil {
    http.Error(w, fmt.Sprintf("%v", err), 500)
  }

  j, err := json.Marshal(resp)
  if err != nil {
    http.Error(w, fmt.Sprintf("%v", err), 500)
  }

  w.Write(j)
}

func init() {
  // since we can run with go run, we want to handle the ephemeral binary
  f, err := os.Executable()
  if err != nil {
    panic(err)
  }

  installPath = f
  
  if port = os.Getenv("PORT"); port == "" {
    port = DEFAULT_PORT
  }
}

func main() {
  http.HandleFunc("/", ResponseHandler)
  http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
