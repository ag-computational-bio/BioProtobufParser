package main

import (
	"../gbparse"
	"compress/gzip"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	//genbank()
	fasta()
}

func fasta() {
	var wg sync.WaitGroup
	wg.Add(1)
	file, err := os.Open("/home/basti/Schreibtisch/testdata/complete.1.1.genomic(1).fna")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	parser := gbparse.FASTAParser{}
	parser.Init()
	go parser.ReadAndParseFile(file, &wg)
	wg.Wait()
}

func genbank() {
	var wg sync.WaitGroup
	start := time.Now()
	wg.Add(1)
	parser := gbparse.GBParser{}
	parser.Init()
	resp, err := http.Get("https://ftp.ncbi.nih.gov/refseq/release/complete/complete.1.genomic.gbff.gz")
	resp.Close = true
	handleError(err)
	gz, err := gzip.NewReader(resp.Body)
	handleError(err)
	defer resp.Body.Close()
	defer gz.Close()
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
