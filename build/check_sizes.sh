#!/bin/sh
# Fails when a production Go file grows past the 60-line ceiling.
#
# Test files are exempt, and so is anything generated. So is a file that is one
# cohesive block rather than code to read: a package doc, a type, an interface,
# a const block or a data table. Those are long because they are explained or
# because they hold a lot of data, not because they do many things — see
# AGENTS.md. They are listed here so the exemption is deliberate: a file that
# grows past the ceiling for any other reason gets no such excuse.
set -eu

limit=60

# One cohesive block each: docs, a type or interface, a const block, a bitmap.
exempt='^\./template/css/doc\.go$|^\./backend/android/doc\.go$|^\./font/font8x16\.go$|^\./svg/style\.go$|^\./svg/image\.go$|^\./svg/node\.go$|^\./svg/pattern\.go$|^\./backend/linux/driver\.go$|^\./backend/android/apk/config\.go$|^\./backend/android/apk/checklist\.go$|^\./backend/android/ndk/key_codes\.go$|^\./template/template\.go$|^\./template/style\.go$|^\./backend/android/apk/overlay\.go$'

bad=0

for f in $(find . -name '*.go' -not -path './.git/*' -not -name '*_test.go' -not -path './build/*'); do
	echo "$f" | grep -qE "$exempt" && continue
	n=$(wc -l < "$f")
	if [ "$n" -gt "$limit" ]; then
		echo "$n $f"
		bad=1
	fi
done

if [ "$bad" -ne 0 ]; then
	echo "files above $limit lines (production only; tests, generated and cohesive blocks are exempt)"
	exit 1
fi