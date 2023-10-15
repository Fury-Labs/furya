#!/bin/bash


# Make sure that dlv exists under GOPATH
go_path=$(go env GOPATH)
if [ -z "${go_path}" ]; then
    echo "GOPATH is not in go env"
    exit 1
fi


# Make sure that the binary is already running
furya_pid=$(pgrep furya)
if [ -z "${furya_pid}" ]; then
    echo "furyad is not running, cannot find its process ID"
    exit 1
fi


$go_path/bin/dlv attach --headless --listen=:2345 $furya_pid