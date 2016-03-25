[![Build Status](https://travis-ci.org/fd0/grobi.svg?branch=master)](https://travis-ci.org/fd0/grobi)

# grobi

This program watches for changes in the available outputs (e.g. when a monitor
is connected/disconnected) and will automatically configure the currently used
outputs via RANDR according to configurable profiles.

# Installation

Grobi requires Go version 1.3 or newer to compile. To build grobi, run the
following command:

```shell
$ go run build.go
```

Afterwards please find a binary of grobi in the current directory:
```
$ ./grobi --help
Usage:
  grobi [OPTIONS] <command>

Application Options:
  -v, --verbose   Be verbose (false)
  -C, --config=   Read config from this file
  -n, --dry-run   Only print what commands would be executed without actually runnig them
  -i, --interval= Number of seconds between polls, set to zero to disable polling (5)
  -p, --pause=    Number of seconds to pause after a change was executed (2)

Help Options:
  -h, --help      Show this help message

Available commands:
  apply    apply a rule
  update   update outputs
  version  display version
  watch    watch for changes
```

# Development

Grobi is developed using the build tool [gb](https://getgb.io). It needs at
least Go 1.4 and can be installed by running the following command:

```shell
$ go get github.com/constabulary/gb/...
```

The program can be compiled using `gb` as follows:

```shell
$ gb build
```

# Compatibility

Grobi follows [Semantic Versioning](http://semver.org) to clearly define which
versions are compatible. The configuration file and command-line parameters and
user-interface are considered the "Public API" in the sense of Semantic
Versioning.
