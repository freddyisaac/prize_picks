#!/bin/sh

CAGE_ID=1
CAGE_STATUS=ACTIVE

while [[ $# -gt 0 ]]; do
	case $1 in
		-id)
			CAGE_ID="$2"
			shift
			shift
			;;
		-s|--status)
			CAGE_STATUS="$2"
			shift
			shift
			;;
		-h)
			echo "cage_status.sh -id <cage id> -s <ACTIVE|DOWN>"
			shift
			;;
	esac
done

curl -X POST http://localhost:8000/cage/${CAGE_ID}/status/${CAGE_STATUS}

