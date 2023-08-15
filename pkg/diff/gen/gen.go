package gen

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
	"github.com/cygnetdigital/roachdiff/pkg/diff/model"
)

// Generator converts cockroachdb tree nodes into strings
type Generator struct {
	PrettyCfg tree.PrettyCfg

	// comment out potentially dangerous statements with a warning
	Warnings bool
}

// Statement encapsulates a SQL statement, and whether it should generate a warning
type Statement struct {
	SQL        string
	HasWarning bool
}

// NewGenerator configures a generator with sensible defaults
func NewGenerator() *Generator {
	cfg := tree.DefaultPrettyCfg()

	cfg.LineWidth = 120
	cfg.TabWidth = 2
	cfg.UseTabs = false
	cfg.Align = tree.PrettyAlignAndDeindent

	return &Generator{PrettyCfg: cfg}
}

func (g *Generator) string(stm tree.NodeFormatter) Statement {
	return Statement{
		SQL:        fmt.Sprintf("%s;", g.PrettyCfg.Pretty(stm)),
		HasWarning: false,
	}
}

func (g *Generator) stringWithWarning(stm tree.NodeFormatter) Statement {
	stmt := g.string(stm)
	stmt.HasWarning = true

	if g.Warnings {
		// add hypen comments to each new line
		s := strings.ReplaceAll(stmt.SQL, "\n", "\n--")

		warning := `WARNING: This is a destructive operation`
		stmt.SQL = fmt.Sprintf("-- %s\n-- %s", warning, s)
	}

	return stmt
}

func (g *Generator) NewCreateTable(table *model.Table) Statement {
	return g.string(table.Tree)
}

func (g *Generator) NewDropTable(table *model.Table) Statement {

	dt := &tree.DropTable{
		Names: tree.TableNames{table.Tree.Table},
	}

	return g.stringWithWarning(dt)
}

func (g *Generator) NewAlterTableAddColumn(col *model.Column) Statement {
	at := &tree.AlterTable{
		Table: col.Table.Tree.Table.ToUnresolvedObjectName(),
		Cmds: tree.AlterTableCmds{
			&tree.AlterTableAddColumn{
				ColumnDef: col.Tree,
			},
		},
	}

	return g.string(at)
}

func (g *Generator) NewAlterTableDropColumn(col *model.Column) Statement {
	dc := &tree.AlterTable{
		Table: col.Table.Tree.Table.ToUnresolvedObjectName(),
		Cmds: tree.AlterTableCmds{
			&tree.AlterTableDropColumn{
				Column: col.Tree.Name,
			},
		},
	}

	return g.stringWithWarning(dc)
}

func (g *Generator) NewCreateIndex(idx *model.Index) Statement {
	return g.string(idx.Tree)
}

func (g *Generator) NewDropIndex(idx *model.Index) Statement {
	ci := &tree.DropIndex{
		IndexList: tree.TableIndexNames{
			&tree.TableIndexName{
				Table: idx.Table.Tree.Table,
				Index: tree.UnrestrictedName(idx.Tree.Name),
			},
		},
	}

	return g.stringWithWarning(ci)
}

func (g *Generator) NewAlterTableAddConstraint(cons *model.Constraint) Statement {
	at := &tree.AlterTable{
		Table: cons.Table.Tree.Table.ToUnresolvedObjectName(),
		Cmds: tree.AlterTableCmds{
			&tree.AlterTableAddConstraint{
				ConstraintDef: cons.Tree,
			},
		},
	}

	return g.stringWithWarning(at)
}

func (g *Generator) NewAlterTableDropConstraint(cons *model.Constraint) Statement {
	at := &tree.AlterTable{
		Table: cons.Table.Tree.Table.ToUnresolvedObjectName(),
		Cmds: tree.AlterTableCmds{
			&tree.AlterTableDropConstraint{
				Constraint: tree.Name(cons.Name),
			},
		},
	}
	return g.stringWithWarning(at)
}

func (g *Generator) NewAlterTableAlterColumnSetNotNull(col *model.Column) Statement {
	stm := &tree.AlterTable{
		Table: col.Table.Tree.Table.ToUnresolvedObjectName(),
		Cmds: tree.AlterTableCmds{
			&tree.AlterTableSetNotNull{
				Column: col.Tree.Name,
			},
		},
	}

	return g.stringWithWarning(stm)
}

func (g *Generator) NewAlterTableAlterColumnDropNotNull(col *model.Column) Statement {
	stm := &tree.AlterTable{
		Table: col.Table.Tree.Table.ToUnresolvedObjectName(),
		Cmds: tree.AlterTableCmds{
			&tree.AlterTableDropNotNull{
				Column: col.Tree.Name,
			},
		},
	}

	return g.stringWithWarning(stm)
}

func (g *Generator) NewCreateEnum(enum *model.Enum) Statement {
	return g.string(enum.Enum)
}

func (g *Generator) NewDropEnum(enum *model.Enum) Statement {
	de := &tree.DropType{
		Names: []*tree.UnresolvedObjectName{enum.Enum.TypeName},
	}

	return g.stringWithWarning(de)
}

func (g *Generator) NewAlterTypeAddValue(enum *model.Enum, value string, placement *tree.AlterTypeAddValuePlacement) Statement {
	at := &tree.AlterType{
		Type: enum.Enum.TypeName,
		Cmd: &tree.AlterTypeAddValue{
			NewVal:    tree.EnumValue(value),
			Placement: placement,
		},
	}

	return g.string(at)
}

func (g *Generator) NewAlterTypeDropValue(enum *model.Enum, value string) Statement {
	at := &tree.AlterType{
		Type: enum.Enum.TypeName,
		Cmd: &tree.AlterTypeDropValue{
			Val: tree.EnumValue(value),
		},
	}

	return g.stringWithWarning(at)
}
