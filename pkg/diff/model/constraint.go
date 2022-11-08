package model

import (
	"fmt"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
)

type ConstraintList struct {
	table         *Table
	constraints   []*Constraint
	constraintMap map[string]*Constraint
}

type Constraint struct {
	Table *Table
	Tree  tree.ConstraintTableDef
	Name  string
}

func newConstraintList(table *Table) *ConstraintList {
	return &ConstraintList{
		table:         table,
		constraints:   []*Constraint{},
		constraintMap: map[string]*Constraint{},
	}
}

func (il *ConstraintList) add(cons tree.ConstraintTableDef, name string) error {
	if name == "\"\"" {
		return fmt.Errorf("all constraints must be named")
	}

	wrapped := &Constraint{
		Table: il.table,
		Tree:  cons,
		Name:  name,
	}

	il.constraintMap[name] = wrapped
	il.constraints = append(il.constraints, wrapped)
	return nil
}

func (il *ConstraintList) All() []*Constraint {
	return il.constraints
}

func (il *ConstraintList) Find(name string) (*Constraint, bool) {
	col, ok := il.constraintMap[name]
	return col, ok
}
