package stool

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

type ViewExplainer struct {
	RootPath string
}

func (this *ViewExplainer) GetAllVariablesFrom(viewName string) (map[string]int, error) {
	c, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.blade.php", this.RootPath, viewName))

	contents := string(c)

	if (err != nil) {
		return nil, err
	}

	re := regexp.MustCompile(`\$[a-zA-Z_\x80-\xff][a-zA-Z0-9_\x80-\xff]*`)

	result := re.FindAllString(contents, -1)

	//for r := range result {
	//
	//}

	fmt.Println(result)

	return map[string]int{}, nil
}
