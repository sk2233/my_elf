section .data  ; 定义数据段   这个是测试文件，功能可能不正常，不要用于使用
        text db "Hello World!",10   ; 使用 text（地址） 标记定义的byte数组  10 是回车符
section .text  ; 定义代码段
        global _start
_start:    ; 链接时使用的入口地址
        mov rax,1
        mov rdi,1
        mov rsi,text
        mov rdx,14
        syscall

        mov rax,60
        mov rdi,0
        syscall

        jmp _start ; 跳转到指定位置

        cmp rax,23   ; 对比的结果存储在  Flag 中
        je _start    ; 根据 cmp 对比存储在 Flag 中的值进行判断

        mov rax,[rbx] ; [] 取地址进行赋值

        call test ; 调用方法，call 与 jmp 最大的区别就是 call 需要恢复现场

test:
    mov ax,100
    ret
