[TOC]



## Shell入门



## 条件判断



## 循环控制



## 变量的高级用法

### 变量替换

|            语法            |                     说明                     |
| :------------------------: | :------------------------------------------: |
|      ${变量#匹配规则}      |            从头开始匹配，最短删除            |
|     ${变量##匹配规则}      |            从头开始匹配，最长删除            |
|      ${变量%匹配规则}      |            从尾开始匹配，最短删除            |
|     ${变量%%匹配规则}      |            从尾开始匹配，最长删除            |
| ${变量/旧字符串/新字符串}  | 替换变量内的旧字符串为新字符串，只替换第一个 |
| ${变量//旧字符串/新字符串} |   替换变量内的旧字符串为新字符串，全部替换   |

```bash
variable_1="I love you,Do you love me"

var1=${variable_1#*ov}
#e you,Do you love me
var2=${variable_1##*ov}
#e me
var3=${variable_1%ov*}
# I love you,Do you l
var4=${variable_1%%ov*}
#I l

var5=${PATH/bin/BIN}
# /usr/lib64/qt-3.3/BIN:/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin:/root/bin
var6=${PATH//bin/BIN}
# /usr/lib64/qt-3.3/BIN:/usr/local/sBIN:/usr/local/BIN:/sBIN:/BIN:/usr/sBIN:/usr/BIN:/root/BIN
```

### 字符串处理

#### 计算字符串长度
|        |        语法         |
| ------ | :-----------------: |
| 方法一 |     ${#string}      |
| 方法二 | expr length $string |
```bash
var1="Hello World"
len=${#var1}
len=`expr length "$var1"`
# 11
# 字符串var1中有空格，必须加双引号，不然报语法错误
```
#### 获取字符索引位置
```bash
expr index "$string" substr
```
```bash
var1="quickstart is a app"
ind=`expr index "$var1" start`
# 6
# 索引从1开始查找
#查找的是字符索引位置，不是子串
```
#### 获取子串长度
```bash
expr match "$string" substr
```
```bash
var1="quickstart is a app"
sub_len=`expr match "$var1" app`
# 0
# 从头开始匹配，若未匹配到则返回0
```
#### 字符串截取

|                 语法                  |              说明               |
| :-----------------------------------: | :-----------------------------: |
|          ${string:position}           |     从string的position开始      |
|       ${string:position:length}       | 从string的position开始到length  |
|         ${string: -position}          |         从右边开始匹配          |
|         ${string:(position)}          |         从左边开始匹配          |
| expr substr $string $position $length | 从postion开始，匹配长度为length |

```bash
var1="kafka hadoop yarn mapreduce"

substr1=${var1:10}
# op yarn mapreduce
substr2=${var1:10:6}
# op yar
substr3=${var1: -5}
# educe
substr4=${var1: -10:4}
# map
substr5=`expr substr "$var1" 5 10`
# a hadoop y
# 注意：使用expr，索引计数是从1开始计算；使用${string:position}，索引计数是从0开始计数
```

### 字符串处理练习

#### 需求描述

```tex
	变量string="Bigdata process framework is Hadoop,Hadoop is an open source project"
	执行脚本后，打印输出string字符串变量，并给出用户以下选项：

	(1)、打印string长度
	(2)、删除字符串中所有的Hadoop
	(3)、替换第一个Hadoop为Mapreduce
	(4)、替换全部Hadoop为Mapreduce
	
	用户输入数字1|2|3|4，可以执行对应项的功能；输入q|Q则退出交互模式
```
#### 思路分析

##### 1、将不同的功能模块划分，并编写函数

	function print_tips
	function len_of_string
	function del_hadoop
	function rep_hadoop_mapreduce_first
	function rep_hadoop_mapreduce_all
##### 2、实现第一步所定义的功能函数

	function print_tips
	{
	    echo "******************************"
	    echo "(1) 打印string长度"
	    echo "(2) 删除字符串中所有的Hadoop"
	    echo "(3) 替换第一个Hadoop为Mapreduce"
	    echo "(4) 替换全部Hadoop为Mapreduce"
	    echo "******************************"
	}
	
	function len_of_string
	{	
	    echo "${#string}"
	}
	
	function del_hadoop
	{
	    echo "${string//Hadoop/}"
	}
	
	function rep_hadoop_mapreduce_first
	{
	    echo "${string/Hadoop/Mapreduce}"
	}
	
	function rep_hadoop_mapreduce_all
	{
	    echo "${string//Hadoop/Mapreduce}"
	}
##### 3、程序主流程的设计

#### 代码实现

```bash
#!/bin/bash
#

string="Bigdata process framework is Hadoop,Hadoop is an open source project"

function print_tips
{
	echo "******************************"
	echo "(1) 打印string长度"
	echo "(2) 删除字符串中所有的Hadoop"
	echo "(3) 替换第一个Hadoop为Mapreduce"
	echo "(4) 替换全部Hadoop为Mapreduce"
	echo "******************************"
}

function len_of_string
{	
	echo "${#string}"
}

function del_hadoop
{
	echo "${string//Hadoop/}"
}

function rep_hadoop_mapreduce_first
{
	echo "${string/Hadoop/Mapreduce}"
}

function rep_hadoop_mapreduce_all
{
	echo "${string//Hadoop/Mapreduce}"
}

while true
do
	echo "【string=$string】"
	echo
	print_tips
	read -p "Pls input your choice(1|2|3|4|q|Q): " choice
	
case $choice in
	1)
		len_of_string
		;;
	2)
		del_hadoop
		;;
	3)
		rep_hadoop_mapreduce_first
		;;
	4)
		rep_hadoop_mapreduce_all
		;;
	q|Q)
		exit
		;;
	*)
		echo "Error,input only in {1|2|3|4|q|Q}"
		;;
esac
done
```

### 命令替换
|      | 语法 |
| ---- | :-------: |
| 方法一 | &#96;command` |
| 方法二 | $(command) |

#### 示例1：获取系统得所有用户并输出	

```bash
#!/bin/bash
#
index=1
for user in `cat /etc/passwd | cut -d ":" -f 1`
do
    echo "This is $index user: $user"
    index=$(($index + 1))
done
```
#### 示例2：根据系统时间计算今年或明年

```bash
echo "This is $(date +%Y) year"
echo "This is $(($(date +%Y) + 1)) year"
# Bash中的(())用来支持算术表达式

# C语言规则运算
# $((exp))，exp为符合C语言规则的运算表达式
```
#### 示例3：根据系统时间获取今年还剩下多少星期，已经过了多少星期

```bash
date +%j
echo "This year have passed $(date +%j) days"
echo "This year have passed $(($(date +%j)/7)) weeks"

echo "There is $((365 - $(date +%j))) days before new year"
echo "There is $(((365 - $(date +%j))/7)) days before new year"
```
#### 示例4：判定nginx进程是否存在，若不存在则自动拉起该进程

```bash
#!/bin/bash
#
nginx_process_num=$(ps -ef | grep nginx | grep -v grep | wc -l)
if [ $nginx_process_num -eq 0 ];then
	systemctl start nginx
fi
```
#### 总结

	``和$()两者是等价的，但推荐初学者使用$()，易于掌握；缺点是极少数UNIX可能不支持，但``都是支持的
	$(())主要用来进行整数运算，包括加减乘除,引用变量前面可以加$，也可以不加$
	
	$(( (100 + 30) / 13 ))
	
	num1=20;num2=30
	((num++));
	((num--))
	$(($num1+$num2*2))

### 有类型变量

#### 声明变量为只读类型

```bash
declare -r
```

```bash
declare -r var="hello"
var="world"				
#test.sh: line 5: var: readonly variable
```

#### 声明变量类型为整型

```bash
declare -i
```

```bash
num1=2001
num2=$num1+1
echo $num2

declare -i num2
num2=$num1+1
echo $num2
```
#### 显示定义的函数和内容

```bash
declare -f
```

#### 显示定义的函数

```bash
declare -F
```

#### 显示定义的数组

```bash
declare -a
```

```bash
array=("jones" "mike" "kobe" "jordan")
```

##### 输出数组内容

```bash
echo ${array[@]}	
# 输出全部内容
echo ${array[1]}	
# 输出下标索引为1的内容
```

##### 获取数组长度

```bash
echo ${#array}		
# 数组内元素个数
echo ${#array[2]}	
# 数组内下标索引为2的元素长度						
```

##### 给数组某个下标赋值

```bash
array[0]="lily"				
# 给数组下标索引为1的元素赋值为lily
array[20]="hanmeimei"		
# 在数组尾部添加一个新元素
```

##### 删除元素

```bash
unset array[2]		
# 清除元素
unset array			
# 清空整个数组
```
##### 分片访问

```bash
${array[@]:1:4}		
# 显示数组下标索引从1开始到3的3个元素，不显示索引为4的元素
```

##### 内容替换

```bash
${array[@]/an/AN}	
# 将数组中所有元素内包含kobe的子串替换为mcgrady
```

##### 数组遍历

```bash
for v in ${array[@]}
do
	echo $v
done
```

##### 声明为环境变量

```bash
declare -x
```

##### 取消已声明的变量

```bash
declare +r
declare +i
declare +a
declar
```

### bash数学运算之expr

```bash
num1=20
num2=100

expr $num1 \| $num2
expr $num1 \& $num2
expr $num1 \< $num2
expr $num1 \< $num2
expr $num1 \<= $num2
expr $num1 \> $num2
expr $num1 \>= $num2
expr $num1 = $num2
expr $num1 != $num2
expr $num1 + $num2
expr $num1 - $num2
expr $num1 \* $num2
expr $num1 / $num2
expr $num1 % $num2
```

#### 示例1：提示用户输入一个正整数num，然后计算1+2+3+...+sum的值；必须对num是否为正整数做判断，不符合应当允许再此输入

```bash
#!/bin/bash
#
while true
do
	read -p "Pls enter a positive integer(num>0): " num

	expr $num + 1 &> /dev/null
	if [ $? -ne 0 ];then
		echo "Error,You must input a interger"
		continue
	else
		if [ `expr $num \> 0` -ne 1 ];then
			echo "Error,You must input a postive interger"
			continue
		else
			sum=0
			for((i=0;i<=$num;i++))
			do
				sum=`expr $sum + $i`
			done
			echo "1+2+3+4+5+...+$num=$sum"
		fi
	fi
done
```
### bash数学运算之bc

```bash
echo "options;expression" | bc

num1=23.5
num2=50

var1=`echo "scale=2;$num1 * $num2" | bc`
```
## 函数的高级用法

### 函数定义和使用

#### 函数定义两种语法

```bash
name()
{
    command1
    command2
    .....
    commandn
}
```
```bash
function name
{
    command1
    command2
    .....
    commandn
}
```
#### 函数使用

调用函数直接使用函数名即可，相当于一条命令

##### 示例1

```bash
function hello
{
	echo "Hello,Zhangsan"
}

hello
```
##### 示例2

```bash
function print_num
{
	for((i=0;i<=10;i++))
	do
		echo -n "$i "
	done
}
```
##### 示例3

需求描述：写一个监控nginx的脚本；如果Nginx服务宕掉，则该脚本可以检测到并将进程启动；如果正常运行，则不做任何处理

```bash
#!/bin/bash
#

this_pid=$$

function nginx_daemon
{
	status=$(ps -ef | grep -v $this_pid | grep nginx | grep -v grep &> /dev/null)
	if [ $? -eq 1 ];then
		systemctl start nginx && echo "Start Nginx Successful" || echo "Failed To Start Nginx"
	else
		echo "Nginx is RUNNING Well"
		sleep 5
	fi
}

while true
do
	nginx_daemon
done

```

##### 示例4

###### 终端命令行定义函数

在脚本中定义好disk_usage函数，然后直接使用. test.sh，再使用declare -F查看，是否可以列出disk_usage函数

```bash
function disk_usage
{
	if [ $# -eq 0 ];then
		df
	else
		case $1 in
			-h)
				df -h
				;;
			-i)
				df -i
				;;
			-ih|-hi)
				df -ih
				;;
			-T)
				df -T
				;;
			*)
				echo "Usage: $0 { -h|-i|-ih|-T }"
				;;
		esac
	fi
}	
```

### 向函数传递参数

函数传参和给脚本传参类似，都是使用$1 $2 $3 $4 $5 $6 $7这种方式

#### 示例1

需求描述：写一个脚本，该脚本可以实现计算器的功能，可以进行+-*/四种计算。

```bash
#!/bin/bash
#

function calculate
{
	case "$2" in
		+)
			echo "$1 + $3 = $(expr $1 + $3)"
			;;
		-)
			echo "$1 + $3 = $(expr $1 - $3)"
			;;
		\*)			
			echo "$1 * $3 = $(expr $1 \* $3)"
			;;
		/)
			echo "$1 / $3 = $(expr $1 / $3)"
			;;
	esac
}

calculate $1 $2 $3
```

### 函数返回值

```
return 返回的是函数退出状态码，函数结束return不返回函数返回值，可以在前面使用echo返回函数返回值
echo  返回的是函数返回值，函数退出状态码是函数最后一条命令的退出状态码 
```

#### 使用return返回退出状态码

##### 示例1：测试nginx是否在运行

```bash
#!/bin/bash
#

this_pid=$$

function is_nginx_running
{

	ps -ef | grep nginx | grep -v $this_pid | grep -v grep > /dev/null 2>&1
	if [ $? -eq 0 ];then
		return 0
	else
		return 1
	fi
}

is_nginx_running && echo "Nginx is running" || echo "Nginx is stopped"
```

#### 使用echo返回值

##### 示例1：两数字相加

```bash
#!/bin/bash
#

function add
{
	echo "`expr $1 \+ $2`"
	# echo `expr $1 \+ $2`
	# echo $(($1 + $2))
}

sum=`add $1 $2`

echo "$1 + $2 = $sum"
```

##### 示例2：返回Linux上所有的不可登陆用户

```bash
#!/bin/bash
#

function get_users
{
	echo `cat /etc/passwd | awk -F: '/\/sbin\/nologin/{print $1}'`
}

index=1
for user in `get_users`;do
	echo "The $index user is $user"
	index=$(expr $index + 1)
done	

echo
echo "System have $index users(do not login)"
```

### 局部变量和全局变量

		要点1：Shell脚本中，默认所有变量都是全局变量；即使函数内部定义的变量，一旦函数调用后，改变了就将一直存在，直到脚本执行完毕
		要点2：定义局部变量，使用local关键字；
		要点3：函数内部，变量会自动覆盖外部变量
#### 示例1

```bash
#!/bin/bash
#

variable_1="Global Variable"

function local_func
{
	variable_2="Local Variable"
}

echo "variable_1=$variable_1"
# variable_1=Global Variable
echo "variable_2=$variable_2"
# variable_2=

local_func

echo "variable_1=$variable_1"
# variable_1=Global Variable
echo "variable_2=$variable_2"
# variable_2=Local Variable

function test_local
{
	echo "variable_2=$variable_2"
	# variable_2=Local Variable
}

test_local

```

#### 编程习惯原则

```
1、尽量在函数内部使用local关键字，将变量的作用于限制在函数内部
2、命名变量名时尽可能遵循实义性的，尽量做到见名知意
```


## Shell编程中的常用工具

### 文件查找之find命令

#### 语法格式

```bash
find [路径] [选项] [操作]
```

#### 常用选项

| 选项      | 参数 | 描述                           | 使用示例                                                     |
| --------- | ---- | ------------------------------ | ------------------------------------------------------------ |
| -name     |      | 文件名称                       | 查找/etc目录下以conf结尾的文件<br/>find /etc -name '*conf'   |
| -iname    |      |                                | 查找当前目录下文件名为aa的文件，不区分大小写<br/>find . -iname aa |
| -user     |      |                                | 查找文件属主为hdfs的所有文件<br/>find . -user hdfs           |
| -group    |      |                                | 查找文件属组为yarn的所有文件<br/>find . -group yarn          |
| -type     | f    | 文件                           | find . -type f                                               |
|           | d    | 目录                           | find . -type d                                               |
|           | c    | 字符设备文件                   | find . -type c                                               |
|           | b    | 块设备文件                     | find . -type b                                               |
|           | l    | 链接文件                       | find . -type l                                               |
|           | p    | 管道文件                       | find . -type p                                               |
| -size     | -n   | 小于n的文件                    | 查找/etc目录下小于10000字节的文件<br/>find /etc -size -10000c |
|           | +n   | 大于n的文件                    | 查找/etc目录下大于1M的文件<br/>find /etc -size +1M           |
|           | n    | 等于n的文件                    |                                                              |
| -mtime    | -n   | n天以内修改的文件              | 查找/etc目录下5天之内修改且以conf结尾的文件<br/>find /etc -mtime -5 -name '*.conf' |
|           | +n   | n天以外修改的文件              | 查找/etc目录下10天之前修改且属主为root的文件<br/>find /etc -mtime +10 -user root |
|           | n    | 正好n天修改的文件              |                                                              |
| -mmin     | -n   | n分钟以内修改的文件            | 查找/etc目录下30分钟之前修改的文件<br/>find /etc -mmin +30   |
|           | +n   | n分钟以外修改的文件            | 查找/etc目录下30分钟之内修改的目录<br/>find /etc -mmin -30 -type d |
| -mindepth | n    | 表示从n级子目录<br/>开始搜索   | 在/etc下的3级子目录开始搜索<br/>find /etc -mindepth 3        |
| -maxdepth | n    | 表示最多搜索到<br/>n-1级子目录 | 在/etc下搜索符合条件的文件,<br>但最多搜索到2级子目录<br/>find /etc -maxdepth 3 -name '*.conf' |

#### 操作

| 操作   | 说明                                                         | 使用示例                                                     |
| ------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| -print | 打印输出                                                     |                                                              |
| -exec  | 对搜索到的文件执行特定的操作，<br>格式为-exec 'command' {} \\; | 搜索/etc下的文件(非目录)，文件名以conf结尾，<br/>且大于10k，然后将其删除<br/>find ./etc/ -type f -name '*.conf'<br/> -size +10k -exec rm -f {} \\; |
|        |                                                              | 将/var/log/目录下以log结尾的文件，<br/>且更改时间在7天以上的删除<br/>find /var/log/ -name '*.log' -mtime +7 -exec rm -rf {} \\; |
|        |                                                              | 搜索条件和例子1一样，只是不删除，<br/>而是将其复制到/root/conf目录下<br/>find ./etc/ -size +10k -type f -name '*.conf' <br/>-exec cp {} /root/conf/ \\; |
| -ok    |                                                              | 和exec功能一样，只是每次操作都会给用户提示                   |

#### 逻辑运算符

| 运算符     | 描述 |
| ---------- | ---- |
| -a         | 与   |
| -o         | 或   |
| -not \| ！ | 非   |

##### 示例1：查找当前目录下，属主不是hdfs的所有文件

```bash
find . -not -user hdfs 	|	find . ! -user hdfs
```

##### 示例2：查找当前目录下，属主属于hdfs，且大小大于300字节的文件

```bash
find . -type f -a -user hdfs -a -size +300c
```

##### 示例3：查找当前目录下的属主为hdfs或者以conf结尾的普通文件

```bash
find . -type f -a -user hdfs -o -name '*.conf'
```

#### 脚本练习

##### 需求描述

```
1、提示用户输入一个目录，然后继续提示用户输入一个搜索文件的查询条件(文件名、文件大小)，然后脚本可以将符合搜索条件的文件打印出来
2、继续提示用户是拷贝或删除这些文件，如果删除，则执行删除操作，同时将删除的文件记录到一个remove.list文件中；
  如果是拷贝，则继续提示用户输入一个目标目录，然后执行拷贝动作
```

##  文本处理三剑客之grep

### grep语法格式

#### 形式一

```bash
grep [option] [pattern] [file1,file2...]
```

#### 形式二

```bash
command | grep [option] [pattern]
```

### 常用选项

| 选项 | 含义                                         |
| ---- | -------------------------------------------- |
| -v   | 显示不匹配pattern的行                        |
| -i   | 搜索时忽略大小写                             |
| -n   | 显示行号                                     |
| -E   | 支持扩展的正则表达式                         |
| -F   | 不支持正则表达式，按字符串的字面意思进行匹配 |
| -r   | 递归搜索                                     |
| -c   | 只输出匹配行的数量，不显示具体内容           |
| -w   | 匹配整词                                     |
| -x   | 匹配整行                                     |
| -l   | 只列出匹配的文件名，不显示具体匹配行内容     |

### grep和egrep的区别

```
grep默认不支持扩展正则表达式，只支持基础正则表达式

使用grep -E可以支持扩展正则表达式

使用egrep可以支持扩展正则表达式，与grep -E等价
```


## 文本处理三剑客之sed

### sed的工作模式

#### 对标准输出或文件逐行处理

#### 语法格式

##### 形式一

```bash
stdout | set [option] "pattern command"
```

##### 形式二

```bash
set [option] "pattern command" file
```

### sed的选项

| 选项 | 含义                               |
| ---- | ---------------------------------- |
| -n   | 只打印模式匹配行                   |
| -e   | 直接在命令行进行sed编辑，默认选项  |
| -f   | 编辑动作报错在文件中，指定文件执行 |
| -r   | 支持扩展正则表达式                 |
| -i   | 直接修改文件内容                   |

### sed的匹配模式

| 匹配模式                     | 含义                                           |
| ---------------------------- | ---------------------------------------------- |
| 10command                    | 匹配到第10行                                   |
| 10,20command                 | 匹配从第10行开始，到第20行结束                 |
| 10,+5command                 | 匹配从第10行开始，到第16行结束                 |
| /pattern1/command            | 匹配到pattern1行                               |
| /pattern1/,/pattern2/command | 匹配从pattern1行开始，到匹配到pattern2的行结束 |
| 10,/pattern1/command         | 匹配从第10行开始，到匹配到pattern1的行结束     |
| /pattern1/,10command         | 匹配从第pattern1的行开始，到匹配到第10行结束   |

### sed的编辑命令

| 类别 | 编辑命令     | 含义                                |
| ---- | ------------ | ----------------------------------- |
| 查询 | p            | 打印                                |
| 增加 | a            | 行后追加                            |
|      | i            | 行前追加                            |
|      | r            | 外部文件读入，行后追加              |
|      | w            | 匹配行写入外部文件                  |
| 删除 | d            | 删除                                |
| 修改 | s/old/new    | 将行内第一个old替换为new            |
|      | s/old/new/g  | 将行内全部的old替换为new            |
|      | s/old/new/2g | 将行内前2个old替换为new             |
|      | s/old/new/ig | 将行内全部的old替换为new,忽略大小写 |


### 利用sed查找文件内容

#### 示例

```bash
#!/bin/bash
#

FILE_NAME=my.cnf

function get_all_segments
{
	echo "`sed -n '/\[.*\]/p' $FILE_NAME  | sed -e 's/\[//g' -e 's/\]//g'`"
	# 思路：先匹配[.*]行，然后删除[]
}

function count_items_in_segment
{
	items=`sed -n '/\['$1'\]/,/\[.*\]/p' $FILE_NAME | grep -v "^#" | grep -v "^$" | grep -v "\[.*\]"`
	# 思路：先匹配当前行到[.*]，然后去掉注释^#,去掉空格^$，去掉最后一行[.*]
	
	index=0
	for item in $items
	do
		index=`expr $index + 1`
	done

	echo $index
}

number=0

for segment in `get_all_segments`
do
	number=`expr $number + 1`
	items_count=`count_items_in_segment $segment`
	echo "$number: $segment  $items_count"
done
```

#### my.conf
```properties
# this is read by the standalone daemon and embedded servers
[client]
port=3306
socket=/tmp/mysql.socket

#ThisSegmentForserver
[server]
innodb_buffer_pool_size=91750M
innodb_buffer_pool_instances=8

#thisisonlyforthemysqldstandalonedaemon
[mysqld]
port=3306
socket=/tmp/mysql.sock

#ThisSegmentFormysqld_safe
[mysqld_safe]
log-error=/var/log/mariadb/mariadb.log
pid-file=/var/run/mariadb/mariadb.pid

#thisisonlyforembeddedserver
[embedded]
gtid_mode=on
enforce_gtid_consistency=1

#usethisgroupforoptionsthatolderserversdon'tunderstand
[mysqld-5.5]
key_buffer_size=32M
read_buffer_size=8M
```

### 利用sed删除文件内容

#### 用法总结

```tex
1、1d
2、5,10d
3、10,+10d
4、/pattern1/d
5、/pattern1/,/pattern2/d
6、/pattern1/,20d
7、15,/pattern1/d
```
#### 示例1

```bash
# 1、删除/etc/passwd中的第15行	
	sed -i '15d' /etc/passwd
# 2、删除/etc/passwd中的第8行到第14行的所有内容	
	sed -i '8,14d' passwd
# 3、删除/etc/passwd中的不能登录的用户(筛选条件：/sbin/nologin)	
	sed -i '/\/sbin\/nologin/d' passwd
# 4、删除/etc/passwd中以mail开头的行，到以yarn开头的行的所有内容		
	sed -i '/^mail/,/^yarn/d' passwd
# 5、删除/etc/passwd中第一个不能登录的用户，到第13行的所有内容
	sed -i '/\/sbin\/nologin/,13d' passwd
# 6、删除/etc/passwd中第5行到以ftp开头的所有行的内容
	sed -i '5,/^ftp/d' passwd
# 7、删除/etc/passwd中以yarn开头的行到最后行的所有内容	
	sed -i '/^yarn/,$d' passwd
```
##### paasword

```
root:x:0:0:root:/root:/bin/bash
bin:x:1:1:bin:/bin:/sbin/nologin
daemon:x:2:2:daemon:/sbin:/sbin/nologin
adm:x:3:4:adm:/var/adm:/sbin/nologin
lp:x:4:7:lp:/var/spool/lpd:/sbin/nologin
```

#### 示例2：典型需求

```bash
# 1、删除配置文件中的所有注释行和空行
	sed -i '/[:blank:]*#/d;/^$/d' nginx.conf
# 2、在配置文件中所有不以#开头的行前面添加*符号，注意：以#开头的行不添加
	sed -i 's/^[^#]/\*&/g' nginx.conf
	# 备注：使用&反向引用匹配到行的内容
```
##### nginx.conf

```nginx
# For more information on configuration, see:
#   * Official English Documentation: http://nginx.org/en/docs/
#   * Official Russian Documentation: http://nginx.org/ru/docs/

user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log;
pid /run/nginx.pid;

# Load dynamic modules. See /usr/share/nginx/README.dynamic.
include /usr/share/nginx/modules/*.conf;

events {
    worker_connections 1024;
}

http {
    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;

    sendfile            on;
    tcp_nopush          on;
    tcp_nodelay         on;
    keepalive_timeout   65;
    types_hash_max_size 2048;

    include             /etc/nginx/mime.types;
    default_type        application/octet-stream;

    upstream www_servers {
	server 192.168.1.70:80 weight=1;
	server 192.168.1.80:80 weight=1;
    }

    # Load modular configuration files from the /etc/nginx/conf.d directory.
    # See http://nginx.org/en/docs/ngx_core_module.html#include
    # for more information.
    include /etc/nginx/conf.d/*.conf;

    server {
        listen       192.168.1.129:80 default_server;
        server_name  _;
        root         /usr/share/nginx/html;

        # Load configuration files for the default server block.
        include /etc/nginx/default.d/*.conf;

        location / {
	    proxy_pass http://www_servers;
	    proxy_set_header Host $host;
	    proxy_set_header X-Forwarded-For $remote_addr;
        }

        error_page 404 /404.html;
            location = /40x.html {
        }

        error_page 500 502 503 504 /50x.html;
            location = /50x.html {
        }
    }

# Settings for a TLS enabled server.
#
#    server {
#        listen       443 ssl http2 default_server;
#        listen       [::]:443 ssl http2 default_server;
#        server_name  _;
#        root         /usr/share/nginx/html;
#
#        ssl_certificate "/etc/pki/nginx/server.crt";
#        ssl_certificate_key "/etc/pki/nginx/private/server.key";
#        ssl_session_cache shared:SSL:1m;
#        ssl_session_timeout  10m;
#        ssl_ciphers HIGH:!aNULL:!MD5;
#        ssl_prefer_server_ciphers on;
#
#        # Load configuration files for the default server block.
#        include /etc/nginx/default.d/*.conf;
#
#        location / {
#        }
#
#        error_page 404 /404.html;
#            location = /40x.html {
#        }
#
#        error_page 500 502 503 504 /50x.html;
#            location = /50x.html {
#        }
#    }

}
```

#### sed中引入变量时注意事项

```tex
1、匹配模式中存在变量，则建议使用双引号
2、匹配模式中存在变量，如果外面使用单引号，则自定义变量也必须使用单引号
3、反向引用匹配到的内容使用&或\1
```

### 利用sed修改文件内容

#### 用法总结

```tex
1、1s/old/new/
2、5,10s/old/new/
3、10,+10s/old/new/
4、/pattern1/s/old/new/
5、/pattern1/,/pattern2/s/old/new/
6、/pattern1/,20s/old/new/
7、15,/pattern1/s/old/new/
```

#### 示例1

```bash
# 1、修改/etc/passwd中第1行中第1个root为ROOT	
	sed -i '1s/root/ROOT/' passwd
# 2、修改/etc/passwd中第5行到第10行中所有的/sbin/nologin为/bin/bash
	sed -i '5,10s/\/sbin\/nologin/\/bin\/bash/g' passwd
# 3、修改/etc/passwd中匹配到/sbin/nologin的行，将匹配到行中的login改为大写的LOGIN
	sed -i '/\/sbin\/nologin/s/login/LOGIN/g' passwd
# 4、修改/etc/passwd中从匹配到以root开头的行，到匹配到行中包含mail的所有行。修改内为将这些所有匹配到的行中的bin改为HADOOP
	sed -i '/^root/,/mail/s/bin/HADOOP/g' passwd
# 5、修改/etc/passwd中从匹配到以root开头的行，到第15行中的所有行，修改内容为将这些行中的nologin修改为SPARK
	sed -i '/^root/,15s/nologin/SPARK/g' passwd
# 6、修改/etc/passwd中从第15行开始，到匹配到以yarn开头的所有航，修改内容为将这些行中的bin换位BIN
	sed -i '15,/^yarn/s/bin/BIN/g' passwd
```
### 利用sed追加文件内容

#### 用法总结

```tex
1、a					在匹配行后面追加		
2、i					在匹配行前面追加
3、r					将文件内容追加到匹配行后面
4、w					将匹配行写入指定文件
```

#### 示例详解

##### 在匹配行后面追加

```bash
# (1)、passwd文件第10行后面追加"Add Line Behind"		
	sed -i '10a Add Line Begind' passwd
# (2)、passwd文件第10行到第20行，每一行后面都追加"Test Line Behind"
	sed -i '10,20a Test Line Behind' passwd
# (3)、passwd文件匹配到/bin/bash的行后面追加"Insert Line For /bin/bash Behind"
	sed -i '/\/bin\/bash/a Insert Line For /bin/bash Behind' passwd
```
##### 在匹配行前面追加

```bash
# (1)、passwd文件匹配到以yarn开头的行，在匹配航前面追加"Add Line Before"
	sed -i '/^yarn/i Add Line Before' passwd
# (2)、passwd文件每一行前面都追加"Insert Line Before Every Line"
	sed -i 'i Insert Line Before Every Line' passwd
```
##### 将文件内容追加到匹配行后面

```bash
# (1)、将/etc/fstab文件的内容追加到passwd文件的第20行后面
	sed -i '20r /etc/fstab' passwd
# (2)、将/etc/inittab文件内容追加到passwd文件匹配/bin/bash行的后面
	sed -i '/\/bin\/bash/r /etc/inittab' passwd
# (3)、将/etc/vconsole.conf文件内容追加到passwd文件中特定行后面，匹配以ftp开头的行，到第18行的所有行
	sed -i '//,18r /etc/vconsole.conf' passwd
```
##### 将匹配行写入指定文件

```bash
# (1)、将passwd文件匹配到/bin/bash的行追加到/tmp/sed.txt文件中
	sed -i '/\/bin\/bash/w /tmp/sed.txt' passwd
# (2)、将passwd文件从第10行开始，到匹配到hdfs开头的所有行内容追加到/tmp/sed-1.txt
	sed -i '10,/^hdfs/w /tmp/sed-1.txt' passwd
```
##  文本处理三剑客之awk

### awk的工作模式

#### awk是一个文本处理工具，通常用于处理数据并生成结果报告

#### 语法格式

##### 形式一

```bash
awk 'BEGIN{} pattern {commands} END{}' file_name
```

##### 形式二

```bash
stout | awk 'BEGIN{} pattern {commands} END{}'
```

### awk的内置变量

| 内置变量 | 全称                   | 含义                                      |
| -------- | ---------------------- | ----------------------------------------- |
| $0       |                        | 打印行所有信息                            |
| $1~$n    |                        | 打印行的第1到n个字段的信息                |
| NF       | Number Field           | 处理行的字段个数                          |
| NR       | Number Row             | 处理行的行号                              |
| FNR      | File Number Row        | 多文件处理时，每个文件单独记录行号        |
| FS       | Field Separator        | 字段分割符，不指定时默认以空格或tab键分割 |
| RS       | Row Separator          | 行分隔符，不指定时以回车分割\n            |
| OFS      | Output Filed Separator | 输出字段分隔符                            |
| ORS      | Output Row Separator   | 输出行分隔符                              |
| FILENAME |                        | 处理文件的文件名                          |
| ARGC     |                        | 命令行参数个数                            |
| ARGV     |                        | 命令行参数数组                            |

### awk格式输出之printf

#### 用法总结

##### 格式符

| 格式符 | 含义                     |
| ------ | ------------------------ |
| %s     | 打印字符串               |
| %d     | 打印10进制数             |
| %f     | 打印浮点数               |
| %x     | 打印16进制数             |
| %o     | 打印8进制数              |
| %e     | 打印数字的科学计数法格式 |
| %c     | 打印单个字符的ASCII码    |


##### 修饰符

| 修饰符 | 含义                                     |
| ------ | ---------------------------------------- |
| -      | 左对齐                                   |
| +      | 右对齐                                   |
| #      | 显示8进制在前面加0，显示16进制在前面加0x |

##### 示例1

```bash
# 1、以字符串格式打印/etc/passwd中的第7个字段，以":"作为分隔符
	awk 'BEGIN{FS=":"} {printf "%s",$7}' /etc/passwd
	
# 2、以10进制格式打印/etc/passwd中的第3个字段，以":"作为分隔符
	awk 'BEGIN{FS=":"} {printf "%d\n",$3}' /etc/passwd

# 3、以浮点数格式打印/etc/passwd中的第3个字段，以":"作为分隔符
	awk 'BEGIN{FS=":"} {printf "%0.3f\n",$3}' /etc/passwd

# 4、以16进制数格式打印/etc/passwd中的第3个字段，以":"作为分隔符
	awk 'BEGIN{FS=":"} {printf "%#x\n",$3}' /etc/passwd

# 5、以8进制数格式打印/etc/passwd中的第3个字段，以":"作为分隔符
	awk 'BEGIN{FS=":"} {printf "%#o\n",$3}' /etc/passwd

# 6、以科学计数法格式打印/etc/passwd中的第3个字段，以":"作为分隔符
	awk 'BEGIN{FS=":"} {printf "%e\n",$3}' /etc/passwd
```

### awk模式匹配两种用法

#### 方法一：RegExp

```bash
# 匹配/etc/passwd文件行中含有root字符串的所有行
	awk 'BEGIN{FS=":"}/root/{print $0}' /etc/passwd

# 匹配/etc/passwd文件行中以yarn开头的所有行
	awk 'BEGIN{FS=":"}/^yarn/{print $0}' /etc/passwd
```

#### 方法二：运算符匹配

##### 关系运算符匹配

| 运算符 | 含义             |
| ------ | ---------------- |
| <      | 小于             |
| >      | 大于             |
| <=     | 小于等于         |
| >=     | 大于等于         |
| ==     | 等于             |
| !=     | 不等于           |
| ~      | 匹配正则表达式   |
| ！~    | 不匹配正则表达式 |

```bash
# (1)、以:为分隔符，匹配/etc/passwd文件中第3个字段小于50的所有行信息
	awk 'BEGIN{FS=":"}$3<50{print $0}' /etc/passwd
	
# (2)、以:为分隔符，匹配/etc/passwd文件中第3个字段大于50的所有行信息
	awk 'BEGIN{FS=":"}$3>50{print $0}' /etc/passwd
	
# (3)、以:为分隔符，匹配/etc/passwd文件中第7个字段为/bin/bash的所有行信息
	awk 'BEGIN{FS=":"}$7=="/bin/bash"{print $0}' /etc/passwd

# (4)、以:为分隔符，匹配/etc/passwd文件中第7个字段不为/bin/bash的所有行信息
	awk 'BEGIN{FS=":"}$7!="/bin/bash"{print $0}' /etc/passwd

# (5)、以：为分隔符，匹配/etc/passwd中第3个字段包含3个以上数字的所有行信息
	awk 'BEGIN{FS=":"}$3~/[0-9]{3,}/{print $0}' /etc/passwd
```

##### 布尔运算符匹配

| 运算符 | 含义 |
| ------ | ---- |
| \|\|   | 或   |
| &&     | 与   |
| ！     | 非   |

```bash
# (1)、以:为分隔符，匹配/etc/passwd文件中包含hdfs或yarn的所有行信息
	awk 'BEGIN{FS=":"}$1=="hdfs" || $1=="yarn" {print $0}' /etc/passwd

# (2)、以:为分隔符，匹配/etc/passwd文件中第3个字段小于50并且第4个字段大于50的所有行信息
	awk 'BEGIN{FS=":"}$3<50 && $4>50 {print $0}' /etc/passwd
```

### awk中表达式的用法

#### 用法总结

| 算数运算符 | 含义                      |
| ---------- | ------------------------- |
| +          | 加                        |
| -          | 减                        |
| *          | 乘                        |
| /          | 除                        |
| %          | 取模                      |
| ^或**      | 乘方                      |
| ++x        | 在返回x变量之前，x变量加1 |
| x++        | 在返回x变量之后，x变量加1 |
| --x        | 在返回x变量之前，x变量减1 |
| x--        | 在返回x变量之后，x变量减1 |

##### 示例1

使用awk计算/etc/services中的空白行数量

```bash
awk '/^$/{sum++}END{print sum}' /etc/services
```
##### 示例2

计算学生课程分数平均值，学生课程文件内容如下：
```tex
Allen	80	90	96	98
Mike	93	98	92	91
Zhang	78	76	87	92
Jerry	86	89	68	92
Han		85	95	75	90
Li		78	88	98	100
```
```bash
awk 'BEGIN{printf "%-8s%-8s%-8s%-8s%-8s%s\n","Name","Yuwen","Math","English","Pysical","Average"}{total=$2+$3+$4+$5;AVG=total/4;printf "%-8s%-8d%-8d%-8d%-8d%0.2f\n",$1,$2,3,$4,$5,AVG}' student.txt

# 思路：在BEGIN中先打印title，在command中计算并打印每一列
```

### awk中条件及循环语句

#### 条件语句

```bash
if(条件表达式1)
	动作
else if(条件表达式2)
	动作
else
	动作
```

##### 示例1

以:为分隔符，只打印/etc/passwd中第3个字段的数值在50-100范围内的行信息

```bash
BEGIN{
	FS=":"
}

{
	if($3<50)
	{
		printf "%-20s%-25s%-5d\n","UID<50",$1,$3
	}
	else if($3>50 && $3<100)
	{
		printf "%-20s%-25s%-5d\n","50<UID<100",$1,$3
	}
	else
	{
		printf "%-20s%-25s%-5d\n","UID>100",$1,$3
	}
}
```

#### 循环条件语句

##### for循环语句

```bash
for(初始化计数器;计数器测试;计数器变更)
	动作
```

##### while循环语句

```bash
while(条件表达式)
	动作
```

##### do...while语句

```bash
do
	动作
while(条件表达式)
```

###### 示例1

计算下列每个同学的平均分数，并且只打印平均分数大于90的同学姓名和分数信息

```tex
Name	Chinese		English		Math		Physical	Average
Allen	80			90			96			98
Mike	93			98			92			91
Zhang	78			76			87			92
Jerry	86			89			68			92
Han		85			95			75			90
Li		78			88			98			100	
```

```bash
BEGIN{
	printf "%-10s%-10s%-10s%-10s%-10s%-10s\n","Name","Chinese","English","Math","Physical","Average"
}

{
	total=$2+$3+$4+$5
	avg=total/4
	if(avg>90)
	{
		printf "%-10s%-10d%-10d%-10d%-10d%-0.2f\n",$1,$2,$3,$4,$5,avg
		score_chinese+=$2
        	score_english+=$3
        	score_math+=$4
        	score_physical+=$5
	}
}

END{
	printf "%-10s%-10d%-10d%-10d%-10d\n","",score_chinese,score_english,score_math,score_physical
}
```

###### 示例2

计算1+2+3+4+...+100的和，请使用while、do while、for三种循环方式实现

```bash
BEGIN{
	while(i<=100)
	{
		sum+=i
		i++
	}
	print sum
}
```

```bash
BEGIN{
	do
	{
		sum+=i
		i++
	}while(i<=100)

	print sum
}
```

```bash
BEGIN{
	for(i=0;i<=100;i++)
	{
		sum+=i
	}
	print sum
}
```

### awk中的常用选项

| 选项 | 含义            |
| ---- | --------------- |
| -v   | 定义或引用变量  |
| -f   | 指定awk命令文件 |
| -F   | 指定分隔符      |
| -V   | 查看awk的版本号 |

### awk中数组的用法

#### Shell中数组的用法

```bash
array=("Allen" "Mike" "Messi" "Jerry" "Hanmeimei" "Wang")
打印元素：			echo ${array[2]}
打印元素个数:		   echo ${#array[@]}
打印元素长度：		  echo ${#array[3]}
给元素赋值：		   array[3]="Li"
删除元素：		    unset array[2];unset array
分片访问：			echo ${array[@]:1:3}
元素内容替换：		  ${array[@]/e/E}	只替换第一个e;${array[@]//e/E}	替换所有的e
数组的遍历：
                    for a in ${array[@]}
                    do
                        echo $a
                    done
```

#### awk中数组的用法

在awk中，使用数组时，不仅可以使用1.2..n作为数组下标，也可以使用字符串作为数组下标

##### 形式一

当使用1.2.3..n时，直接使用array[2]访问元素;需要遍历数组时，使用以下形式

```bash
str="Allen Jerry Mike Tracy Jordan Kobe Garnet"
split(str,array)
for(i=1;i<=length(array);i++)
	print array[i]
```

##### 形式二

当使用字符串作为数组下标时，需要使用array[str]形式访问元素;遍历数组时，使用以下形式

```bash
array["var1"]="Jin"
array["var2"]="Hao"
array["var3"]="Fang"

for(a in array)
	print array[a]
```
##### 典型应用

###### 示例1

统计主机上所有的TCP连接状态数，按照每个TCP状态分类

```bash
netstat -an | grep tcp | awk '{array[$6]++}END{for(a in array) print a,array[a]}'
```

###### 示例2（字典）

计算横向数据总和，计算纵向数据总和

```tex
allen	80	90	87	91	348
mike	78	86	93	96	256
Kobe	66	92	82	78	232
Jerry	98	74	66	54  356
Wang	87	21	100	43  322
		234 342 451 456 342
```
```bash
BEGIN {
	printf "%-10s%-10s%-10s%-10s%-10s%-10s\n","Name","Yuwen","Math","English","Physical","Total"
}

{
	total=$2+$3+$4+$5
	yuwen_sum+=$2
	math_sum+=$3
	eng_sum+=$4
	phy_sum+=$5
	printf "%-10s%-10d%-10d%-10d%-10d%-10d\n",$1,$2,$3,$4,$5,total
}
END {
	printf "%-10s%-10d%-10d%-10d%-10d\n","",yuwen_sum,math_sum,eng_sum,phy_sum
}

```

### 复杂的awk处理生产数据的示例

#### 需求描述

利用awk处理日志，并生成结果报告

#### 生成数据脚本

```bash
#!/bin/bash
#

function create_random()
{
	min=$1
	max=$(($2-$min+1))
	num=$(date +%s%N)
	echo $(($num%$max+$min))
}

INDEX=1

while true
do
	for user in Allen Mike Jerry Tracy Hanmeimei Lilei
	do
		COUNT=$RANDOM
		NUM1=`create_random 1 $COUNT`
		NUM2=`expr $COUNT - $NUM1`		
		echo "`date '+%Y-%m-%d %H:%M:%S'` $INDEX Batches:$user INSERT $COUNT DATA INTO database.table 'test',Insert $NUM1 Records Successfully,Failed Insert $NUM2 Records" >> /root/db.log.`date +%Y%m%d`
		INDEX=`expr $INDEX + 1`
	done
done
```

#### 数据格式如下

```tex
2019-01-29 00:58:30 1 Batches: user allen insert 22498 records into database:product table:detail, insert 20771 records successfully,failed 1727 records
2019-01-29 00:58:30 2 Batches: user mike insert 29378 records into database:product table:detail, insert 21426 records successfully,failed 7952 records
2019-01-29 00:58:30 3 Batches: user jerry insert 22779 records into database:product table:detail, insert 9397 records successfully,failed 13382 records
2019-01-29 00:58:30 4 Batches: user tracy insert 25232 records into database:product table:detail, insert 21255 records successfully,failed 3977 records
```
#### 统计每个人分别插入多少数据，多少成功，多少失败，并且要格式化输出，加上标题

##### 输出结果

```tex
User      Total               Sucess              Failed              
tracy     7472277             3945659             3526618             
allen     7390330             3597157             3793173             
mike      7226579             3679395             3547184
		  21384945
```
```bash
BEGIN {
	printf "%-10s%-20s%-20s%-20s\n","User","Total","Sucess","Failed"
}

{
	TOTAL[$6]+=$8
	SUCCESS[$6]+=$14
	FAIL[$6]+=$17
}

END {
	for(u in SUCCESS)
	{
		total+=TOTAL[u]
		success+=SUCCESS[u]
		fail+=FAIL[u]
		printf "%-10s%-20d%-20d%-20d\n",u,TOTAL[u],SUCCESS[u],FAIL[u]
	}

	printf "%-10s%-20d%-20d%-20d\n","",total,success,fail
}
```

#### 查找丢失数据的现象，也就是成功+失败的记录数，不等于一共插入的记录数。找出这些数据并显示行号和对应行的日志信息

##### 输出结果

```bash
awk '{if($8!=$14+$17) print NR,$0}' db.log.201901
```

##  Shell脚本操作数据库实战

### 测试数据导入
