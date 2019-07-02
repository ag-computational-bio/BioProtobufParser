package gbparse

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"sync"
)

type GBParser struct {
	output chan *Genbank
}

func (gb *GBParser) Init() {
	gb.output = make(chan *Genbank, 1000000)
}

func (gb GBParser) ReadAndParseFile(reader io.Reader, mainwg *sync.WaitGroup) {

	scanner := bufio.NewScanner(reader)

	var lines []string
	currentLine := 0
	recordStart := 0
	featureStart := 0
	sequenceStart := 0
	hasSequence := false

	// Iterate over Lines
	for scanner.Scan() {

		// Handle occured errors
		if scanner.Err() != nil {
			log.Fatal(scanner.Err())
		}

		line := scanner.Text()
		// DEBUG: fmt.Println(line)
		lines = append(lines, line)
		if strings.HasPrefix(line, "FEATURES") {
			hasSequence = false
			featureStart = currentLine
		} else if strings.HasPrefix(line, "ORIGIN") {
			sequenceStart = currentLine
			hasSequence = true
		} else if strings.HasPrefix(line, "//") {
			if hasSequence {
				go parseGBRecord(&lines, recordStart, featureStart, sequenceStart, currentLine, mainwg, gb.output)
			} else {
				go parseGBRecord(&lines, recordStart, featureStart, currentLine, currentLine, mainwg, gb.output)
			}
			recordStart = currentLine
		}
		currentLine++
	}
	// Waitgroup -> Done
	mainwg.Done()
}

func parseGBRecord(lines *[]string, startpoint int, startpointqual int, startpointseq int, startpointnext int, wg *sync.WaitGroup, output chan *Genbank) {
	wg.Add(1)
	// DEBUG: fmt.Println(startpoint, startpointqual, startpointseq, startpointnext)
	currentGenbankRecord := &Genbank{}
	parseHeader((*lines)[startpoint:startpointqual], currentGenbankRecord)
	parseQualifier((*lines)[startpointqual:startpointseq], currentGenbankRecord)
	if startpointseq != startpointnext {
		parseSequence((*lines)[startpointseq:startpointnext], currentGenbankRecord)
	}
	output <- currentGenbankRecord
	fmt.Println(len(output))
	wg.Done()
}

func parseHeader(lines []string, gbRecord *Genbank) {

	currentRef := &Reference{}
	var RefList []*Reference

	beforeCategory := ""
	currentReference := 0

	for _, line := range lines {
		if len(line) > 12 {
			switch line[0:12] {
			case "LOCUS       ":
				beforeCategory = "LOCUS"
				gbRecord.LOCUS = line[12:]
			case "DEFINITION  ":
				beforeCategory = "DEFINITION"
				gbRecord.DEFINITION = line[12:]
			case "ACCESSION   ":
				beforeCategory = "ACCESSION"
				gbRecord.ACCESSION = line[12:]
			case "VERSION     ":
				beforeCategory = "VERSION"
				gbRecord.VERSION = line[12:]
			case "DBLINK      ":
				beforeCategory = "DBLINK"
				gbRecord.DBLINK = line[12:]
			case "KEYWORDS    ":
				beforeCategory = "KEYWORDS"
				gbRecord.KEYWORDS = line[12:]
			case "SOURCE      ":
				beforeCategory = "SOURCE"
				gbRecord.SOURCE = line[12:]
			case "  ORGANISM  ":
				beforeCategory = "ORGANISM"
				gbRecord.ORGANISM = line[12:]
			case "REFERENCE   ":
				beforeCategory = "REFERENCE"
				if currentReference >= 1 {
					RefList = append(RefList, currentRef)
				}
				currentRef = &Reference{}
				currentReference++
				currentRef.Number = int32(currentReference)
			case "  AUTHORS   ":
				beforeCategory = "  AUTHORS"
				currentRef.AUTHORS = line[12:]
			case "  TITLE     ":
				beforeCategory = "  TITLE"
				currentRef.TITLE = line[12:]
			case "  JOURNAL   ":
				beforeCategory = "  JOURNAL"
				currentRef.JOURNAL = line[12:]
			case "  PUBMED    ":
				beforeCategory = "  PUBMED"
				currentRef.PUBMED = line[12:]
			case "COMMENT     ":
				beforeCategory = "COMMENT"
				gbRecord.COMMENT = line[12:]
			default:
				switch beforeCategory {
				case "COMMENT":
					gbRecord.COMMENT += line[11:]
				case "  AUTHORS":
					currentRef.AUTHORS += line[12:]
				case "  TITLE":
					currentRef.TITLE += line[12:]
				case "  JOURNAL":
					currentRef.JOURNAL += line[12:]
				case "  PUBMED":
					currentRef.PUBMED += line[12:]
				case "LOCUS":
					gbRecord.LOCUS += line[11:]
				case "DEFINITION":
					gbRecord.DEFINITION += line[11:]
				case "ACCESSION":
					gbRecord.ACCESSION += line[11:]
				case "VERSION":
					gbRecord.VERSION += line[11:]
				case "DBLINK":
					gbRecord.DBLINK += line[11:]
				case "KEYWORDS":
					gbRecord.KEYWORDS += line[11:]
				case "SOURCE":
					gbRecord.SOURCE += line[11:]
				case "ORGANISM":
					gbRecord.ORGANISM += line[11:]
				}
			}
		}
	}
	// AddLastReference
	RefList = append(RefList, currentRef)

	gbRecord.REFERENCES = RefList
}

