package stool

import (
	"fmt"
	"os"
	"strings"
)

type ViewFinder struct {
	Root string
}

func (this *ViewFinder) getFile (fileName string) (*os.File, error) {
	filepath := strings.ReplaceAll(fileName, ".", "/")

	file, err := os.Open(fmt.Sprintf("%s/%s.blade.php", this.Root, filepath))

	if err != nil {
		return nil, err
	}

	return file, nil
}
