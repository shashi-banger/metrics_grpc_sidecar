
# Metrics GRPC Sidecar

This repository implements an application with GRPC interface to register 
metrics and to instrument metrics. The application provides a `/metrics` 
endpoint where the metrics can be pulled. This application uses
 [prometheus golang client](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus#section-readme)


## Build

```make all```

The above produces to applications `metrics_collector` and `metrics_producer`. 

- `metrics_collector`: Application with grpc interface which can be used as side car for collecting metrics. It serves a `/metrics` end point where metrics can be fetched or scraped by prometheus server.
- `metrics_producer`: An example application making grpc calls to `media_collector` for producing metrics

## Example run

```
# Terminal 1
./metrics_collector

# Terminal 2
./metrics_producer

# Terminal 3
curl http://localhost:8080/metrics
```
 


