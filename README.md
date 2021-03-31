# Embed

> Help to build embedded resource for Go's project.
>
> According to the resource configuration, the file resources are packaged into the program for embedded use.

## Installation

> To install embed command line program, use the following:
```
go get -u github.com/thecxx/embed/cmd/embed
```

## Default configuration

> Command: `embed init`

```
# embed build [-f embed.yaml]
---

# The package name
pkg: "embed"

# Source file path
path: "embed/embed.go"

# The compression method for embedded resource, it can be set to [no] [gz|gzip]
compress: "gz"

# Archive the data to a separate file
# archive: "embed/archive.go"

#
# Use the resource like:
#
# embed.Demo1.Size()
# embed.Demo1.Bytes()
# embed.Demo1.NewReader()
# ...
#
# Resource list
items:
  - name: "Demo1"
    file: "tests/demo1.txt"
    comment: "A demo file, \njust for test."

  - name: "Demo2"
    file: "tests/demo2.txt"
    comment: "Another demo file, \njust for test."

```

## Build

> Command: `embed build`