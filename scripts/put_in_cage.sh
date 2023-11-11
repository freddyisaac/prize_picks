#!/bin/sh

CAGE_ID=1
DINO_FILE="dino_c.json"

while [[ $# -gt 0 ]]; do
        case $1 in
                -id)
                        CAGE_ID="$2"
                        shift
                        shift
                        ;;
                -f)
                        DINO_FILE="$2"
                        shift
                        shift
                        ;;
                -h)
                        echo "put_in_cage.sh -id <cage id> -f <file with dino data>"
                        shift
                        ;;
        esac
done

curl -X POST -d @${DINO_FILE} http://localhost:8000/cage/${CAGE_ID}/add_dino
