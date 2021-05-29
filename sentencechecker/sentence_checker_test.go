package sentencechecker

import (
	"os"
	"testing"
)

func TestSentenceChecker_Validate(t *testing.T) {
	file, err := os.Open("../file-to-check.txt")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	sentenceChecker, err := New("../dictionary.txt", file)
	if err != nil {
		t.Error(err)
	}
	sentenceChecker.Validate()
}
