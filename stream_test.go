package gonativeextractor_test

import (
	"testing"

	"github.com/SpongeData-cz/gonativeextractor"
)

func TestFileStream(t *testing.T) {
	t.Run("isstream", func(t *testing.T) {
		_, err := gonativeextractor.NewFileStream("./fixtures/create_stream.txt")
		if err != nil {
			t.Errorf("Error unexpected")
		}
	})

	t.Run("nostream", func(t *testing.T) {
		_, err := gonativeextractor.NewFileStream("./fixtures/nosuchfile.txt")
		if err == nil {
			t.Errorf("Error expected but not set")
		}
	})

	t.Run("closestream", func(t *testing.T) {
		s, _ := gonativeextractor.NewFileStream("./fixtures/create_stream.txt")
		if s.Close() != nil {
			t.Errorf("Error unexpected")
		}
		if s.Close() == nil {
			t.Errorf("Error expected but not set")
		}
	})
}
