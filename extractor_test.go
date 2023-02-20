package gonativeextractor_test

import (
	"testing"

	"github.com/SpongeData-cz/gonativeextractor"
)

func TestExtractor(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		e := gonativeextractor.NewExtractor(-1, -1, 0)
		if e == nil {
			t.Errorf("Should return non nil pointer to Extractor")
		}

		if e.Close() != nil {
			t.Errorf("Problem in extractor close.")
		}
		if e.Close() == nil {
			t.Errorf("Not throwing error when closed the second time.")
		}
	})

}
