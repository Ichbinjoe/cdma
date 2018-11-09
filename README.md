# cdma

This executable simulates CDMA. This was written for a class.

## Compiling

CDMA uses no non-standard golang libraries. Compiling this program requires a
standard [golang installation](https://golang.org/doc/install).

To get and build this repository, run the following:

```
go get github.com/ichbinjoe/cdma
go build github.com/ichbinjoe/cdma
```

This will create a `cdma` executable within your current directory.

## Running

`cdma` takes one optional flag, -streams. It defaults to three.

Next, enter in the CDMA codes for each stream. After those are entered, enter
the messages you wish to be passed. Once entered, `cdma` simulates encoding and
decoding a simple no-interference CDMA stream.
