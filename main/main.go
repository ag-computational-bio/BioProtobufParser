package main

import (
	"fmt"
	"gbparsertest2/gbparse"
	"gbparsertest2/generators"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	//genbank()
	//fasta()
	genbank_read()
}

func genbank_write(parser *gbparse.GBParser) {
	f, err := os.Create("./testdata/test.gbff")
	if err != nil {
		fmt.Println(err)
		return
	}
	for record := range parser.Output {
		_, err = f.WriteString(generators.GenerateGBfromproto(record))
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func fasta_read() {
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
	close(parser.Output)
	fastaWritefile(&parser)
}

func fastaWritefile(parser *gbparse.FASTAParser) {
	f, err := os.Create("./testdata/test.fasta")
	if err != nil {
		fmt.Println(err)
		return
	}
	for record := range parser.Output {
		_, err = f.WriteString(generators.GenerateFastafromproto(record))
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func genbank_read() {
	var wg sync.WaitGroup
	start := time.Now()
	wg.Add(1)
	parser := gbparse.GBParser{}
	parser.Init()
	//resp, err := http.Get("https://ftp.ncbi.nih.gov/refseq/release/complete/complete.1.genomic.gbff.gz")
	gz, err := os.Open("/home/basti/Schreibtisch/testdata/gbffrecords/xx00.gbff")
	//resp.Close = true
	if err != nil {
		log.Fatal(err)
	}
	//gz, err := gzip.NewReader(resp.Body)
	//defer resp.Body.Close()
	defer gz.Close()
	go parser.ReadAndParseFile(gz, &wg)
	wg.Wait()
	t := time.Now()
	close(parser.Output)
	genbank_write(&parser)
	elapsed := t.Sub(start)
	fmt.Println(elapsed, " seconds")
}
