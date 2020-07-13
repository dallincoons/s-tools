package stool

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type OrphanFinder struct {
	Root string
}

func (this *OrphanFinder) GetOrphans() []string {
	indexer := ViewIndexer{
		RootDir:    this.Root,
		Explainer:  &VariableCollector{},
		ViewFinder: &ViewFinder{},
		Writer:     bufio.NewWriter(os.Stdout),
	}

	controllerUsages := make(map[string]bool)

	filepath.Walk(this.Root, func(path string, info os.FileInfo, err error) error {
		if (!strings.HasSuffix(info.Name(), "Controller.php") && !strings.HasSuffix(info.Name(), "DataTable.php")) {
			return nil
		}

		contents, _ := ioutil.ReadFile(path)

		re := regexp.MustCompile("(?:view|render|route)\\('(.+?)'[),]")

		for _, usage := range re.FindAllStringSubmatch(string(contents), -1) {
			controllerUsages[usage[1]]	= true
		}

		return nil
	})

	views := indexer.IndexViews(this.Root)

	orphanCandidates := []string{}
	for _, view := range views {
		controllerUsageFound := controllerUsages[view.Name]
		if controllerUsageFound {
			continue
		}
		if len(view.Parents) == 0 && len(view.Children) == 0 {
			orphanCandidates = append(orphanCandidates, view.Name)
		}
	}

	return orphanCandidates
}