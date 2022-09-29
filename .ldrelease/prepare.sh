#!/bin/bash

set -uex

yum install -y helm git

mkdir ~/.ssh
ssh-keyscan github.com > ~/.ssh/known_hosts
