#!/bin/bash

basedir=$( cd $(dirname $0) && pwd)
cd $basedir > /dev/null

bin=etclient

ps x | grep ${bin} | grep -v grep | awk -F ' ' '{print $1}' | xargs kill -9

cd bin

user=100

for ((user=100000;user<100100;user++))
do
  for ((i=0;i<5;i++))
  do
    key=`expr 110 + $i`
    akey=`printf "%x" $key`
    tmp=`echo $akey | perl -pe 's/([0-9a-f]{2})/chr hex $1/gie'`
    ./${bin} login,write -u $user -p $user -k $tmp -v $tmp
  done
done

cd -