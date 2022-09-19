#! /usr/bin/bash

go build -o clientGrep client/*.go
go build -o master server/master/*.go
go build -o mapper server/worker/mapper/*.go
go build -o reducer server/worker/reducer/*.go
