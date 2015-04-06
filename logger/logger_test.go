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
	if l.GetLogLevel() != Debug {
		t.Errorf("Invalid Log Level Set.")
	}
	err := l.SetLogLevel(Info)
	if err != nil {
		t.Errorf("Set log level func should have been called correctly for value.")
	}
	if l.level != Info {
		t.Errorf("Set log level func should have set new value correctly.")
	}

	err = l.SetLogLevel(UseDefault)
	if err != nil {
		t.Errorf("Set log level func should have been called correctly for value.")
	}

	if l.level != Info {
		t.Errorf("Set default log level should have set new value correctly.")
	}

	err = l.SetLogLevel(UseDefault - 1)
	if err == nil {
		t.Errorf("Low param value was not tested properly.")
	}
	err = l.SetLogLevel(Debug + 1)
	if err == nil {
		t.Errorf("High param value was not tested properly.")
	}
}

func TestDefaultSetLogLevel(t *testing.T) {
	l := New(UseDefault, false)
	if l.GetLogLevel() != Info {
		t.Errorf("Invalid default Log Level Set.")
	}
}

func TestSetErrorFunc(t *testing.T) {
	l := New(Debug, false)
	if err := l.SetExitFunc(nil); err == nil {
		t.Errorf("Invalid set exit function with nil.")
	}

	if err := l.SetExitFunc(func(code int) {}); err != nil {
		t.Errorf("Invalid set exit function with vald value.")
	}
}

func TestSetColourLabels(t *testing.T) {
	l := New(UseDefault, true)
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
			t.Errorf("Invalid colour label\nExpected:%s\nActual:%s", expected, actual)
		}
	}
}

func TestEmergencyf(t *testing.T) {
	t.Parallel()
	test_message := "Emergencyf"
	expectOutput(t, func() {
		l := New(Debug, false) // Mock the exit so coverage can complete.
		l.exit = func(code int) {}
		l.Emergencyf(test_message)
	}, fmt.Sprintf("%s%s\n", labels[Emergency], test_message))
}

func TestAlertf(t *testing.T) {
	t.Parallel()
	test_message := "Alertf"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.Alertf(test_message)
	}, fmt.Sprintf("%s%s\n", labels[Alert], test_message))
}

func TestCriticalf(t *testing.T) {
	t.Parallel()
	test_message := "Criticalf"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.Criticalf(test_message)
	}, fmt.Sprintf("%s%s\n", labels[Critical], test_message))
}

func TestErrorf(t *testing.T) {
	t.Parallel()
	test_message := "Errorf"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.Errorf(test_message)
	}, fmt.Sprintf("%s%s\n", labels[Error], test_message))
}

func TestWarningf(t *testing.T) {
	t.Parallel()
	test_message := "Warningf"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.Warningf(test_message)
	}, fmt.Sprintf("%s%s\n", labels[Warning], test_message))
}

func TestNoticef(t *testing.T) {
	t.Parallel()
	test_message := "Noticef"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.Noticef(test_message)
	}, fmt.Sprintf("%s%s\n", labels[Notice], test_message))
}

func TestInfof(t *testing.T) {
	t.Parallel()
	test_message := "Infof"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.Infof(test_message)
	}, fmt.Sprintf("%s%s\n", labels[Info], test_message))
}

func TestDebugf(t *testing.T) {
	t.Parallel()
	test_message := "Debugf"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.Debugf(test_message)
	}, fmt.Sprintf("%s%s\n", labels[Debug], test_message))
}

func TestOutputf(t *testing.T) {
	t.Parallel()
	test_label := "[OUTPUT] "
	test_message := "Output"
	expectOutput(t, func() {
		l := New(Debug, false)
		l.Output(-1, test_label, test_message)
	}, fmt.Sprintf("%s%s\n", test_label, test_message))
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
		t.Errorf("Expected '%s', received '%s'.", expected, out)
	}
}
