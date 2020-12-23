package server

import (
  "github.com/ugosan/logshark/v1/config"
  "fmt"
  "log"
  "io/ioutil"
  "net/http"
  "encoding/json"
  "strings"
  "bytes"
  "time"
)


type Stats struct {
  Events int
  EpsT0 int
  Eps int
  MaxEvents int
}

var configflags config.Config

var currentStats = Stats{0, 0, 0, 0}
var channel = make(chan map[string]interface{})

func addEvent(jsonBody string){

  if(currentStats.Events < configflags.MaxEvents){
  
    var obj map[string]interface{}
    json.Unmarshal([]byte(jsonBody), &obj)

    channel <- obj

  }

  currentStats.Events += 1

}

func updateEps() {
  ticker := time.NewTicker(time.Second).C

  for {
    select {
    case <-ticker:
      currentStats.Eps = currentStats.Events - currentStats.EpsT0
      currentStats.EpsT0 = currentStats.Events
      currentStats.MaxEvents = configflags.MaxEvents
    }
  }
}


func home(w http.ResponseWriter, r *http.Request) {

  body, err := ioutil.ReadAll(r.Body)

  if err != nil {
    log.Printf("Error reading body: %v", err)
      http.Error(w, "can't read body", http.StatusBadRequest)
      return
  }

  switch r.Method {
    case "GET":
      fmt.Fprintf(w, "{	\"name\" : \"instance-000000001\",	\"cluster_name\" : \"dummy-cluster\",	\"cluster_uuid\" : \"yaVi2rdIQT-v-qN9v4II9Q\",	\"version\" : {		\"number\" : \"6.8.3\",		\"build_flavor\" : \"default\",		\"build_type\" : \"tar\",		\"build_hash\" : \"0c48c0e\",		\"build_date\" : \"2019-08-29T19:05:24.312154Z\",		\"build_snapshot\" : false,		\"lucene_version\" : \"7.7.0\",		\"minimum_wire_compatibility_version\" : \"5.6.0\",		\"minimum_index_compatibility_version\" : \"5.0.0\"	},	\"tagline\" : \"You Know, for Search\"}")
    case "POST":

      addEvent(string(body))

    default:
        fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
    }

}


func bulk(w http.ResponseWriter, r *http.Request) {

  body, err := ioutil.ReadAll(r.Body)

  if err != nil {
    log.Printf("Error reading body: %v", err)
    http.Error(w, "can't read body", http.StatusBadRequest)
    return
  }

  switch r.Method {

    case "POST":

      var splits = strings.Split(string(body), "\n")

      for i := 0; i < len(splits); i++ {
        if i%2 == 0 {
            continue
        }

        if(currentStats.Events < configflags.MaxEvents){
          addEvent(splits[i])
        }else{
          currentStats.Events += 1
        }

      }

      fmt.Fprintf(w, "{\"errors\": false}")

    default:
      fmt.Fprintf(w, "Sorry, only POST method is supported.")
    }

}

func SendTestRequest(){

  var testJson = fmt.Sprintf("{	\"sequence\": %d, \"hola\": \"hola\",\"obj\": {\"a\": 1, \"string\": \"stringsss\", \"array\": [\"one\",\"two\",\"three\"],\"float\": 3.14}, \"name\" : \"instance-000000001\",	\"cluster_name\" : \"dummy-cluster\",	\"cluster_uuid\" : \"yaVi2rdIQT-v-qN9v4II9Q\",	\"version\" : {		\"number\" : \"6.8.3\",		\"build_flavor\" : \"default\",		\"build_type\" : \"tar\",		\"build_hash\" : \"0c48c0e\",		\"build_date\" : \"2019-08-29T19:05:24.312154Z\",		\"build_snapshot\" : false,		\"lucene_version\" : \"7.7.0\",		\"minimum_wire_compatibility_version\" : \"5.6.0\",		\"minimum_index_compatibility_version\" : \"5.0.0\"	},	\"tagline\" : \"You Know, for Search\", \"test\": \" <>()weird \\\"aaa\\\": 1 \"}", currentStats.Events)

  
  resp, err := http.Post(
    fmt.Sprintf("http://%s:%s", configflags.Host, configflags.Port),
    "application/json",
    bytes.NewBuffer([]byte(testJson)))
  if err != nil {
    print(err)
  }

  defer resp.Body.Close()
}

func GetStats() Stats {
  return currentStats
}

func ResetStats() Stats {
  currentStats = Stats{0, 0, 0, configflags.MaxEvents}
  return currentStats
}


func Start(c chan map[string]interface{}, config config.Config) {
  
  configflags = config

  

  log.Print("Listening "+config.Host+":"+config.Port)

  http.HandleFunc("/", home)
  http.HandleFunc("/_bulk", bulk)

 
  channel = c

  go updateEps()

  err := http.ListenAndServe(config.Host+":"+config.Port, nil)
  if err != nil {
    log.Fatal("Error Starting the HTTP Server : ", err)
    return
  }
}