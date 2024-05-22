package main

/*
//#include <stdio.h>
//#include <errno.h>
#include <stdlib.h>
//#include <string.h>
//#include <stdarg.h>
//#include <sys/resource.h>
//#include <sys/syscall.h>
//#include <unistd.h>
#include <bpf/libbpf.h>
*/
import "C"
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func main() {
	module := NewModuleFromFile("/root/libbpf-go/main.bpf.o")
	defer module.Close()
	err := module.LoadAllMap()
	if err != nil {
		panic(err)
	}
	err = module.LoadAllProgram()
	if err != nil {
		panic(err)
	}
	err = module.LoadAllIntoKernel()
	if err != nil {
		panic(err)
	}
	module.NewPerfBuffer("events")
	module.PerfStart("events")
	module.NewPerfBuffer("close_events")
	module.PerfStart("close_events")
	select {}
}

func perfOutput(mapName string, cpu int, data []byte) {
	fmt.Println(mapName)
	res := BindSocketRet{}
	err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &res)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.String())
	//var a FuncEntryRecord
	//if err := binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &a); err != nil {
	//	panic(err)
	//}
	//fmt.Printf("GOID:%d TIME:%d FuncAddr:%x\n", a.GoId, a.KTime, a.FuncAddr)
}

func ringBufferOutput(data []byte) {

}

type FuncEntryRecord struct {
	GoId     uint64
	KTime    uint64
	FuncAddr uint32
}

// SchedSwitch
/*
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
*/

type SchedSwitch struct {
	Pad       int64
	PrevComm  [16]byte
	PrevPid   int32
	PrevPrio  int32
	PrevState int64
	NextComm  [16]byte
	NextPid   int32
	NextPrio  int32
}

type BindSocketRet struct {
	//Pad      int64
	//Nr       int32
	//_        int32
	//Fd       int64
	//SaFamily uint16
	//_        uint16
	//_        uint32
	SaData [14]byte
	//_        [2]byte
	//Addrlen  int
}

func (ret BindSocketRet) String() string {
	// 解析 SaData 中的 IP 地址和端口号
	// 解析端口号（sa_data 的前两个字节）
	port := uint16(ret.SaData[0])<<8 | uint16(ret.SaData[1])
	// 解析 IP 地址（sa_data 的后四个字节）
	ip := net.IPv4(ret.SaData[2], ret.SaData[3], ret.SaData[4], ret.SaData[5])
	// 将 IP 地址和端口号格式化为字符串
	address := fmt.Sprintf("%s:%d", ip.String(), port)
	return fmt.Sprintf("%s Address: %s", time.Now().Local().String(), address)
}
