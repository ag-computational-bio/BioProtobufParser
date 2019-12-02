package main

import (
	"log"
	"os"
	"sync"

	"git.computational.bio.uni-giessen.de/sbeyvers/golanggbffparser/gbparse"
)

func main() {
	parser := gbparse.GBParser{}
	parser.Init()

	file, err := os.Open("/home/marius/Downloads/archaea.7.rna.gbff")
	if err != nil {
		log.Println(err.Error())
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go parser.ReadAndParseFile(file, &wg)

	go func() {
		for foo := range parser.Output {
			log.Println(foo.GetACCESSION())
		}
	}()
	wg.Wait()

}
