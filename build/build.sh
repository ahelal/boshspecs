#!/bin/bash
set -e

# Get rid of existing binaries
rm -f dist/boshspecs*

# Check if VERSION variable set and not empty, otherwise set to default value
if [ -z "$VERSION" ]; then
  VERSION="0.0.1-dev"
fi
echo "Building application version $VERSION"

echo "Building default binary"
CGO_ENABLED=0 go build -ldflags "-s -w" -ldflags "-X main.version=${VERSION}" -o "dist/boshspecs" $PKG_SRC

if [ ! -z "${BUILD_ONLY_DEFAULT}" ]; then
    echo "Only default binary was requested to build"
    exit 0
fi

# Build binaries
OS_PLATFORM_ARG=(linux darwin windows)
OS_ARCH_ARG=(amd64)
for OS in ${OS_PLATFORM_ARG[@]}; do
  for ARCH in ${OS_ARCH_ARG[@]}; do
    echo "Building binary for $OS/$ARCH..."
    GOARCH=$ARCH GOOS=$OS CGO_ENABLED=0 go build -ldflags "-s -w" -ldflags "-X main.version=${VERSION}" -o "dist/boshspecs_$OS-$ARCH" $PKG_SRC
  done
done
