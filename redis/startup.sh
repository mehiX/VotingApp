#!/bin/bash

mkdir -p /usr/local/etc/redis

envsubst < /tmp/redis.conf > /usr/local/etc/redis/redis.conf && exec redis-server /usr/local/etc/redis/redis.conf