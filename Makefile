PHONY: all
all: metrics_collector metrics_producer


.PHONY: metrics_collector
metrics_collector:
	echo ">> Building metrics_collector"
	go build -o metrics_collector ./cmd/main.go

metrics_producer:
	echo ">> Building metrics_producer"
	go build -o metrics_producer ./examples/metrics_producer.go

clean:
	rm metrics_collector metrics_producer