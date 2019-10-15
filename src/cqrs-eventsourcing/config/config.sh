#!/bin/bash
sudo docker run --rm -d -P -p 127.0.0.1:5432:5432 -e POSTGRES_PASSWORD="1234" --name pg postgres
echo "wait 5s"
sleep 5
psql postgresql://postgres:1234@localhost:5432/postgres -c "create database eventstore"
