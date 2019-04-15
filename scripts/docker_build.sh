#!/bin/bash

branch=`git branch|grep '*'|awk '{print $2}'`

docker build -t pttai-signal-server:${branch} .