func parseQualifier(lines []string, gbRecord *Genbank) {

	currentFeature := &Feature{}
	var FeatureList []*Feature
	qualMap := make(map[string]string)

	wordRegEx, _ := regexp.Compile("[^\\s]+")
	qualifier, _ := regexp.Compile("^[/].*[=]?")

	initialized := false
	currentType := ""
	for _, line := range lines {

		switch line[0:6] {
		case "      ":
			if qualifier.MatchString(line[21:]) {
				splits := strings.Split(line, "=")
				currentType = splits[0][21:]
				if len(splits) == 2 {
					qualMap[currentType] = splits[1]
				} else if len(splits) == 1 {
					qualMap[currentType] = ""
				}
			} else {
				qualMap[currentType] += line[21:]
			}
		case "CONTIG":
			gbRecord.CONTIG = line[12:]
		default:
			currentType = wordRegEx.FindString(line[0:21])
			if currentType != "FEATURES" {
				if initialized {
					currentFeature.QUALIFIERS = qualMap
					FeatureList = append(FeatureList, currentFeature)
					currentFeature = &Feature{}
					qualMap = make(map[string]string)
				}
				currentFeature.TYPE = currentType
				x, y := getPositionFormat(line[21:])
				currentFeature.IsCompliment = x
				currentFeature.START = y[0]
				currentFeature.STOP = y[1]
				initialized = true
			}
		}
	}
	currentFeature.QUALIFIERS = qualMap
	FeatureList = append(FeatureList, currentFeature)
	gbRecord.FEATURES = FeatureList
}

func parseSequence(lines []string, gbRecord *Genbank) {

	seqRegEx, _ := regexp.Compile("[a-zA-Z]+")
	sequence := ""
	for _, line := range lines[1:] {
		sequence += strings.Join(seqRegEx.FindAllString(line, -1), "")
	}
	gbRecord.SEQUENCE = sequence
}

func getPositionFormat(line string) (bool, []string) {
	regxComp, _ := regexp.Compile("complement")
	regxFromTo, _ := regexp.Compile("[<]?[0-9]+[>]?")
	isComplement := false
	if regxComp.MatchString(line) {
		isComplement = true
	}
	return isComplement, regxFromTo.FindAllString(line, -1)
}

func printRecord(gbRecord *Genbank) {
	fmt.Println(gbRecord.LOCUS)
	fmt.Println(gbRecord.ACCESSION)
	fmt.Println(gbRecord.DEFINITION)
	fmt.Println(gbRecord.VERSION)
	fmt.Println(gbRecord.DBLINK)
	fmt.Println(gbRecord.KEYWORDS)
	fmt.Println(gbRecord.ORGANISM)
	fmt.Println(gbRecord.COMMENT)
	fmt.Println(len(gbRecord.SEQUENCE))

	for _, x := range gbRecord.REFERENCES {
		fmt.Println(x.Number)
		fmt.Println(x.TITLE)
		fmt.Println(x.AUTHORS)
		fmt.Println(x.JOURNAL)
	}

	for _, y := range gbRecord.FEATURES {
		fmt.Println(y.TYPE, y.IsCompliment, y.START, y.STOP)
		for k, v := range y.QUALIFIERS {
			fmt.Println(k + "=" + v)
		}
	}

}
