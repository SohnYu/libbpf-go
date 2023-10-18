package libbpfgo

/*
#include <bpf/libbpf.h>
#include <bpf/bpf.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"syscall"
	"unsafe"
)

func (m *BPFMap) PinMap(pinPath string) error {
	cPinPath := C.CString(pinPath)
	defer C.free(unsafe.Pointer(cPinPath))
	err := C.bpf_map__pin(m.bpfMap, cPinPath)
	if err < 0 {
		err := fmt.Errorf("failed to pin map to path %s %w", pinPath, syscall.Errno(-err))
		panic(err.Error())
	}
	return nil
}

func (m *BPFMap) UnPinMap(unpinPath string) error {
	cPinPath := C.CString(unpinPath)
	defer C.free(unsafe.Pointer(cPinPath))
	err := C.bpf_map__unpin(m.bpfMap, cPinPath)
	if err < 0 {
		err := fmt.Errorf("failed to unpin map to path %s %w", unpinPath, syscall.Errno(-err))
		panic(err.Error())
	}
	return nil
}

func mapKeySize(cm *C.struct_bpf_map) int {
	return int(C.bpf_map__key_size(cm))
}

func mapValueSize(cm *C.struct_bpf_map) int {
	return int(C.bpf_map__value_size(cm))
}

func mapType(cm *C.struct_bpf_map) int {
	bpfMapType := C.bpf_map__type(cm)
	return int(bpfMapType)
}

func (m *BPFPinnedMap) MapLookupElem(key unsafe.Pointer) []byte {
	value := make([]byte, m.valueSize)
	valuePtr := unsafe.Pointer(&value[0])

	errC := C.bpf_map_lookup_elem_flags(m.fd, key, valuePtr, C.ulonglong(MapFlagUpdateAny))
	if errC != 0 {
		err := fmt.Errorf("failed to lookup value %v in map %s: %w", key, m.name, syscall.Errno(-errC))
		panic(err.Error())
	}
	return value
}

func (m *BPFPinnedMap) MapDelElem(key unsafe.Pointer) error {
	ret := C.bpf_map_delete_elem(m.fd, key)
	if ret != 0 {
		err := fmt.Errorf("failed to del map elem %d from map %s: %w", key, m.name, syscall.Errno(-ret))
		panic(err.Error())
	}
	return nil
}

func (m *BPFPinnedMap) MapUpdateElem(key, value unsafe.Pointer) {
	errC := C.bpf_map_update_elem(m.fd, key, value, C.ulonglong(MapFlagUpdateAny))
	if errC != 0 {
		err := fmt.Errorf("failed to update map : %w", syscall.Errno(-errC))
		panic(err.Error())
	}
	return
}
