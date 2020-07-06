#!/usr/bin/env bash

echo "==> Checking documentation terraform blocks are formatted..."

files=$(find ./website -type f -name "*.html.markdown")
error=false

for f in $files; do
  terrafmt diff -c -q "$f" || error=true
done

if ${error}; then
  exit 1
fi

exit 0
