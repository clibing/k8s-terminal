#!/usr/bin/env bash

set -e

DIST_PREFIX="k8s-terminal"
TARGET_DIR="dist"
PLATFORMS="darwin/amd64 darwin/arm64 windows/amd64 windows/386 windows/arm"
# PLATFORMS_LINUX="linux/amd64 linux/arm linux/arm64"
PLATFORMS_LINUX="linux/amd64"

BUILD_VERSION=$(cat version)
BUILD_DATE=$(date "+%F %T")
COMMIT_SHA1=$(git rev-parse HEAD)

rm -rf ${TARGET_DIR}
mkdir ${TARGET_DIR}

for pl in ${PLATFORMS}; do
    export GOOS=$(echo ${pl} | cut -d'/' -f1)
    export GOARCH=$(echo ${pl} | cut -d'/' -f2)
    export TARGET=${TARGET_DIR}/${DIST_PREFIX}_${GOOS}_${GOARCH}
    if [ "${GOOS}" == "windows" ]; then
        export TARGET=${TARGET_DIR}/${DIST_PREFIX}_${GOOS}_${GOARCH}.exe
    fi

    echo "build => ${TARGET}"
    go build -trimpath -o ${TARGET} \
            -ldflags    "-X 'main.version=${BUILD_VERSION}' \
                        -X 'main.buildDate=${BUILD_DATE}' \
                        -X 'main.commitID=${COMMIT_SHA1}'\
                        -w -s"

done

for pl in ${PLATFORMS_LINUX}; do
    export GOOS=$(echo ${pl} | cut -d'/' -f1)
    export GOARCH=$(echo ${pl} | cut -d'/' -f2)
    export TARGET=${TARGET_DIR}/${DIST_PREFIX}_${GOOS}_${GOARCH}

    if [ "${GOARCH}" == "amd64" ]; then
        export CC=x86_64-linux-musl-gcc  
        export CXX=x86_64-linux-musl-g++
    elif [ "${GOARCH}" == "arm" ]; then
        export  CC=arm-linux-gnueabi-gcc
        export CXX=arm-linux-gnueabi-g++
    elif [ "${GOARCH}" == "arm64" ]; then
        export  CC=arm64-linux-gnueabi-gcc
        export CXX=arm64-linux-gnueabi-g++
    fi

    echo "build => ${TARGET}"
    go build -trimpath -o ${TARGET} \
            -ldflags    "-X 'main.version=${BUILD_VERSION}' \
                        -X 'main.buildDate=${BUILD_DATE}' \
                        -X 'main.commitID=${COMMIT_SHA1}'\
                        -w -s"
done
