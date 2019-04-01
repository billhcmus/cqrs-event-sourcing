#!/bin/bash
docker run --rm -d -P -p 127.0.0.1:5432:5432 -e POSTGRES_PASSWORD="1234" --name pg postgres

echo "wait..."
sleep 5

#psql postgresql://postgres:1234@localhost:5432/postgres -c "create database eventstore"

docker exec -it pg createdb -U postgres eventstore

docker run --rm -d --name nats-es -p 4222:4222 -p 6222:6222 -p 8222:8222 nats

docker run --rm -d --name redis-cache -p 6379:6379 redis
