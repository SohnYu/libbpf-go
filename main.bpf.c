// SPDX-License-Identifier: GPL-2.0
// Copyright (c) 2017 Facebook

#include <vmlinux.h>
#include <bpf/bpf_helpers.h>
#include "src/common.h"
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

struct bind_socket {
	  unsigned long long pad;
	  int nr;
    long fd;
    struct sockaddr *umyaddr;
    int addrlen;
};

struct sys_enter_close {
	  unsigned long long pad;
	  int nr;
    long fd;
};

struct bind_socket_ret {
//	  unsigned long long pad;
//	  int nr;
//    long fd;
//    sa_family_t sa_family;
    char sa_data[14];
//    int addrlen;
//    int _pad;
};

struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(key_size, sizeof(u32));
    __uint(value_size, sizeof(u32));
} events SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(key_size, sizeof(u32));
    __uint(value_size, sizeof(u32));
} close_events SEC(".maps");

SEC("tracepoint/syscalls/sys_enter_bind")
int sys_enter_bind(struct bind_socket *ctx)
{
    struct bind_socket_ret e = {};
//    e.prev_pid = ctx->prev_pid;
//    e.prev_prio = ctx->prev_prio;
//    e.prev_state = ctx->prev_state;
//    e.next_pid = ctx->next_pid;
//    e.next_prio = ctx->next_prio;
//    bpf_probe_read(&e.prev_comm, sizeof(e.prev_comm), ctx->prev_comm);
    bpf_probe_read(&e.sa_data, sizeof(e.sa_data), ctx->umyaddr->sa_data);

    bpf_perf_event_output(ctx, &events, BPF_F_CURRENT_CPU, &e, sizeof(struct bind_socket_ret));

//  bpf_printk("xsk_redirect: %d %s\n", ctx->next_pid, ctx->next_comm);
	return 0;
}

SEC("tracepoint/syscalls/sys_enter_close")
int sys_exit_bind(struct sys_enter_close *ctx)
{
    struct sys_enter_close e = {};
//    e.prev_pid = ctx->prev_pid;
//    e.prev_prio = ctx->prev_prio;
//    e.prev_state = ctx->prev_state;
//    e.next_pid = ctx->next_pid;
//    e.next_prio = ctx->next_prio;
//    bpf_probe_read(&e.prev_comm, sizeof(e.prev_comm), ctx->prev_comm);
//    bpf_probe_read(&e.sa_data, sizeof(e.sa_data), ctx->umyaddr->sa_data);
    e.fd = ctx->fd;
    bpf_perf_event_output(ctx, &close_events, BPF_F_CURRENT_CPU, &e, sizeof(struct sys_enter_close));

//  bpf_printk("xsk_redirect: %d %s\n", ctx->next_pid, ctx->next_comm);
	return 0;
}

char _license[] SEC("license") = "GPL";