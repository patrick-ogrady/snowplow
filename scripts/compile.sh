#!/bin/sh
# Copyright (c) 2021 patrick-ogrady
#
# Permission is hereby granted, free of charge, to any person obtaining a copy of
# this software and associated documentation files (the "Software"), to deal in
# the Software without restriction, including without limitation the rights to
# use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
# the Software, and to permit persons to whom the Software is furnished to do so,
# subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
# FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
# COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
# IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
# CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

VERSION=$1;

xgo -go go-1.15.5 --targets=darwin/*,windows/*,linux/* -out
"bin/avalanche-runner-${VERSION}" .;

# Rename some files
mv "bin/avalanche-runner-${VERSION}-darwin-10.6-amd64" "bin/avalanche-runner-${VERSION}-darwin-amd64" 
mv "bin/avalanche-runner-${VERSION}-windows-4.0-amd64.exe" "bin/avalanche-runner-${VERSION}-windows-amd64"

# Tar all files
cd bin || exit;
for i in *; do tar -czf "$i.tar.gz" "$i" && rm "$i"; done