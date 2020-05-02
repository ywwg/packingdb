# packingdb

The #1 Enterprise solution for making sure you remember all your shit.

packingdb is web scale.

## Go env setup

Setting up Go is a little more complex than it needs to be.  First, set the GOPATH environment
variable to point to the root of this repository.

```shell
export GOPATH="$(pwd)"
```

You'll need the promptui module: (TODO: make this versioned!)

```shell
go get -v github.com/manifoldco/promptui
```

Subsequently, if you want to update modules:

```shell
go get -u all
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

packingdb has two binaries.  An older non-interactive program and a newer promptui-based interactive program.  The prompt-based program is much easier to use.

To start a new packing list, just specify a filename where the data will be stored:

```shell
./packingdb mytrip.csv
```

packingdb will ask you some basic questions about the trip, like the number of nights and the minimum and maximum predicted temperatures. Then you'll be presented with a list of all the configured Properties, and you can scroll through the list and select which ones apply to the current trip.

You can change all of these parameters later from the main menu.

You can also edit the packfile yourself.  It's simple csv.
