# Overview
A release tool built in go.

This was originally an experiment in another project to do all Go code
everywhere.  No shell, no external binaries, nothing.

There's a little work left to do here.

- Don't invoke `go build` anymore.  This executable would run afterward.
- Don't require the `release` subcommand. Always do that codepath.
- Always build all support platform executables.

# Dependencies
[Go](https://golang.org/doc/install)


# Installation
```
go get -u code.groovestomp.com/go-release-pkg
```

# Usage
```
go-release-pkg --help
```
