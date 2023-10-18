package libbpfgo

import "C"
import (
	"unsafe"
)

//go:linkname main_ringBufferOutput main.ringBufferOutput
func main_ringBufferOutput(data []byte)

//export RingBufferFunc
func RingBufferFunc(ctx unsafe.Pointer, data unsafe.Pointer, size C.int) C.int {
	value := C.GoBytes(data, size)
	main_ringBufferOutput(value)
	return C.int(0)
}
