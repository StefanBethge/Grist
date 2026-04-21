package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// TestGristRunExampleEndToEnd copies the examples/minimal fixture into a
// temp directory, invokes the run subcommand, and asserts the produced
// output matches what the recipe should have emitted.
func TestGristRunExampleEndToEnd(t *testing.T) {
	dir := t.TempDir()

	recipePath := filepath.Join(dir, "minimal.yaml")
	copyFile(t, "../../examples/minimal.yaml", recipePath)
	copyFile(t, "../../examples/minimal.in.csv", filepath.Join(dir, "minimal.in.csv"))

	var stdout, stderr bytes.Buffer
	if err := run([]string{"run", recipePath}, &stdout, &stderr); err != nil {
		t.Fatalf("run: %v (stderr=%q)", err, stderr.String())
	}

	got, err := os.ReadFile(filepath.Join(dir, "minimal.out.csv"))
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	want := "name;email\nAda;ada@example.com\nGrace;grace@example.com\n"
	if string(got) != want {
		t.Errorf("output mismatch\n got: %q\nwant: %q", string(got), want)
	}
}

func copyFile(t *testing.T, src, dst string) {
	t.Helper()
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read %s: %v", src, err)
	}
	if err := os.WriteFile(dst, data, 0o600); err != nil {
		t.Fatalf("write %s: %v", dst, err)
	}
}
