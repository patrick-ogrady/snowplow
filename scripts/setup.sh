#!/bin/bash
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


# Update apt
sudo apt update;

# Provision swap
sudo fallocate -l 16G /swapfile;
sudo chmod 600 /swapfile;
sudo mkswap /swapfile;
sudo swapon /swapfile;
echo '/swapfile swap swap defaults 0 0' | sudo tee -a /etc/fstab;

# Install Google Cloud Monitoring agent
cd ..;
curl -sSO https://dl.google.com/cloudagents/add-monitoring-agent-repo.sh;
sudo bash add-monitoring-agent-repo.sh;
sudo apt update;
sudo apt install stackdriver-agent -y;
sudo service stackdriver-agent start;
cd snowplow;

# Install git
sudo apt install git -y;

# Install golang
sudo apt install golang -y;

# Install docker
sudo apt install apt-transport-https ca-certificates curl software-properties-common -y;
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -;
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable";
sudo apt update;
sudo apt install docker-ce -y;

# Install make
sudo apt install make -y;

# Install zip
sudo apt install zip -y;

# Install snowplow
make install;

echo "snowplow installation succeeded!";
