#include <bpf/bpf.h>
#include <bpf/libbpf.h>
#include <linux/bpf.h> // uapi
#include <errno.h>
extern int RingBufferFunc(void *ctx, void *data, size_t size);
extern void PerfBufferFunc(void *ctx, int cpu, void *data, __u32 size);
extern void PerfBufferLostFunc(void *ctx, int cpu, __u64 cnt);

struct ring_buffer *init_ring_buf(int map_fd, uintptr_t ctx) {
  struct ring_buffer *rb = NULL;

  rb = ring_buffer__new(map_fd, RingBufferFunc, (void *)ctx, NULL);
  if (!rb) {
    int saved_errno = errno;
    fprintf(stderr, "Failed to initialize ring buffer: %s\n", strerror(errno));
    errno = saved_errno;
    return NULL;
  }

  return rb;
}

struct perf_buffer *init_perf_buf(int map_fd, int page_cnt, void* ctx) {
  struct perf_buffer_opts pb_opts = {};
  struct perf_buffer *pb = NULL;

  pb_opts.sz = sizeof(struct perf_buffer_opts);

  pb = perf_buffer__new(map_fd, page_cnt, PerfBufferFunc, PerfBufferLostFunc,
                        ctx, &pb_opts);
  if (!pb) {
    int saved_errno = errno;
    fprintf(stderr, "Failed to initialize perf buffer: %s\n", strerror(errno));
    errno = saved_errno;
    return NULL;
  }

  return pb;
}