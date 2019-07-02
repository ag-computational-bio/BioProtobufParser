package gbparse

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"strings"
	"sync"
)

type FASTAParser struct {
	output chan *Fasta
}

func (fastaparser *FASTAParser) Init() {
	fastaparser.output = make(chan *Fasta, 10000000)
}

func (fastaparser FASTAParser) ReadAndParseFile(reader io.Reader, mainwg *sync.WaitGroup) {

	scanner := bufio.NewScanner(reader)
	header := ""
	sequence := ""

	// Iterate over Lines
	for scanner.Scan() {

		// Handle occured errors
		if scanner.Err() != nil {
			log.Fatal(scanner.Err())
		}

		line := scanner.Text()
		if strings.HasPrefix(line, ">") {
			if header != "" {
				go parseFastaRecord(header, sequence, mainwg, fastaparser.output)
			}
			header = line
			sequence = ""
		} else {
			sequence += line
		}
	}
	// Letztes Record parsen
	go parseFastaRecord(header, sequence, mainwg, fastaparser.output)
	// Waitgroup -> Done
	mainwg.Done()
}

func parseFastaRecord(header string, sequence string, wg *sync.WaitGroup, output chan *Fasta) {
	wg.Add(1)
	currentFastaRecord := &Fasta{}
	regxaccession, _ := regexp.Compile("[A-Z]{2}_[A-Z0-9]+")
	regxaccessionVersion, _ := regexp.Compile("[A-Z]{2}_[A-Z0-9]+[.]?[0-9]+")
	currentFastaRecord.ACCESSION = regxaccession.FindAllString(header, -1)[0]
	currentFastaRecord.VERSION = regxaccessionVersion.FindAllString(header, -1)[0]
	currentFastaRecord.SEQUENCE = sequence
	currentFastaRecord.HEADER = header
	output <- currentFastaRecord
	wg.Done()
}
