package gonativeextractor_test

import (
	"testing"

	"github.com/SpongeData-cz/gonativeextractor"
)

func TestExtractor(t *testing.T) {

	t.Run("glob", func(t *testing.T) {

		// Creating extractor
		e := gonativeextractor.NewExtractor(-1, -1, 0)

		// Adding miner
		err := e.AddMinerSo(gonativeextractor.DEFAULT_MINERS_PATH+"/glob_entities.so", "match_glob", []byte("world\x00"))
		if err != nil {
			t.Fatal(err.Error())
		}

		// Adding non-existent miner
		err = e.AddMinerSo(gonativeextractor.DEFAULT_MINERS_PATH+"/glob_entities.so", "match_glob_hhh", []byte("world\x00"))
		if err == nil {
			t.Error("Should return unknown symbol error.")
		}

		// Adding miner from non-existent file
		err = e.AddMinerSo(gonativeextractor.DEFAULT_MINERS_PATH+"/glob_entities_hhh.so", "match_glob", []byte("world\x00"))
		if err == nil {
			t.Error("Should return unknown path to the .so.")
		}

		// Creating and setting stream
		st, err := gonativeextractor.NewBufferStream([]byte("Hello world byte\x00"))
		if err != nil {
			t.Fatal(err.Error())
		}
		err = e.SetStream(st)
		if err != nil {
			t.Fatal(err.Error())
		}

		// Output validation
		for !e.Eof() {
			r, err := e.Next()
			if err != nil {
				t.Fatal(err.Error())
			}
			if r.Eof() {
				t.Fatal("Should find an occurrence.")
			}
			if r.Str() != "world" {
				t.Error("Str should be 'world'.")
			}
			if r.Pos() != 6 {
				t.Error("Pos should be 6.")
			}
			if r.Upos() != 6 {
				t.Error("Upos should be 6.")
			}
			if r.Len() != 5 {
				t.Error("Len should be 5.")
			}
			if r.Ulen() != 5 {
				t.Error("Ulen should be 5.")
			}
			if r.Label() != "Glob" {
				t.Error("Label should be 'Glob'.")
			}
			if r.Prob() != 1 {
				t.Error("Prob should be 1.")
			}
			r.Next()
			if !r.Eof() {
				t.Error("Should find only one occurrence.")
			}
		}

		// Meta validation
		meta := e.Meta()
		if len(meta) != 1 {
			t.Fatal("Meta should have exactly one element.")
		}
		if meta[0].Ldpath != gonativeextractor.DEFAULT_MINERS_PATH+"/glob_entities.so" {
			t.Error("Wrong path in meta.")
		}
		if meta[0].Ldptr == nil {
			t.Error("Nil pointer in meta.")
		}
		if meta[0].Ldsymb != "match_glob" {
			t.Error("Ldsymb should be 'match_glob'.")
		}
		if meta[0].Params != "world" {
			t.Error("Params should be 'world'.")
		}
		if len(meta[0].Meta) != 1 {
			t.Fatal("Meta[0].Meta should have exactly one element.")
		}
		if meta[0].Meta[0] != "match_glob" {
			t.Error("Meta[0].Meta[0] should be 'match_glob'.")
		}

		if e.GetLastError() == nil {
			t.Error("No errors but there should be.")
		}

		// Setting and unsetting flags
		err = e.SetFlags(gonativeextractor.E_SORT_RESULTS)
		if err != nil {
			t.Error(err.Error())
		}
		err = e.UnsetFlags(gonativeextractor.E_SORT_RESULTS)
		if err != nil {
			t.Error(err.Error())
		}

		// Destroying the extractor
		if e.Destroy() != nil {
			t.Error(err.Error())
		}

		if e.Destroy() == nil {
			t.Error("Not throwing error when destroyed the second time.")
		}

	})

	t.Run("globFileStream", func(t *testing.T) {

		// Creating extractor
		e := gonativeextractor.NewExtractor(-1, -1, 0)

		// Adding miner
		err := e.AddMinerSo(gonativeextractor.DEFAULT_MINERS_PATH+"/glob_entities.so", "match_glob", []byte("world\x00"))
		if err != nil {
			t.Fatal(err.Error())
		}

		// Creating and setting file stream
		st, err := gonativeextractor.NewFileStream("./fixtures/hworld.txt")
		if err != nil {
			t.Fatal(err.Error())
		}
		err = e.SetStream(st)
		if err != nil {
			t.Fatal(err.Error())
		}

		// Output validation
		for !e.Eof() {
			r, err := e.Next()
			if err != nil {
				t.Fatal(err.Error())
			}
			for !r.Eof() {
				if r.Str() != "world" {
					t.Error("Str should be 'world'.")
				}
				if r.Pos() != 6 {
					t.Error("Pos should be 6.")
				}
				if r.Upos() != 6 {
					t.Error("Upos should be 6.")
				}
				if r.Len() != 5 {
					t.Error("Len should be 5.")
				}
				if r.Ulen() != 5 {
					t.Error("Ulen should be 5.")
				}
				if r.Label() != "Glob" {
					t.Error("Label should be 'Glob'.")
				}
				if r.Prob() != 1 {
					t.Error("Prob should be 1.")
				}
				r.Next()
			}
		}

		// Meta validation
		meta := e.Meta()
		if len(meta) != 1 {
			t.Fatal("Meta should have exactly one element.")
		}
		if meta[0].Ldpath != gonativeextractor.DEFAULT_MINERS_PATH+"/glob_entities.so" {
			t.Error("Wrong path in meta.")
		}
		if meta[0].Ldptr == nil {
			t.Error("Nil pointer in meta.")
		}
		if meta[0].Ldsymb != "match_glob" {
			t.Error("Ldsymb should be 'match_glob'.")
		}
		if meta[0].Params != "world" {
			t.Error("Params should be 'world'.")
		}
		if len(meta[0].Meta) != 1 {
			t.Fatal("Meta[0].Meta should have exactly one element.")
		}
		if meta[0].Meta[0] != "match_glob" {
			t.Error("Meta[0].Meta[0] should be 'match_glob'.")
		}

		err = e.GetLastError()
		if err != nil {
			t.Error("Unexpected error: ", err.Error())
		}

		// Destroying the extractor
		if e.Destroy() != nil {
			t.Error(err.Error())
		}

	})

	t.Run("Un/setStream", func(t *testing.T) {

		// Creating extractor
		e := gonativeextractor.NewExtractor(-1, -1, 0)

		// Creating and setting stream
		st, err := gonativeextractor.NewBufferStream([]byte("Hello world byte\x00"))
		if err != nil {
			t.Fatal(err.Error())
		}
		err = e.SetStream(st)
		if err != nil {
			t.Fatal(err.Error())
		}

		// Closing stream
		err = st.Close()
		if err != nil {
			t.Error(err.Error())
		}

		// Unsetting stream
		e.UnsetStream()
		if !e.Eof() {
			t.Error("Stream is set but should not be.")
		}
		_, err = e.Next()
		if err == nil {
			t.Error("Stream is set but should not be.")
		}

		// Destroying the extractor
		if e.Destroy() != nil {
			t.Error(err.Error())
		}

	})

}
