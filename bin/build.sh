#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

echo "go formatting..."
gofmt -s -w $SCRIPT_DIR/.. 2>grep reflex

LAST_ERR=$?
if [ $LAST_ERR -ne 0 ]; then
  echo "go format error: $LAST_ERR"
  exit 1
fi

$SCRIPT_DIR/update-imports.sh

LAST_ERR=$?
if [ $LAST_ERR -ne 0 ]; then
  echo "update import error: $LAST_ERR"
  exit 1
fi

$SCRIPT_DIR/compile.sh

LAST_ERR=$?
if [ $LAST_ERR -ne 0 ]; then
  echo "compile error: $LAST_ERR"
  exit 1
fi

echo "pre-build checks passed successfully."