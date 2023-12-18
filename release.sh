#!/bin/sh

set -ex 

if [ -z ${RELEASE_VERSION+x} ]
then
    echo 'Error! $RELEASE_VERSION is required.'
    exit 64
fi

echo $RELEASE_VERSION

goreleaser check

tag_and_push() {
    local component="$1"
    git tag -a "$component/$RELEASE_VERSION" -m "release $component/$RELEASE_VERSION"
    git push origin $component/$RELEASE_VERSION
}


tag_and_push "sdk"

# main
git tag -a $RELEASE_VERSION -m "release $RELEASE_VERSION"
git push origin $RELEASE_VERSION

make pulumi go_sdk nodejs_sdk python_sdk
