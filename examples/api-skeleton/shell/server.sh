#!/bin/sh
echo "============`date +%F' '%T`==========="

file=/project/bin/program
cmd=api

getpid()
{
  docmd=`ps aux | grep ${file} | grep ${cmd} | grep -v 'grep' | grep -v '\.sh' | awk '{print $2}' | xargs`
  echo $docmd
}

start()
{
  pidstr=`getpid`
  if [ -n "$pidstr" ];then
    echo "running with pids $pidstr"
  else
     $file $cmd > /dev/null 2>&1 &
     sleep 1
     pidstr=`getpid`
     echo "start with pids $pidstr"
  fi
}

stop()
{
  pidstr=`getpid`
  if [ ! -n "$pidstr" ];then
     echo "not executed!"
     return
  fi
  echo "kill $pidstr"
  kill $pidstr
}

restart()
{
  stop
  sleep 1
  start
}

case "$1" in
  start)
    start
    ;;
  stop)
    stop
    ;;
  restart)
    restart
    ;;
esac
