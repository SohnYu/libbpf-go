//+build ignore
#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_tracing.h>
#ifdef asm_inline
#undef asm_inline
#define asm_inline asm
#endif
#define BPF_F_CURRENT_CPU 0xffffffffULL
#define GO_PT_REGS_IP(x) ((x)->ip)
#define GOID_OFFSET 152 // go 1.19

struct func_entry_record
{
    u64 goid;
    u64 time;
    uintptr_t func_addr;
};

struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(key_size, sizeof(u32));
    __uint(value_size, sizeof(u32));
} events SEC(".maps");

static __always_inline
        u64 get_goid(struct pt_regs *ctx)
{
//    struct task_struct *task = (struct task_struct *)bpf_get_current_task();
//    unsigned long fsbase = 0;
//    void *g = NULL;
    u64 goid = 0;
//    bpf_probe_read(&fsbase, sizeof(fsbase), &task->thread.fsbase);
//    bpf_probe_read(&g, sizeof(g), (void*)fsbase-8);
    bpf_probe_read(&goid, sizeof(goid), (void *)(ctx->r14+GOID_OFFSET));
    return goid;
}

static __always_inline
        int process_func_ctx(struct pt_regs *ctx)
{
    struct func_entry_record result = {};
    int goid = get_goid(ctx);
    uintptr_t func_addr = (uintptr_t)GO_PT_REGS_IP(ctx);
    u64 ktime = bpf_ktime_get_ns();

    result.goid = goid;
    result.time = ktime;
    result.func_addr = func_addr;

    bpf_perf_event_output(ctx, &events, BPF_F_CURRENT_CPU, &result, sizeof(struct func_entry_record));
    return 1;
}

char LICENSE[] SEC("license") = "GPL";

SEC("uprobe/main.GoroutineTest")
int gindemo(struct pt_regs *ctx)
{
    return process_func_ctx(ctx);
}

SEC("uretprobe/3624724")
int gindemoret(struct pt_regs *ctx)
{
    return process_func_ctx(ctx);
}
