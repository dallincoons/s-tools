package stool

import (
	"log"
	"testing"
)

func TestGetNamesOfOrphanFiles(t *testing.T) {
	finder := newOrphanFinder()
	orphanMap := make(map[string]bool)

	orphans := finder.GetOrphans()

	for _, orphan := range orphans {
		orphanMap[orphan] = true
	}

	dir1view1found := orphanMap["dir1.view3"]
	dir2view2found := orphanMap["dir2.view2"]
	usedbycontrollerfound := orphanMap["used_by_controller"]
	usedbydatatable := orphanMap["used_by_data_table"]

	if !dir1view1found {
		log.Fatalf("orphan dir1.view3 not found")
	}

	if !dir2view2found {
		log.Fatalf("orphan dir2.view2 not found")
	}

	if usedbycontrollerfound {
		log.Fatalf("non orphan usedbycontrollerfound found")
	}

	if usedbydatatable {
		log.Fatalf("non orphan usedbydatatable found")
	}
}

func newOrphanFinder() OrphanFinder {
	return OrphanFinder{
		Root: "../test_fixtures/views2",
	}
}
