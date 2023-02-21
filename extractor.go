package gonativeextractor

/*
#cgo CFLAGS: -I/usr/include/nativeextractor
#cgo LDFLAGS: -lnativeextractor -lglib-2.0 -ldl
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
	"runtime"
	"unsafe"
)

type Occurrence struct {
	Str   string
	Pos   uint64
	Upos  uint64
	Len   uint32
	Ulen  uint32
	Label string
	Prob  float64
}

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
- E_NO_ENCLOSED_OCCURRENCES disables enclosed occurrence feature.
- E_SORT_RESULTS enables results ascending sorting of occurrence records.
*/
const (
	E_NO_ENCLOSED_OCCURRENCES = 1 << 0
	E_SORT_RESULTS            = 1 << 1
)

/* Default path to .so libs representing miners. */
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
Creates a new Extractor.
Parameters:
  - batch: number of logical symbols to be analyzed in the stream
  - threads: number of threads for miners to run on
  - flags: initial flags

Returns: pointer to a new instance of Extractor
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
Destroys the Extractor
Returns: error if the extractor has been already closed, nil otherwise
*/
func (ego *Extractor) Close() error {
	if ego.extractor == nil {
		return fmt.Errorf("Extractor has been already closed.")
	}
	C.extractor_c_destroy(ego.extractor)
	ego.extractor = nil
	return nil
}

/*
Sets stream to the Extractor.
Parameters:
  - stream an instance of Streamer interface.

Returns:
  - error.
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
Unsets stream attached to the Extractor.
*/
func (ego *Extractor) UnsetStream() {
	C.extractor_c_unset_stream(ego.extractor)
	ego.stream = nil
}

/*
Set NativeExtractor flags.
Parameters:
  - flags:
  - valid flags are these E_NO_ENCLOSED_OCCURRENCES, E_SORT_RESULTS.

Returns:
  - error.
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
Unsets flags.
Parameters:
  - flags:
  - valid flags are these E_NO_ENCLOSED_OCCURRENCES, E_SORT_RESULTS.

Returns:
  - error.
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
  - sodir a path to the shared object,
  - symbol shared object symbol,
  - params (Optional) may be empty array or nil.
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

func (ego *Extractor) GetLastError() error {
	err := C.GoString(C.extractor_get_last_error(ego.extractor))
	if err == "" {
		return nil
	}
	return fmt.Errorf(err)
}

func (ego *Extractor) Eof() bool {
	if ego.stream == nil {
		return true
	}
	return ego.extractor.stream.state_flags&C.STREAM_EOF != 0
}

func (ego *Extractor) Meta() []DlSymbol {

	dlSymbols := ego.extractor.dlsymbols
	result := make([]DlSymbol, 0)

	for i := 0; true; i++ {
		elem := *(**C.dl_symbol_t)(unsafe.Add(unsafe.Pointer(dlSymbols), i*int(unsafe.Sizeof(dlSymbols))))
		if elem == nil {
			break
		}
		meta := make([]string, 0)
		for j := 0; true; j++ {
			str := *(**C.char)(unsafe.Add(unsafe.Pointer(elem.meta), j*int(unsafe.Sizeof(elem.meta))))
			if str == nil {
				break
			}
			meta = append(meta, C.GoString(str))
		}
		result = append(result, DlSymbol{
			Ldpath: C.GoString(elem.ldpath),
			Ldsymb: C.GoString(elem.ldsymb),
			Meta:   meta,
			Params: C.GoString(elem.params),
			Ldptr:  unsafe.Pointer(elem.ldptr),
		})
	}

	return result

}

func (ego *Extractor) Next() ([]Occurrence, error) {

	if ego.stream == nil {
		return nil, fmt.Errorf("Stream is not set.")
	}

	occurrences := C.next(ego.extractor, C.uint(ego.batch))
	result := make([]Occurrence, 0)

	for i := 0; true; i++ {
		occ := *(**C.struct_occurrence_t)(unsafe.Add(unsafe.Pointer(occurrences), i*int(unsafe.Sizeof(occurrences))))
		if occ == nil {
			break
		}
		result = append(result, Occurrence{
			Str:   C.GoString(occ.str),
			Pos:   uint64(occ.pos),
			Upos:  uint64(occ.upos),
			Len:   uint32(occ.len),
			Ulen:  uint32(occ.ulen),
			Label: C.GoString(occ.label),
			Prob:  float64(occ.prob),
		})
	}

	return result, nil

}
