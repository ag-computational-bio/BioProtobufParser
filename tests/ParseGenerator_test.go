package tests

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"

	"git.computational.bio.uni-giessen.de/sbeyvers/golanggbffparser/gbparse"
	"git.computational.bio.uni-giessen.de/sbeyvers/golanggbffparser/generators"
)

func TestGBFFParserAndGenerator(t *testing.T) {

	// Add Waitgroup
	var wg sync.WaitGroup
	wg.Add(1)

	// Initialize Parser
	parser := gbparse.GBParser{}
	parser.Init()

	// Open testfile
	gbff, err := os.Open("../testfiles/Test50Entries.gbff")
	// Read file as bytebuffer for comparison
	filecontent, err := ioutil.ReadFile("../testfiles/Test50Entries.gbff")

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Parsing gbff file...")
	defer gbff.Close()
	go parser.ReadAndParseFile(gbff, &wg)
	wg.Wait()

	log.Println("Parsing complete, reading protobuf...")
	result := ""
	// Close Output channel before reading
	close(parser.Output)
	for record := range parser.Output {
		result += generators.GenerateGBfromproto(record)
	}

	// compare resultstring from protobuf object against raw string from file
	if result != string(filecontent) {
		t.Errorf("Parsed and generated file not equal!")
	}
}

func TestFASTAParserAndGenerator(t *testing.T) {

	// Add Waitgroup
	var wg sync.WaitGroup
	wg.Add(1)

	// Initialize Parser
	parser := gbparse.FASTAParser{}
	parser.Init()

	// Open testfile
	gbff, err := os.Open("../testfiles/Test50Entries.fasta")
	// Read file as bytebuffer for comparison
	filecontent, err := ioutil.ReadFile("../testfiles/Test50Entries.fasta")

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Parsing fasta file...")
	defer gbff.Close()
	go parser.ReadAndParseFile(gbff, &wg)
	wg.Wait()

	log.Println("Parsing complete, reading protobuf...")
	result := ""
	// Close Output channel before reading
	close(parser.Output)
	for record := range parser.Output {
		result += generators.GenerateFastafromproto(record)
	}

	// compare resultstring from protobuf object against raw string from file
	if result != string(filecontent) {
		t.Errorf("Parsed and generated file not equal!")
	}

}
