package cmd

import (
	"testing"
	"bytes"
	"io"
)

func Test_RootCmd(t *testing.T) {
	b := bytes.NewBufferString("")
	cmd := NewRootCmd()
	cmd.SetOut(b)
	cmd.Execute()
	out, err := io.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "Welcome to PFT!" {
		t.Fatalf("expected \"%s\" got \"%s\"", "Welcome to PFT!", string(out))
	}
}