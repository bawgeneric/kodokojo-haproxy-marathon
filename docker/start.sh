#!/bin/bash

haproxy -D -f /usr/local/etc/haproxy/haproxy.cfg -p /tmp/haproxy.pid

kodokojo-haproxy-marathon "$@"

