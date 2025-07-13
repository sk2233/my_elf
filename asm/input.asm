section .data
        text1 db "What is your name?"
        text2 db "Hello "

section .bss
        name resb 16  ; 设置 16 位的 bss 段

section .text
        global _start
_start:
        mov rsi,text1
        mov rdx,18
        call _print  ; 调用方法，把部分操作系统调用发生变动的参数放在外面

        call _input

        mov rsi,text2
        mov rdx,6
        call _print

        mov rsi,name
        mov rdx,16
        call _print

        mov rax,60
        mov rdi,0
        syscall

_print:
        mov rax,1
        mov rdi,1
        syscall
        ret

_input:
        mov rax,0   ; 准备读取的系统调用
        mov rdi,0
        mov rsi,name
        mov rdx,16
        syscall
        ret