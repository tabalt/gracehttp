#!/bin/bash

action=$1

ps_result=$(ps aux | awk '/gracehttpdemo/{print $2}')
echo -e $ps_result

if [ "$1" = "" ] 
then
	kill -9 $ps_result
else
	kill -$action $ps_result
fi


