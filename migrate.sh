#!/usr/bin/env bash

export GO111MODULE=on

ENV=${1:-development}
export CONFIG_JSON=config.$ENV.json

echo "-----------------------------"
echo " Environment : $ENV"
echo " Config      : $CONFIG_JSON"
echo "-----------------------------"

echo "Building..."
rm -rf bin/app
go build -o bin/app .

if [ $? -ne 0 ]; then
  echo "Build failed!"
  exit 1
fi

echo "Running migrations..."
./bin/app migrate up
