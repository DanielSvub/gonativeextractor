package gonativeextractor

/*
   #cgo LDFLAGS: -lglib-2.0 -ldl
   #include <dlfcn.h>
   #include <nativeextractor/common.h>
   #include <nativeextractor/stream.h>

   void stream_c_destroy_bridge(void * f, stream_c * self)
   {
      return ((void (*)(stream_c *))f)(self);
   }

   stream_buffer_c * stream_buffer_c_new_bridge(void * f, unsigned char * buffer, unsigned long length)
   {
      return ((stream_buffer_c * (*) (unsigned char *, unsigned long))f)(buffer, length);
   }


   stream_file_c * stream_file_c_new_bridge(void * f, const char* path)
   {
      return ((stream_file_c * (*) (const char*))f)( path);
   }


*/
import "C"
import (
	"fmt"
	"io"
	"os"
	"unsafe"
)

/*
Interface for streams.
*/
type Streamer interface {
	GetStream() *C.struct_stream_c
	Check() bool
	io.Closer
}

/*
Structure representing stream from file.
*/
type FileStream struct {
	Ptr       *C.struct_stream_file_c
	Path      string
	dlHandler unsafe.Pointer
}

/*
Gets the inner stream structure.

Returns:
  - pointer to the C struct stream_c.
*/
func (ego *FileStream) GetStream() *C.struct_stream_c {
	return &ego.Ptr.stream
}

/*
Checks if an error occurred in a FileStream.

Returns:
  - true if an error occurred, false otherwise.
*/
func (ego *FileStream) Check() bool {
	return ego.Ptr.stream.state_flags&C.STREAM_FAILED == 0
}

/*
Creates a new FileStream.

Parameters:
  - path - path to a file.

Returns:
  - pointer to a new instance of FileStream.
  - error if any occurred, nil otherwise.
*/
func NewFileStream(path string) (*FileStream, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist")
	}
	out := FileStream{Path: path}

	nativeextractorpath := C.CString(DEFAULT_NATIVEEXTRACOTR_PATH + "/libnativeextractor.so")
	defer C.free(unsafe.Pointer(nativeextractorpath))
	out.dlHandler = C.dlopen(nativeextractorpath, C.RTLD_NODELETE|C.RTLD_LAZY)
	if out.dlHandler == nil {
		panic("Can not dlopen libnativeextractor.so")
	}

	fName := C.CString("stream_file_c_new")
	defer C.free(unsafe.Pointer(fName))
	fPtr := C.dlsym(out.dlHandler, fName)
	filePath := C.CString(path)
	defer C.free(unsafe.Pointer(filePath))
	out.Ptr = C.stream_file_c_new_bridge(fPtr, filePath) //C.stream_file_c_new(C.CString(path))
	if !out.Check() {
		return nil, fmt.Errorf("unable to create FileStream")
	}

	return &out, nil
}

/*
Closes a FileStream.

Returns:
  - error if the stream has been already closed, nil otherwise.
*/
func (ego *FileStream) Close() error {
	if ego.Ptr == nil {
		return fmt.Errorf("fileStream has been already closed")
	}

	fName := C.CString("stream_c_destroy")
	defer C.free(unsafe.Pointer(fName))
	fPtr := C.dlsym(ego.dlHandler, fName)
	C.stream_c_destroy_bridge(fPtr, &ego.Ptr.stream) //C.stream_c_destroy(&ego.Ptr.stream)
	C.free(unsafe.Pointer(ego.Ptr))
	ego.Ptr = nil
	C.dlclose(ego.dlHandler)
	return nil

}

/*
Structure representing buffer over heap memory.
*/
type BufferStream struct {
	Buffer    []byte
	Ptr       *C.struct_stream_buffer_c
	dlHandler unsafe.Pointer
}

/*
Creates a new BufferStream.

Parameters:
  - buffer - byte array for stream initialization (has to be terminated with \x00).

Returns:
  - pointer to a new instance of BufferStream.
  - error if any occurred, nil otherwise.
*/
func NewBufferStream(buffer []byte) (*BufferStream, error) {
	if buffer == nil {
		return nil, fmt.Errorf("nil buffer given")
	}
	out := BufferStream{Buffer: buffer}

	nativeextractorpath := C.CString(DEFAULT_NATIVEEXTRACOTR_PATH + "/libnativeextractor.so")
	defer C.free(unsafe.Pointer(nativeextractorpath))
	out.dlHandler = C.dlopen(nativeextractorpath, C.RTLD_NODELETE|C.RTLD_LAZY)
	if out.dlHandler == nil {
		panic("Can not dlopen libnativeextractor.so")
	}

	fName := C.CString("stream_buffer_c_new")
	defer C.free(unsafe.Pointer(fName))
	fPtr := C.dlsym(out.dlHandler, fName)
	out.Ptr = C.stream_buffer_c_new_bridge(fPtr, (*C.uchar)(&buffer[0]), C.ulong(len(buffer))) //C.stream_buffer_c_new((*C.uchar)(&buffer[0]), C.ulong(len(buffer)))
	if !out.Check() {
		return nil, fmt.Errorf("unable to create BufferStream")
	}

	return &out, nil
}

/*
Gets the inner stream structure.

Returns:
  - pointer to the C struct stream_c.
*/
func (ego *BufferStream) GetStream() *C.struct_stream_c {
	return &ego.Ptr.stream
}

/*
Checks if an error occurred in a BufferStream.

Returns:
  - true if an error occurred, false otherwise.
*/
func (ego *BufferStream) Check() bool {
	return ego.Ptr.stream.state_flags&C.STREAM_FAILED == 0
}

/*
Closes a BufferStream.

Returns:
  - error if the stream has been already closed, nil otherwise.
*/
func (ego *BufferStream) Close() error {
	if ego.Ptr == nil {
		return fmt.Errorf("bufferStream has been already closed")
	}
	fName := C.CString("stream_c_destroy")
	defer C.free(unsafe.Pointer(fName))
	fPtr := C.dlsym(ego.dlHandler, fName)
	C.stream_c_destroy_bridge(fPtr, &ego.Ptr.stream)
	//	C.stream_c_destroy(&ego.Ptr.stream)
	C.free(unsafe.Pointer(ego.Ptr))
	ego.Ptr = nil
	C.dlclose(ego.dlHandler)
	return nil
}
