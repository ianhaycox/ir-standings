#!/bin/bash
if [ -z $COV_CUTOVER ]; then
	COV_CUTOVER="50"
fi

FNAME=coverage.nomocks.out
if [ $1 ]; then
	FNAME=$1
fi

COV=$(go tool cover -func $FNAME | tail -n 1 | awk '{print $3}' | cut -d "." -f 1)
echo "Test coverage at $COV %"
if [[ "$COV" -lt "$COV_CUTOVER" ]]; then echo "Failing due to low test coverage"; exit 1; fi

