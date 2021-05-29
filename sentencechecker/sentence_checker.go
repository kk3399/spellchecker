package sentencechecker

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/jdkato/prose/v2"

	"spellchecker/sentencereader"
	"spellchecker/spellsuggestor"
)

type SentenceEnding byte

const (
	SentenceEndingDot         SentenceEnding = '.'
	SentenceEndingQuestion    SentenceEnding = '?'
	SentenceEndingExclamation SentenceEnding = '!'

	ProseNounByte        = 'N' // noun types
	ProseNounLabelPerson = "PERSON"
	ProseNounLabelGPE    = "GPE" // geographical/political Entities
)

type SentenceChecker struct {
	spellSuggestor *spellsuggestor.SpellSuggestor
	sentenceReader *sentencereader.SentenceReader
}

func New(dictionaryPath string, fileReader io.Reader) (*SentenceChecker, error) {
	spellSuggestor, err := spellsuggestor.New(dictionaryPath)
	if err != nil {
		return nil, err
	}

	sentenceReader := sentencereader.New(fileReader, []byte{byte(SentenceEndingDot), byte(SentenceEndingExclamation), byte(SentenceEndingQuestion)})
	sentenceChecker := &SentenceChecker{spellSuggestor, sentenceReader}
	return sentenceChecker, nil
}

func (sentenceChecker *SentenceChecker) Validate() {
	for {
		if success, sentence := sentenceChecker.sentenceReader.Read(); success {
			sentenceNounMap := map[string]bool{}

			// not sure how big of a file we are parsing, if the file is small we could parse all sentences and load into the document
			// we could also choose the size of data we want to process in one go
			sentenceText := sentence.Text.String()
			doc, err := prose.NewDocument(sentenceText)
			if err != nil {
				log.Fatal(err)
			}

			for _, tok := range doc.Tokens() {
				if tok.Tag[0] == ProseNounByte || tok.Label == ProseNounLabelPerson || tok.Label == ProseNounLabelGPE { // noun type
					sentenceNounMap[tok.Text] = true
				}
			}

			for i, word := range sentence.Words {
				if sentenceNounMap[word] {
					continue
				}

				word = trimSpecialCharsToRight(word)

				if !isWord(word) {
					continue
				}

				if sentenceChecker.spellSuggestor.IsCorrect(strings.ToLower(word)) {
					continue
				}

				printSuggestions(word, sentence.Positions[i], sentenceText, sentenceChecker.spellSuggestor.Suggest(word))
			}
		} else {
			return
		}
	}
}

func trimSpecialCharsToRight(word string) string {
	r := len(word)
	for i := range word {
		if isValidAlphabet(word[len(word)-1-i]) {
			break
		}
		r--
	}
	return word[:r]
}

// pretty naive at the moment
func isWord(word string) bool {
	if len(word) == 0 {
		return false
	}
	for i := range word {
		if isValidAlphabet(word[i]) {
			continue
		}
		return false
	}
	return true
}

func isValidAlphabet(b byte) bool {
	return ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z')
}

func printSuggestions(word string, position []int, sentenceText string, suggestions []string) {
	fmt.Printf("Line: %d, Position: %d, Misspelt: '%s'\n", position[0], position[1], word)
	fmt.Printf("Context '%s'\n", sentenceText)
	fmt.Printf("Suggestions %v\n", suggestions)
	fmt.Println()
}
