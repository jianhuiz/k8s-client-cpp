package main

import (
	"C"
	"unsafe"
)

func no_copy_slice_from_c_array(bytes unsafe.Pointer, size C.int) []byte {
        return (*[1 << 30]byte)(bytes)[:size:size]
}

