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

func (il *IndexList) add(ci *tree.CreateIndex) error {
	if ci.Name.String() == "\"\"" {
		return fmt.Errorf("index name is required")
	}

	idx := newIndex(il.table, ci)

	il.indexMap[ci.Name.String()] = idx
	il.indexes = append(il.indexes, idx)
	return nil
}

func (il *IndexList) addTableDef(idxdef *tree.IndexTableDef) error {
	if idxdef.Name.String() == "\"\"" {
		return fmt.Errorf("all indicies must be named")
	}

	idx := newIndex(il.table, &tree.CreateIndex{
		Table:            il.table.Tree.Table,
		Name:             idxdef.Name,
		Inverted:         idxdef.Inverted,
		Columns:          idxdef.Columns,
		Sharded:          idxdef.Sharded,
		Storing:          idxdef.Storing,
		PartitionByIndex: idxdef.PartitionByIndex,
		StorageParams:    idxdef.StorageParams,
		Predicate:        idxdef.Predicate,
		NotVisible:       idxdef.NotVisible,
		Unique:           false,
		IfNotExists:      false,
		Concurrently:     false,
	})

	il.indexMap[idxdef.Name.String()] = idx
	il.indexes = append(il.indexes, idx)
	return nil
}

type Index struct {
	Table *Table
	Name  string
	Tree  *tree.CreateIndex
}

func newIndex(tbl *Table, idx *tree.CreateIndex) *Index {
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
