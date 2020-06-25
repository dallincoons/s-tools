package stool

import (
	"strings"
)

var pathLookup = make(map[string]string)
var nameLookup = make(map[string]string)

type Blade struct {
	RootDir string
}

func (this *Blade) AddPath(p string) {
	base := strings.Trim(strings.TrimSuffix(p, ".blade.php"), "/")

	dottedBase := strings.Replace(base, "/", ".", -1)

	pathLookup[dottedBase] = p
	nameLookup[p] = dottedBase
}

func (this *Blade) GetPath(path string) (string, bool) {
	name, ok := pathLookup[path]

	return name, ok
}

func (this *Blade) GetName(name string) (string, bool) {
	path, ok := nameLookup[name]

	return path, ok
}
