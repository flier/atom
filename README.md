# Overview ![Travis Build](https://travis-ci.org/flier/atom.svg?branch=master)

Package atom provides integer codes (also known as atoms) for a fixed set of frequently occurring strings

An atom is an unique ID of a frequently occurring name, it could be more effective in store and compare in O(1).
``` go
hello := atom.New("hello")
```
The atom value will be unique in the same process base on the creating order.
``` go
hello.Value()
```
The atom name/value mapping could be saved as a snapshot and serialized to the disk / database.
``` go
data, cache := atom.Save()
```
When new process start next time or on the remote, load those will restore the atom name/value mapping.
``` go
atom.Load(data, cache)
```

# Install

Get and build package source.
```
go get -u github.com/flier/atom
```
Get and install generate atom command.
```
go install github.com/flier/atom/cmd/genatoms
```

# Internal

An atom is a 4 bytes uint32, which contains the offset and length of the atom name in a global buffer, or the embedded short name (less than 4 bytes) in the value.

The highest bit set to 0 means the atom contains a long string in the global buffer.
```
+--------+--------------------------+
|0| len  |          offset          |
+--------+--------------------------+
|01234567|01234567|01234567|01234567|
+--------+--------------------------+
```
The hightest bit set to 1 means the atom embedded a short string in the value. The first byte must less than 0x80.
```
+--------+--------------------------+
|1|str[0]| str[1] | str[2] | str[3] |
+--------+--------------------------+
|1|  3   | str[0] | str[1] | str[2] |
+--------+--------------------------+
|1|  2   |        | str[0] | str[1] |
+--------+--------------------------+
|1|  1   |        |        | str[0] |
+--------+--------------------------+
|01234567|01234567|01234567|01234567|
+--------+--------------------------+
```
It means the maximum length of atom name is 127 bytes.

`Atom.IsEmbedded` will reports whether the atom embedded the name

# Pregenerated Atoms

A build-in atoms buffer and cache could be generated with command:
```
$ genatoms -i atoms.txt -o atom.go -p atom -test
```
It will scan and extract all the Golang identifier from input file `atoms.txt`, the atom data and cache will be save to the output file `atom.go` with package name `atom`.
```
-case-insensitive
      case-insensitive atom (default true)
-format
      format the generated code (default true)
-i string
      read atom from the input file (default STDIN)
-o string
      write atom table to the output file (default STDOUT)
-p string
      generated package name (default "atom")
-test
      generate test table for the atom data
The extracted atom will be case insensitive by default.
```
