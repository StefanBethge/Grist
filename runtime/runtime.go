package runtime

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/stefanbethge/grist/recipe"
	"github.com/stefanbethge/gseq-table/csv"
	"github.com/stefanbethge/gseq-table/table"
)

// Run loads the source declared by r, applies each step in order, and writes
// the result to the declared sink. The recipe is assumed to have been
// validated (see recipe.Load).
func Run(r recipe.Recipe) error {
	t, err := loadSource(r.Source, r.BaseDir)
	if err != nil {
		return fmt.Errorf("load source: %w", err)
	}
	for i, step := range r.Steps {
		t, err = applyStep(t, step)
		if err != nil {
			return fmt.Errorf("step %d (%s): %w", i+1, step.Kind, err)
		}
	}
	if err := writeSink(r.Sink, t, r.BaseDir); err != nil {
		return fmt.Errorf("write sink: %w", err)
	}
	return nil
}

func loadSource(s recipe.Source, baseDir string) (table.Table, error) {
	opts, err := csvReaderOptions(s)
	if err != nil {
		return table.Table{}, err
	}
	res := csv.New(opts...).ReadFile(resolvePath(baseDir, s.Path))
	if res.IsErr() {
		return table.Table{}, res.UnwrapErr()
	}
	return res.Unwrap(), nil
}

func writeSink(s recipe.Sink, t table.Table, baseDir string) error {
	opts, err := csvWriterOptions(s)
	if err != nil {
		return err
	}
	return csv.NewWriter(opts...).WriteFile(resolvePath(baseDir, s.Path), t)
}

// resolvePath makes p absolute relative to baseDir when it is not already
// absolute. An empty baseDir leaves p unchanged (callers like tests may pass
// absolute paths directly).
func resolvePath(baseDir, p string) string {
	if baseDir == "" || filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(baseDir, p)
}

func csvReaderOptions(s recipe.Source) ([]csv.Option, error) {
	var opts []csv.Option
	if s.Sep != "" {
		sep, err := singleRune(s.Sep, "source.sep")
		if err != nil {
			return nil, err
		}
		opts = append(opts, csv.WithSeparator(sep))
	}
	if s.Header != nil && !*s.Header {
		opts = append(opts, csv.WithNoHeader())
	}
	return opts, nil
}

func csvWriterOptions(s recipe.Sink) ([]csv.WriterOption, error) {
	var opts []csv.WriterOption
	if s.Sep != "" {
		sep, err := singleRune(s.Sep, "sink.sep")
		if err != nil {
			return nil, err
		}
		opts = append(opts, csv.WithWriteSeparator(sep))
	}
	if s.Header != nil && !*s.Header {
		opts = append(opts, csv.WithoutHeader())
	}
	return opts, nil
}

func singleRune(s, what string) (rune, error) {
	runes := []rune(s)
	if len(runes) != 1 {
		return 0, fmt.Errorf("%s must be a single character, got %q", what, s)
	}
	return runes[0], nil
}

func applyStep(t table.Table, s recipe.Step) (table.Table, error) {
	switch p := s.Params.(type) {
	case recipe.TrimParams:
		for _, col := range p.Columns {
			t = t.Map(col, strings.TrimSpace)
		}
		return t, nil
	case recipe.RenameParams:
		return t.Rename(p.From, p.To), nil
	default:
		return t, fmt.Errorf("unhandled params type %T", s.Params)
	}
}
