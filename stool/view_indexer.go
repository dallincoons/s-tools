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

	re := regexp.MustCompile("@include\\('(.+)'\\)")

	results := re.FindAllStringSubmatch(string(contents), -1)

	for _, result := range results {
		viewNames = append(viewNames, result[1])
	}

	return viewNames, nil
}
