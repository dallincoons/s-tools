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
	path := strings.TrimRight(p, ".blade.php")
	base := strings.Trim(path, "/")

	pathLookup[strings.Replace(base, "/", ".", -1)] = p
	nameLookup[p] = strings.Replace(base, "/", ".", -1)
}

func (this *Blade) GetPath(path string) (string, bool) {
	name, ok := pathLookup[path]

	return name, ok
}

func (this *Blade) GetName(name string) (string, bool) {
	path, ok := nameLookup[name]

	return path, ok
}
