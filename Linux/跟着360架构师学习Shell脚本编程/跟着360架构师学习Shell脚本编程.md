[TOC]



## Shell入门



## 条件判断



## 循环控制



## 变量的高级用法

### 变量替换和测试

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

