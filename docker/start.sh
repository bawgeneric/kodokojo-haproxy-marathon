#!/bin/bash

haproxy -D -f /usr/local/etc/haproxy/haproxy.cfg -p /tmp/haproxy.pid

kodokojo-haproxy-marathon -httpPort=4444 -marathonUrl=$MARATHON_URL -marathonCallbackUrl=$MARATHON_URL_CALLBACK

