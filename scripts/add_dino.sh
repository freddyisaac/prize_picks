#!/bin/sh

if [ "$1" == "" ]; then
DINO_FILE="dino_c.json"
else
DINO_FILE=$1
fi

curl -v -X POST -d @$DINO_FILE http://localhost:8000/dino/add

