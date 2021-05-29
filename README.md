# spellchecker
basic spell checker

## how to run
- `go get .`
- `go build`
- `./spellchecker dictionary.txt file-to-check.txt`

## layout
.
|____sentencechecker                # sentencereader and spellsuggestor to check sentences
| |____sentence_checker.go
| |____sentence_checker_test.go
|____sentencereader                 # reads the input file as sentences
| |____sentence_reader.go
| |____sentence_reader_test.go
|____spellsuggestor                 # suggests words for misspelt word
| |____spell_suggestor.go
| |____spell_suggestor_test.go
|____main.go
|____news.txt
|____file-to-check.txt
|____dictionary.txt
|____README.md


