#!/usr/bin/env bash

(cd 7000/ && redis-server ./redis.conf &)
(cd 7001/ && redis-server ./redis.conf &)
(cd 7002/ && redis-server ./redis.conf &)

sleep 3
redis-cli --cluster create 127.0.0.1:7000 127.0.0.1:7001 127.0.0.1:7002 --cluster-yes

sleep infinity
