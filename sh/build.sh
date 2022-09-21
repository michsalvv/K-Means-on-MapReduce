#! /usr/bin/bash

#go build -o client src/client/*.go
for mapper in $(docker ps -a --format "{{.Names}}" --filter "name=mapper")
do 
    sudo docker exec $mapper bash -c "go env -w GO111MODULE=off ; go build -o bin/mapper src/kmeans-MR/workers/mapper/*.go"
done

for reducer in $(docker ps -a --format "{{.Names}}" --filter "name=reducer")
do 
    sudo docker exec $reducer bash -c "go env -w GO111MODULE=off ; go build -o bin/reducer src/kmeans-MR/workers/reducer/*.go"

done

sudo docker exec km-master bash -c "go env -w GO111MODULE=off ; go build -o bin/master src/kmeans-MR/master/*.go"
