# BioProtobufParser

This tool parses Gbff and FASTA files and converts them to a [protocol buffer](https://developers.google.com/protocol-buffers) representation.
The corresponding protobuf files and schemas can be found [here](https://github.com/ag-computational-bio/BioProtobufSchemas).

## Usage


1. The first need [Go](https://golang.org/) installed (**version 1.14+ is required**), then you can use the below Go command to install BioProtobufParser.

```sh
$ go get -u github.com/ag-computational-bio/BioProtobufParser
```

2. Import it in your code:

```go
import "github.com/ag-computational-bio/BioProtobufParser"
```

## Quick Start

For 

```go
package main

import (

"github.com/ag-computational-bio/BioProtobufParser/gbparse"
schemas "github.com/ag-computational-bio/BioProtobufSchemas/go"

)

func main() {
    parser := gbparse.FASTAParser{}
    
    output := make(chan *schemas.Fasta)

    file := bufio.
    parser.Rea

}
```
