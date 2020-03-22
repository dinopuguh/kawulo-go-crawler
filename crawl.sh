#!/bin/bash
BASEDIR=$(dirname $0)

if [ $# -eq 0 ]
  then
    go run "$BASEDIR/main.go"
  else
    go run "$BASEDIR/main.go" -data $1
fi
