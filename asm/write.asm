section .data
        file db "test.txt",0 ; 文件名称必须以 0 结尾
        text db "My Input"

section .text
        global _start

_start:
        mov rax,2     ; 打开文件
        mov rdi,file
        mov rsi,65
        mov rdx,0644o  ; 必须使用 o 表示 8进制
        syscall
        push rax

        mov rdi,rax  ; 写入文件
        mov rax,1
        mov rsi,text
        mov rdx,8
        syscall

        mov rax,3  ; 关闭文件
        pop rdi
        syscall

        mov rax,60
        mov rdi,0
        syscall