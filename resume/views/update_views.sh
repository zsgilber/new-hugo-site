#!/usr/bin/env bash

for view in *.jq; do
    echo "processing view for $view"
    < ../resume.json jq -f "$view" >"${view%.*}".json
done
