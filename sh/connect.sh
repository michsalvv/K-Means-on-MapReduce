#! /usr/bin/bash

for mapper in $(docker ps -a --format "{{.Names}}" --filter "name=mapper")
do 
    terminator --new-tab --title "mapper $mapper" -x "docker exec -it $mapper ./bin/worker 172.17.0.1:9001 mapper"
done
