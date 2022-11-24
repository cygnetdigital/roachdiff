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

func (c *Constraint) Equal(n *Constraint) bool {
	fmt1 := tree.NewFmtCtx(tree.FmtSimple)
	c.Tree.Format(fmt1)
	fmt2 := tree.NewFmtCtx(tree.FmtSimple)
	n.Tree.Format(fmt2)
	return fmt1.String() == fmt2.String()
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

func (il *ConstraintList) addByCol(c *tree.ColumnTableDef) error {
	if c == nil {
		return nil
	}

	if c.References.Table == nil {
		return nil
	}

	name := fmt.Sprintf("%s_%s_fkey", il.table.Name, c.Name)

	cons := &tree.ForeignKeyConstraintTableDef{
		Name:     tree.Name(name),
		Table:    *c.References.Table,
		FromCols: tree.NameList{tree.Name(c.Name)},
		ToCols:   tree.NameList{c.References.Col},
		Actions:  c.References.Actions,
	}

	return il.add(cons, name)
}

func (il *ConstraintList) All() []*Constraint {
	return il.constraints
}

func (il *ConstraintList) Find(name string) (*Constraint, bool) {
	col, ok := il.constraintMap[name]
	return col, ok
}
