# Simple TCP Server [![Build Status](https://travis-ci.org/Negaihoshi/simple-tcp-server.svg?branch=master)](https://travis-ci.org/Negaihoshi/simple-tcp-server)

A simple TCP server that support mutiple connections.

***

### setup mock api

the mock server using mockingjay server

Build application
```
$ go get github.com/quii/mockingjay-server
$ cd $GOPATH/src/github.com/quii/mockingjay-server
$ ./build.sh
$ mockingjay-server -config=fake.yaml -monkeyConfig=monkey.yaml
```
more detail: https://github.com/quii/mockingjay-server

***

### set tcp server

command
```
$ go build
$ ./server
```

### set tcp client

command
```
$ go build
$ ./client
```

***

### status
you can check server status when server is up.

default url: http://127.0.0.1:1337/status

![2018-12-10 3 38 08](https://user-images.githubusercontent.com/1733006/49701903-17edfa00-fc2d-11e8-998f-cca6d8ad5844.png)
