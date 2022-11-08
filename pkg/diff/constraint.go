package diff

import "github.com/cygnetdigital/roachdiff/pkg/diff/model"

func CreateConstraints(d *Diff, a, b *model.Model) error {

	// look for new constraints by looking for table matches
	for _, tblNew := range b.Tables.All() {

		tblOld, hasOldTable := a.Tables.Find(tblNew.Name)

		// now compare all constraints
		for _, consNew := range tblNew.Constraints.All() {

			// skip constraint if old table doesn't exist yet
			if !hasOldTable {
				continue
			}

			// find the matching column, if it exists do nothing
			if _, ok := tblOld.Constraints.Find(consNew.Name); ok {
				continue

			}

			d.append(d.gen.NewAlterTableAddConstraint(consNew))
		}
	}

	return nil
}

func DropConstraints(d *Diff, a, b *model.Model) error {

	for _, tblNew := range b.Tables.All() {

		// skip if no matching table
		tblOld, ok := a.Tables.Find(tblNew.Name)
		if !ok {
			continue
		}

		// look through all the constraints on the old, to see
		// if they are present on the new.
		for _, consOld := range tblOld.Constraints.All() {

			// if the column is found on the new, no need to drop
			if _, ok := tblNew.Constraints.Find(consOld.Name); ok {
				continue

			}

			d.append(d.gen.NewAlterTableDropConstraint(consOld))
		}
	}

	return nil
}
