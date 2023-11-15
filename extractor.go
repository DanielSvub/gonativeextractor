package gonativeextractor

/*
   #cgo LDFLAGS: -lglib-2.0 -ldl
   #include <dlfcn.h>
   #include <string.h>
   //#include <nativeextractor/common.h>
   //#include <nativeextractor/extractor.h>
   //#include <nativeextractor/stream.h>

   void cleanup_fun_bridge(void * f, extractor_c * self)
   {
      return ((void (*)(extractor_c *))f)(self);
   }
   bool extractor_c_set_stream_bridge(void * f, extractor_c * self, stream_c * stream)
   {
      return ((bool (*)(extractor_c *, stream_c *))f)(self, stream);
   }
   bool flags_fun_bridge(void * f, extractor_c * self, unsigned flags)
   {
      return ((bool (*)(extractor_c *, unsigned))f)(self, flags);
   }
   occurrence_t** next_bridge(void * f, extractor_c * self, unsigned batch)
   {
      return ((occurrence_t** (*)(extractor_c *, unsigned))f)(self, batch);
   }
   bool extractor_c_add_miner_from_so_bridge(void *f, extractor_c * self, const char * miner_so_path, const char * miner_name, void * params)
   {
     return ((bool (*)(extractor_c *, const char *, const char *, void* ))f)(self, miner_so_path, miner_name, params);
   }
   const char * extractor_get_last_error_bridge(void * f, extractor_c * self)
   {
     return ((const char * (*)(extractor_c *))f)(self);
   }
   extractor_c * extractor_c_new_bridge(void * f, int threads, miner_c ** miners)
   {
     return ((extractor_c * (*)(int, miner_c **))f)(threads, miners);
   }





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
Default path to libnativeextractor.so.
*/
const DEFAULT_NATIVEEXTRACOTR_PATH = "/usr/lib"

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
	dlHandler unsafe.Pointer
}

/*
Creates a new Extractor. Has to be deallocated with 'Destroy' method after use.

Parameters:
  - batch - number of logical symbols to be analyzed in the stream (if negative, defaults to 2^16),
  - threads - number of threads for miners to run on (if negative, defaults to maximum threads available),
  - flags - initial flags.

Returns:
  - pointer to a new instance of Extractor.
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
	nativeextractorpath := C.CString(DEFAULT_NATIVEEXTRACOTR_PATH + "/libnativeextractor.so")
	defer C.free(unsafe.Pointer(nativeextractorpath))
	fmt.Println("Dlopening nativeextractor")
	out.dlHandler = C.dlopen(nativeextractorpath, C.RTLD_LAZY)
	if out.dlHandler == nil {
		panic("Can not dlopen libnativeextractor.so")
	}
	fName := C.CString("extractor_c_new")
	defer C.free(unsafe.Pointer(fName))
	fPtr := C.dlsym(out.dlHandler, fName)
	out.extractor = C.extractor_c_new_bridge(fPtr, C.int(out.threads), miners) // C.extractor_c_new(C.int(out.threads), miners)

	return &out
}

/*
Destroys the Extractor.

Returns:
  - error if the extractor has been already closed, nil otherwise.
*/
func (ego *Extractor) Destroy() error {
	if ego.extractor == nil {
		return fmt.Errorf("Extractor has been already closed")
	}

	if ego.stream != nil {
		err := ego.stream.Close()
		ego.stream = nil
		if err != nil {
			return err
		}
	}

	fName := C.CString("extractor_c_destroy")
	defer C.free(unsafe.Pointer(fName))
	fPtr := C.dlsym(ego.dlHandler, fName)
	C.cleanup_fun_bridge(fPtr, ego.extractor) //C.extractor_c_destroy(ego.extractor

	C.free(unsafe.Pointer(ego.extractor))
	ego.extractor = nil

	C.dlclose(ego.dlHandler)
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
	fName := C.CString("extractor_c_set_stream")
	defer C.free(unsafe.Pointer(fName))
	fPtr := C.dlsym(ego.dlHandler, fName)
	ok := C.extractor_c_set_stream_bridge(fPtr, ego.extractor, stream.GetStream()) //C.extractor_c_set_stream(ego.extractor, stream.GetStream())
	if !ok {
		return fmt.Errorf("unable to set stream")
	}
	ego.stream = stream
	return nil
}

/*
Dettaches the stream from the Extractor.
*/
func (ego *Extractor) UnsetStream() {
	fName := C.CString("extractor_c_unset_stream")
	defer C.free(unsafe.Pointer(fName))
	fPtr := C.dlsym(ego.dlHandler, fName)
	C.cleanup_fun_bridge(fPtr, ego.extractor) //C.extractor_c_unset_stream(ego.extractor)
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
	fName := C.CString("extractor_set_flags")
	defer C.free(unsafe.Pointer(fName))
	fPtr := C.dlsym(ego.dlHandler, fName)
	ok := C.flags_fun_bridge(fPtr, ego.extractor, C.uint(flags)) //C.extractor_set_flags(ego.extractor, C.uint(flags))
	if !ok {
		return fmt.Errorf("unable to set flags")
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
	fName := C.CString("extractor_unset_flags")
	defer C.free(unsafe.Pointer(fName))
	fPtr := C.dlsym(ego.dlHandler, fName)
	ok := C.flags_fun_bridge(fPtr, ego.extractor, C.uint(flags)) // C.extractor_unset_flags(ego.extractor, C.uint(flags))
	if !ok {
		return fmt.Errorf("unable to unset flags")
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

Returns:
  - error if any occurred, nil otherwise.
*/
func (ego *Extractor) AddMinerSo(sodir string, symbol string, params []byte) error {
	var data unsafe.Pointer
	if params == nil {
		data = nil
	} else {
		data = unsafe.Pointer(&params[0])
	}
	//fmt.Println("Loading miner: ", sodir, symbol)
	fName := C.CString("extractor_c_add_miner_from_so")
	defer C.free(unsafe.Pointer(fName))
	fPtr := C.dlsym(ego.dlHandler, fName)
	extractorAdded := C.extractor_c_add_miner_from_so_bridge(fPtr, ego.extractor, C.CString(sodir), C.CString(symbol), data)

	if extractorAdded { //C.extractor_c_add_miner_from_so(ego.extractor, C.CString(sodir), C.CString(symbol), data) {
		//	fmt.Println("OK")
		return nil
	}
	return ego.GetLastError()
}

/*
Gives the last error which occurred in Extractor.

Returns:
  - error, nil if no error occurred.
*/
func (ego *Extractor) GetLastError() error {
	fName := C.CString("extractor_get_last_error")
	defer C.free(unsafe.Pointer(fName))
	fPtr := C.dlsym(ego.dlHandler, fName)

	err := C.GoString(C.extractor_get_last_error_bridge(fPtr, ego.extractor))
	//	err := C.GoString(C.extractor_get_last_error(ego.extractor))
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
		return nil, fmt.Errorf("stream is not set")
	}
	fName := C.CString("next")
	defer C.free(unsafe.Pointer(fName))
	fPtr := C.dlsym(ego.dlHandler, fName)

	result := C.next_bridge(fPtr, ego.extractor, C.uint(ego.batch)) //C.next(ego.extractor, C.uint(ego.batch))

	return &Occurrence{result}, nil

}
