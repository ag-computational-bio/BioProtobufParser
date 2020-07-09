package gbparse

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	bioproto "github.com/ag-computational-bio/BioProtobufSchemas/go"
)

type GBParser struct {
}

func (parser *GBParser) ReadAndParseFile(reader io.Reader, output chan *bioproto.Genbank) {

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
				parseGBRecord(&lines, recordStart, featureStart, sequenceStart, currentLine, output)
			} else {
				//mainwg.Add(1)
				parseGBRecord(&lines, recordStart, featureStart, currentLine, currentLine, output)
			}
			recordStart = currentLine
		}
		currentLine++
	}
}

func parseGBRecord(lines *[]string, startpoint int, startpointqual int, startpointseq int, startpointnext int, output chan *bioproto.Genbank) {
	// DEBUG: fmt.Println(startpoint, startpointqual, startpointseq, startpointnext)
	currentGenbankRecord := &bioproto.Genbank{}
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

func parseHeader(lines []string, gbRecord *bioproto.Genbank) {

	currentRef := &bioproto.Reference{}
	var RefList []*bioproto.Reference

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
				gbRecord.ACCESSION = findAccesions(line[12:])
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
				currentRef = &bioproto.Reference{}
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
					gbRecord.ACCESSION = append(gbRecord.ACCESSION, findAccesions(line[11:])...)
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

func findAccesions(line string) (accessions []string) {
	accregex, _ := regexp.Compile("[A-Z0-9_]+")
	return accregex.FindAllString(line, -1)
}

func parseQualifier(lines []string, gbRecord *bioproto.Genbank) {

	currentFeature := &bioproto.Feature{}
	var FeatureList []*bioproto.Feature
	var QualList []*bioproto.Qualifier

	currentQual := &bioproto.Qualifier{}

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
					currentQual = &bioproto.Qualifier{}
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
					currentQual = &bioproto.Qualifier{}
					currentFeature.QUALIFIERS = QualList
					FeatureList = append(FeatureList, currentFeature)
					currentFeature = &bioproto.Feature{}
					QualList = []*bioproto.Qualifier{}
				}
				currentFeature.TYPE = currentType
				isComplement, isJoined, isOrdered, minPos, maxPos, locations := parseAllLocations(line[21:])
				currentFeature.IsCompliment = isComplement
				currentFeature.IsJoined = isJoined
				currentFeature.IsOrdered = isOrdered
				currentFeature.MINPOSITION = minPos
				currentFeature.MAXPOSITION = maxPos
				currentFeature.LOCATIONS = locations
				initialized = true
			}
		}
	}
	QualList = append(QualList, currentQual)
	currentQual = &bioproto.Qualifier{}
	currentFeature.QUALIFIERS = QualList
	FeatureList = append(FeatureList, currentFeature)
	gbRecord.FEATURES = FeatureList
}

func parseSequence(lines []string, gbRecord *bioproto.Genbank) {

	seqRegEx, _ := regexp.Compile("[a-zA-Z]+")
	sequence := ""
	for _, line := range lines[1:] {
		sequence += strings.Join(seqRegEx.FindAllString(line, -1), "")
	}
	gbRecord.SEQUENCE = sequence
}

func parseSingleLocation(line string) (loc *bioproto.Location) {
	var startpos, stoppos int

	hasexternal := strings.ContainsRune(line, ':')

	// Regexes
	positions, _ := regexp.Compile("[0-9]+")
	unknownsingleregx, _ := regexp.Compile("[>0-9<]+\\.[>0-9<]+")

	newLocation := bioproto.Location{}

	if hasexternal {
		idsplit := strings.Split(line, ":")
		newLocation.EXTERNALREFERENCE = idsplit[0]
		line = idsplit[1]
	}

	newLocation.SITEBETWEEN = strings.ContainsRune(line, '^')
	newLocation.UNKNOWNLB = strings.ContainsRune(line, '<')
	newLocation.UNKNOWNUB = strings.ContainsRune(line, '>')
	newLocation.UNKNOWNSINGLESITE = unknownsingleregx.MatchString(line)

	locStrings := positions.FindAllString(line, -1)

	if len(locStrings) == 1 {
		startpos, _ = strconv.Atoi(locStrings[0])
		stoppos = startpos
	} else if len(locStrings) == 2 {
		startpos, _ = strconv.Atoi(locStrings[0])
		stoppos, _ = strconv.Atoi(locStrings[1])
	} else {
		fmt.Println(line)
		panic("regex for position finding failed, 0 or more than 2 results detected: abort")
	}

	newLocation.START = int64(startpos)
	newLocation.STOP = int64(stoppos)

	return &newLocation
}

func parseAllLocations(line string) (isComplement, isJoined, order bool, maxloc, minloc int64, locations []*bioproto.Location) {
	regxComp, _ := regexp.Compile("complement")
	regxJoin, _ := regexp.Compile("join")
	regxOrder, _ := regexp.Compile("order")
	ismulti := strings.ContainsRune(line, ',')

	if ismulti {
		multiEntrys := strings.Split(line, ",")
		for _, entry := range multiEntrys {
			locations = append(locations, parseSingleLocation(entry))
		}
	} else {
		locations = append(locations, parseSingleLocation(line))
	}

	isComplement = false
	isJoined = false
	order = false
	if regxComp.MatchString(line) {
		isComplement = true
	}
	if regxJoin.MatchString(line) {
		isJoined = true
	}

	if regxOrder.MatchString(line) {
		order = true
	}

	return isComplement, isJoined, order, locations[0].START, locations[len(locations)-1].STOP, locations
}

func printRecord(gbRecord *bioproto.Genbank) {
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
		fmt.Println(y.TYPE, y.IsCompliment, y.MINPOSITION, y.MAXPOSITION)
		for _, v := range y.QUALIFIERS {
			fmt.Println(v.Key + "=" + v.Value)
		}
	}

}
