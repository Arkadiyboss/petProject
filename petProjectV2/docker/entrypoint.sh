#!/bin/sh

/usr/local/bin/docker-entrypoint.sh postgres &

sleep 5

/usr/local/bin/myapp