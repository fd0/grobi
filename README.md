[![Build Status](https://travis-ci.org/fd0/grobi.svg?branch=master)](https://travis-ci.org/fd0/grobi)

# grobi
Watch for changes in outputs for the i3 window manager and react accordingly

# Requirements

Grobi requires Go version 1.4 or newer and the build tool [gb](https://getgb.io). The latter can be installed by running the following command:
```shell
$ go get github.com/constabulary/gb/...
```

# Installation

Get all dependencies:
```shell
$ gb vendor restore
```

Compile the program:
```shell
$ gb build
```
