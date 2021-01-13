#!/bin/bash

# Update apt
sudo apt update;

# Install git
sudo apt install git -y;

# Install golang
sudo apt install golang -y;
export PATH=$PATH:~/go/bin;

# Install docker
sudo apt install apt-transport-https ca-certificates curl software-properties-common -y;
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -;
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable";
sudo apt update;
sudo apt install docker-ce -y;

# Install make
sudo apt install make -y;

# Install snowplow
make install;

echo "snowplow installation succeeded!";
