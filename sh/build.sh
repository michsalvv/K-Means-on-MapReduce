#! /usr/bin/bash

#go build -o client src/client/*.go
for mapper in $(docker ps -a --format "{{.Names}}" --filter "name=mapper")
do 
    docker exec $mapper bash -c "go env -w GO111MODULE=off ; go build -o bin/worker src/kmeans-MR/workers/*.go"
done

for reducer in $(docker ps -a --format "{{.Names}}" --filter "name=reducer")
do 
    docker exec $reducer bash -c "go env -w GO111MODULE=off ; go build -o bin/worker src/kmeans-MR/workers/*.go"

done

docker exec km-master bash -c "go env -w GO111MODULE=off ; go build -o bin/master src/kmeans-MR/master/*.go"

echo " -- BUILDED --"
