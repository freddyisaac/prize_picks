#!/bin/sh

PARAM=""
if [ "$1" != "" ]; then
	PARAM="?species="$1
fi

curl -v http://localhost:8000/v1/dino/list${PARAM} | jq

