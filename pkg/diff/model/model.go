package model

import (
	"fmt"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
)

// Model is the declarative representation of a sql schema
type Model struct {
	Tables *TableList
	Enums  *EnumList
}

// NewModel parses a set of `CREATE (TABLE|ENUM)` statements
// to produce a model
func NewModel(sql string) (*Model, error) {
	statements, err := parser.Parse(sql)
	if err != nil {
		return nil, err
	}

	model := &Model{
		Tables: newTableList(),
		Enums:  newEnumList(),
	}

	for _, stmt := range statements {

		switch stm := stmt.AST.(type) {
		case *tree.CreateTable:
			if err := model.addTable(stm); err != nil {
				return nil, err
			}

		case *tree.CreateIndex:
			if err := model.addIndex(stm); err != nil {
				return nil, err
			}

		case *tree.CreateType:
			if err := model.addEnum(stm); err != nil {
				return nil, err
			}

		default:
			return nil, fmt.Errorf("do not support type %T", stm)
		}

	}

	return model, nil
}

func (m *Model) addTable(ct *tree.CreateTable) error {
	return m.Tables.add(ct)
}

func (m *Model) addIndex(ci *tree.CreateIndex) error {
	table := ci.Table.String()

	tbl, ok := m.Tables.tableMap[table]
	if !ok {
		return fmt.Errorf("table '%s' not found", table)
	}

	return tbl.Indexes.add(ci)
}

func (m *Model) addEnum(ct *tree.CreateType) error {
	return m.Enums.add(ct)
}

type TableList struct {
	tables   []*Table
	tableMap map[string]*Table
}

func newTableList() *TableList {
	return &TableList{
		tables:   []*Table{},
		tableMap: map[string]*Table{},
	}
}

func (tl *TableList) add(ct *tree.CreateTable) error {
	t, err := NewTable(ct)
	if err != nil {
		return fmt.Errorf("unable to parse CREATE TABLE %s statement: %w", ct.Table.String(), err)
	}

	tl.tableMap[t.Name] = t
	tl.tables = append(tl.tables, t)
	return nil
}

func (tl *TableList) All() []*Table {
	return tl.tables
}

// Find table by name
func (tl *TableList) Find(name string) (*Table, bool) {
	t, ok := tl.tableMap[name]
	return t, ok
}

type EnumList struct {
	enums   []*Enum
	enumMap map[string]*Enum
}

func newEnumList() *EnumList {
	return &EnumList{
		enumMap: make(map[string]*Enum),
	}
}

func (el *EnumList) add(ce *tree.CreateType) error {
	e, err := NewEnum(ce)
	if err != nil {
		return fmt.Errorf("unable to parse CREATE TYPE %s statement: %w", ce.TypeName.String(), err)
	}

	el.enumMap[e.Name] = e
	el.enums = append(el.enums, e)
	return nil
}

func (el *EnumList) All() []*Enum {
	return el.enums
}

// Find enum by name
func (el *EnumList) Find(name string) (*Enum, bool) {
	e, ok := el.enumMap[name]
	return e, ok
}
