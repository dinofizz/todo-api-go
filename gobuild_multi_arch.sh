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

echo "$(pwd)" > /log
contains $1 "amd64"
echo "$(ls -al)" >> /log
if [ $? -eq 0 ]
then
  go build -o todo-api -v
  return 0
fi

contains $1 "arm/v7"
if [ $? -eq 0 ]
then
  env CC=arm-linux-gnueabi-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 go build -o todo-api -v
  return 0
fi

exit 1
