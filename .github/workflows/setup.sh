#!/bin/sh

set -eux

mkdir -p /home/runner/go/bin
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
echo "::add-path::/home/runner/go/bin"
mkdir -p /home/runner/go/src/github/ngs/
ln -s /home/runner/work/akerun-slack-notifier/akerun-slack-notifier /home/runner/go/src/github/ngs/akerun-slack-notifier
