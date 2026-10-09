#!/bin/sh
# Fails when a production Go file grows past the 60-line ceiling.
# Test files are exempt; so is anything generated.
set -eu

limit=60
bad=0

for f in $(find . -name '*.go' -not -path './.git/*' -not -name '*_test.go' -not -path './backend/android/apk/overlay.go'); do
	n=$(wc -l < "$f")
	if [ "$n" -gt "$limit" ]; then
		echo "$n $f"
		bad=1
	fi
done

if [ "$bad" -ne 0 ]; then
	echo "files above $limit lines (production only; tests and generated files are exempt)"
	exit 1
fi
