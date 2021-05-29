package sentencereader

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type SentenceReader struct {
	scanner        *bufio.Scanner
	lineNumber     int
	cursorPosition int
	sentences      []*Sentence
	terminators    []byte
	done           bool
}

type Sentence struct {
	Words     []string
	Positions [][]int
	Text      strings.Builder
	done      bool
}

func New(fileReader io.Reader, terminators []byte) *SentenceReader {
	scanner := bufio.NewScanner(fileReader)
	sentenceReader := &SentenceReader{
		scanner:        scanner,
		lineNumber:     0,
		cursorPosition: 0,
		sentences:      nil,
		terminators:    terminators,
		done:           false,
	}
	return sentenceReader
}

func (sentenceReader *SentenceReader) Read() (bool, *Sentence) {
	if sentenceReader.done {
		return false, nil
	}
	for len(sentenceReader.sentences) > 0 && sentenceReader.sentences[0].done {
		sentence := sentenceReader.sentences[0]
		sentenceReader.sentences = sentenceReader.sentences[1:]
		return true, sentence
	}

	var lastSentence *Sentence
	if len(sentenceReader.sentences) > 0 {
		lastSentence = sentenceReader.sentences[len(sentenceReader.sentences)-1]
		sentenceReader.sentences = sentenceReader.sentences[:len(sentenceReader.sentences)-1]
	} else {
		lastSentence = &Sentence{
			Words:     nil,
			Positions: nil,
			Text:      strings.Builder{},
			done:      false,
		}
	}

	for sentenceReader.scanner.Scan() {
		sentenceReader.lineNumber++
		sentenceReader.cursorPosition = 1
		line := sentenceReader.scanner.Text()
		for _, word := range strings.Fields(line) {
			if len(word) == 0 {
				continue
			}
			lastSentence.Positions = append(lastSentence.Positions, []int{sentenceReader.lineNumber, sentenceReader.cursorPosition})
			sentenceReader.cursorPosition += len(word) + 1 // assuming single space
			if sentenceReader.isSentenceEnding(word[len(word)-1]) {
				lastSentence.Words = append(lastSentence.Words, word[:len(word)-1])
				lastSentence.Text.WriteString(word)
				lastSentence.done = true
				sentenceReader.sentences = append(sentenceReader.sentences, lastSentence)
				lastSentence = &Sentence{
					Words:     nil,
					Positions: nil,
					Text:      strings.Builder{},
					done:      false,
				}
			} else {
				lastSentence.Words = append(lastSentence.Words, word)
				lastSentence.Text.WriteString(fmt.Sprintf("%s ", word))
			}
		}
		if len(sentenceReader.sentences) > 0 && sentenceReader.sentences[0].done {
			firstSentence := sentenceReader.sentences[0]
			sentenceReader.sentences = sentenceReader.sentences[1:]
			return true, firstSentence
		}
	}
	sentenceReader.done = true
	// hoping the sentence & file would end in . or ? or !

	return false, nil
}

func (sentenceReader *SentenceReader) isSentenceEnding(c byte) bool {
	for _, b := range sentenceReader.terminators {
		if b == c {
			return true
		}
	}
	return false
}
