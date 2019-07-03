package generators

import (
	"bytes"
	"gbparsertest2/gbparse"
	"strings"
)

func GenerateGBfromproto(record *gbparse.Genbank) {

}

func generateHeaderString(record *gbparse.Genbank) {
	var buffer bytes.Buffer
	buffer.WriteString("LOCUS       " + record.LOCUS + "\n")
	buffer.WriteString(formatStringWithNewlineChars("DEFINITION  "+record.DEFINITION, "            ", true) + "\n")
	buffer.WriteString("ACCESSION   " + record.ACCESSION + "\n")
	buffer.WriteString("VERSION     " + record.ACCESSION + "\n")
	buffer.WriteString("DBLINK      " + addSpacesSpecialHeader(record.DBLINK) + "\n")
	buffer.WriteString("KEYWORDS    " + record.KEYWORDS + "\n")
	buffer.WriteString("SOURCE      " + record.SOURCE + "\n")
	buffer.WriteString("  ORGANISM  " + addSpacesSpecialHeader(record.ORGANISM) + "\n")

}

func generateQualifierString(record *gbparse.Genbank) {

}

func generateSequenceString(record *gbparse.Genbank) {

}

func addSpacesSpecialHeader(inputString string) (Output string) {
	var returnbuffer bytes.Buffer
	for _, char := range inputString {
		if char == '\n' {
			returnbuffer.WriteString("            ")
		}
		returnbuffer.WriteRune(char)
	}
	return returnbuffer.String()
}

func formatStringWithNewlineChars(Splittedstring string, newlineinsertion string, hasKeyword bool) (result string) {
	var buffer bytes.Buffer
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
			if currentlength >= 80-len(newlineinsertion)-1 {
				buffer.WriteString(newlineinsertion + Splittedstring[lastsplitindex:lastspaceindex] + "\n")
				lastsplitindex = lastspaceindex + 1
				currentlength = 0
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
	return keyword + buffer.String()[len(keyword):]
}
