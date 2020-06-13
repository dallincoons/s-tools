package stool

import (
	"fmt"
	"log"
	"testing"
)

func TestFindViewFromDotDelimitedString(t *testing.T) {
	finder := getViewFinder()

	file, err := finder.getFile("dir1.dir2.view1")

	if err != nil {
		log.Fatal("could not find file")
	}

	if file.Name() != "../test_fixtures/views/dir1/dir2/view1.blade.php" {
		log.Fatal(
			fmt.Sprintf("file not found, got %s, expected %s", file.Name(), "../test_fixtures/views/dir1/dir2/view1.blade.php"),
		)
	}
}

func getViewFinder() *ViewFinder {
	return &ViewFinder{
		Root: "../test_fixtures/views",
	}
}
