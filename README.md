# logshark

Logshark is a debugger for JSON logs.


<kbd>![](./_doc/demo.gif)</kbd>

Logshark works by listening for logs on an HTTP port, it mimicks the Elasticsearch protocol so as to receive data from Beats (Filebeat, Metricbeat, Heartbeat, etc.) and Logstash using the standard elasticsearch output. 

Features:
- Terminal UI 
- Navigable list of logs 
- Highlightable, pretty printed JSON
- 🎨 Colorful
- Beats/Logstash integration
- Stats such as *Events per second* and *Average size* in bytes per event - useful for calculating bulk/batch size

## Download

Releases [here](https://github.com/ugosan/logshark/releases)

## How to use

### binary

```perl
./logshark --host 0.0.0.0 --port 9200 --max 1000
```

### docker

```perl
docker run -p 9200:9200 -it ugosan/logshark -host 0.0.0.0 -port 9200
```