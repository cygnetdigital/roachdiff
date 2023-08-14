package diff

import (
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
	"github.com/cygnetdigital/roachdiff/pkg/diff/model"
)

func CreateEnum(d *Diff, a, b *model.Model) error {
	// look for new enums by looping through new, and comparing against old
	for _, enumNew := range b.Enums.All() {
		// skip if enum already exists
		if _, ok := a.Enums.Find(enumNew.Name); ok {
			continue
		}

		d.append(d.gen.NewCreateEnum(enumNew))
	}

	return nil
}

func DropEnum(d *Diff, a, b *model.Model) error {
	// look for enums to drop by looping through old, and comparing against the new
	for _, enumOld := range a.Enums.All() {
		// skip if enum still exists
		if _, ok := b.Enums.Find(enumOld.Name); ok {
			continue
		}

		d.append(d.gen.NewDropEnum(enumOld))
	}

	return nil
}

// CreateEnumValues
func CreateEnumValues(d *Diff, a, b *model.Model) error {
	// for each enum that exists in bpoth old and new, insert any new values
	for _, enumNew := range b.Enums.All() {
		enumOld, ok := a.Enums.Find(enumNew.Name)
		if !ok {
			continue
		}

		// shortcut if no old values to insert relative to.
		if len(enumOld.Values.All()) == 0 {
			for _, newVal := range enumNew.Values.All() {
				d.append(d.gen.NewAlterTypeAddValue(enumNew, newVal, nil))
			}
			continue
		}

		// All values in `new` are either to be created, or already exist.
		// So we can always place relative to the first value in `old`, or the last seen new value.

		// old:      bar, bix,      qix
		// new: foo, bar,      baz, qix, qux

		// when begin is true, we are inserting before the first existing (i.e. old) value.
		// when begin is false, we are inserting after the last seen value.

		lastSeen, begin := enumOld.Values.All()[0], true
		for _, newVal := range enumNew.Values.All() {
			if enumOld.Values.Find(newVal) {
				lastSeen, begin = newVal, false
				continue
			}

			// newVal is not in enumOld, so we need to insert it.
			d.append(d.gen.NewAlterTypeAddValue(enumNew, newVal, &tree.AlterTypeAddValuePlacement{
				Before:      begin,
				ExistingVal: tree.EnumValue(lastSeen),
			}))

			if !begin {
				lastSeen = newVal
			}
		}

	}
	return nil
}

// DropEnumValues
func DropEnumValues(d *Diff, a, b *model.Model) error {
	// for each enum that exists in both old and new, drop any values that no longer exist
	for _, enumNew := range b.Enums.All() {
		enumOld, ok := a.Enums.Find(enumNew.Name)
		if !ok {
			continue
		}

		for _, oldValue := range enumOld.Values.All() {
			if enumNew.Values.Find(oldValue) {
				continue
			}

			d.append(d.gen.NewAlterTypeDropValue(enumNew, oldValue))
		}
	}
	return nil
}
