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

	t.Run("glob", func(t *testing.T) {
		e := gonativeextractor.NewExtractor(-1, -1, 0)

		err := e.AddMinerSo(gonativeextractor.DEFAULT_MINERS_PATH+"/glob_entities.so", "match_glob", []byte("world"))
		if err != nil {
			t.Errorf(err.Error())
		}

		err = e.AddMinerSo(gonativeextractor.DEFAULT_MINERS_PATH+"/glob_entities.so", "match_glob_hhh", []byte("world"))
		if err == nil {
			t.Errorf("Should return unknown symbol error.")
		}

		err = e.AddMinerSo(gonativeextractor.DEFAULT_MINERS_PATH+"/glob_entities_hhh.so", "match_glob", []byte("world"))
		if err == nil {
			t.Errorf("Should return unknown path to the .so.")
		}

		st, err := gonativeextractor.NewBufferStream([]byte("Hello world byte"))
		if err != nil {
			t.Errorf(err.Error())
		}
		e.SetStream(st)

		for !e.Eof() {
			_, err := e.Next()
			if err != nil {
				t.Errorf(err.Error())
			}
		}

		e.Meta()

		err = e.GetLastError()
		if err == nil {
			t.Errorf("No errors but there should be.")
		}

		e.UnsetStream()

		for !e.Eof() {
			_, err := e.Next()
			if err == nil {
				t.Errorf("Stream is set but should not be.")
			}
		}

		e.SetFlags(1 << 1)
		e.UnsetFlags(1 << 1)

		e.Close()

	})

}
