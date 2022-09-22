#! /usr/bin/bash

terminator --new-tab -x "docker exec -it km-master ./bin/master; /usr/bin/zsh"

for mapper in $(docker ps -a --format "{{.Names}}" --filter "name=mapper")
do 
    terminator --new-tab -x "docker exec -it $mapper ./bin/worker 172.17.0.1:9001 mapper"
done


for reducer in $(docker ps -a --format "{{.Names}}" --filter "name=reducer")
do 
    terminator --new-tab -x "docker exec -it $reducer ./bin/worker 172.17.0.1:9001 reducer"
done
