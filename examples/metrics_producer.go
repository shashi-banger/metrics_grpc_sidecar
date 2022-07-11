package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	pb "github.com/shashi-banger/metrics_grpc_sidecar/metricspb/metricspb_v1"
	"google.golang.org/grpc"
)

func random_metrics_producer(ctx context.Context, container_name string, client pb.MetricsClient) {
	counterParams := pb.CreateCounterParams{Name: "num_requests",
		Labels: []string{"container_name", "stage"},
		Help:   "Number of requests at various stages"}
	// Create a counter metric
	client.CreateCounter(ctx, &counterParams)

	// Create a Gauge metric
	gaugeParams := pb.CreateGaugeParams{Name: "active_requests",
		Labels: []string{"container_name"},
		Help:   "Number of requests at various stages"}
	// Create a counter metric
	client.CreateGauge(ctx, &gaugeParams)

	// Create a Histo metrics

	for {
		update_params := pb.UpdateParams{Name: "num_requests",
			LabelValues: map[string]string{"container_name": container_name, "stage": "input"}}
		resp, err := client.CounterInc(context.Background(), &update_params)
		fmt.Println("Calling counter")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(resp.StatusCode)
		}

		update_gauge := pb.UpdateParams{Name: "active_requests",
			LabelValues: map[string]string{"container_name": container_name}, Value: float64(rand.Intn(15))}
		client.GaugeSet(context.Background(), &update_gauge)
		time.Sleep(time.Second)
	}

}

func main() {
	var address = flag.String("a", "", "getafix service address")
	flag.Parse()

	conn, err := grpc.Dial(*address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Printf("Grpc dial failed\n")
		panic("Grpc dial failed")
	}
	defer conn.Close()

	client := pb.NewMetricsClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	go random_metrics_producer(ctx, "tarang", client)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}
