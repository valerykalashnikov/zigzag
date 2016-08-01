package pearson

// #include "pearson.h"
import "C"

import (
	"strconv"
	"unsafe"
)

func CreatePearson16(key []byte, length, hexlength int) uint64 {
	var (
		hex string
	)

	cKey := C.malloc(C.size_t(unsafe.Sizeof(key)))
	defer C.free(cKey)

	phex := C.CString(hex)
	defer C.free(unsafe.Pointer(phex))

	C.Pearson16((*C.uchar)(cKey), (C.size_t)(length), phex, (C.size_t)(hexlength))

	hex = C.GoString(phex)
	result, _ := strconv.ParseUint(hex, 16, 64)

	return result
}
