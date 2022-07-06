#!/bin/bash
cd metricspb/v1
protoc -I ./ --go_out=./ --go_grpc_out=./ metrics.proto 
cd .. && cd ..