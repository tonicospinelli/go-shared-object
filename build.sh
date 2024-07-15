#!/bin/sh

export LD_LIBRARY_PATH="${LD_LIBRARY_PATH}:."

apk update
apk add --no-cache build-base go

make libcsv
