package gonativeextractor

/*
#cgo CFLAGS: -I/usr/include/nativeextractor
#cgo LDFLAGS: -lnativeextractor -lglib-2.0 -ldl
#include <string.h>
#include <nativeextractor/common.h>
#include <nativeextractor/extractor.h>
*/
import "C"
import (
	"math/bits"
	"unsafe"
)

/*
Size of pointer in bytes.
*/
const POINTER_SIZE = bits.UintSize / 8

/*
Interface alowing an access to the found occurrence.
*/
type Occurrencer interface {
	Next()
	Eof() bool
	Str() string
	Pos() uint64
	Upos() uint64
	Len() uint32
	Ulen() uint32
	Label() string
	Prob() float64
}

/*
Struct implementing Occurrencer, contains a pointer to the current occurrence (C struct).
*/
type Occurrence struct {
	ptr **C.struct_occurrence_t
}

/*
Moves the pointer to the next occurrence. If EOF, does nothing.
*/
func (ego *Occurrence) Next() {
	if !ego.Eof() {
		ego.ptr = (**C.struct_occurrence_t)(unsafe.Add(unsafe.Pointer(ego.ptr), POINTER_SIZE))
	}
}

/*
Checks if all occurrences have been read.
Returns:
  - true if there is nothing to read (current pointer is nil), false otherwise.
*/
func (ego *Occurrence) Eof() bool {
	return *ego.ptr == nil
}

/*
Private method, causes panic when EOF.
*/
func (ego *Occurrence) check() {
	if ego.Eof() {
		panic("Attempt to access a nil pointer.")
	}
}

/*
Creates a string containing found occurrence.
Returns:
  - found occurrence.
*/
func (ego *Occurrence) Str() string {
	ego.check()
	cstr := C.strndup((*ego.ptr).str, C.size_t((*ego.ptr).len))
	retVal := C.GoString(cstr)
	C.free(unsafe.Pointer(cstr))
	return retVal
}

/*
Casts position of the found occurrence to Go integer type.
Returns:
  - position of the found occurrence (in bytes).
*/
func (ego *Occurrence) Pos() uint64 {
	ego.check()
	return uint64((*ego.ptr).pos)
}

/*
Casts UTF position of the found occurrence to Go integer type.
Returns:
  - position of the found occurrence (in unicode characters).
*/
func (ego *Occurrence) Upos() uint64 {
	ego.check()
	return uint64((*ego.ptr).upos)
}

/*
Casts length of the found occurrence to Go integer type.
Returns:
  - length of the found occurrence (in bytes).
*/
func (ego *Occurrence) Len() uint32 {
	ego.check()
	return uint32((*ego.ptr).len)
}

/*
Casts UTF length of the found occurrence to Go integer type.
Returns:
  - length of the found occurrence (in unicode characters).
*/
func (ego *Occurrence) Ulen() uint32 {
	ego.check()
	return uint32((*ego.ptr).ulen)
}

/*
Casts label of the found occurrence to Go string.
Returns:
  - label of the found entity.
*/
func (ego *Occurrence) Label() string {
	ego.check()
	return C.GoString((*ego.ptr).label)
}

/*
Casts probability of the found occurrence to Go float type.
Returns:
  - probability of the occurrence.
*/
func (ego *Occurrence) Prob() float64 {
	ego.check()
	return float64((*ego.ptr).prob)
}
