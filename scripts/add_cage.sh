#!/bin/sh

CAGE_DIET="H"
CAGE_PARAM=""

usage_message() {
		echo "add_cage.sh -diet <H|C> -cap <cage dinosaur capacity>"
}

if [[ $# -eq 0 ]]; then
	usage_message
	exit
fi

while [[ $# -gt 0 ]]; do
        case $1 in
		-diet)
			CAGE_DIET="$2"
			shift
			shift
			;;
		-cap)
			CAGE_CAP="$2"
			shift
			shift
			CAGE_PARAM="?cap="$CAGE_CAP
			;;
		-h)
			usage_message
			shift
			exit
			;;
		esac
done

curl -v -X POST http://localhost:8000/v1/cage/${CAGE_DIET}/add${CAGE_PARAM}

