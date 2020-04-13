# todo-api-go
Yet Another ToDo App - I'm using this one to learn building web services in Go

## Build

```shell script
$ go build -o todo-api
```

## Test

```shell script
$ go test
```

## Run

You will need to export the following environment variables:

* `GORM_DIALECT` : currently only `sqlite3` is supported.
* `GORM_CONNECTION_STRING` : for `sqlite3` this will be the path to the database file. It will be created if it does not exist.
* `HOST_ADDRESS` : This should the IP address and port on in the form of `<HOST_IP>:<PORT>`

Example:

```shell script
$ export GORM_DIALECT=sqlite3
$ export GORM_CONNECTION_STRING=test.db
$ export HOST_ADDRESS=127.0.0.1:8000
```

The binary (once built) can then be run:

```shell script
$ ./todo-api
2020/04/13 16:54:18 Starting web server on 127.0.0.1:8000

```
