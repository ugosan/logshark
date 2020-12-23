package main


import (
  "github.com/ugosan/logshark/v1/ui"
  "github.com/ugosan/logshark/v1/config"
  "flag"
  "log"
  "os"
)


func main() {

  config := config.Config{}

  flag.StringVar(&config.Host, "host", "0.0.0.0", "Specify host. Default is 0.0.0.0")
  flag.StringVar(&config.Port, "port", "8080", "Specify port. Default is 8080")
  flag.IntVar(&config.MaxEvents, "max", 1000, "Specify max events. Default is 1000")
  flag.StringVar(&config.LogFile, "log", "/dev/null", "Specify a log file (debug). Default is /dev/null")

  flag.Parse()  

  f, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil {
    log.Println(err)
  }
  defer f.Close()

  log := log.New(f, "[main] ", log.LstdFlags)
  log.Println(config)

  ui.Start(config)

}