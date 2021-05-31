#!/bin/bash

basedir=$( cd $(dirname $0) && pwd)
cd $basedir > /dev/null

bin=etclient

ps x | grep ${bin} | grep -v grep | awk -F ' ' '{print $1}' | xargs kill -9

cd bin && ./${bin} $@
cd -