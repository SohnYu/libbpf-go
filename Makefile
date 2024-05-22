ROOT_PATH = $(abspath ./)
LIB_PATH = $(ROOT_PATH)/lib
ARCH := $(shell uname -m | sed 's/x86_64/amd64/g; s/aarch64/arm64/g')
OS_VERSION := $(shell uname -r)
CGO_CFLAGS_DYN = "-I. -I/usr/include/"
CGO_LDFLAGS_DYN = "-lelf -lz -lbpf"

build-go-dynamic:
	CC=clang \
		CGO_CFLAGS=$(CGO_CFLAGS_DYN) \
		CGO_LDFLAGS=$(CGO_LDFLAGS_DYN) \
		go build .

build-go:
	CC=clang \
		CGO_CFLAGS="-I$(LIB_PATH)/" \
		CGO_LDFLAGS="-lelf -lz $(LIB_PATH)/libbpf.a" \
		GOOS=linux GOARCH=$(ARCH) \
		go build \
		-tags netgo -ldflags  '-w -extldflags "-static"' \
		-o libbpf-go .

build-libbpf:
	CC="gcc" CFLAGS="-g -O2 -Wall -fpie" \
       /usr/bin/make -C $(ROOT_PATH)/libbpf/src \
    	BUILD_STATIC_ONLY=y \
    	OBJDIR=$(LIB_PATH)/libbpf \
    	DESTDIR=$(LIB_PATH) \
    	INCLUDEDIR= LIBDIR= UAPIDIR= install

build-ebpf:
	clang -g -O2 -I$(abspath ./) -c -target bpf -o main.bpf.o main.bpf.c

gen-vmlinux:
	bpftool btf dump file /sys/kernel/btf/vmlinux format c > vmlinux.h
	if [ $? -eq 0 ]; then
		echo "vmlinux gen success!"
	else
		./extract-vmlinux /boot/vmlinuz-$(OS_VERSION) > vmlinux
		bpftool btf dump file ./vmlinux format c > vmlinux.h
	fi


