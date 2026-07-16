#!/bin/sh

for i in $(seq 1 4); do
    curl -k -i -s https://shop.local/api/v1/clients | grep -E "HTTP/|X-Proxy-Cache|X-Upstream-Addr"
    echo
done