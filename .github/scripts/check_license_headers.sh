#!/bin/sh

SCRIPT_PATH="$(dirname "$0")"
find . -type f -name "*.go" \( -exec "$SCRIPT_PATH"/check_license_header.sh {} \; -o -print \) | awk 'BEGIN{count=0} {count++; print "Error: File "$0" is missing a correct license header"} END{if(count != 0){ exit 1 }}' >&2
