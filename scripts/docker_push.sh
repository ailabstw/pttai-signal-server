#!/bin/bash

branch=`git branch|grep '*'|awk '{print $2}'`

docker tag pttai-signal-server:${branch} ailabstw/pttai-signal-server:latest
docker push ailabstw/pttai-signal-server:latest
