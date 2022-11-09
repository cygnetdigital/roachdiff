package model

import (
	"fmt"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
)

// Table represents the declarative state of a `CREATE TABLE` statement
type Table struct {
	Name        string
	Tree        *tree.CreateTable
	Columns     *ColumnList
	Indexes     *IndexList
	Constraints *ConstraintList
}

// NewTable parses a `CREATE TABLE` statement to produce a declarative
// representation of the table state
func NewTable(ct *tree.CreateTable) (*Table, error) {
	tbl := &Table{
		Name: ct.Table.Table(),
		Tree: ct,
	}

	tbl.Columns = newColumnList(tbl)
	tbl.Indexes = newIndexList(tbl)
	tbl.Constraints = newConstraintList(tbl)

	for _, def := range ct.Defs {

		switch t := def.(type) {
		case *tree.ColumnTableDef:
			tbl.Columns.add(t)

		case *tree.IndexTableDef:
			if err := tbl.Indexes.addTableDef(t); err != nil {
				return nil, fmt.Errorf("unable to add index: %w", err)
			}

		case *tree.UniqueConstraintTableDef:
			if err := tbl.Constraints.add(t, t.Name.String()); err != nil {
				return nil, fmt.Errorf("unable to add unique constraint: %w", err)
			}

		case *tree.CheckConstraintTableDef:
			if err := tbl.Constraints.add(t, t.Name.String()); err != nil {
				return nil, fmt.Errorf("unable to add check constraint: %w", err)
			}

		case *tree.ForeignKeyConstraintTableDef:
			if err := tbl.Constraints.add(t, t.Name.String()); err != nil {
				return nil, fmt.Errorf("unable to add foreign key constraint: %w", err)
			}

		default:
			return nil, fmt.Errorf("unable to parse tabledef: %T", t)
		}
	}

	return tbl, nil
}
