#!/bin/sh

CAGE_DIET="H"
CAGE_PARAM=""

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
			echo "add_cage.sh -diet <H|C> -cap <cage dinosaur capacity>"
			shift
			;;
		esac
done

curl -v -X POST http://localhost:8000/cage/${CAGE_DIET}/add${CAGE_PARAM}

