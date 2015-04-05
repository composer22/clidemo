package server

import (
	"regexp"
	"testing"
)

const (
	v4UUIDRegExpFmt = "^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$"
)

func TestUtilsCreateV4UUID(t *testing.T) {
	r, _ := regexp.Compile(v4UUIDRegExpFmt)

	for i := 0; i < 10; i++ {
		uuid := createV4UUID()
		if !r.MatchString(uuid) {
			t.Errorf("UUID not V4 standard.\n")
			break
		}
	}
	uuid1 := createV4UUID()
	uuid2 := createV4UUID()
	if uuid1 == uuid2 {
		t.Errorf("UUID not being created uniquely.\n")
	}

}
