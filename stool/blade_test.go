package stool

import (
	"log"
	"testing"
)

func TestStoresViewPathAsBladeName(t *testing.T) {
	b := newBlade()

	b.AddPath("/hot/topic/is/not/punkrock.blade.php")
	path, _ := b.GetPath("hot.topic.is.not.punkrock")

	if (path != "/hot/topic/is/not/punkrock.blade.php") {
		log.Fatalf("expected /hot/topic/is/not/punkrock.blade.php, got %s", path)
	}
}

func newBlade() *Blade {
	return &Blade{
		RootDir: "hot/topic",
	}
}
