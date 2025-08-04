#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

cd $SCRIPT_DIR/..

echo "go compile..."

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o reflex.bin

go test -c .

cd - > /dev/null
