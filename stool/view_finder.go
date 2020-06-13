package stool

import (
	"fmt"
	"os"
	"strings"
)

type ViewFinder struct {
	Root string
}

func (this *ViewFinder) getFile (filename string) (*os.File, error) {
	filepath := this.GetFilePath(filename)

	file, err := os.Open(filepath)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (this *ViewFinder) GetFilePath(view_name string) string {
	filename := strings.ReplaceAll(view_name, ".", "/")

	return fmt.Sprintf("%s/%s.blade.php", this.Root, filename)
}
