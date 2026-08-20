#!/bin/sh

openssl req -x509 -nodes -days 365 \
    -newkey rsa:2048 \
    -keyout docker/nginx/certs/shop.local.key \
    -out docker/nginx/certs/shop.local.crt \
    -config docker/nginx/certs/shop.local.cnf \
    -extensions v3_req

echo "Generated local HTTPS certificate for shop.local"