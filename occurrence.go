package gonativeextractor

/*
#cgo CFLAGS: -I/usr/include/nativeextractor
#cgo LDFLAGS: -lnativeextractor -lglib-2.0 -ldl
#include <string.h>
#include <nativeextractor/common.h>
#include <nativeextractor/extractor.h>
*/
import "C"
import "unsafe"

/*
Interface representing one found occurrence.
*/
type Occurrencer interface {
	Str() string
	Pos() uint64
	Upos() uint64
	Len() uint32
	Ulen() uint32
	Label() string
	Prob() float64
}

/*
Struct implementing Occurrencer, contains a pointer to C struct.
*/
type Occurrence struct {
	ptr *C.struct_occurrence_t
}

/*
Creates a string containing found occurrence.
Returns:
  - found occurrence.
*/
func (ego *Occurrence) Str() string {
	cstr := C.strndup(ego.ptr.str, C.size_t(ego.ptr.len))
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
	return uint64(ego.ptr.pos)
}

/*
Casts UTF position of the found occurrence to Go integer type.
Returns:
  - position of the found occurrence (in unicode characters).
*/
func (ego *Occurrence) Upos() uint64 {
	return uint64(ego.ptr.upos)
}

/*
Casts length of the found occurrence to Go integer type.
Returns:
  - length of the found occurrence (in bytes).
*/
func (ego *Occurrence) Len() uint32 {
	return uint32(ego.ptr.len)
}

/*
Casts UTF length of the found occurrence to Go integer type.
Returns:
  - length of the found occurrence (in unicode characters).
*/
func (ego *Occurrence) Ulen() uint32 {
	return uint32(ego.ptr.ulen)
}

/*
Casts label of the found occurrence to Go string.
Returns:
  - label of the found entity.
*/
func (ego *Occurrence) Label() string {
	return C.GoString(ego.ptr.label)
}

/*
Casts probability of the found occurrence to Go float type.
Returns:
  - probability of the occurrence.
*/
func (ego *Occurrence) Prob() float64 {
	return float64(ego.ptr.prob)
}
