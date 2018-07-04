#! /bin/bash
export logfile=./logs/leisure.log
export binaryName=leisure

function reload() { 
    pid=`ps -ef|grep $binaryName |grep -v grep|awk '{print $2}'`
    kill -HUP $pid
    sleep 1
    newpid=`ps -ef|grep $binaryName|grep -v grep|awk '{print $2}'` 
    echo "reload..., pid=$newpid"
}

function start(){
   nohup ./$binaryName &
}

function stop() {
    pid=`ps -ef|grep $binaryName |grep -v grep|awk '{print $2}'`
    echo $pid 
   kill -9 $pid
   echo "loggerydraw stop"
}

function restart() {
    stop
    sleep 1
    start 
}

function build() {
    GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $binaryName
    echo "build success"
}

function iosbuild() {
    GO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o $binaryName
    echo "build success"
}

function tailf() {
   tail -f $logfile
}

function help() {
    echo "$0 start|stop|restart"
}

if [ "$1" == "" ]; then
    help
elif [ "$1" == "start" ];then
    start
elif [ "$1" == "stop" ];then
    stop
elif [ "$1" == "restart" ];then
    restart
elif [ "$1" == "reload" ];then
    reload
elif [ "$1" == "build" ];then
    build
elif [ "$1" == "iosbuild" ];then
    iosbuild    		
elif [ "$1" == "tail" ];then
    tailf
else
    help
fi
