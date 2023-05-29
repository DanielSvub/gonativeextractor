package gonativeextractor

/*
#cgo CFLAGS: -I/usr/include/nativeextractor
#cgo LDFLAGS: -lnativeextractor -lglib-2.0 -ldl
#include <string.h>
#include <nativeextractor/common.h>
#include <nativeextractor/extractor.h>
#include <nativeextractor/stream.h>
bool extractor_c_add_miner_from_so(extractor_c * self, const char * miner_so_path, const char * miner_name, void * params );
const char * extractor_get_last_error(extractor_c * self);
void extractor_c_destroy(extractor_c * self);
bool extractor_c_set_stream(extractor_c * self, stream_c * stream);
void extractor_c_unset_stream(extractor_c * self);
bool extractor_set_flags(extractor_c * self, unsigned flags);
bool extractor_unset_flags(extractor_c * self, unsigned flags);
occurrence_t** next(extractor_c * self, unsigned batch);
const char * extractor_get_last_error(extractor_c * self);
*/
import "C"
import (
	"fmt"
	"math/bits"
	"runtime"
	"unsafe"
)

/*
Structure with information about a miner.
*/
type DlSymbol struct {
	// Path to the .so library
	Ldpath string
	// Miner function name
	Ldsymb string
	// Meta info about miner functions and labels
	Meta []string
	// Miner params
	Params string
	// Pointer to the loaded .so library
	Ldptr unsafe.Pointer
}

/*
Constants for valid NativeExtractor flags.
*/
const (
	// Disables enclosed occurrence feature.
	E_NO_ENCLOSED_OCCURRENCES = 1 << 0
	// Enables results ascending sorting of occurrence records.
	E_SORT_RESULTS = 1 << 1
)

/*
Default path to .so libs representing miners.
*/
const DEFAULT_MINERS_PATH = "/usr/lib/nativeextractor_miners"

/*
Analyzes next batch with miners.
*/
type Extractor struct {
	miners    []C.struct_miner_c
	batch     int
	stream    Streamer
	flags     uint32
	threads   int
	extractor *C.struct_extractor_c
}

/*
Creates a new Extractor. Has to be deallocated with 'Destroy' method after use.
Parameters:
  - batch - number of logical symbols to be analyzed in the stream (if negative, defaults to 2^16),
  - threads - number of threads for miners to run on (if negative, defaults to maximum threads available),
  - flags - initial flags.
Returns: pointer to a new instance of Extractor.
*/
func NewExtractor(batch int, threads int, flags uint32) *Extractor {
	out := Extractor{}
	out.flags = flags

	out.threads = threads
	if out.threads <= 0 {
		out.threads = runtime.NumCPU()
	}

	out.batch = batch
	if out.batch <= 0 {
		out.batch = 1 << 16
	}

	miner := &C.struct_miner_c{}
	miners := (**C.struct_miner_c)(C.calloc(1, C.ulong(unsafe.Sizeof(&miner))))
	out.extractor = C.extractor_c_new(C.int(out.threads), miners)

	return &out
}

/*
Destroys the Extractor.
Returns: error if the extractor has been already closed, nil otherwise.
*/
func (ego *Extractor) Destroy() error {
	if ego.extractor == nil {
		return fmt.Errorf("Extractor has been already closed.")
	}

	if ego.stream != nil {
		err := ego.stream.Close()
		ego.stream = nil
		if err != nil {
			return err
		}
	}

	C.extractor_c_destroy(ego.extractor)
	C.free(unsafe.Pointer(ego.extractor))
	ego.extractor = nil

	return nil
}

/*
Sets a stream to the Extractor.
Parameters:
  - stream - an instance of Streamer interface.
Returns:
  - error if any occurred, nil otherwise.
*/
func (ego *Extractor) SetStream(stream Streamer) error {
	ok := C.extractor_c_set_stream(ego.extractor, stream.GetStream())
	if !ok {
		return fmt.Errorf("Unable to set stream.")
	}
	ego.stream = stream
	return nil
}

