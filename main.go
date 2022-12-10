package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "path"
  "net"
  "strings"
)

var (
  root string
  index string
  port string
)

func main() {
  port = "8082"

  //url := GetLocalIPv4()+":"+port
  //local := "localhost:"+port

  //root = "/home/avery/goserve/public/"
  index = "index.html"

  http.HandleFunc("/", handler)

  finish := make(chan bool)

  //go serveSSL(url)
  go serve(":"+port)

  <-finish
}

func GetLocalIPv4() string {
    netInterfaceAddresses, err := net.InterfaceAddrs()
    rerr := "Cannot find IPv4 address"
    if err != nil { return rerr }
    for _, netInterfaceAddress := range netInterfaceAddresses {
        networkIp, ok := netInterfaceAddress.(*net.IPNet)
        if ok && !networkIp.IP.IsLoopback() && networkIp.IP.To4() != nil {
            ip := networkIp.IP.String()
            fmt.Println("Resolved Host IP: " + ip)
            return ip
        }
    }
    return rerr
}

func serveSSL(url string) {
  fmt.Println("Starting server on", url)
  err := http.ListenAndServeTLS(url, "server.crt", "server.key", nil)
  if err != nil {
    fmt.Println(err)
    return
  }
}

func serve(url string) {
  fmt.Println("Starting server on", url)
  err := http.ListenAndServe(url, nil)
  if err != nil {
    fmt.Println(err)
    return
  }
}

func loadFile(rt string, file string) string {
  if file == "/" {
    file = index
  }
  bytes, err := ioutil.ReadFile(path.Join(rt, file))
  if err != nil {
    if file == index {
      return "<html><center style='padding:40px'>My Go Web Server</center><title>My Go Web Server</title></html>"
    } else {
      return "404 - " + file + " - file not found"
    }
  }
  str := string(bytes)
  return str
}

func loadConf(file string) string {
  bytes, err := ioutil.ReadFile(path.Join("sites/", file+".conf"))
  if err != nil {
    bytes, _ = ioutil.ReadFile("sites/default.conf")
  }
  str := string(bytes)
  return str
}

func handler(w http.ResponseWriter, r *http.Request) {
  url := r.URL.String()
  fmt.Println("Requested URL:", url)
  host := strings.Replace(strings.Replace(r.Host, ":", "", 1), port, "", 1)
  host = strings.Replace(host, "www.", "", 1)
  fmt.Println("Requested Host:", host)
  rt := strings.Replace(loadConf(host), "\n", "", 1)
  fmt.Println("Root:", rt)
  if r.URL.String() != "" {
    fmt.Fprintf(w, loadFile(rt, url))
  }
  fmt.Println("")
}
