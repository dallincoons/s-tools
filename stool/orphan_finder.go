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
		RootDir:    filepath.Join(this.Root, "resources", "views"),
		Explainer:  &VariableCollector{},
		ViewFinder: &ViewFinder{},
		Writer:     bufio.NewWriter(os.Stdout),
	}

	controllerUsages := make(map[string]bool)

	filepath.Walk(filepath.Join(this.Root, "app"), func(path string, info os.FileInfo, err error) error {
		if (!strings.HasSuffix(info.Name(), "Controller.php") && !strings.HasSuffix(info.Name(), "DataTable.php") && !strings.Contains(path, "Mail")) {
			return nil
		}

		contents, _ := ioutil.ReadFile(path)

		re := regexp.MustCompile("(?:view|render|make|route)\\('(.+?)'[),]")

		for _, usage := range re.FindAllStringSubmatch(string(contents), -1) {
			controllerUsages[usage[1]]	= true
		}

		return nil
	})

	views := indexer.IndexViews(filepath.Join(this.Root, "resources", "views"))

	orphanFound := make(map[string]bool)
	orphanCandidates := []string{}
	for _, view := range views {
		controllerUsageFound := controllerUsages[view.Name]
		orphanAlreadyFound := orphanFound[view.Name]
		if controllerUsageFound || orphanAlreadyFound {
			continue
		}
		if len(view.Parents) == 0 && len(view.Children) == 0 {
			orphanCandidates = append(orphanCandidates, view.Name)
			orphanFound[view.Name] = true
		}
	}

	return orphanCandidates
}
