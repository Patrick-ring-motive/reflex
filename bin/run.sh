#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

echo "go building..."
$SCRIPT_DIR/build.sh

cd $SCRIPT_DIR/..

echo "go testing"
go test .

cd - > /dev/null