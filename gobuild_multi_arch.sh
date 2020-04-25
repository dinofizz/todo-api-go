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

contains $1 "amd64"
if [ $? -eq 0 ]
then
  go build -o todo-api -v
  exit
fi

contains $1 "arm/v7"
if [ $? -eq 0 ]
then
  env CC=arm-linux-gnueabihf-gcc CXX=arm-linux-gnueabihf-g++ \
      CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 \
      go build -v -o todo-api
  exit
fi

exit 1
