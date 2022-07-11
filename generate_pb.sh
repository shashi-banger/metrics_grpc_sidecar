#!/bin/bash
cd metricspb/metricspb_v1
#protoc --plugin=/Users/shashidhar/go/bin/protoc-gen-go-grpc --go_out=. --go-grpc_out=.  metrics.proto
protoc --go_out=. --go-grpc_out=.  metricspbV1.proto
cd .. && cd ..
