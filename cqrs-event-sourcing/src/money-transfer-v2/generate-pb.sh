#!/usr/bin/env bash

protoc -I pb/ pb/*.proto --go_out=pb