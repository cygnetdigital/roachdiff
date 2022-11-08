package diff

import "github.com/cygnetdigital/roachdiff/pkg/diff/model"

func CreateColumns(d *Diff, a, b *model.Model) error {

	// look for new columns by looking for table matches
	for _, tblNew := range b.Tables.All() {

		// skip if no matching table
		tblOld, ok := a.Tables.Find(tblNew.Name)
		if !ok {
			continue
		}

		// now compare all columns
		for _, colNew := range tblNew.Columns.All() {

			// find the matching column, if it exists do nothing
			if _, ok := tblOld.Columns.Find(colNew.Name); ok {
				continue

			}

			d.append(d.gen.NewAlterTableAddColumn(colNew))
		}
	}

	return nil
}

func DropColumns(d *Diff, a, b *model.Model) error {

	for _, tblNew := range b.Tables.All() {

		// skip if no matching table
		tblOld, ok := a.Tables.Find(tblNew.Name)
		if !ok {
			continue
		}

		// look through all the columns on the old, to see
		// if they are present on the new.
		for _, colOld := range tblOld.Columns.All() {

			// if the column is found on the new, no need to drop
			if _, ok := tblNew.Columns.Find(colOld.Name); ok {
				continue

			}

			d.append(d.gen.NewAlterTableDropColumn(colOld))
		}
	}

	return nil
}
