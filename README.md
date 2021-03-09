# logshark

Logshark is a debugger for JSON logs.

<kbd>![](./_doc/demo.gif)</kbd>



Features:
- Terminal UI 
- Beats/Logstash integration
- Small (<2Mb)
- Provides stats such as events per second and average document size, useful for benchmarking

Logshark works by listening for logs on a specific TCP port, it mimicks the Elasticsearch protocol so as to receive data from Beats (Filebeat, Metricbeat, Heartbeat, etc.) using the standard elasticsearch output.

## How to use

```bash
./logshark --host 127.0.0.1 --port 9200 --max 1000
```