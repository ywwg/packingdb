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
