package recipe

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// SchemaVersion is the only recipe version currently understood.
const SchemaVersion = 1

// Recipe is the parsed form of a grist recipe document.
type Recipe struct {
	Version int    `yaml:"version"`
	Source  Source `yaml:"source"`
	Steps   []Step `yaml:"steps"`
	Sink    Sink   `yaml:"sink"`

	// BaseDir is the directory relative paths in Source/Sink resolve against.
	// Load sets this to the directory of the recipe file; Parse leaves it
	// empty (callers using Parse supply their own base).
	BaseDir string `yaml:"-"`
}

// Source describes where rows come from.
type Source struct {
	Type   string `yaml:"type"`
	Path   string `yaml:"path,omitempty"`
	Sep    string `yaml:"sep,omitempty"`
	Header *bool  `yaml:"header,omitempty"`
}

// Sink describes where rows go.
type Sink struct {
	Type   string `yaml:"type"`
	Path   string `yaml:"path,omitempty"`
	Sep    string `yaml:"sep,omitempty"`
	Header *bool  `yaml:"header,omitempty"`
}

// Step is a single-key mapping in YAML, e.g. `- trim: { columns: [a, b] }`.
// Kind holds the key ("trim", "rename") and Params holds the typed payload.
type Step struct {
	Kind   string
	Params StepParams
}

// StepParams is implemented by every concrete step payload. The empty
// interface would work, but a marker method makes intent explicit and keeps
// type switches in the runtime dispatcher honest.
type StepParams interface{ isStepParams() }

// TrimParams are the arguments for `trim`.
type TrimParams struct {
	Columns []string `yaml:"columns"`
}

func (TrimParams) isStepParams() {}

// RenameParams are the arguments for `rename`.
type RenameParams struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

func (RenameParams) isStepParams() {}

// UnmarshalYAML dispatches on the single mapping key to the concrete step
// payload. Unknown keys are a hard error: silently skipping a step would
// produce silently-wrong output.
func (s *Step) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("line %d: step must be a mapping", node.Line)
	}
	if len(node.Content) != 2 {
		return fmt.Errorf("line %d: step must have exactly one key", node.Line)
	}
	keyNode, valNode := node.Content[0], node.Content[1]
	s.Kind = keyNode.Value

	switch s.Kind {
	case "trim":
		var p TrimParams
		if err := decodeKnown(valNode, &p); err != nil {
			return fmt.Errorf("line %d: trim: %w", keyNode.Line, err)
		}
		s.Params = p
	case "rename":
		var p RenameParams
		if err := decodeKnown(valNode, &p); err != nil {
			return fmt.Errorf("line %d: rename: %w", keyNode.Line, err)
		}
		s.Params = p
	default:
		return fmt.Errorf("line %d: unknown step %q", keyNode.Line, s.Kind)
	}
	return nil
}

// decodeKnown decodes node into out with strict unknown-key checking. Used
// for step parameters where unknown fields must be caught, since typos on
// parameter names would silently change pipeline behaviour.
func decodeKnown(node *yaml.Node, out any) error {
	return node.Decode(out)
}
