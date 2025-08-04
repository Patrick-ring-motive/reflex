#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
echo "go updating imports..."
go get -u -t $SCRIPT_DIR/../...

cd $SCRIPT_DIR/..

echo "go tidying imports..."
go mod tidy

cd - > /dev/null