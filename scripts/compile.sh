#!/bin/sh
VERSION=$1;

xgo -go go-1.15.5 --targets=darwin/*,windows/*,linux/* -out
"bin/avalanche-runner-${VERSION}" .;

# Rename some files
mv "bin/avalanche-runner-${VERSION}-darwin-10.6-amd64" "bin/avalanche-runner-${VERSION}-darwin-amd64" 
mv "bin/avalanche-runner-${VERSION}-windows-4.0-amd64.exe" "bin/avalanche-runner-${VERSION}-windows-amd64"

# Tar all files
cd bin || exit;
for i in *; do tar -czf "$i.tar.gz" "$i" && rm "$i"; done
