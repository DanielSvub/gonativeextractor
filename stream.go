package gonativeextractor

/*
#cgo CFLAGS: -I/usr/include/nativeextractor
#cgo LDFLAGS: -lnativeextractor -lglib-2.0 -ldl
#include <nativeextractor/common.h>
#include <nativeextractor/extractor.h>
#include <nativeextractor/stream.h>
bool extractor_c_add_miner_from_so(extractor_c * self, const char * miner_so_path, const char * miner_name, void * params );
const char * extractor_get_last_error(extractor_c * self);
void stream_c_destroy(stream_c * self);
*/
import "C"
import (
	"fmt"
	"io"
	"unsafe"
)

type Streamer interface {
	GetStream() *C.struct_stream_c
	Check() bool
	io.Closer
}

type FileStream struct {
	Ptr  *C.struct_stream_file_c
	Path string
}

/*
Gets the inner stream structure.
Returns: pointer to the C struct stream_c
*/
func (ego *FileStream) GetStream() *C.struct_stream_c {
	return &ego.Ptr.stream
}

/*
Checks if an error occurred in a FileStream.
Returns: whether an error occurred or not
*/
func (ego *FileStream) Check() bool {
	return ego.Ptr.stream.state_flags&C.STREAM_FAILED == 0
}

/*
NewFileStream returns a pointer to created file stream from file.
Parameters:
  - path: path to a file

Returns: pointer to a new instance of FileStream
*/
func NewFileStream(path string) (*FileStream, error) {
	out := FileStream{Path: path}
	out.Ptr = C.stream_file_c_new(C.CString(path))
	if !out.Check() {
		return nil, fmt.Errorf("Unable to create FileStream.")
	}

	return &out, nil
}

/*
Closes a FileStream.
Returns: error if the stream has been already closed, nil otherwise
*/
func (ego *FileStream) Close() error {
	if ego.Ptr == nil {
		return fmt.Errorf("FileStream has been already closed.")
	}
	C.stream_c_destroy(&ego.Ptr.stream)
	C.free(unsafe.Pointer(ego.Ptr))
	ego.Ptr = nil
	return nil

}

/*
Buffer over heap memory.
*/
type BufferStream struct {
	Buffer []byte
	Ptr    *C.struct_stream_buffer_c
}

/*
Creates a new BufferStream.
Returns: pointer to a new instance of BufferStream
*/
func NewBufferStream(buffer []byte) (*BufferStream, error) {
	if buffer == nil {
		return nil, fmt.Errorf("Nil buffer given.")
	}
	out := BufferStream{Buffer: buffer}
	out.Ptr = C.stream_buffer_c_new((*C.uchar)(&buffer[0]), C.ulong(len(buffer)))
	if !out.Check() {
		return nil, fmt.Errorf("Unable to create BufferStream.")
	}
	return &out, nil
}

/*
Gets the inner stream structure.
Returns: pointer to the C struct stream_c
*/
func (ego *BufferStream) GetStream() *C.struct_stream_c {
	return &ego.Ptr.stream
}

/*
Checks if an error occurred in a BufferStream.
Returns: whether an error occurred or not
*/
func (ego *BufferStream) Check() bool {
	return ego.Ptr.stream.state_flags&C.STREAM_FAILED == 0
}

/*
Closes a BufferStream.
Returns: error if the stream has been already closed, nil otherwise
*/
func (ego *BufferStream) Close() error {
	if ego.Ptr == nil {
		return fmt.Errorf("BufferStream has been already closed.")
	}
	C.stream_c_destroy(&ego.Ptr.stream)
	C.free(unsafe.Pointer(ego.Ptr))
	ego.Ptr = nil
	return nil
}
