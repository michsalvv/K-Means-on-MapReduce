#! /usr/bin/bash

go build -o client src/client/*.go
sudo docker exec km-master bash -c "go build -o /bin/master src/kmeans-MR/src/server/master/*.go"
sudo docker exec km-mapper bash -c "go build -o /bin/mapper src/kmeans-MR/src/server/workers/mapper/*.go"
sudo docker exec km-reducer bash -c "go build -o /bin/reducer src/kmeans-MR/src/server/workers/reducer/*.go"
