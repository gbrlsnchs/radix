#!/bin/bash

EXCLUDE=$(find . -type f -name '*.go' -not -path './vendor/*')
GOIMPORTS_LIST=$(goimports -l $EXCLUDE)

if [ -n "$GOIMPORTS_LIST" ]; then
  echo "goimports error:"
  goimports -d $EXCLUDE
  exit 1
fi
