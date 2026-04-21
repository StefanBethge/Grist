package recipe

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// knownTopLevel lists the recognised top-level keys. Anything else parses
// into a warning so newer recipes stay forward-compatible.
var knownTopLevel = map[string]struct{}{
	"version": {},
	"source":  {},
	"steps":   {},
	"sink":    {},
}

// Load parses the recipe file at path and validates it. Warnings are
// non-fatal messages (e.g. unknown top-level keys) that the caller should
// surface to the user. A non-nil error means the recipe could not be loaded.
func Load(path string) (Recipe, []string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Recipe{}, nil, err
	}
	r, warnings, err := Parse(data)
	if err != nil {
		return r, warnings, err
	}
	r.BaseDir = filepath.Dir(path)
	return r, warnings, nil
}

// Parse is Load for an already-read byte slice, useful for tests.
func Parse(data []byte) (Recipe, []string, error) {
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return Recipe{}, nil, fmt.Errorf("parse yaml: %w", err)
	}

	warnings := collectTopLevelWarnings(&root)

	var r Recipe
	if err := root.Decode(&r); err != nil {
		return Recipe{}, warnings, fmt.Errorf("decode recipe: %w", err)
	}
	if err := r.Validate(); err != nil {
		return Recipe{}, warnings, err
	}
	return r, warnings, nil
}

func collectTopLevelWarnings(root *yaml.Node) []string {
	mapping := mappingContent(root)
	if mapping == nil {
		return nil
	}
	var warnings []string
	for i := 0; i < len(mapping); i += 2 {
		key := mapping[i]
		if _, ok := knownTopLevel[key.Value]; !ok {
			warnings = append(warnings, fmt.Sprintf("line %d: ignoring unknown top-level field %q", key.Line, key.Value))
		}
	}
	return warnings
}

// mappingContent returns the key/value pairs of the first mapping node under
// root, or nil if root is not a mapping-bearing document.
func mappingContent(root *yaml.Node) []*yaml.Node {
	n := root
	if n.Kind == yaml.DocumentNode && len(n.Content) > 0 {
		n = n.Content[0]
	}
	if n.Kind != yaml.MappingNode {
		return nil
	}
	return n.Content
}

// Validate checks required fields and the schema version. It runs after
// Decode so type errors are reported first.
func (r Recipe) Validate() error {
	if r.Version == 0 {
		return fmt.Errorf("missing required field %q", "version")
	}
	if r.Version != SchemaVersion {
		return fmt.Errorf("unsupported recipe version %d (this build understands %d)", r.Version, SchemaVersion)
	}
	if err := r.Source.validate("source"); err != nil {
		return err
	}
	if err := r.Sink.validate("sink"); err != nil {
		return err
	}
	for i, step := range r.Steps {
		if err := validateStep(step); err != nil {
			return fmt.Errorf("step %d (%s): %w", i+1, step.Kind, err)
		}
	}
	return nil
}

func (s Source) validate(what string) error {
	if s.Type == "" {
		return fmt.Errorf("%s.type is required", what)
	}
	if s.Type != "csv" {
		return fmt.Errorf("%s.type %q is not supported yet", what, s.Type)
	}
	if s.Path == "" {
		return fmt.Errorf("%s.path is required", what)
	}
	return nil
}

func (s Sink) validate(what string) error {
	return Source(s).validate(what)
}

func validateStep(s Step) error {
	switch p := s.Params.(type) {
	case TrimParams:
		if len(p.Columns) == 0 {
			return fmt.Errorf("columns must not be empty")
		}
	case RenameParams:
		if p.From == "" || p.To == "" {
			return fmt.Errorf("from and to are required")
		}
	default:
		return fmt.Errorf("unhandled params type %T", s.Params)
	}
	return nil
}
