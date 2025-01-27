package stool

import (
	"bufio"
	"io/ioutil"
	"os"
	"regexp"
)

type VariableCollector struct {
}

func (this *VariableCollector) GetAllVariablesFrom(path string) (map[string]int, error) {
	variablesCounts := make(map[string]int)
	c, err := ioutil.ReadFile(path)

	contents := string(c)

	if (err != nil) {
		return nil, err
	}

	re := regexp.MustCompile(`\$[a-zA-Z_\x80-\xff][a-zA-Z0-9_\x80-\xff]*`)

	result := re.FindAllString(contents, -1)

	for _, r := range result {
		variablesCounts[r]++
	}

	return variablesCounts, nil
}

func GetExplainer(view_root string) ViewExplainer {
	finder := &ViewFinder{
		view_root,
	}

	return ViewExplainer{ViewIndexer:ViewIndexer{
		RootDir:    view_root,
		Explainer:  &VariableCollector{},
		ViewFinder: finder,
		Writer:     bufio.NewWriter(os.Stdout),
	}}
}
