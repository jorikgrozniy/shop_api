#!/bin/sh

for i in $(seq 1 16); do
    curl -s http://localhost/api/v1/whoami
    echo
done