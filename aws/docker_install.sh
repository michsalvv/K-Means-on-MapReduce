#!/bin/bash
service docker start
usermod -a -G docker ec2-user
id ec2-user
chmod 666 /var/run/docker.sock
curl -L https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m) -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
