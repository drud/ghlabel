#!/bin/bash

set -o errexit

ARTIFACTS=$1
BASE_DIR=$PWD

sudo mkdir $ARTIFACTS && sudo chmod 777 $ARTIFACTS
export VERSION=$(git describe --tags --always)

# Generate OSX tarball/zipball
cd $BASE_DIR/bin/darwin/darwin_amd64
tar -czf $ARTIFACTS/ghlabel_osx.$VERSION.tar.gz ghlabel
zip $ARTIFACTS/ghlabel_osx.$VERSION.zip ghlabel

# Generate linux tarball/zipball
cd $BASE_DIR/bin/linux
tar -czf $ARTIFACTS/ghlabel_linux.$VERSION.tar.gz ghlabel
zip $ARTIFACTS/ghlabel_linux.$VERSION.zip ghlabel

# generate windows tarball/zipball
cd $BASE_DIR/bin/windows/windows_amd64
tar -czf $ARTIFACTS/ghlabel_windows.$VERSION.tar.gz ghlabel.exe
zip $ARTIFACTS/ghlabel_windows.$VERSION.zip ghlabel.exe

# Create the sha256 files
cd $ARTIFACTS
for item in *.tar.gz *.zip; do
  sha256sum $item > $item.sha256.txt
done