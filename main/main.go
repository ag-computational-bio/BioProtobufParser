package main

import (
	"../gbparse"
	"compress/gzip"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	start := time.Now()
	wg.Add(1)
	parser := gbparse.GBParser{}
	parser.Init()
	resp, err := http.Get("https://ftp.ncbi.nih.gov/refseq/release/complete/complete.1.genomic.gbff.gz")
	handleError(err)
	gz, err := gzip.NewReader(resp.Body)
	handleError(err)
	go parser.ReadAndParseFile(gz, &wg)
	wg.Wait()
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(elapsed, " seconds")
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
