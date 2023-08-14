package model

import (
	"fmt"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
)

// Enum represents the declarative state of a `CREATE TYPE _ AS ENUM` statement
type Enum struct {
	Name   string
	Values EnumValueList
	Enum   *tree.CreateType
}

// NewEnum parses a CREATE TYPE statement and returns an Enum
func NewEnum(ct *tree.CreateType) (*Enum, error) {
	if ct.Variety != tree.Enum {
		return nil, fmt.Errorf("CREATE TYPE %s is not supported", ct.Variety.String())
	}

	e := &Enum{
		Name: ct.TypeName.Object(),
		Enum: ct,
	}

	e.Values = *newEnumValueList(e)

	for _, v := range ct.EnumLabels {
		e.Values.add(string(v))
	}

	return e, nil
}

type EnumValueList struct {
	enum     *Enum
	values   []string
	valueMap map[string]struct{}
}

func newEnumValueList(e *Enum) *EnumValueList {
	return &EnumValueList{
		enum:     e,
		values:   nil,
		valueMap: make(map[string]struct{}),
	}
}

func (ev *EnumValueList) add(v string) {
	ev.values = append(ev.values, v)
	ev.valueMap[v] = struct{}{}
}

func (ev *EnumValueList) All() []string {
	return ev.values
}

func (ev *EnumValueList) Find(v string) bool {
	_, ok := ev.valueMap[v]
	return ok
}
