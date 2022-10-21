package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ugosan/logshark/v1/config"
	"github.com/ugosan/logshark/v1/logging"
)

type Stats struct {
	Events     int
	EpsT0      int
	Eps        int
	MaxEvents  int
	TotalBytes int
	AvgBytes   int
}

var (
	configflags config.Config
	currentStats = Stats{0, 0, 0, 0, 0, 0}
	eventChannel = make(chan string)
	statsChannel = make(chan Stats)
)

func addEvent(jsonBody string) {

	if currentStats.Events < configflags.MaxEvents {

		eventChannel <- jsonBody

	}

	currentStats.Events++
	currentStats.TotalBytes += cap([]byte(jsonBody))
	currentStats.AvgBytes = currentStats.TotalBytes / currentStats.Events
	statsChannel <- currentStats
}

func updateEps() {
	ticker := time.NewTicker(time.Second).C

	for range ticker {
			currentStats.Eps = currentStats.Events - currentStats.EpsT0
			currentStats.EpsT0 = currentStats.Events
			currentStats.MaxEvents = configflags.MaxEvents
			statsChannel <- currentStats
	}
}

func home(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "{\"name\":\"instance-0000000000\",\"cluster_name\":\"5b4b73e528774cf9a8fe60b7909adcf0\",\"cluster_uuid\":\"d4XqGxciQxqAFle4qYNWIQ\",\"version\":{\"number\":\"7.13.0\",\"build_flavor\":\"default\"},\"tagline\":\"You Know, for Search\"}")
	case "POST":

		addEvent(string(body))

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}

func license(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	switch r.Method {
		case "GET":
			fmt.Fprintf(w, "{\"license\":{\"status\":\"active\",\"uid\":\"fb7a9815-8b0a-4608-b786-049fbec4a4a8\",\"type\":\"platinum\",\"issue_date\":\"2020-03-24T00:00:00.000Z\",\"issue_date_in_millis\":1585008000000,\"expiry_date\":\"2222-06-30T00:00:00.000Z\",\"expiry_date_in_millis\":1656547200000,\"max_nodes\":100000,\"issued_to\":\"ElasticCloud\",\"issuer\":\"API\",\"start_date_in_millis\":1614556800000}}")
		case "POST":
			addEvent(string(body))

		default:
			fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}

func bulk(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)

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

			if currentStats.Events < configflags.MaxEvents {
				addEvent(splits[i])
			} else {
				currentStats.Events++
			}

		}

		fmt.Fprintf(w, "{\"errors\": false}")

	default:
		fmt.Fprintf(w, "Sorry, only POST method is supported.")
	}

}

func SendTestRequest() {

	var testJSON = fmt.Sprintf("{	\"sequence\": %d, \"hola\": \"hola\",\"obj\": {\"a\": 1, \"string\": \"stringsss\", \"array\": [\"one\",\"two\",\"three\"],\"float\": 3.14}, \"name\" : \"instance-000000001\",	\"cluster_name\" : \"<dummy-cluster\",	\"cluster_uuid\" : \"yaVi2rdIQT-v-qN9v4II9Q\",	\"version\" : {		\"number\" : \"6.8.3\",		\"build_flavor\" : \"default\",		\"build_type\" : \"tar\",		\"build_hash\" : \"0c48c0e\",		\"build_date\" : \"2019-08-29T19:05:24.312154Z\",		\"build_snapshot\" : false,		\"lucene_version\" : \"7.7.0\",		\"minimum_wire_compatibility_version\" : \"5.6.0\",		\"minimum_index_compatibility_version\" : \"5.0.0\"	},	\"xml\": \"<car><mirror>XML</mirror></car>\", \"tagline\" : \"You Know, for Search\", \"test\": \" <19> ()weird \\\"aaa\\\": 1 \"}", currentStats.Events)

	resp, err := http.Post(
		fmt.Sprintf("http://%s:%s", configflags.Host, configflags.Port),
		"application/json",
		bytes.NewBuffer([]byte(testJSON)))
	if err != nil {
		print(err)
		logging.GetManager().Log(err)
	}

	defer resp.Body.Close()
}

func GetStats() Stats {
	return currentStats
}

func ResetStats() Stats {
	currentStats = Stats{0, 0, 0, configflags.MaxEvents, 0, 0}
	return currentStats
}

func Start(_eventChannel chan string, _statsChannel chan Stats, config config.Config) {

	configflags = config

	logs := logging.GetManager()

	logs.Log("Listening on " + config.Host + ":" + config.Port)

	http.HandleFunc("/", home)
	http.HandleFunc("/_bulk", bulk)
	http.HandleFunc("/_license", license)

	eventChannel = _eventChannel
	statsChannel = _statsChannel

	go updateEps()

	err := http.ListenAndServe(config.Host+":"+config.Port, nil)
	if err != nil {
		log.Fatal("Error Starting the HTTP Server : ", err)
		return
	}
}
