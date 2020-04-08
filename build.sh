#!/bin/bash

CGO_ENABLED=0

if [ "${1}" == "windows" ]
then
    echo "The Windows operating system is not supported!"
    exit;
elif [ "${1}" == "osx" ]
then
    go build -v -ldflags "-s -w" -o builds/image2webp-darwin-${2}
else
    go build -v -ldflags "-s -w" -o builds/image2webp-${1}-${2}
fi

cp ./config.json ./builds/

echo "build done!"
ls builds