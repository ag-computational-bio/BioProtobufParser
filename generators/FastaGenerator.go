package generators

import (
	"bytes"
	"gbparsertest2/gbparse"
)

func generateFastafromproto(record *gbparse.Fasta) (fastarecord string) {
	return record.HEADER + "\n" + insertNth(record.SEQUENCE, 81)
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
