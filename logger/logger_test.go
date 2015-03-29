package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestSetAndGetLogLevel(t *testing.T) {
	l := New(Debug, false)
	l.SetLogLevel(Debug)
	if l.GetLogLevel() != Debug {
		t.Fatalf("Invalid Log Level Set.\n")
	}
	err := l.SetLogLevel(Emergency - 1)
	if err == nil {
		t.Fatalf("Low param value was not tested properly.\n")
	}
	err = l.SetLogLevel(Debug + 1)
	if err == nil {
		t.Fatalf("High param value was not tested properly.\n")
	}
}

func TestDefaultSetLogLevel(t *testing.T) {
	l := New(-1, false)
	if l.GetLogLevel() != Info {
		t.Fatalf("Invalid default Log Level Set.\n")
	}
}

func TestSetColourLabels(t *testing.T) {
	l := New(-1, true)
	for i, actual := range l.labels {
		var colour int
		switch i {
		case Emergency, Alert, Critical, Error:
			colour = foregroundRed
		case Warning:
			colour = foregroundYellow
		case Notice:
			colour = foregroundGreen
		case Debug:
			colour = foregroundBlue
		default:
			colour = foregroundDefault
		}
		expected := fmt.Sprintf(colourFormat, colour, labels[i])
		if expected != actual {
			t.Fatalf("Invalid colour label\nExpected:%s\nActual:%s\n", expected, actual)
		}
	}
}

func TestEmergencyf(t *testing.T) {
	test_message := "Emergencyf"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.SetLogLevel(Debug)
		l.Emergencyf(false, test_message)
	}, fmt.Sprintf("%s%s\n", labels[Emergency], test_message))
}

func TestAlertf(t *testing.T) {
	test_message := "Alertf"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.SetLogLevel(Debug)
		l.Alertf(test_message)
	}, fmt.Sprintf("%s%s\n", labels[Alert], test_message))
}

func TestCriticalf(t *testing.T) {
	test_message := "Criticalf"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.SetLogLevel(Debug)
		l.Criticalf(test_message)
	}, fmt.Sprintf("%s%s\n", labels[Critical], test_message))
}

func TestErrorf(t *testing.T) {
	test_message := "Errorf"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.SetLogLevel(Debug)
		l.Errorf(test_message)
	}, fmt.Sprintf("%s%s\n", labels[Error], test_message))
}

func TestWarningf(t *testing.T) {
	test_message := "Warningf"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.SetLogLevel(Debug)
		l.Warningf(test_message)
	}, fmt.Sprintf("%s%s\n", labels[Warning], test_message))
}

func TestNoticef(t *testing.T) {
	test_message := "Noticef"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.SetLogLevel(Debug)
		l.Noticef(test_message)
	}, fmt.Sprintf("%s%s\n", labels[Notice], test_message))
}

func TestInfof(t *testing.T) {
	test_message := "Infof"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.SetLogLevel(Debug)
		l.Infof(test_message)
	}, fmt.Sprintf("%s%s\n", labels[Info], test_message))
}

func TestDebugf(t *testing.T) {
	test_message := "Debugf"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.SetLogLevel(Debug)
		l.Debugf(test_message)
	}, fmt.Sprintf("%s%s\n", labels[Debug], test_message))
}

func TestOutputf(t *testing.T) {
	test_label := "[OUTPUT] "
	test_message := "Output"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.SetLogLevel(Debug)
		l.Output(-1, test_label, test_message)
	}, fmt.Sprintf("%s%s\n", "[OUTPUT] ", test_message))
}

// expectOutput is a helper function that repipes or mocks out stdout and allows error messages to be tested
// against the pipe.
func expectOutput(t *testing.T, f func(), expected string) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	os.Stdout.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC
	if !strings.Contains(out, expected) {
		t.Fatalf("Expected '%s', received '%s'\n", expected, out)
	}
}
