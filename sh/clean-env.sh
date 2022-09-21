#!/bin/bash
docker-compose stop
docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)
