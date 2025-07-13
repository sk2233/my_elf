section .data
        file db "test.txt",0 ; 文件名称必须以 0 结尾

section .bss
        text resb 7

section .text
        global _start

_start:
        mov rax,2     ; 打开文件
        mov rdi,file
        mov rsi,0
        mov rdx,0644o  ; 必须使用 o 表示 8进制
        syscall
        push rax

        mov rdi,rax  ; 读取文件
        mov rax,0
        mov rsi,text
        mov rdx,7
        syscall

        mov rax,3  ; 关闭文件
        pop rdi
        syscall

        mov rax,1  ; 打印读取的内容
        mov rdi,1
        mov rsi,text
        mov rdx,7
        syscall

        mov rax,60
        mov rdi,0
        syscall