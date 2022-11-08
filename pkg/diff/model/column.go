package model

import "github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"

type ColumnList struct {
	table  *Table
	cols   []*Column
	colMap map[string]*Column
}

type Column struct {
	Table *Table
	Name  string
	Tree  *tree.ColumnTableDef
}

func newColumn(table *Table, col *tree.ColumnTableDef) *Column {
	return &Column{
		Table: table,
		Name:  col.Name.String(),
		Tree:  col,
	}
}

func newColumnList(table *Table) *ColumnList {
	return &ColumnList{
		table:  table,
		cols:   []*Column{},
		colMap: map[string]*Column{},
	}
}

func (cd *ColumnList) add(col *tree.ColumnTableDef) {
	c := newColumn(cd.table, col)
	cd.colMap[col.Name.String()] = c
	cd.cols = append(cd.cols, c)
}

func (cd *ColumnList) All() []*Column {
	return cd.cols
}

func (cd *ColumnList) Find(name string) (*Column, bool) {
	col, ok := cd.colMap[name]
	return col, ok
}
