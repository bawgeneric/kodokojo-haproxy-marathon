#!/bin/bash

haproxy -D -f /usr/local/etc/haproxy/haproxy.cfg -p /tmp/haproxy.pid

kodokojo-haproxy-marathon -httpPort=$PORT -marathonUrl=$MARATHON_URL -marathonCallbackUrl=$MARATHON_URL_CALLBACK

