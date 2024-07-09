#!/bin/bash

set -e

trap "killall distributeddb" SIGINT

cd "$(dirname "$0")"

killall distributeddb || true
sleep 0.1
 
go build


./distributeddb -db-location=delhi.db -http-addr=127.0.0.1:8080 -config-file=sharding.toml -shard=Delhi &
./distributeddb -db-location=moscow.db -http-addr=127.0.0.1:8081 -config-file=sharding.toml -shard=Moscow &
./distributeddb -db-location=london.db -http-addr=127.0.0.1:8082 -config-file=sharding.toml -shard=London


wait