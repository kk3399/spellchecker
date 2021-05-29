package spellsuggestor

import (
	"bufio"
	"log"
	"os"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
)

type SpellSuggestor struct {
	dictionary map[string]bool
}

func New(dictionaryFile string) (*SpellSuggestor, error) {
	SpellSuggestor := &SpellSuggestor{make(map[string]bool)}
	if err := SpellSuggestor.loadDictionary(dictionaryFile); err != nil {
		return nil, err
	}
	return SpellSuggestor, nil
}

func (SpellSuggestor *SpellSuggestor) loadDictionary(dictionaryFile string) error {
	file, err := os.Open(dictionaryFile)
	if err != nil {
		return err
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		SpellSuggestor.dictionary[scanner.Text()] = true
	}
	return nil
}

func (SpellSuggestor *SpellSuggestor) IsCorrect(word string) bool {
	return SpellSuggestor.dictionary[word]
}

func (SpellSuggestor *SpellSuggestor) Suggest(word string) []string {
	suggestionsMap := make(map[string]bool)
	if SpellSuggestor.dictionary[word] {
		return nil
	}

	editDistance1 := SpellSuggestor.getWordsOneEditAway(word)
	unMatchedAtEdit1 := make([]string, 0)
	for _, editDistance1Word := range editDistance1 {
		if SpellSuggestor.dictionary[editDistance1Word] {
			suggestionsMap[editDistance1Word] = true
		} else {
			unMatchedAtEdit1 = append(unMatchedAtEdit1, editDistance1Word)
		}
	}
	if len(suggestionsMap) == 0 { // assuming suggestions with shorter edit distance are preferred
		for _, editDistance1Word := range unMatchedAtEdit1 {
			for _, editDistance2Word := range SpellSuggestor.getWordsOneEditAway(editDistance1Word) {
				if SpellSuggestor.dictionary[editDistance2Word] {
					suggestionsMap[editDistance2Word] = true
				}
			} // O(n^k) if we go further
		}
	}

	suggestions := make([]string, 0)
	for suggestion := range suggestionsMap {
		suggestions = append(suggestions, suggestion)
	}
	return suggestions
}

func (SpellSuggestor *SpellSuggestor) getWordsOneEditAway(word string) []string {
	if len(word) == 1 {
		return nil
	}
	suggestions := make([]string, 0)
	wordBytes := []byte(word) // assuming valid alphabet
	for i, chari := range wordBytes {
		// add new char
		tempBytes := make([]byte, len(wordBytes)+1)
		copy(tempBytes, wordBytes[:i])
		copy(tempBytes[i+1:], wordBytes[i:])
		for _, c := range []byte(alphabet) {
			tempBytes[i] = c
			suggestions = append(suggestions, string(tempBytes))
		}

		copy(tempBytes[i:], tempBytes[i+1:])
		tempBytes = tempBytes[:len(tempBytes)-1]

		// swap with adjascent char
		if i > 0 {
			tempBytes[i-1], tempBytes[i] = tempBytes[i], tempBytes[i-1]
			suggestions = append(suggestions, string(tempBytes))
			tempBytes[i-1], tempBytes[i] = tempBytes[i], tempBytes[i-1]
		}

		// replace with other char from alphabet
		for _, c := range []byte(alphabet) {
			if chari == c {
				continue
			}
			tempBytes[i] = c
			suggestions = append(suggestions, string(tempBytes))
		}
		tempBytes[i] = chari

		// delete the char
		copy(tempBytes[i:], tempBytes[i+1:])
		tempBytes = tempBytes[:len(tempBytes)-1]
		suggestions = append(suggestions, string(tempBytes))
	}
	// add new char at the end
	tempBytes := make([]byte, len(wordBytes)+1)
	copy(tempBytes, wordBytes)
	for _, c := range []byte(alphabet) {
		tempBytes[len(wordBytes)] = c
		suggestions = append(suggestions, string(tempBytes))
	}

	return suggestions
}
