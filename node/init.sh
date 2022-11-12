#! /bin/bash
go env -w GO111MODULE=off
go get gopkg.in/yaml.v2 
go get github.com/kelseyhightower/envconfig

if [ "$1" = "master" ]; then
    go env -w GO111MODULE=off ; go build -o bin/master src/kmeans-MR/master/*.go ; ./bin/master
elif [ "$1" = "mapper" ]; then
    go env -w GO111MODULE=off ; go build -o bin/worker src/kmeans-MR/workers/*.go ; ./bin/worker 172.17.0.1:9001 mapper
elif [ "$1" = "reducer" ]; then
    go env -w GO111MODULE=off ; go build -o bin/worker src/kmeans-MR/workers/*.go ; ./bin/worker 172.17.0.1:9001 reducer
else
    echo "Please use 'sh init.sh master/mapper/reducer'"
fi