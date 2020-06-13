package stool

import (
	"io/ioutil"
	"regexp"
)

type ViewIndexer struct {
	Path string
}

func (this *ViewIndexer) FindAllIncludes() ([]string, error) {
	var viewNames []string

	contents, err := ioutil.ReadFile(this.Path)

	if err != nil {
		return nil, err
	}

	return this.getViewNames(contents, viewNames), nil
}

func (this *ViewIndexer) getViewNames(contents []byte, viewNames []string) []string {
	re := regexp.MustCompile("@include\\('(.+)'\\)")

	results := re.FindAllStringSubmatch(string(contents), -1)

	for _, result := range results {
		viewNames = append(viewNames, result[1])
	}

	return viewNames
}
