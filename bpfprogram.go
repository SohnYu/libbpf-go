package main

/*
#include <stdlib.h>
#include "libbpfgo.h"
*/
import "C"
import (
	"debug/elf"
	"fmt"
	"path"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type BaseProgram struct {
	name       string
	program    *C.struct_bpf_program
	module     *Module
	pinnedPath string
}

type Program interface {
	LoadIntoKernel() error
	Unload() error
}

type TracePointProgram struct {
	BaseProgram
	category string
	funcName string
}

func (t *TracePointProgram) LoadIntoKernel() error {
	tpCategory := C.CString(t.category)
	defer C.free(unsafe.Pointer(tpCategory))
	tpName := C.CString(t.funcName)
	defer C.free(unsafe.Pointer(tpName))

	link, errno := C.bpf_program__attach_tracepoint(t.program, tpCategory, tpName)
	if link == nil {
		return fmt.Errorf("failed to attach tracepoint %s to program %s: %w", t.name, t.name, errno)
	}
	return nil
}

func (t *TracePointProgram) Unload() error {
	return nil
}

type KprobeProgram struct {
	BaseProgram
}

func (t *KprobeProgram) LoadIntoKernel() error {
	cs := C.CString(t.name)
	cbool := C.bool(false)
	link, errno := C.bpf_program__attach_kprobe(t.program, cbool, cs)
	C.free(unsafe.Pointer(cs))
	if link == nil {
		return fmt.Errorf("failed to attach %s k(ret)probe to program %s: %w", t.name, t.name, errno)
	}
	return nil
}

func (t *KprobeProgram) Unload() error {
	return nil
}

type UProbeProgram struct {
	BaseProgram
	programPath string
	offset      uint32
	resolved    bool
}

func (t *UProbeProgram) LoadIntoKernel() error {
	if !t.resolved {
		panic("you need add offset by yourself!")
	}
	retCBool := C.bool(false)
	pidCint := C.int(-1)
	pathCString := C.CString(t.module.elfPath)
	offsetCsizet := C.size_t(t.offset)
	link, errno := C.bpf_program__attach_uprobe(t.program, retCBool, pidCint, pathCString, offsetCsizet)
	if link == nil {
		return fmt.Errorf("failed to attach u(ret)probe to program %s:%d with pid %d: %w ", t.module.elfPath, t.offset, 0, errno)
	}
	C.free(unsafe.Pointer(pathCString))
	return nil
}

func (t *UProbeProgram) Unload() error {
	return nil
}

func (t *UProbeProgram) AddOffset(offset string) error {
	val, err := strconv.ParseUint(offset, 0, 32)
	if err != nil {
		return err
	}
	t.offset = uint32(val)
	t.resolved = true
	return nil
}

func (t *UProbeProgram) AddOffsetByFuncName(funcName string) error {
	t.offset = symbolOffset(t.module.elf, funcName)
	t.resolved = true
	return nil
}

type TCProgram struct {
	BaseProgram
}

func (t *TCProgram) LoadIntoKernel() error {
	fmt.Println(t.name, "load suceess")
	return nil
}

func (t *TCProgram) Unload() error {
	return nil
}

type XDPProgram struct {
	BaseProgram
}

func (t *XDPProgram) LoadIntoKernel() error {
	fmt.Println(t.name, "load suceess")
	return nil
}

func (t *XDPProgram) Unload() error {
	return nil
}

func NewModuleFromFile(bpfFilePath string) *Module {
	opts := C.struct_bpf_object_open_opts{}
	opts.sz = C.sizeof_struct_bpf_object_open_opts

	bpfFile := C.CString(bpfFilePath)
	defer C.free(unsafe.Pointer(bpfFile))

	obj, errno := C.bpf_object__open_file(bpfFile, &opts)
	if obj == nil {
		panic(errno)
	}
	m := &Module{
		obj:          obj,
		elf:          nil,
		loaded:       false,
		deferFunc:    make([]func(), 0),
		bpfProgram:   make(map[string]Program, 0),
		bpfPinnedMap: make(map[string]*BPFPinnedMap, 0),
		maps:         make(map[string]*BPFMap, 0),
	}
	err := m.BPFLoadObject()
	if err != nil {
		panic(err)
	}
	return m
}

func (m *Module) AddDeferFunc(function func()) {
	m.deferFunc = append(m.deferFunc, function)
}

func (m *Module) BPFLoadObject() error {
	ret := C.bpf_object__load(m.obj)
	if ret != 0 {
		return fmt.Errorf("failed to load BPF object: %w", syscall.Errno(-ret))
	}
	m.loaded = true
	m.AddDeferFunc(func() {
		C.bpf_object__close(m.obj)
	})
	return nil
}

func (m *Module) GetAllProgram() map[string]Program {
	return m.bpfProgram
}

func (p BaseProgram) GetProgramSecName() string {
	return C.GoString(C.bpf_program__section_name(p.program))
}

func (m *Module) LoadAllProgram() error {
	var bpfProgram *C.struct_bpf_program
	for {
		bpfProgram = C.bpf_object__next_program(m.obj, bpfProgram)
		if bpfProgram == nil {
			break
		}
		programName := C.bpf_program__name(bpfProgram)
		generationProgram := m.GenerationProgram(bpfProgram)
		m.bpfProgram[C.GoString(programName)] = generationProgram
	}
	return nil
}

func (m *Module) LoadAllIntoKernel() error {
	for _, program := range m.bpfProgram {
		err := program.LoadIntoKernel()
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func (m *Module) GenerationProgram(bpfProgram *C.struct_bpf_program) Program {
	var res Program
	programName := C.bpf_program__name(bpfProgram)
	baseProgram := BaseProgram{
		name:       C.GoString(programName),
		program:    bpfProgram,
		module:     m,
		pinnedPath: "",
	}
	secName := baseProgram.GetProgramSecName()
	if secName == "" {
		panic(baseProgram.name + " has no sec name!")
	}
	secNames := strings.Split(secName, "/")
	if secNames == nil {
		panic("nil slice!")
	}
	if len(secNames) < 2 {
		panic("invalid sec name:" + secName)
	}
	switch secNames[0] {
	case TracePoint:
		res = &TracePointProgram{
			BaseProgram: baseProgram,
		}
	case KProbe:
		res = &KprobeProgram{
			BaseProgram: baseProgram,
		}
	case UProbe, URetProbe:
		u := &UProbeProgram{
			BaseProgram: baseProgram,
		}
		res = u

		// add by func addr
		err := u.AddOffset(secNames[1])
		if err == nil {
			break
		}

		// add by funcName (addr resolved by elf)
		err = u.AddOffsetByFuncName(secNames[1])
		if err != nil {
			panic(err)
		}

	case XDP:
		res = &XDPProgram{
			BaseProgram: baseProgram,
		}
	case TC:
		res = &TCProgram{
			BaseProgram: baseProgram,
		}
	default:
		panic("no such ebpf program sec named:" + secName)
	}
	return res
}

func (m *Module) AddElfFile(path string) {
	m.elfPath = path
	f, err := elf.Open(path)
	if err != nil {
		panic(err)
	}
	m.elf = f
	m.AddDeferFunc(func() {
		f.Close()
	})
}

func (m *Module) NewPerfBuffer(mapName string) *C.struct_perf_buffer {
	bpfMap := m.maps[mapName]
	if mapType(bpfMap.bpfMap) != BPF_MAP_TYPE_PERF_EVENT_ARRAY {
		panic("map type error, need BPF_MAP_TYPE_PERF_EVENT_ARRAY")
	}

	pb := C.init_perf_buf(bpfMap.fd, C.int(1), nil)
	return pb
}

func (m *Module) GetAllMapName() (res []string) {
	var bpfMap *C.struct_bpf_map
	for {
		bpfMap = C.bpf_object__next_map(m.obj, bpfMap)
		if bpfMap == nil {
			break
		}
		mapName := C.bpf_map__name(bpfMap)
		res = append(res, C.GoString(mapName))
	}
	return
}

func (m *Module) LoadAllMap() error {
	for _, mapName := range m.GetAllMapName() {
		cs := C.CString(mapName)
		bpfMap, errno := C.bpf_object__find_map_by_name(m.obj, cs)
		C.free(unsafe.Pointer(cs))
		if bpfMap == nil {
			return fmt.Errorf("failed to load BPF map %s: %w", mapName, errno)
		}
		b := &BPFMap{
			bpfMap: bpfMap,
			module: m,
			fd:     C.bpf_map__fd(bpfMap),
		}
		m.maps[mapName] = b
	}
	return nil
}

func (m *Module) LoadPinnedMapFromPath(name string) {
	bpfPath := path.Join(PinnedFilePath, name)
	CPath := C.CString(bpfPath)
	defer C.free(unsafe.Pointer(CPath))
	mapFd := C.bpf_obj_get(CPath)

	heapMap := C.CString(name)
	defer C.free(unsafe.Pointer(heapMap))

	mapReuse := C.bpf_object__find_map_by_name(m.obj, heapMap)
	ret := C.bpf_map__reuse_fd(mapReuse, mapFd)
	if ret != 0 {
		err := fmt.Errorf("failed to reuse map object: %w", syscall.Errno(-ret))
		panic(err)
	}
	newMap := &BPFPinnedMap{
		path: bpfPath,
		fd:   mapFd,
		bPFMap: BPFMap{
			bpfMap: mapReuse,
		},
		keySize:   mapKeySize(mapReuse),
		valueSize: mapValueSize(mapReuse),
	}
	delete(m.maps, name)
	m.bpfPinnedMap[name] = newMap
}

func (m *Module) Close() {
	for _, f := range m.deferFunc {
		f()
	}
}
