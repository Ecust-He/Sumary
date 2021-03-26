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

#### 例子

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

##### 方法一
```bash
${#string}
```

##### 方法二
```bash
expr length $string
```

##### 例子

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
##### 例子
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
##### 例子
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

##### 例子

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

#### 两种方法

##### 方法一

```bash
`command`
```
##### 方法二

```bash
$(command)
```

#### 例子1：获取系统得所有用户并输出	

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
#### 例子2：根据系统时间计算今年或明年

```bash
echo "This is $(date +%Y) year"
echo "This is $(($(date +%Y) + 1)) year"
# Bash中的(())用来支持算术表达式

# C语言规则运算
# $((exp))，exp为符合C语言规则的运算表达式
```
#### 例3：根据系统时间获取今年还剩下多少星期，已经过了多少星期

```bash
date +%j
echo "This year have passed $(date +%j) days"
echo "This year have passed $(($(date +%j)/7)) weeks"

echo "There is $((365 - $(date +%j))) days before new year"
echo "There is $(((365 - $(date +%j))/7)) days before new year"
```
#### 例4：判定nginx进程是否存在，若不存在则自动拉起该进程

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

##### 测试nginx是否在运行

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

##### 两数字相加

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

##### 返回Linux上所有的不可登陆用户

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
