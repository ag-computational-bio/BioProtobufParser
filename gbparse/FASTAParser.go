package gbparse

import (
	"bufio"
	"git.computational.bio.uni-giessen.de/sbeyvers/protobuffiles/gocompiled"
	"io"
	"log"
	"regexp"
	"strings"
	"sync"
)

type FASTAParser struct {
	Output chan *gbparse.Fasta
}

func (fastaparser *FASTAParser) Init() {
	fastaparser.Output = make(chan *gbparse.Fasta, 10000000)
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
				parseFastaRecord(header, sequence, fastaparser.Output)
			}
			header = line
			sequence = ""
		} else {
			sequence += line
		}
	}
	// Letztes Record parsen
	parseFastaRecord(header, sequence, fastaparser.Output)
	// Waitgroup -> Done
	mainwg.Done()
}

func parseFastaRecord(header string, sequence string, output chan *gbparse.Fasta) {
	currentFastaRecord := &gbparse.Fasta{}
	regxaccession, _ := regexp.Compile("[A-Z]{2}_[A-Z0-9]+")
	regxaccessionVersion, _ := regexp.Compile("[A-Z]{2}_[A-Z0-9]+[.]?[0-9]+")
	currentFastaRecord.ACCESSION = regxaccession.FindAllString(header, -1)[0]
	currentFastaRecord.VERSION = regxaccessionVersion.FindAllString(header, -1)[0]
	currentFastaRecord.SEQUENCE = sequence
	currentFastaRecord.HEADER = header
	output <- currentFastaRecord
}
