package gbparse

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"strings"

	gbparse "git.computational.bio.uni-giessen.de/sbeyvers/protobuffiles/go"
)

type FASTAParser struct {
}

func (fastaparser FASTAParser) ReadAndParseFile(reader io.Reader, output chan *gbparse.Fasta) {

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
				parseFastaRecord(header, sequence, output)
			}
			header = line
			sequence = ""
		} else {
			sequence += line
		}
	}
	// Letztes Record parsen
	parseFastaRecord(header, sequence, output)
	// Waitgroup -> Done
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
