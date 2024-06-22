package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func Test_main(t *testing.T) {
	stdOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	_ = w.Close()

	res, _ := io.ReadAll(r)
	output := string(res)

	os.Stdout = stdOut

	if !strings.Contains(output, "33660") {
		t.Errorf("Expected 34320, got %s", output)
	}
}
