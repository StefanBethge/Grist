package recipe

import (
	"strings"
	"testing"
)

const minimalYAML = `
version: 1
source:
  type: csv
  path: in.csv
  sep: ";"
steps:
  - trim: { columns: [name, email] }
  - rename: { from: "E-Mail", to: email }
sink:
  type: csv
  path: out.csv
`

func TestParseMinimal(t *testing.T) {
	r, warnings, err := Parse([]byte(minimalYAML))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %v", warnings)
	}
	if r.Version != 1 {
		t.Errorf("version: want 1, got %d", r.Version)
	}
	if r.Source.Type != "csv" || r.Source.Path != "in.csv" || r.Source.Sep != ";" {
		t.Errorf("source mismatch: %+v", r.Source)
	}
	if len(r.Steps) != 2 {
		t.Fatalf("steps: want 2, got %d", len(r.Steps))
	}

	trim, ok := r.Steps[0].Params.(TrimParams)
	if !ok || r.Steps[0].Kind != "trim" {
		t.Fatalf("step 0: want trim TrimParams, got %s %T", r.Steps[0].Kind, r.Steps[0].Params)
	}
	if len(trim.Columns) != 2 || trim.Columns[0] != "name" || trim.Columns[1] != "email" {
		t.Errorf("trim columns: %v", trim.Columns)
	}

	rn, ok := r.Steps[1].Params.(RenameParams)
	if !ok || r.Steps[1].Kind != "rename" {
		t.Fatalf("step 1: want rename RenameParams, got %s %T", r.Steps[1].Kind, r.Steps[1].Params)
	}
	if rn.From != "E-Mail" || rn.To != "email" {
		t.Errorf("rename: %+v", rn)
	}
}

func TestParseUnknownTopLevelFieldWarns(t *testing.T) {
	y := minimalYAML + "\nfuture_feature: ok\n"
	_, warnings, err := Parse([]byte(y))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(warnings) != 1 || !strings.Contains(warnings[0], "future_feature") {
		t.Errorf("expected one warning about future_feature, got %v", warnings)
	}
}

func TestParseUnknownStepIsError(t *testing.T) {
	y := `
version: 1
source: { type: csv, path: in.csv }
steps:
  - explode_universe: {}
sink: { type: csv, path: out.csv }
`
	_, _, err := Parse([]byte(y))
	if err == nil {
		t.Fatal("expected error for unknown step")
	}
	if !strings.Contains(err.Error(), "explode_universe") {
		t.Errorf("error should mention the step: %v", err)
	}
}

func TestParseMissingVersion(t *testing.T) {
	y := `
source: { type: csv, path: in.csv }
steps: []
sink: { type: csv, path: out.csv }
`
	_, _, err := Parse([]byte(y))
	if err == nil || !strings.Contains(err.Error(), "version") {
		t.Fatalf("expected missing-version error, got %v", err)
	}
}

func TestParseUnsupportedVersion(t *testing.T) {
	y := `
version: 99
source: { type: csv, path: in.csv }
steps: []
sink: { type: csv, path: out.csv }
`
	_, _, err := Parse([]byte(y))
	if err == nil || !strings.Contains(err.Error(), "unsupported recipe version") {
		t.Fatalf("expected unsupported-version error, got %v", err)
	}
}

func TestValidateRejectsEmptyTrimColumns(t *testing.T) {
	y := `
version: 1
source: { type: csv, path: in.csv }
steps:
  - trim: { columns: [] }
sink: { type: csv, path: out.csv }
`
	_, _, err := Parse([]byte(y))
	if err == nil || !strings.Contains(err.Error(), "columns") {
		t.Fatalf("expected empty-columns error, got %v", err)
	}
}

func TestValidateRejectsIncompleteRename(t *testing.T) {
	y := `
version: 1
source: { type: csv, path: in.csv }
steps:
  - rename: { from: a }
sink: { type: csv, path: out.csv }
`
	_, _, err := Parse([]byte(y))
	if err == nil || !strings.Contains(err.Error(), "from and to") {
		t.Fatalf("expected from/to error, got %v", err)
	}
}
