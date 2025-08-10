section .bss
    temp resb 10
section .text
    global _start
_start:
    mov rax,2233
    call _print

    mov rax,60 ; exit
    mov rdi,0
    syscall
_print:
    mov rbx,temp
    mov rcx,0
_loop1:
    mov rdx,0 ; 存储余数
    mov rsi,10
    div rsi  ;  rax=rax/rsi
    add rdx,48  ; 余数+'0'
    mov [rbx],rdx  ; 进行暂存
    inc rbx   ; 移动计算与偏移位置
    inc rcx
    cmp rax,0
    jne _loop1

    mov rax,1  ; 按数目输出
    mov rdi,1
    mov rsi,temp
    mov rdx,rcx
    syscall

    ret