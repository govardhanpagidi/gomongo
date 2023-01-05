#!/usr/bin/env bash

set -x

cd schema-gen
go build

cd -


FILE=diff.json
rm -rf $FILE

./schema-gen/schema-gen compare "$FILE"

if test -f "$FILE"; then
 cat $FILE
fi

