#!/usr/bin/env bash

cd ~/go
mkdir -p amazegardenai
cd amazegardenai
cp -r /var/app/current/* .

# Install dependencies.
go get ./...

# Build app
go build