// SPDX-License-Identifier: GPL-2.0
// Copyright (c) 2017 Facebook

#include <vmlinux.h>
#include <bpf/bpf_helpers.h>
#define TASK_COMM_LEN 16

/* taken from /sys/kernel/tracing/events/sched/sched_switch/format */
struct sched_switch_args {
	unsigned long long pad;
	char prev_comm[TASK_COMM_LEN];
	int prev_pid;
	int prev_prio;
	long long prev_state;
	char next_comm[TASK_COMM_LEN];
	int next_pid;
	int next_prio;
};
#define BPF_F_CURRENT_CPU 0xffffffffULL
struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(key_size, sizeof(u32));
    __uint(value_size, sizeof(u32));
} events SEC(".maps");

SEC("tracepoint/sched/sched_switch")
int oncpu(struct sched_switch_args *ctx)
{
    struct sched_switch_args e = {};
    e.prev_pid = ctx->prev_pid;
    bpf_probe_read(&e.prev_comm, sizeof(e.prev_comm), ctx->prev_comm);

    bpf_perf_event_output(ctx, &events, BPF_F_CURRENT_CPU, &e, sizeof(struct sched_switch_args));

//  bpf_printk("xsk_redirect: %d %s\n", 1111, ctx->next_comm);
	return 0;
}

char _license[] SEC("license") = "GPL";