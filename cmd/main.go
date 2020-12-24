package main


import (
  "github.com/ugosan/logshark/v1/ui"
  "github.com/ugosan/logshark/v1/config"
  "github.com/ugosan/logshark/v1/logging"

  "flag"
)


func main() {

  config := config.Config{}

  flag.StringVar(&config.Host, "host", "0.0.0.0", "Specify host. Default is 0.0.0.0")
  flag.StringVar(&config.Port, "port", "8080", "Specify port. Default is 8080")
  flag.IntVar(&config.MaxEvents, "max", 1000, "Specify max events. Default is 1000")
  flag.StringVar(&config.LogFile, "log", "/dev/null", "Specify a log file (debug). Default is /dev/null")

  flag.Parse()  

  if(config.LogFile != "/dev/null"){
    logs := logging.GetManager()
    logs.InitLogger(config)
    logs.Log(config)    
  }



  ui.Start(config)
}