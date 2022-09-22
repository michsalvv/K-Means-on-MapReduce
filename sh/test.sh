#! /usr/bin/bash

echo master

for mapper in $(docker ps -a --format "{{.Names}}" --filter "name=mapper")
do 
    echo $mapper
done


for reducer in $(docker ps -a --format "{{.Names}}" --filter "name=reducer")
do 
    echo $reducer
done
