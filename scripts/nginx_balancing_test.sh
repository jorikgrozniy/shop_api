#!/bin/sh

for i in $(seq 1 16); do
    curl -k -i -s https://shop.local/api/v1/swagger/index.html | grep -E "HTTP/|X-Proxy-Cache|X-Upstream-Addr"
    echo
done