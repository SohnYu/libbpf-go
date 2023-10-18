ROOT_PATH = $(abspath ./)
LIB_PATH = $(ROOT_PATH)/lib
ARCH := $(shell uname -m | sed 's/x86_64/amd64/g; s/aarch64/arm64/g')
OS_VERSION := $(shell uname -r)

build-go:
	CC=clang \
		CGO_CFLAGS="-I$(LIB_PATH)/" \
		CGO_LDFLAGS="-lelf -lz $(LIB_PATH)/libbpf.a" \
		GOOS=linux GOARCH=$(ARCH) \
		go build \
		-tags netgo -ldflags  '-w -extldflags "-static"' \
		-o libbpf-go ./cmd/libbpf-go.go

build-libbpf:
	CC="gcc" CFLAGS="-g -O2 -Wall -fpie" \
       /usr/bin/make -C $(ROOT_PATH)/libbpf/src \
    	BUILD_STATIC_ONLY=y \
    	OBJDIR=$(LIB_PATH)/libbpf \
    	DESTDIR=$(LIB_PATH) \
    	INCLUDEDIR= LIBDIR= UAPIDIR= install

build-ebpf:
	clang -g -O2 -c -target bpf -o main.bpf.o main.bpf.c

gen-vmlinux:
	./extract-vmlinux /boot/vmlinuz-$(OS_VERSION) > vmlinux
	bpftool btf dump file ./vmlinux format c > vmlinux.h

