#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

mmd4_git="https://github.com/fletcher/MultiMarkdown-4"
mmd_latex_support_git="https://github.com/fletcher/peg-multimarkdown-latex-support"

pwd=$(pwd)

sudo -v

echo "==> Upgrading System"
sudo apt-get update
sudo apt-get dist-upgrade --assume-yes


echo "==> Installing utils"
sudo apt-get install --assume-yes git htop golang


echo "==> Installing LaTex"
sudo apt-get install --assume-yes \
     texlive-latex-base \
     texlive-latex-extra texlive-fonts-extra texlive-math-extra \
     texlive-xetex latex-xcolor \
     latexmk


echo "==> Installing MultiMarkdown"
sudo apt-get install --assume-yes libglib2.0-dev
mmd4="./mmd4"
git clone --recursive $mmd4_git $mmd4
cd $mmd4
     git submodule init
     git submodule update

     make

     sudo make install
     sudo make install-scripts
cd "$pwd"

echo "==> Adding MultiMarkdown LaTeX support files to texmf"
mmd_latex_support="./mmd_latex_support"
texmf="/etc/texmf/tex/latex"
git clone $mmd_latex_support_git $mmd_latex_support

sudo mkdir -p $texmf
sudo mv $mmd_latex_support "${texmf}/mmd"
