#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

sudo -v

echo "==> Installing multimardown"
sudo apt-get install --assume-yes git libglib2.0-dev

git clone --recursive git://github.com/fletcher/peg-multimarkdown.git
cd peg-multimarkdown

./update_submodules.sh

make

sudo mv ./multimarkdown /usr/local/bin

