package diff

import "github.com/cygnetdigital/roachdiff/pkg/diff/model"

func CreateIndexes(d *Diff, a, b *model.Model) error {

	// look for new indexes by looking for table matches
	for _, tblNew := range b.Tables.All() {

		tblOld, hasOldTable := a.Tables.Find(tblNew.Name)

		// now compare all indexes
		for _, idxNew := range tblNew.Indexes.All() {

			if !hasOldTable {
				d.append(d.gen.NewCreateIndex(idxNew))
				continue
			}

			// find the matching index, if it exists do nothing
			if _, ok := tblOld.Indexes.Find(idxNew.Name); ok {
				continue

			}

			d.append(d.gen.NewCreateIndex(idxNew))
		}
	}

	return nil
}

func DropIndexes(d *Diff, a, b *model.Model) error {

	for _, tblNew := range b.Tables.All() {

		// skip if no matching table
		tblOld, ok := a.Tables.Find(tblNew.Name)
		if !ok {
			continue
		}

		// look through all the indexes on the old, to see
		// if they are present on the new.
		for _, colOld := range tblOld.Indexes.All() {

			// if the column is found on the new, no need to drop
			if _, ok := tblNew.Indexes.Find(colOld.Name); ok {
				continue

			}

			d.append(d.gen.NewDropIndex(colOld))
		}
	}

	return nil
}
