package model

import (
	"fmt"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
)

// IndexList ,,,
type IndexList struct {
	table    *Table
	indexes  []*Index
	indexMap map[string]*Index
}

func newIndexList(table *Table) *IndexList {
	return &IndexList{
		table:    table,
		indexes:  []*Index{},
		indexMap: map[string]*Index{},
	}
}

func (il *IndexList) add(idxdef *tree.IndexTableDef) error {
	if idxdef.Name.String() == "\"\"" {
		return fmt.Errorf("all indicies must be named")
	}

	idx := newIndex(il.table, idxdef)

	il.indexMap[idxdef.Name.String()] = idx
	il.indexes = append(il.indexes, idx)
	return nil
}

type Index struct {
	Table *Table
	Name  string
	Tree  *tree.IndexTableDef
}

func newIndex(tbl *Table, idx *tree.IndexTableDef) *Index {
	return &Index{
		Table: tbl,
		Tree:  idx,
		Name:  idx.Name.String(),
	}
}

func (il *IndexList) All() []*Index {
	return il.indexes
}

func (il *IndexList) Find(name string) (*Index, bool) {
	col, ok := il.indexMap[name]
	return col, ok
}
