section .data
        text db "Hello",10,0  ; 一定要以 0 结尾
section .text
        global _start
_start:
        mov rax,text
        call _print

        mov rax,60
        mov rdi,0
        syscall

_print:
        push rax ; 保存 rax 的值
        mov rbx,0
_loop:
        inc rax
        inc rbx ; 计算字符串的长度
        mov cl,[rax]
        cmp cl,0
        jne _loop

        mov rax,1  ; 调用打印函数
        mov rdi,1
        pop rsi
        mov rdx,rbx
        syscall
        ret