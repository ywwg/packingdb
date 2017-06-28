# packingdb

A thing that helps me figure out what to pack.

## Go env setup

Setting up Go is a little more complex than it needs to be.  First, set the GOPATH environment
variable to point to the root of this repository.

```shell
export GOPATH="$(pwd)"
```

Then, you should be able to build and run the program:

```shell
go install github.com/ywwg/packingdb && ./bin/packingdb
```

I created a little shell script that does this, hopefully it works!

## Other helpful go utils

You may also want Guru and Godef if you want to set up Eclipse:

```shell
export GOPATH="$(pwd)"
go get golang.org/x/tools/cmd/guru
go get github.com/npat-efault/godef
go install golang.org/x/tools/cmd/guru
go install github.com/npat-efault/godef
```

Then you can point your eclipse project to those binaries.
