# isUTF8

Detect whether a file is well-formed UTF-8 or not.

isUTF8 is written in [Go](https://golang.org/) and uses [memory mapped files](https://en.wikipedia.org/wiki/Memory-mapped_file) to run as quickly as possible.  It uses the [golang.org/x/sys/unix](https://godoc.org/golang.org/x/sys/unix) package and will probably run only on Unix-like systems (e.g., MacOS, Linux).

On a 2016 MacBook Pro, isUTF8 checked a 1GB file in around 1 second, about 30% faster than a nearly identical C program compiled with gcc's ‑O3 flag (run times will vary depending on the system and how much of the file is already in memory cache).

For information about well-formed UTF-8 see [The Unicode Standard](http://www.unicode.org/versions/Unicode9.0.0/), [Chapter 3 Conformance](http://www.unicode.org/versions/Unicode9.0.0/ch03.pdf), Table 3-7 Well-Formed UTF-8 Byte Sequences.

## Prerequisites
[Go](https://golang.org/) programming language.

[golang.org/x/sys/unix](https://godoc.org/golang.org/x/sys/unix) package.  Not part of the standard Go installation so it must be installed separately.

```
go get golang.org/x/sys/unix
```

## Building
```
git clone https://github.com/mfuhr/isUTF8.git
cd isUTF8
go build
go test
```

## Examples
```
$ ./isUTF8 testdata/test_utf8.txt
true testdata/test_utf8.txt
$ echo $?
0
$ ./isUTF8 testdata/test_latin1.txt
false testdata/test_latin1.txt
$ echo $?
1
```

## Status
In active development (May 2017).  Behavior, especially the output, subject to change.

