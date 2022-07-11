package main

import (
	"context"
	"log"
	"net"
	"strings"
	"time"

	"github.com/shashi-banger/metrics_grpc_sidecar/metrics"
	v1 "github.com/shashi-banger/metrics_grpc_sidecar/metricspb/metricspb_v1"
	"google.golang.org/grpc"
)

type server struct {
	v1.UnimplementedMetricsServer
}

var (
	timeout = 3 * time.Second
)

func handleError(err error) *v1.Response {
	var response v1.Response
	if strings.HasPrefix(err.Error(), metrics.ErrorAlreadyExists) {
		response.StatusCode = 409
		response.Message = err.Error()
		return &response
	} else if strings.HasPrefix(err.Error(), metrics.ErrorUnknownCollector) {
		response.StatusCode = 404 // Metric resource does not exists. Needs to be created
		response.Message = err.Error()
		return &response
	}
	response.StatusCode = 500
	response.Message = err.Error()
	return &response
}

func (*server) CreateCounter(ctx context.Context, req *v1.CreateCounterParams) (*v1.Response, error) {
	var response v1.Response
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := metrics.CreateCounter(req)

	if err != nil {
		response := handleError(err)
		return response, err
	}

	response.StatusCode = 200 //Ok
	return &response, nil
}

func (*server) CounterInc(ctx context.Context, req *v1.UpdateParams) (*v1.Response, error) {
	var response v1.Response
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := metrics.CounterInc(req.Name, req.LabelValues)
	if err != nil {
		response := handleError(err)
		return response, err
	}

	response.StatusCode = 200 //Ok
	return &response, nil
}

func (*server) CreateGauge(ctx context.Context, req *v1.CreateGaugeParams) (*v1.Response, error) {
	var response v1.Response
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := metrics.CreateGauge(req)

	if err != nil {
		response := handleError(err)
		return response, err
	}

	response.StatusCode = 200 //Ok
	return &response, nil
}

func (*server) GaugeSet(ctx context.Context, req *v1.UpdateParams) (*v1.Response, error) {
	var response v1.Response
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := metrics.GaugeSet(req.Name, req.LabelValues, req.Value)
	if err != nil {
		response := handleError(err)
		return response, err
	}

	response.StatusCode = 200 //Ok
	return &response, nil
}

func (*server) CreateHistogram(ctx context.Context, req *v1.CreateHistogramParams) (*v1.Response, error) {
	var response v1.Response
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := metrics.CreateHistogram(req)

	if err != nil {
		response := handleError(err)
		return response, err
	}

	response.StatusCode = 200 //Ok
	return &response, nil
}

func (*server) HistogramObserve(ctx context.Context, req *v1.UpdateParams) (*v1.Response, error) {
	var response v1.Response
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := metrics.HistogramObserve(req.Name, req.LabelValues, req.Value)
	if err != nil {
		response := handleError(err)
		return response, err
	}

	response.StatusCode = 200 //Ok
	return &response, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Println("ERROR:", err.Error())
	}
	s := grpc.NewServer()
	v1.RegisterMetricsServer(s, &server{})

	log.Printf("Server started at %v", lis.Addr().String())

	go func() {
		metrics.StartMetricsEndpointServer()
		s.Stop()
	}()

	err = s.Serve(lis)
	if err != nil {
		log.Println("ERROR:", err.Error())
	}
}
