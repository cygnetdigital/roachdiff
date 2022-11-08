package diff

import (
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
	"github.com/cygnetdigital/roachdiff/pkg/diff/model"
)

func ColumnNullability(d *Diff, a, b *model.Model) error {

	// look for new columns by looking for table matches
	for _, tblNew := range b.Tables.All() {

		// skip if no matching table
		tblOld, ok := a.Tables.Find(tblNew.Name)
		if !ok {
			continue
		}

		// now compare all columns
		for _, colNew := range tblNew.Columns.All() {

			// find the matching column, if it doesnt exists do nothing
			colOld, ok := tblOld.Columns.Find(colNew.Name)
			if !ok {
				continue
			}

			colNewNull := colNew.Tree.Nullable.Nullability
			colOldNull := colOld.Tree.Nullable.Nullability

			// fix cases where nullabilty not been explicitly set
			if colNewNull == tree.SilentNull {
				colNewNull = tree.Null
			}
			if colOldNull == tree.SilentNull {
				colOldNull = tree.Null
			}

			// if nullability is the same do nothing
			if colNewNull == colOldNull {
				continue
			}

			if colNewNull == tree.NotNull {
				d.append(d.gen.NewAlterTableAlterColumnSetNotNull(colNew))
				continue
			}

			if colNewNull == tree.Null {
				d.append(d.gen.NewAlterTableAlterColumnDropNotNull(colNew))
				continue
			}

		}
	}

	return nil
}
