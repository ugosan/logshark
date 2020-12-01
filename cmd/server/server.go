// Golang HTTP Server
package server

import (
  "fmt"
  "log"
  "io/ioutil"
  "net/http"
  "encoding/json"
  "github.com/TylerBrock/colorjson"
  "strings"
)

const (
  Host = "0.0.0.0"
  Port = "8080"
)

var channel = make(chan map[string]interface{})

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

			log.Printf(string(body))
      
    
      var obj map[string]interface{}
      json.Unmarshal([]byte(body), &obj)
      f := colorjson.NewFormatter()
      f.Indent = 4
      //s, _ := f.Marshal(obj)
      //fmt.Println(string(s))
    
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
			//log.Printf("POST _bulk:")
      
      var splits = strings.Split(string(body), "\n")
      
      for i := 0; i < len(splits); i++ {
        if i%2 == 0 { 
            continue
        }
        var obj map[string]interface{}
				json.Unmarshal([]byte(splits[i]), &obj)
				channel <- obj

        //f := colorjson.NewFormatter()
        //f.Indent = 4
        //s, _ := f.Marshal(obj)
        //log.Printf(string(s))
      }

      fmt.Fprintf(w, "{\"errors\": false}")
			
    default:
      fmt.Fprintf(w, "Sorry, only POST method is supported.")
    }
    


  /*
  var prettyJSON bytes.Buffer
  error := json.Indent(&prettyJSON, body, "", "  ")
  if error != nil {
      log.Println("JSON parse error: ", error)
      return
  }

  log.Println("\n", string(prettyJSON.Bytes()))
*/

}




func Start(c chan map[string]interface{}) {
  log.Print("Listening "+Host+":"+Port)
  
  http.HandleFunc("/", home)
	http.HandleFunc("/_bulk", bulk)
	
	channel = c

  err := http.ListenAndServe(Host+":"+Port, nil)
  if err != nil {
    log.Fatal("Error Starting the HTTP Server : ", err)
    return
  }
}

