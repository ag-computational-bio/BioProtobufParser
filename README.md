# BioProtobufParser

This tool parses Gbff and FASTA files and converts them to a [protocol buffer](https://developers.google.com/protocol-buffers) representation.
The corresponding protobuf files and schemas can be found [here](https://github.com/ag-computational-bio/BioProtobufSchemas).

## Usage


1. Preparations: [Go](https://golang.org/) installed (**version 1.14+ is required**), then you can use the below Go command to install BioProtobufParser.

```sh
$ go get -u github.com/ag-computational-bio/BioProtobufParser
```

2. Import it in your code:

```go
import "github.com/ag-computational-bio/BioProtobufParser"
```

## Quick Start 
For Fasta files:
```go
package main

import (

"fmt"
"github.com/ag-computational-bio/BioProtobufParser/gbparse"
schemas "github.com/ag-computational-bio/BioProtobufSchemas/go"
"os"

)

func main() {
    
    # Create a parser
    parser := gbparse.FASTAParser{}
    
    # Create output channel
    output := make(chan *schemas.Fasta)

    # Open fasta file
    fasta, err := os.Open("FASTAFILENAME")
    defer fasta.Close()

    # Read and serialize fasta records in protobuf output channel
    parser.ReadAndParseFile(fasta, output)
    if err != nil {
        fmt.Errorf(err.Error())
    }
    
    # Use the serialized protobuf messages in the output channel

    # Close it afterwards
    close(output)


}
```

For GBFF files:
```go
package main

import (

"fmt"
"github.com/ag-computational-bio/BioProtobufParser/gbparse"
schemas "github.com/ag-computational-bio/BioProtobufSchemas/go"
"os"

)

func main() {
    # Create a parser
    parser := gbparse.GBParser{}
    
    # Create output channel
    output := make(chan *schemas.Genbank)
    
    # Open gbff file
    gbff, err := os.Open("GBFFFIlENAME")
    defer gbff.Close()
    # Read and serialize gbff records in protobuf output channel
    parser.ReadAndParseFile(gbff, output)
    if err != nil {
        fmt.Errorf(err.Error())
    }
    
    # Use the serialized protobuf messages in the output channel
    
    # Close it afterwards
    close(output)


}
```