package generators

import (
	"bytes"

	bioproto "github.com/ag-computational-bio/BioProtobufSchemas/go"
)

//GenerateFastafromproto Fasta protobuf to fasta

func GenerateFastafromproto(record *bioproto.Fasta) string {
	return record.HEADER + "\n" + insertNth(record.SEQUENCE, 80) + "\n"
}

func insertNth(s string, n int) string {
	var buffer bytes.Buffer
	var n1 = n - 1
	var l1 = len(s) - 1
	for i, rune := range s {
		buffer.WriteRune(rune)
		if i%n == n1 && i != l1 {
			buffer.WriteRune('\n')
		}
	}
	return buffer.String()
}
