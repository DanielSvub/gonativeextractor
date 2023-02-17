package gonativeextractor

/*
#cgo CFLAGS: -I/usr/include/nativeextractor
#cgo LDFLAGS: -lnativeextractor -lglib-2.0 -ldl
#include <nativeextractor/common.h>
#include <nativeextractor/extractor.h>
#include <nativeextractor/stream.h>
bool extractor_c_add_miner_from_so(extractor_c * self, const char * miner_so_path, const char * miner_name, void * params );
const char * extractor_get_last_error(extractor_c * self);
*/
import "C"
import (
	"fmt"
	"io"
)

type Streamer interface {
	Open()
	Check() bool
	GetStream() *C.struct_stream_c
	io.Closer
}

type FileStream struct {
	Ptr  *C.struct_stream_file_c
	Path string
}

func (ego *FileStream) GetStream() *C.struct_stream_c {
	return &ego.Ptr.stream
}

func (ego *FileStream) Check() bool {
	return ego.Ptr.stream.state_flags&C.STREAM_FAILED == 0
}

func NewFileStream(path string) (*FileStream, error) {
	out := FileStream{Path: path}
	out.Ptr = C.stream_file_c_new(C.CString(path))
	if !out.Check() {
		return nil, fmt.Errorf("Seek out of bounds")
	}

	return &out, nil
}
