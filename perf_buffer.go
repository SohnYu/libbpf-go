package main

import "C"
import (
	"fmt"
	"unsafe"
)

//go:linkname main_perfOutput main.perfOutput
func main_perfOutput(mapName string, cpu int, data []byte)

//export PerfBufferFunc
func PerfBufferFunc(ctx unsafe.Pointer, cpu C.int, data unsafe.Pointer, size C.int) {
	// 将 unsafe.Pointer 转换回 *string
	strPtr := (*string)(ctx)
	value := C.GoBytes(data, size)
	main_perfOutput(*strPtr, int(cpu), value)
	//fmt.Printf("cpu %d output: %d\n", int(cpu), binary.LittleEndian.Uint32(value))
	return
}

//export PerfBufferLostFunc
func PerfBufferLostFunc(ctx unsafe.Pointer, cpu C.int, cnt C.ulonglong) {
	fmt.Printf("cpu %d err: %d", int(cpu), int(cnt))
	return
}
