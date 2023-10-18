package libbpfgo

import "C"
import "debug/elf"

type Module struct {
	obj          *C.struct_bpf_object
	elf          *elf.File
	elfPath      string
	loaded       bool
	deferFunc    []func()
	bpfProgram   map[string]Program
	maps         map[string]*BPFMap
	bpfPinnedMap map[string]*BPFPinnedMap
	bpfPerfEvent map[string]*PerfEvent
}

type BPFProgram struct {
	name       string
	program    *C.struct_bpf_program
	module     *Module
	pinnedPath string
}

type MapFlag uint32

const (
	MapFlagUpdateAny     MapFlag = iota // create new element or update existing
	MapFlagUpdateNoExist                // create new element if it didn't exist
	MapFlagUpdateExist                  // update existing element
	MapFlagFLock                        // spin_lock-ed map_lookup/map_update
)

const (
	TracePoint = "tracepoint"
	KProbe     = "kprobe"
	KRetProbe  = "kretprobe"
	UProbe     = "uprobe"
	URetProbe  = "uretprobe"
	XDP        = "xdp"
	TC         = "tc"
	LSM        = "lsm"
)

const (
	PinnedFilePath = "/sys/fs/bpf/"
)

const (
	BPF_MAP_TYPE_UNSPEC = iota
	BPF_MAP_TYPE_HASH
	BPF_MAP_TYPE_ARRAY
	BPF_MAP_TYPE_PROG_ARRAY
	BPF_MAP_TYPE_PERF_EVENT_ARRAY
	BPF_MAP_TYPE_PERCPU_HASH
	BPF_MAP_TYPE_PERCPU_ARRAY
	BPF_MAP_TYPE_STACK_TRACE
	BPF_MAP_TYPE_CGROUP_ARRAY
	BPF_MAP_TYPE_LRU_HASH
	BPF_MAP_TYPE_LRU_PERCPU_HASH
	BPF_MAP_TYPE_LPM_TRIE
	BPF_MAP_TYPE_ARRAY_OF_MAPS
	BPF_MAP_TYPE_HASH_OF_MAPS
	BPF_MAP_TYPE_DEVMAP
	BPF_MAP_TYPE_SOCKMAP
	BPF_MAP_TYPE_CPUMAP
	BPF_MAP_TYPE_XSKMAP
	BPF_MAP_TYPE_SOCKHASH
	BPF_MAP_TYPE_CGROUP_STORAGE
	BPF_MAP_TYPE_REUSEPORT_SOCKARRAY
	BPF_MAP_TYPE_PERCPU_CGROUP_STORAGE
	BPF_MAP_TYPE_QUEUE
	BPF_MAP_TYPE_STACK
	BPF_MAP_TYPE_SK_STORAGE
	BPF_MAP_TYPE_DEVMAP_HASH
	BPF_MAP_TYPE_STRUCT_OPS
	BPF_MAP_TYPE_RINGBUF
	BPF_MAP_TYPE_INODE_STORAGE
	BPF_MAP_TYPE_TASK_STORAGE
	BPF_MAP_TYPE_BLOOM_FILTER
)

type BPFMap struct {
	bpfMap *C.struct_bpf_map
	fd     C.int
	module *Module
}

type PerfEvent struct {
	perfBuffer *C.struct_perf_buffer
	stopChan   chan struct{}
}

type BPFPinnedMap struct {
	fd        C.int
	keySize   int
	valueSize int
	path      string
	name      string
	bPFMap    BPFMap
}
