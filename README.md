# packingdb

The #1 Enterprise solution for making sure you remember all your shit.

packingdb is web scale.

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

## Running

packingdb has two basic modes.  One just prints out a packing list.  The other can keep track of
what you've packed.

A "context" (TODO: change that term, it sucks) is sort of "where" you're traveling but also sort
of "what".  So "Key West" could be a context, but so could "skiing."  Do what you want.

### Simple Mode

For the first mode, just specify the trip context and a number of days:

```shell
./packingdb --context camping --days 5
```

This will generate a packing list for the camping context for a 5 day trip.

### File Mode


#### Creating a new file

To create a persistent file storing your packing data, add the --packfile flag:

```shell
./packingdb --context camping --days 5 --packfile ./mytrip
```

This will still print out the packing list, but will also save packed data to the file.

#### Recording stuff that's been packed

After the file is created, the --context and --days flags are ignored.  To pack an item, just name
it:

```shell
./packingdb --packfile ./mytrip --pack "boots"
```

If you mess up the name, the program will quit with an error.

Note that if you change your item .go files, subsequent runs of packingdb on your old packfile
*will* pick up the changes!

You can also edit the packfile yourself.  It's simple csv.
