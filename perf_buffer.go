package libbpfgo

import "C"
import (
	"fmt"
	"unsafe"
)

//go:linkname main_perfOutput main.perfOutput
func main_perfOutput(cpu int, data []byte)

//export PerfBufferFunc
func PerfBufferFunc(ctx unsafe.Pointer, cpu C.int, data unsafe.Pointer, size C.int) {
	value := C.GoBytes(data, size)
	main_perfOutput(int(cpu), value)
	//fmt.Printf("cpu %d output: %d\n", int(cpu), binary.LittleEndian.Uint32(value))
	return
}

//export PerfBufferLostFunc
func PerfBufferLostFunc(ctx unsafe.Pointer, cpu C.int, cnt C.ulonglong) {
	fmt.Printf("cpu %d err: %d", int(cpu), int(cnt))
	return
}
