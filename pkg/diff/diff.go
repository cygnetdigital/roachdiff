package diff

import (
	"fmt"
	"strings"

	"github.com/cygnetdigital/roachdiff/pkg/diff/gen"
	"github.com/cygnetdigital/roachdiff/pkg/diff/model"
)

// Differ between two schemas
type Differ struct {
	Generator *gen.Generator
	left      string
	right     string
}

// NewDiffer with strings
func NewDiffer(left, right string) *Differ {
	return &Differ{
		Generator: gen.NewGenerator(),
		left:      left,
		right:     right,
	}
}

type diffFunc func(d *Diff, a, b *model.Model) error

// Run the differ to produce a diff
func (d *Differ) Run() (*Diff, error) {

	left, err := model.NewModel(d.left)
	if err != nil {
		return nil, fmt.Errorf("failed to parse left schema: %w", err)
	}

	right, err := model.NewModel(d.right)
	if err != nil {
		return nil, fmt.Errorf("failed to parse right schema: %w", err)
	}

	funcs := []diffFunc{
		CreateEnum,
		DropEnum,
		CreateEnumValues,
		DropEnumValues,
		CreateTable,
		DropTable,
		DropColumns,
		CreateColumns,
		DropIndexes,
		CreateIndexes,
		DropConstraints,
		CreateConstraints,
		AlterConstraints,
		ColumnNullability,
	}

	diff := &Diff{
		gen: d.Generator,
	}

	for _, df := range funcs {
		if err := df(diff, left, right); err != nil {
			return nil, fmt.Errorf("failed to diff with diffFunc %T: %w", df, err)
		}
	}

	return diff, nil
}

type Diff struct {
	sb        strings.Builder
	gen       *gen.Generator
	dangerous bool
}

func (d *Diff) append(s gen.Statement) {
	d.dangerous = d.dangerous || s.HasWarning
	d.sb.WriteString(s.SQL)
	d.sb.WriteString("\n")
}

func (d *Diff) String() string {
	return d.sb.String()
}

// Dangerous returns true if the diff contains dangerous statements
func (d *Diff) Dangerous() bool {
	return d.dangerous
}
