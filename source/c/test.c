const int name = 2233;  // 只读  直接在使用处进行替换
int age=2233;  // 可读可写      data 段  bss 段中的数据

int main(){ // 暂时不接受输入参数  固定入口函数
    char arr ="sdsd"; // 自动分配
    int num=10;
    for(int i=0;i<10;i++){ // 内存都是固定分配的，没有动态申请
        if(i>5){ // 同一作用域中 堆栈大小是可预测的
            continue; // 根据栈顶 for 循环跳转到不同位置
        }else{
            break;
        }
    }
    print(22+33*44/55); // 支持 + - * /
    print("Hello World"); // 绑定的系统调用
    char temp[100]; // 指定使用的内容  唯一支持  [] 的地方
    input(temp); // 绑定的系统调用  要注意自定义函数不能与系统调用同名
    int file=open("test.txt",0);  // mode
    write(file,"Hello World");
    read(file,temp);
    close(file);
    int res=max(22,33); // 传参使用栈，局部变量就不使用栈了
    return 0;  // 所有函数必须有 return
}
// 暂时只支持入参与返回  int  可以使用 int 构造  str 指针进行使用
int max(int num1,int num2){ // 为了支持递归，局部变量必须使用栈
    if(num1>num2){
        return num1;
    }else{
        return num2;
    }
}