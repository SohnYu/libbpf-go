# libbpf-go
Load bpf program in golang
dependence make gcc

### libbpf 
在ubuntu22.04 server版本  
dependence libelf-dev zib

> go tool objdump -s "main.GoroutineTest" main查看main程序中main.GoroutineTest的汇编代码
> TEXT main.GoroutineTest(SB) /home/ysh/Nexus/dockerFile/goroutineIdTest/main.go
main.go:27            0x778d80                4c8d6424f0                      LEAQ -0x10(SP), R12                                                                                     
main.go:27            0x778d85                4d3b6610                        CMPQ 0x10(R14), R12                                                                                     
main.go:27            0x778d89                0f8606010000                    JBE 0x778e95                                                                                            
main.go:27            0x778d8f                4881ec90000000                  SUBQ $0x90, SP                                                                                          
main.go:27            0x778d96                4889ac2488000000                MOVQ BP, 0x88(SP)                                                                                       
main.go:27            0x778d9e                488dac2488000000                LEAQ 0x88(SP), BP                                                                                       
main.go:27            0x778da6                4889442438                      MOVQ AX, 0x38(SP)                                                                                       
main.go:28            0x778dab                e890feffff                      CALL main.goId(SB)                                                                                      
main.go:28            0x778db0                48894c2430                      MOVQ CX, 0x30(SP)                                                                                       
main.go:28            0x778db5                48897c2428                      MOVQ DI, 0x28(SP)                                                                                       
main.go:28            0x778dba                440f117c2440                    MOVUPS X15, 0x40(SP)                                                                                    
main.go:28            0x778dc0                440f117c2450                    MOVUPS X15, 0x50(SP)                                                                                    
main.go:28            0x778dc6                e85536c9ff                      CALL runtime.convTstring(SB)                                                                            
main.go:28            0x778dcb                488d0daefd0200                  LEAQ 0x2fdae(IP), CX                                                                                    
main.go:28            0x778dd2                48894c2440                      MOVQ CX, 0x40(SP)                                                                                       
main.go:28            0x778dd7                4889442448                      MOVQ AX, 0x48(SP)                                                                                       
main.go:28            0x778ddc                488b442430                      MOVQ 0x30(SP), AX                                                                                       
main.go:28            0x778de1                4885c0                          TESTQ AX, AX                                                                                            
main.go:28            0x778de4                7406                            JE 0x778dec                                                                                             
main.go:28            0x778de6                488b4808                        MOVQ 0x8(AX), CX                                                                                        
main.go:28            0x778dea                eb03                            JMP 0x778def                                                                                            
main.go:28            0x778dec                4889c1                          MOVQ AX, CX                                                                                             
main.go:28            0x778def                48894c2450                      MOVQ CX, 0x50(SP)                                                                                       
main.go:28            0x778df4                488b542428                      MOVQ 0x28(SP), DX                                                                                       
main.go:28            0x778df9                4889542458                      MOVQ DX, 0x58(SP)                                                                                       
print.go:294          0x778dfe                488b1dcbfa3d00                  MOVQ os.Stdout(SB), BX                                                                                  
print.go:294          0x778e05                488d0534d41600                  LEAQ go.itab.*os.File,io.Writer(SB), AX                                                                 
print.go:294          0x778e0c                488d4c2440                      LEAQ 0x40(SP), CX                                                                                       
print.go:294          0x778e11                bf02000000                      MOVL $0x2, DI                                                                                           
print.go:294          0x778e16                4889fe                          MOVQ DI, SI                                                                                             
print.go:294          0x778e19                e8a2afd4ff                      CALL fmt.Fprintln(SB)                                                                                   
main.go:29            0x778e1e                488d051b280300                  LEAQ 0x3281b(IP), AX                                                                                    
main.go:29            0x778e25                e8965cc9ff                      CALL runtime.newobject(SB)                                                                              
main.go:29            0x778e2a                66c7006f6b                      MOVW $0x6b6f, 0(AX)                                                                                     
context.go:1023       0x778e2f                440f117c2460                    MOVUPS X15, 0x60(SP)                                                                                    
context.go:1023       0x778e35                440f117c2468                    MOVUPS X15, 0x68(SP)                                                                                    
context.go:1023       0x778e3b                440f117c2478                    MOVUPS X15, 0x78(SP)                                                                                    
context.go:1025       0x778e41                4889442470                      MOVQ AX, 0x70(SP)                                                                                       
context.go:1025       0x778e46                48c744247802000000              MOVQ $0x2, 0x78(SP)                                                                                     
context.go:1025       0x778e4f                48c784248000000002000000        MOVQ $0x2, 0x80(SP)                                                                                     
context.go:1023       0x778e5b                488d05de1a0700                  LEAQ 0x71ade(IP), AX                                                                                    
context.go:1023       0x778e62                488d5c2460                      LEAQ 0x60(SP), BX                                                                                       
context.go:1023       0x778e67                e83433c9ff                      CALL runtime.convT(SB)                                                                                  
context.go:1023       0x778e6c                bbc8000000                      MOVL $0xc8, BX                                                                                          
context.go:1023       0x778e71                488d0dd8ef1600                  LEAQ go.itab.github.com/gin-gonic/gin/render.Data,github.com/gin-gonic/gin/render.Render(SB), CX        
context.go:1023       0x778e78                4889c7                          MOVQ AX, DI                                                                                             
context.go:1023       0x778e7b                488b442438                      MOVQ 0x38(SP), AX                                                                                       
context.go:1023       0x778e80                e85b51ffff                      CALL github.com/gin-gonic/gin.(*Context).Render(SB)                                                     
main.go:30            0x778e85                488bac2488000000                MOVQ 0x88(SP), BP                                                                                       
main.go:30            0x778e8d                4881c490000000                  ADDQ $0x90, SP                                                                                          
main.go:30            0x778e94                c3                              RET      // 这里即是golang方法的返回地址                                                                                          
main.go:27            0x778e95                4889442408                      MOVQ AX, 0x8(SP)                                                                                        
main.go:27            0x778e9a                e841c3ceff                      CALL runtime.morestack_noctxt.abi0(SB)                                                                  
main.go:27            0x778e9f                488b442408                      MOVQ 0x8(SP), AX                                                                                        
main.go:27            0x778ea4                e9d7feffff                      JMP main.GoroutineTest(SB)


**注：**
1. c结构体无属性，在go中不可直接在栈上分配

#### 常用用法
```c
bpf_probe_read(&e.prev_comm, sizeof(e.prev_comm), ctx->prev_comm); // 从ctx->prev_comm读取指定长度的字符串到e.prev_comm
```
#### BPF_MAP_TYPE_PERF_EVENT_ARRAY 定义
```c
struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(key_size, sizeof(u32));
    __uint(value_size, sizeof(u32));
} events SEC(".maps");
```

#### UProbe
```c

```
#### TracePoint
```c

```