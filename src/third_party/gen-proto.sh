#!/bin/bash
protoc -I ../api/proto/v1 ../api/proto/v1/*.proto --go_out=plugins=grpc:../pkg/api/v1