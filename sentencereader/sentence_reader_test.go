package sentencereader

import (
	"os"
	"testing"
)

func TestSentenceReader_Read(t *testing.T) {
	file, err := os.Open("../file-to-check.txt")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	sentenceReader := New(file, []byte{'.', '?', '!'})
	for {
		if success, sentence := sentenceReader.Read(); success {
			t.Log(sentence)
		} else {
			break
		}
	}
}
