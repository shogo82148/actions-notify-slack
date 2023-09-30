#!/bin/bash

set -uex

CURRENT=$(cd "$(dirname "$0")" && pwd)

# if version starts with 'v', remove it
VERSION=${1#v}

MAJOR=$(echo "$VERSION" | cut -d. -f1)
MINOR=$(echo "$VERSION" | cut -d. -f2)
PATCH=$(echo "$VERSION" | cut -d. -f3)
WORKING=$CURRENT/.working

: clone
ORIGIN=$(git remote get-url origin)
rm -rf "$WORKING"
git clone "$ORIGIN" "$WORKING"
cd "$WORKING"

: update the version of package.json
cd "$WORKING/action"
jq ".version=\"$MAJOR.$MINOR.$PATCH\"" < package.json > .tmp.json
mv .tmp.json package.json
jq ".version=\"$MAJOR.$MINOR.$PATCH\"" < package-lock.json > .tmp.json
mv .tmp.json package-lock.json

: publish to GitHub
git tag -a "v$MAJOR.$MINOR.$PATCH" -m "release v$MAJOR.$MINOR.$PATCH"
git push origin "v$MAJOR.$MINOR.$PATCH"

git tag -fa "v$MAJOR" -m "release v$MAJOR.$MINOR.$PATCH"
git push -f origin "v$MAJOR"

cd "$CURRENT"
rm -rf "$WORKING"
