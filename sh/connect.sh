#! /usr/bin/bash

for reducer in $(docker ps -a --format "{{.Names}}" --filter "name=reducer")
do 
    terminator --new-tab -x "docker exec -it $reducer ./bin/reducer 172.17.0.1:9001; bash"
    #gnome-terminal -- bash -c 'docker exec' + $reducer + 'bash -c "./bin/reducer 172.17.0.1:9001" ; bash'
done
