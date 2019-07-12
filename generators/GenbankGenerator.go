package generators

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gbparsertest2/gbparse"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func GenerateGBfromproto(record *gbparse.Genbank) (fastarecord string) {
	var stringbuffer bytes.Buffer

	stringbuffer.WriteString(generateHeaderString(record))
	stringbuffer.WriteString("FEATURES             Location/Qualifiers\n")
	stringbuffer.WriteString(generateQualifierString(record, read_json()))
	if record.FEATURES != nil {

	}
	if record.CONTIG != "" {
		stringbuffer.WriteString("CONTIG      " + record.CONTIG + "\n")
	}
	stringbuffer.WriteString("//\n")
	return stringbuffer.String()
}

func generateHeaderString(record *gbparse.Genbank) (HeadString string) {
	var buffer bytes.Buffer
	buffer.WriteString("LOCUS       " + record.LOCUS + "\n")
	buffer.WriteString(formatStringWithNewlineChars("DEFINITION  "+record.DEFINITION, "            ", true))
	buffer.WriteString("ACCESSION   " + record.ACCESSION + "\n")
	buffer.WriteString("VERSION     " + record.VERSION + "\n")
	if len(record.DBLINK) > 0 {
		for i, line := range record.DBLINK {
			if i == 0 {
				buffer.WriteString("DBLINK      " + line + "\n")
			}
			buffer.WriteString("            " + line + "\n")
		}
	}
	buffer.WriteString("KEYWORDS    " + record.KEYWORDS + "\n")
	buffer.WriteString("SOURCE      " + record.SOURCE + "\n")
	for i, line := range record.ORGANISM {
		if i == 0 {
			buffer.WriteString("  ORGANISM  " + line + "\n")
		}
		buffer.WriteString("            " + line + "\n")
	}

	for _, ref := range record.REFERENCES {
		if ref.ORIGIN != "" {
			buffer.WriteString("REFERENCE   " + ref.ORIGIN + "\n")
			buffer.WriteString(formatStringWithNewlineChars("  AUTHORS   "+ref.AUTHORS, "            ", true))
			if ref.CONSRTM != "" {
				buffer.WriteString(formatStringWithNewlineChars("  CONSRTM   "+ref.CONSRTM, "            ", true))
			}
			buffer.WriteString(formatStringWithNewlineChars("  TITLE     "+ref.TITLE, "            ", true))
			buffer.WriteString(formatStringWithNewlineChars("  JOURNAL   "+ref.JOURNAL, "            ", true))
			if ref.PUBMED != "" {
				buffer.WriteString("   PUBMED   " + ref.PUBMED + "\n")
			}
		}
	}
	b64Decode, _ := base64.StdEncoding.DecodeString(record.COMMENT)
	buffer.WriteString("COMMENT     " + addSpacesSpecialHeader(string(b64Decode)) + "\n")

	return buffer.String()
}
func read_json() (result map[string][]string) {
	jsonFile, err := os.Open("generators/categorys_by_occurence.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func generateQualifierString(record *gbparse.Genbank, jsonmap map[string][]string) (QualString string) {
	var buffer bytes.Buffer
	spacestring := "                "
	for _, feature := range record.FEATURES {
		if feature.IsCompliment {
			buffer.WriteString("     " + feature.TYPE + spacestring[len(feature.TYPE):] + "complement(" + feature.START + ".." + feature.STOP + ")\n")
		} else {
			buffer.WriteString("     " + feature.TYPE + spacestring[len(feature.TYPE):] + feature.START + ".." + feature.STOP + "\n")
		}
		for _, occurence := range jsonmap[feature.TYPE] {
			if val, inMap := feature.QUALIFIERS[occurence]; inMap {
				if occurence == "/pseudo" {
					buffer.WriteString("                     /pseudo\n")
				} else {
					buffer.WriteString(formatStringWithNewlineChars(occurence+"="+val, "                     ", false))
				}
			}
		}

	}
	return buffer.String()
}

func generateSequenceString(record *gbparse.Genbank) {

}

func addSpacesSpecialHeader(inputString string) (Output string) {
	var returnbuffer bytes.Buffer
	for _, char := range inputString {
		returnbuffer.WriteRune(char)
		if char == '\n' {
			returnbuffer.WriteString("            ")
		}
	}
	return returnbuffer.String()
}

func formatStringWithNewlineChars(Splittedstring string, newlineinsertion string, hasKeyword bool) (result string) {
	var buffer bytes.Buffer
	// DEBUG: fmt.Println(Splittedstring)
	keyword := ""
	if hasKeyword {
		keyword = Splittedstring[:len(newlineinsertion)]
		Splittedstring = Splittedstring[len(newlineinsertion):]
	}
	lastsplitindex := 0
	currentlength := 0
	if strings.ContainsRune(Splittedstring, ' ') {
		currentlength := 0
		lastspaceindex := 0
		for i, char := range Splittedstring {
			if char == ' ' {
				lastspaceindex = i
			}
			if currentlength >= 79-len(newlineinsertion) {
				buffer.WriteString(newlineinsertion + Splittedstring[lastsplitindex:lastspaceindex] + "\n")
				lastsplitindex = lastspaceindex + 1
				currentlength = i - lastspaceindex - 1
			}
			if i == len(Splittedstring)-1 {
				buffer.WriteString(newlineinsertion + Splittedstring[lastsplitindex:] + "\n")
			}
			currentlength++
		}
	} else {
		for i := range Splittedstring {
			if currentlength >= 80-len(newlineinsertion)-1 {
				buffer.WriteString(newlineinsertion + Splittedstring[lastsplitindex:i] + "\n")
				lastsplitindex = i
				currentlength = 0
			}
			if i == len(Splittedstring)-1 {
				buffer.WriteString(newlineinsertion + Splittedstring[lastsplitindex:] + "\n")
			}
			currentlength++

		}
	}
	if len(buffer.String()) > 0 {
		return keyword + buffer.String()[len(keyword):]
	} else {
		return keyword + "\n"
	}
}