/*
Dettaches the stream from the Extractor.
*/
func (ego *Extractor) UnsetStream() {
	C.extractor_c_unset_stream(ego.extractor)
	ego.stream = nil
}

/*
Sets NativeExtractor flags.
Parameters:
  - flags - use constants defined above.
Returns:
  - error if any occurred, nil otherwise.
*/
func (ego *Extractor) SetFlags(flags uint32) error {
	ok := C.extractor_set_flags(ego.extractor, C.uint(flags))
	if !ok {
		return fmt.Errorf("Unable to set flags.")
	}
	ego.flags = flags
	return nil
}

/*
Unsets NativeExtractor flags.
Parameters:
  - flags - use constants defined above.
Returns:
  - error if any occurred, nil otherwise.
*/
func (ego *Extractor) UnsetFlags(flags uint32) error {
	ok := C.extractor_unset_flags(ego.extractor, C.uint(flags))
	if !ok {
		return fmt.Errorf("Unable to unset flags.")
	}
	ego.flags = uint32(ego.extractor.flags)
	return nil
}

/*
Loads a Miner from a Shared Object (.so library).
Parameters:
  - sodir - a path to the shared object,
  - symbol - shared object symbol,
  - params - optional (may be empty array or nil, but if present, has to be terminated with \x00).
*/
func (ego *Extractor) AddMinerSo(sodir string, symbol string, params []byte) error {
	var data unsafe.Pointer
	if params == nil {
		data = nil
	} else {
		data = unsafe.Pointer(&params[0])
	}

	if C.extractor_c_add_miner_from_so(ego.extractor, C.CString(sodir), C.CString(symbol), data) {
		return nil
	}

	return fmt.Errorf(C.GoString(C.extractor_get_last_error(ego.extractor)))
}

/*
Gives the last error which occurred in Extractor.
Returns:
  - error, nil if no error occurred.
*/
func (ego *Extractor) GetLastError() error {
	err := C.GoString(C.extractor_get_last_error(ego.extractor))
	if err == "" {
		return nil
	}
	return fmt.Errorf(err)
}

/*
Checks if the stream attached to the Extractor ended.
Returns:
  - true if there is nothing to read (stream ended or no stream set), false otherwise.
*/
func (ego *Extractor) Eof() bool {
	if ego.stream == nil {
		return true
	}
	return ego.extractor.stream.state_flags&C.STREAM_EOF != 0
}

/*
Gives the meta information about Extractor.
Returns:
  - slice of structures with information about miners.
*/
func (ego *Extractor) Meta() []DlSymbol {

	step := bits.UintSize / 8
	result := make([]DlSymbol, 0)

	// Iterating over null terminated array of pointers (size of pointer added to address in each iteration)
	for elem := ego.extractor.dlsymbols; *elem != nil; elem = (**C.struct_dl_symbol_t)(unsafe.Add(unsafe.Pointer(elem), step)) {

		meta := make([]string, 0)

		// Iterating once more, same as the outer loop
		for str := (*elem).meta; *str != nil; str = (**C.char)(unsafe.Add(unsafe.Pointer(&str), step)) {
			meta = append(meta, C.GoString(*str))
		}

		result = append(result, DlSymbol{
			Ldpath: C.GoString((*elem).ldpath),
			Ldsymb: C.GoString((*elem).ldsymb),
			Meta:   meta,
			Params: C.GoString((*elem).params),
			Ldptr:  unsafe.Pointer((*elem).ldptr),
		})

	}

	return result

}

/*
Gives the next batch of found entities.
Returns:
  - the first found occurrence,
  - error, if any occurred.
*/
func (ego *Extractor) Next() (Occurrencer, error) {

	if ego.stream == nil {
		return nil, fmt.Errorf("Stream is not set.")
	}

	result := C.next(ego.extractor, C.uint(ego.batch))

	return &Occurrence{result}, nil

}
