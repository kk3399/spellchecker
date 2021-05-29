# spellchecker
basic spell checker

## how to run
- `go get .`
- `go build`
- `./spellchecker dictionary.txt file-to-check.txt`

## layout
```
├── sentencechecker                  - uses sentencereader and spellsuggestor to check sentences
├── sentencereader                   - reads the input file as sentences
├── spellsuggestor                   - suggests words for misspelt word
main.go                   
```


