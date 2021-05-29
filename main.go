package main

import (
	"fmt"
	"os"

	"spellchecker/sentencechecker"
)

func main() {

	if len(os.Args) > 2 {
		file, err := os.Open(os.Args[2]) // "test.txt"
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		sentenceChecker, err := sentencechecker.New(os.Args[1], file) // "dictionary.txt"
		if err != nil {
			fmt.Println(err)
		}
		sentenceChecker.Validate()
	}

}
