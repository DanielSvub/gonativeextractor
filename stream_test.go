package gonativeextractor_test

import (
	"testing"

	"github.com/DanielSvub/gonativeextractor"
)

func TestFileStream(t *testing.T) {
	t.Run("isstream", func(t *testing.T) {
		_, err := gonativeextractor.NewFileStream("./fixtures/create_stream.txt")
		if err != nil {
			t.Errorf("Stream creation failed.")
		}
	})

	t.Run("nostream", func(t *testing.T) {
		_, err := gonativeextractor.NewFileStream("./fixtures/nosuchfile.txt")
		if err == nil {
			t.Errorf("No error when trying to open non-existent file.")
		}
	})

	t.Run("closestream", func(t *testing.T) {
		s, _ := gonativeextractor.NewFileStream("./fixtures/create_stream.txt")
		if s.Close() != nil {
			t.Errorf("Stream closing failed.")
		}
		if s.Close() == nil {
			t.Errorf("No error when closed the second time.")
		}
	})
}

func TestBufferStream(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		stream, err := gonativeextractor.NewBufferStream([]byte("test\x00"))
		if err != nil {
			t.Errorf("Stream creation failed.")
		}
		err = stream.Close()
		if err != nil {
			t.Errorf("Stream closing failed.")
		}
		err = stream.Close()
		if err == nil {
			t.Errorf("No error when closed the second time.")
		}
	})

	t.Run("nil", func(t *testing.T) {
		_, err := gonativeextractor.NewBufferStream(nil)
		if err == nil {
			t.Errorf("No error when nil given as buffer.")
		}
	})
}
