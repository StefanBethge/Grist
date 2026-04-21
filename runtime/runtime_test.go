package runtime

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stefanbethge/grist/recipe"
)

func TestRunAppliesTrimAndRename(t *testing.T) {
	dir := t.TempDir()
	inPath := filepath.Join(dir, "in.csv")
	outPath := filepath.Join(dir, "out.csv")

	in := "name;E-Mail\n  Ada  ;  ada@example.com\nGrace; grace@example.com \n"
	if err := os.WriteFile(inPath, []byte(in), 0o600); err != nil {
		t.Fatalf("write input: %v", err)
	}

	r := recipe.Recipe{
		Version: recipe.SchemaVersion,
		Source:  recipe.Source{Type: "csv", Path: inPath, Sep: ";"},
		Steps: []recipe.Step{
			{Kind: "trim", Params: recipe.TrimParams{Columns: []string{"name", "E-Mail"}}},
			{Kind: "rename", Params: recipe.RenameParams{From: "E-Mail", To: "email"}},
		},
		Sink: recipe.Sink{Type: "csv", Path: outPath, Sep: ";"},
	}

	if err := Run(r); err != nil {
		t.Fatalf("Run: %v", err)
	}

	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	want := "name;email\nAda;ada@example.com\nGrace;grace@example.com\n"
	if string(got) != want {
		t.Errorf("output mismatch\n got: %q\nwant: %q", string(got), want)
	}
}

func TestRunRejectsMultiCharSeparator(t *testing.T) {
	r := recipe.Recipe{
		Version: recipe.SchemaVersion,
		Source:  recipe.Source{Type: "csv", Path: "in.csv", Sep: "||"},
		Sink:    recipe.Sink{Type: "csv", Path: "out.csv"},
	}
	if err := Run(r); err == nil {
		t.Fatal("expected error for multi-character separator")
	}
}
