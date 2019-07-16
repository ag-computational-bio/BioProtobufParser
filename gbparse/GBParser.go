package gbparse

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"sync"
)

type GBParser struct {
	Output chan *Genbank
}

func (gb *GBParser) Init() {
	gb.Output = make(chan *Genbank, 1000000)
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
				//mainwg.Add(1)
				parseGBRecord(&lines, recordStart, featureStart, sequenceStart, currentLine, mainwg, gb.Output)
			} else {
				//mainwg.Add(1)
				parseGBRecord(&lines, recordStart, featureStart, currentLine, currentLine, mainwg, gb.Output)
			}
			recordStart = currentLine
		}
		currentLine++
	}
	// Waitgroup -> Done
	mainwg.Done()
}

func parseGBRecord(lines *[]string, startpoint int, startpointqual int, startpointseq int, startpointnext int, wg *sync.WaitGroup, output chan *Genbank) {
	// DEBUG: fmt.Println(startpoint, startpointqual, startpointseq, startpointnext)
	currentGenbankRecord := &Genbank{}
	parseHeader((*lines)[startpoint:startpointqual], currentGenbankRecord)
	parseQualifier((*lines)[startpointqual:startpointseq], currentGenbankRecord)
	if startpointseq != startpointnext {
		parseSequence((*lines)[startpointseq:startpointnext], currentGenbankRecord)
	}
	//Encode Comment to b64 before it is written to output-channel
	currentGenbankRecord.COMMENT = base64.RawStdEncoding.EncodeToString([]byte(currentGenbankRecord.COMMENT))
	//printRecord(currentGenbankRecord)
	output <- currentGenbankRecord
	//wg.Done()
}

func parseHeader(lines []string, gbRecord *Genbank) {

	currentRef := &Reference{}
	var RefList []*Reference

	beforeCategory := ""
	currentReference := 0

	for _, line := range lines {
		if len(line) >= 12 {
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
				gbRecord.DBLINK = append(gbRecord.DBLINK, line[12:])
			case "KEYWORDS    ":
				beforeCategory = "KEYWORDS"
				gbRecord.KEYWORDS = line[12:]
			case "SOURCE      ":
				beforeCategory = "SOURCE"
				gbRecord.SOURCE = line[12:]
			case "  ORGANISM  ":
				beforeCategory = "ORGANISM"
				gbRecord.ORGANISM = append(gbRecord.ORGANISM, line[12:])
			case "REFERENCE   ":
				beforeCategory = "REFERENCE"
				if currentReference >= 1 {
					RefList = append(RefList, currentRef)
				}
				currentRef = &Reference{}
				currentRef.ORIGIN = line[12:]
				currentReference++
				currentRef.Number = int32(currentReference)
			case "  AUTHORS   ":
				beforeCategory = "  AUTHORS"
				currentRef.AUTHORS = line[12:]
			case "  CONSRTM   ":
				beforeCategory = "  CONSRTM"
				currentRef.CONSRTM = line[12:]
			case "  TITLE     ":
				beforeCategory = "  TITLE"
				currentRef.TITLE = line[12:]
			case "  JOURNAL   ":
				beforeCategory = "  JOURNAL"
				currentRef.JOURNAL = line[12:]
			case "   PUBMED   ":
				beforeCategory = "  PUBMED"
				currentRef.PUBMED = line[12:]
			case "COMMENT     ":
				beforeCategory = "COMMENT"
				gbRecord.COMMENT = line[12:]
			default:
				switch beforeCategory {
				case "COMMENT":
					gbRecord.COMMENT += "\n" + line[12:]
				case "  AUTHORS":
					currentRef.AUTHORS += line[11:]
				case "  CONSRTM":
					currentRef.CONSRTM += line[11:]
				case "  TITLE":
					currentRef.TITLE += line[11:]
				case "  JOURNAL":
					currentRef.JOURNAL += line[11:]
				case "   PUBMED":
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
					gbRecord.DBLINK = append(gbRecord.DBLINK, line[12:])
				case "KEYWORDS":
					gbRecord.KEYWORDS += line[11:]
				case "SOURCE":
					gbRecord.SOURCE += line[11:]
				case "ORGANISM":
					gbRecord.ORGANISM = append(gbRecord.ORGANISM, line[12:])
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
	var QualList []*Qualifier

	currentQual := &Qualifier{}

	wordRegEx, _ := regexp.Compile("[^\\s]+")
	qualifier, _ := regexp.Compile("^[/].*[=]?")

	currentType := ""

	initialized := false
	for _, line := range lines {

		switch line[0:6] {
		case "      ":
			if qualifier.MatchString(line[21:]) {
				splits := strings.SplitN(line, "=", 2)
				if (currentQual.Key) != "" {
					QualList = append(QualList, currentQual)
					currentQual = &Qualifier{}
				}
				currentQual.Key = splits[0][21:]
				if len(splits) == 2 {
					currentQual.Value = splits[1]
				} else if len(splits) == 1 {
					currentQual.Value = ""
				}
			} else {
				if currentQual.Key != "/translation" {
					currentQual.Value += line[20:]
				} else {
					currentQual.Value += line[21:]
				}
			}
		case "CONTIG":
			gbRecord.CONTIG = line[12:]
		default:
			currentType = wordRegEx.FindString(line[0:21])
			if currentType != "FEATURES" {
				if initialized {
					QualList = append(QualList, currentQual)
					currentQual = &Qualifier{}
					currentFeature.QUALIFIERS = QualList
					FeatureList = append(FeatureList, currentFeature)
					currentFeature = &Feature{}
					QualList = []*Qualifier{}
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
	QualList = append(QualList, currentQual)
	currentQual = &Qualifier{}
	currentFeature.QUALIFIERS = QualList
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
	regxFromTo, _ := regexp.Compile("[>0-9<]+")
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
		for _, v := range y.QUALIFIERS {
			fmt.Println(v.Key + "=" + v.Value)
		}
	}

}
