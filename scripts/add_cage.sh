#!/bin/sh

if [ "$1" == "" ]; then
CAGE_DIET="H"
else
CAGE_DIET=$1
fi

curl -v -X POST http://localhost:8000/cage/${CAGE_DIET}/add

