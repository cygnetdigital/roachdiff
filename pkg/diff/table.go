package diff

import (
	"github.com/cygnetdigital/roachdiff/pkg/diff/model"
)

func CreateTable(d *Diff, a, b *model.Model) error {

	// look for new tables by looping through new, and comparing
	// against old
	for _, tblNew := range b.Tables.All() {

		// skip if table already exists
		if _, ok := a.Tables.Find(tblNew.Name); ok {
			continue
		}

		d.append(d.gen.NewCreateTable(tblNew))
	}

	return nil
}

func DropTable(d *Diff, a, b *model.Model) error {

	// look for tables to drop by looping through old, and comparing
	// against the new
	for _, tblOld := range a.Tables.All() {

		// skip if table still exists
		if _, ok := b.Tables.Find(tblOld.Name); ok {
			continue
		}

		d.append(d.gen.NewDropTable(tblOld))
	}

	return nil
}
