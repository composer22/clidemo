// Package parser implements a mechanism for counting words and sentence locations for a given text source.
package parser

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

const (
	punctuationMarks = ";:,!?.\\/[](){}-\"'`"
)

// Parser represents text source plus a mapping of unique words found in the text with an arrray of sentence ids where the
// words were located.
type Parser struct {
	Words map[string]*wordRef `json:"words"` // Words as key with struct of counts, location.
}

// wordRef represents a word found in the source text, a count on it's use, and which sentences it was found.
type wordRef struct {
	Counter     int   `json:"counter"`     // The number of times the word was found in the text.
	SentenceUse []int `json:"sentenceUse"` // The sentence id where the word was found.
}

// New is a factory function that returns a new parser instance.
func New() *Parser {
	return &Parser{
		Words: make(map[string]*wordRef),
	}
}

// Execute begins the parsing process. The source text is read, words are counted, and unique sentence ids are
// recorded.
func (p *Parser) Execute(source io.Reader) {
	scanner := bufio.NewScanner(source)
	scanner.Split(bufio.ScanWords)

	eos := false
	sentPointer := 0

	// Loop on the text and analyze word usage.
	for scanner.Scan() {
		word := scanner.Text()

		// Check for period in word and mark EOS was found.
		if strings.HasSuffix(word, ".") {
			eos = true
		}

		// Remove beginning and trailing punctuation.
		word = strings.Trim(word, punctuationMarks)

		// Store it as a result.
		if len(word) > 0 {
			key := strings.ToLower(word)
			w, ok := p.Words[key]
			if !ok {
				p.Words[key] = &wordRef{
					Counter:     0,
					SentenceUse: make([]int, 0),
				}
				w = p.Words[key]
			}
			w.Counter++
			w.SentenceUse = append(w.SentenceUse, sentPointer)
		}

		// If a period was found in the word advance the pointer.
		if eos {
			sentPointer++
			eos = false
		}
	}
}

// Reset cleans out the parser and makes it available for another parse job.
func (p *Parser) Reset() {
	p.Words = make(map[string]*wordRef)
}

// String is an implentation of the Stringer interface so the structure is returned as a string to fmt.Print() etc.
func (p *Parser) String() string {
	result, _ := json.Marshal(p)
	return string(result)
}
