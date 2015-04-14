package auth

import (
	"reflect"
	"testing"
)

const (
	tValidToken   = "X3A3E6C4C51F12DF2415682CCF9D18"
	tInvalidToken = "X8A95585DD5B64E33D5BF4C8F4E849"
)

func TestNew(t *testing.T) {
	a := New()
	tp := reflect.TypeOf(a)

	if tp.Kind() != reflect.Ptr {
		t.Fatalf("Auth not created as a pointer.")
	}

	tp = tp.Elem()
	if tp.Kind() != reflect.Struct {
		t.Fatalf("Auth not created as a struct.")
	}
	if tp.Name() != "Auth" {
		t.Fatalf("Auth struct is not named correctly.")
	}
	if !(tp.NumField() > 0) {
		t.Fatalf("Auth struct is empty.")
	}

	fld, ok := tp.FieldByName("Tokens")
	if !ok {
		t.Fatalf("Tokens not found.")
	}
	if fld.Type.String() != "map[string]bool" {
		t.Fatalf("Tokens is not a map[string]bool.")
	}
	if len(a.Tokens) == 0 {
		t.Fatalf("Tokens is empty.")
	}
}

func TestValid(t *testing.T) {
	a := New()
	a.Tokens[tValidToken] = true
	a.Tokens[tInvalidToken] = false

	if !a.Valid(tValidToken) {
		t.Errorf("Valid token validated as false.")
	}
	if a.Valid(tInvalidToken) {
		t.Errorf("Invalid token validated true.")
	}
	if a.Valid("NOT A STORED TOKEN") {
		t.Errorf("Missing token validated true.")
	}
}
