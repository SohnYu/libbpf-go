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

struct openat_args {
	unsigned long long pad;
	int syscall_n;
	int dfd;
	const char *filename;
	int flags;
	umode_t mode;
};

struct sys_enter_openat_struct_ret
{
    int dfd;
    char filename[256];
    int flags;
    umode_t mode;
};

struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(key_size, sizeof(u32));
    __uint(value_size, sizeof(u32));
} events SEC(".maps");

char LICENSE[] SEC("license") = "GPL";

SEC("tracepoint/syscalls/sys_enter_openat")
int tracepoint__syscalls__sys_enter_openat(struct trace_event_raw_sys_enter *ctx)
{
    struct sys_enter_openat_struct_ret ret = {};
//    bpf_printk("%d----------%s", ctx->dfd, ctx->filename);
//        bpf_perf_event_output(ctx, &events, BPF_F_CURRENT_CPU, &result, sizeof(struct func_entry_record));
    return 0;
}
