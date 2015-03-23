package parser

import (
	"bytes"
	"fmt"
	"testing"
)

const (
	TEST_PARSER_TEXT        = "Now is the 'Winter' of our discontent. And then the other dude as well."
	TEST_PARSER_RESULT_JSON = `{"words":{"and":{"counter":1,"sentenceUse":[1]},"as":` +
		`{"counter":1,"sentenceUse":[1]},"discontent":{"counter":1,"sentenceUse":[0]},"dude":` +
		`{"counter":1,"sentenceUse":[1]},"is":{"counter":1,"sentenceUse":[0]},"now":` +
		`{"counter":1,"sentenceUse":[0]},"of":{"counter":1,"sentenceUse":[0]},"other":` +
		`{"counter":1,"sentenceUse":[1]},"our":{"counter":1,"sentenceUse":[0]},"the":` +
		`{"counter":2,"sentenceUse":[0,1]},"then":{"counter":1,"sentenceUse":[1]},"well":` +
		`{"counter":1,"sentenceUse":[1]},"winter":{"counter":1,"sentenceUse":[0]}}}`
)

// TestParserExecute tests the execution of the parser and validates the results.
func TestParserExecute(t *testing.T) {
	p := New()
	r := bytes.NewBufferString(TEST_PARSER_TEXT)
	p.Execute(r)
	result := fmt.Sprint(p)
	if result != TEST_PARSER_RESULT_JSON {
		t.Fatalf("Invalid parser results\nExpected:\n\n%s\n\nResult:\n\n%s\n\n", TEST_PARSER_RESULT_JSON, result)
	}

}
