#!/bin/sh

./cleandocker.sh

if [[ $1 == "inmemory" ]] || [[ $1 == "" ]]; then
    docker build -t inmemory .
    docker run --rm -d --name inmemory -p 8080:8080 inmemory
    clear
    echo "Service: inmemory mode is started"
    sleep 10
elif [ $1 == "database" ]; then
    docker-compose up -d
    clear
    echo "Service: database mode is started"
    sleep 10
else
    echo "Invalid parameter"
    echo "Expected: './app.sh inmemory' or './app.sh database'"
    sleep 10
fi
