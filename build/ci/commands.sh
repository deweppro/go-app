#!/usr/bin/env bash

source `dirname $0`/vars.sh

envup(){
  envdown
  ${dockercmd} up -d
}

envdown(){
  ${dockercmd} down
}

alltests(){
  ${dockercmd} exec app go test ./...
}

onetest(){
  ${dockercmd} exec app go test -run "$2" ./...
}