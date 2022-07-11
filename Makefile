PHONY: all
all: build_metrics_collector


.PHONY: build_metrics_collector
build_metrics_collector:
	echo ">> Building metrics_collector"
	go build -o metrics_collector ./cmd/main.go
