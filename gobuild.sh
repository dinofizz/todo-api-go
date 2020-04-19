#!/bin/sh

contains() {
    string="$1"
    substring="$2"
    if test "${string#*$substring}" != "$string"
    then
        return 0    # $substring is in $string
    else
        return 1    # $substring is not in $string
    fi
}

contains $1 "amd64" && go build -o todo-api -v
contains $1 "arm/v7" && env CC=arm-linux-gnueabi-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 go build -o todo-api -v
exit 0
