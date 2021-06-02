[TOC]

## 1 Python入门导学

### Python的特性

Simple is better than complex.

Now is better than never. Although never is often better than right now

- 易于上手，难于精通
- Python既有动态脚本语言的特点，又有面向对象的特性 

### Python的缺点

- 运行效率和开发效率不可兼得

### Python能做什么

- 爬虫
- 大数据与数据分析
- 自动化运维与自动化测试
- Web开发（Flask、Django）
- 机器学习（TensorFlow）
- 胶水语言（混合其他语言编程）

## 2 基本数据类型

### 什么是代码，什么是写代码

- 代码是现实世界事物在计算机世界中的映射
- 写代码是将现实世界中的事物用计算机语言来描述

### 基本数据类型

- type(1)查看数据类型

#### 数字（Number）

##### 整型（int）

###### 进制

- 二进制（0b）
- 八进制（0o）
- 十进制
- 十六进制（0x）

##### 浮点数（float）

##### 复数（complex）

##### 布尔类型（bool）

- 数字非0、字符串非空、列表不为空、字典不为空都表示True
- None表示False

#### 字符串（str）

- 不可变

##### 单引号

- 单引号内可以使用双引号
- 单引号内使用单引号需要转义

##### 双引号

##### 三引号

- 多行字符串

~~~python
```
hello
world
```
~~~

##### 原始字符串

- 特殊字符串不转义

```python
r''
```

##### 转移字符串

```python
'\nhello\nworld\n'
```

##### 字符串操作

###### 字符串拼接：+

```python
'hello' + 'world'
```

###### 字符串repeat：*

```python
'hello'*3
```

###### 字符串查找：数组索引

- 下标从0开始
- -1代表最后一个字符

```python
'hello world'[0]
```

###### 字符串截取：slice操作

```python
'hello world'[0:2]
```

## 3 python中组的概念与定义

### list列表

#### 定义

#### 操作

- 数组索引
- 切片截取
- 数组合并（+）
- 数组repeat（*）

### tuple元组

- 不可变类型
- type((1))是int类型，type((1,))是tuple类型，type(())是tuple类型

### 序列总结

str、list、tuple

- 数组索引
- 切片操作
- +、*操作
- 成员运算符（in、not in）
- 内置函数len、max、min等

### set集合

- 无序、不重复
- 没有索引和切片操作

#### 操作

- len
- 成员运算符(in、not in)
- 差集（-）、并集（|）、交集（&）
- 内置函数difference、intersection、union等

### dict字典

- key值必须是不可变类型
- type({})是dict类型

## 4 变量

### 命名规则

- 变量名只能是字母、数字和下滑线的组合
- 首字母不能是数字
- 关键字不能作为变量

### 值类型与引用类型

#### 值类型

- 不可变
- int、bool、str、tuple

#### 引用类型

- 可变
- list、set、dict

## 5 运算符

### 算数运算符

- \+
- \-
- \*
- /
- %
- //（保留整数部分）
- **（指数运算，例如2\*\*5=32）

### 赋值运算符

### 比较运算符

- ==
- !=
- <
- <=
- \>
- \>=

#### 不只是数字才能做比较运算

- bool类型可以与数字相加（例如：1 + true值为2）
- 字符串比较（依次比较字符的ASCII码）
- 列表、元组比较（依次比较列表中元素的大小）

### 逻辑运算符

- and
- or
- not

#### 逻辑运算符并不一定返回bool类型

```python
[1] or []
#[1]

'a' and 'b'
# 'b'
'a' or 'b'
# 'a'

0 and 1
# 0
1 and 0
# 0
1 and 2
# 2
0 or 1
# 1
1 or 2
# 1
```



#### 其他数据类型与bool类型之间的转换

- int整型、float浮点型：0被认为是False,非0表示True
- str字符串：空字符串被认为是False
- tuple元组：空元组()被认为是False
- list列表：空列表[]被认为是False
- dict字典：空字典{}被认为是False

### 成员运算符

- in
- not in

**返回值是bool类型**

#### 使用场景

- 序列（str、tuple、list）
- set集合
- dict字典

### 身份运算符

- is
- is not

#### ==与is的区别

- ==比较【值】是否相等
- is 比较【内存地址】是否相等

```python
a = 1
b = 1
a == b
# True
a is b
# True

a = 1
b = 1.0
a == b
# True
a is b
# False

a = {1,2,3}
b = {1,3,2}
a == b
# True
a is b
# False

a = (1,2,3)
b = (1,3,2)
a == b
# False
a is b
# False
```

#### 如何判断变量的值、身份、类型

- 值：==
- 身份：id
- 类型：type（推荐使用isinstance）

```python
a = 'aa'
type(a) == str
# True

isinstance(a, str)
# True

isinstance(a, (str,int))
# True
```

**对象的三个特征id、value、type**

### 位运算符

- &
- |
- ^（按位异或）
- ~（按位取反）
- \>\>
- <<

## 6 表达式

表达式(Expression)是运算符（Operator）和操作数（Operand）所构成的序列

### 优先级

## 7 流程控制

### 条件控制语句

```python
if condition :
	print True
elif
	print False
```

- elif的优点：替换switch分支

#### 嵌套分支

#### 代码块

### 循环控制语句

#### while循环

##### 使用场景

- 递归场景
- while与else一起使用

#### for循环

##### 使用场景

- 遍历序列（str、list、tuple）、set集合、dict字典
- for与else一起使用
- for与range一起使用

```python
nums = [1, 2, 3]
for num in nums:
    print(num, end=' | ')
# 1 | 2 | 3 |     

for num in range(0, 10, 2):
    print(num)
```

## 8 工程组织结构

### 包（package）

- 文件夹名称
- \__init__.py文件

### 模块（module）

.py文件名称

#### import导入模块

#### from import导入变量

- \*是指所有
- as 可以对变量起别名
- **可以使用\__all__=[]指定需要导出的变量**

```python
from module1 import method1 
# from module1 import method1 as m1
# import * from module1
```

#### \__init__.py的用法

- 其他模块导入该包或该包下的模块时首先被调用
- 可以使用\__all__=[]指定需要导出的模块

#### 常见错误

- 避免循环导入

#### 模块内置变量

- dir函数返回模块内所有变量
- \__name__：模块全限定名称
- \__package__：包名称
- \__doc__：模块的注释
- \__file__：跟执行入口文件时的相对路径有关

#### 入口文件

- 你运行哪个文件哪个文件就是入口文件，其他的都是普通模块
- \__name__：\__main__

#### \__name__的经典应用

make a script importable and excutable

**入口文件会执行main方法**

```python
if __name__ == '__main__':
	pass
```

##### 使用场景

- 既可以作为普通的模块，又可以作为程序的入口

#### 相对导入和绝对导入

- 与入口文件的位置有关
- 入口文件不允许使用相对导入

## 9 函数

### 特点

- 功能性
- 隐藏细节
- 避免编写重复的代码

### 返回多个值

- 返回值类型是tuple

#### 序列解包

使用多个变量接收多个返回值

```python
a,b,c = 1,2,3
d = 1,2,3
print(type(d))
# tuple
```

#### 链式赋值

```python
a=b=c=1
```

### 参数类型

#### 必须参数

```python
def add(x, y):
	return x + y
```

#### 关键参数

##### 特点

- 不需要考虑参数的顺序
- 增加代码的可读性

```python
add(y=1, x=2)
```

#### 默认参数

- 必须放在必须参数后面

```python
def add(x=1, y=2):
	return x + y

add()
# 3
```

#### 可变参数

- \*将tuple中元素平铺开
- 必须放在必须参数后面

```python
def demo(*param):
    print(param)
    print(type(param))
    
demo(1,2,3,4,5,6)
# (1,2,3,4,5,6)
#tuple

a=(1,2,3,4,5,6)
demo(*a)
# (1,2,3,4,5,6)
#tuple
```

#### 关键字可变参数

- \*将dict中元素平铺开
- 必须放在必须参数后面

```python
    def demo(**param):
        print(param)
        print(type(param))

    demo(a=1, b=2)
    
    #m = {'a': 1, 'b': 2}
    #demo(**m)
    
    #{'a': 1, 'b': 2}
	#<class 'dict'>
```

### 变量的作用域

#### 作用域链

#### global关键字

**可以将函数内部的局部变量提升为全局变量**

```python
def demo():
    global c
    c = 10
    
demo()
print(c)
```

## 10 面向对象

### 类的定义

- 类的最基本作用：封装
- 建议：对于类单独使用模块定义

### 函数与方法的区别

- 函数：面向过程
- 方法：面向对象、设计层面

### 类与对象

类：类是现实世界或思维世界的实体在计算机中的反映，它将这些数据以及这些数据上的操作封装在一起

对象：类的实例

### 构造函数

- 构造函数默认return None, 且不可修改

```python
class Human:
    pass

# 继承
class Student(Human):
    sum = 0 # 类变量
    name = ''
    age = 0
    __salary = 0 # 私有变量

    # 构造方法
    def __init__(self, name, age):
        self.name = name
        self.age = age
        # self.__class__.sum += 1

    # 实例方法
    def to_string(self):
        print("name = ", self.name)
        print("age = ", self.age)
        print("sum = ", Student.sum)

    # 类方法   
    @classmethod
    def add_one(cls):
        cls.sum += 1
        
    # 静态方法
    @staticmethod
    def add(x, y):
        return x + y    

if __name__ == '__main__':
    stu = Student('kris', 20)
    Student.add_one()
    stu.to_string()
    print(Student.add(1, 2))
```

### 类变量与实例变量

#### 查找顺序

**实例变量 -> 类变量 -> 父类**

### 方法

#### 实例方法

- 显胜于隐：显示声明在方法的第一个入参

##### 在实例方法中访问实例变量和类变量

#### 类方法

- 操作和对象无关的方法，建议使用类方法
- 可以被类和对象调用

#### 静态方法

- 不建议使用静态方法，与面向对象无关，与普通函数无异
- 可以被类和对象调用

### 成员可见性

- 双下划线表示私有的
- 注意给对象动态新增实例变量

### 继承

- 支持多继承
- 子类调用父类的方法：使用super关键字

## 11 正则表达式

正则表达式是一个特殊的字符序列，一个字符串是否与我们所设定的字符序列相匹配。

- 快速检索文本、实现一些替换文本的操作

#### 字符

##### 元字符

| 元字符 | 描述                          |
| ------ | ----------------------------- |
| \d     | 匹配一个数字字符，等价于[0-9] |
| \w     | 匹配单词字符                  |
| \s     | 匹配空白字符                  |

##### 普通字符

一般作为定界使用

#### 字符集

```python
a = 'abc, afc, acd, acc, adc'
result = re.findall('a[cf]c', a)
#['afc', 'acc']

result = re.findall('a[c-f]c', a)
#['afc', 'acc', 'adc']
```

##### 概括字符集

.**可以匹配到除换行符以外任意字符**

#### 数量词

```
a = 'python 11 java89php'
result = re.findall('[a-z]+', a)
#result = re.findall('[a-z]{1,}', a)
print(result)
# ['python', 'java', 'php']
```

##### 匹配0次1次或多次

| 字符 | 描述              |
| ---- | ----------------- |
| ？   | 匹配0次或1次      |
| +    | 至少匹配一次      |
| *    | 匹配0次或任意多次 |

```python
a = 'pytho0python1pythonn2'
result = re.findall('python*', a)
#['pytho', 'python', 'pythonn']

result = re.findall('python+', a)
#['python', 'pythonn']

result = re.findall('python?', a)
#['pytho', 'python', 'python']
```

#### 贪婪与非贪婪

**默认是贪婪匹配**

```python
a = 'python 11 java89php'
# 贪婪模式
result = re.findall('[a-z]{3,6}', a)
print(result)
# ['python', 'java', 'php']

# 非贪婪模式
result = re.findall('[a-z]{3,6}?', a)
print(result)
# ['pyt', 'hon', 'jav', 'php']
```

#### 边界匹配符

| 字符 | 描述       |
| ---- | ---------- |
| ^    | 开始匹配符 |
| $    | 结尾匹配符 |

#### 组

```python
a = 'pythonpythonpythonpython'
result = re.findall('(python){3}', a)
# ['python']
```

#### 匹配模式参数

```python
a = 'PythonC#JavaPhp'
# 忽略字母大小写
result = re.findall('c#', a, re.I)
# ['C#']
```

#### 正则替换

```python
a = 'PythonC#JavaC#PhpC#'
# count 为0（默认值）时，表示替换到所有匹配到的字符串
result = re.sub('C#', "Go", a, count=0)
# PythonGoJavaGoPhpGo

result = re.sub('C#', "Go", a, count=1)
# PythonGoJavaC#PhpC#

def convert(value):
    # 获取匹配到的字符
    matched = value.group()
    return "[" + matched + "]"
# 对匹配到的字符串进行自定义处理
result = re.sub('C#', convert, a, count=0)
# Python[C#]Java[C#]Php[C#]
```

##### 把函数作为参数传递

```python
def convert(value):
    matched = value.group()
    if int(matched) >= 6:
        return '9'
    else:
        return '0'
    
s = 'A81272C23D2983'
result = re.sub('\d', convert, s)
# A90090C00D0990
```

#### search与match函数

```python
s = 'A81272C23D2983'
# 从字符串的开头开始匹配
result = re.match('\d', s)
# None

# 匹配到第一个后立即返回匹配结果
result = re.search('\d', s)
print(result.group())
# 8
```

#### group分组

```python
s = 'life is short, i use python'
result = re.findall('life(.*)python', s)
print(result)
# [' is short, i use ']

result = re.search('life(.*)python', s)
print(result.group(1))
# [' is short, i use ']
```

## 12 JSON

是一种轻量级数据交换格式，数据格式主要包括JSON Object和JSON Array

- 易于阅读
- 易于解析
- 网络传输效率高
- 跨语言交换

### JSON字符串

符合JSON格式的字符串叫JSON字符串

### 序列化

```python
student = {"name": "kris", "age": 20}
json_str = json.dumps(student)
print(type(json_str))
<class 'str'>
```

### 反序列化

将JSON字符串转换为对象的过程

```python
json_str = '{"name":"kris", "age":20}'
student = json.loads(json_str)
print(type(student))
# <class 'dict'>
print(student['name'])
# kris
```

## 13 Python的高级语法和用法

### 枚举

#### 枚举其实是一个类

```python
class VIP(Enum):
    YELLOW = 1
    GREEN = 2
    BLUE = 3
    
print(VIP.GREEN.name) # 枚举名称
#GREEN
print(VIP.GREEN.value) # 枚举值
# 2
```

#### 枚举和普通类相比有什么优势

- 枚举的不可变性
- 枚举的不可重复性

#### 枚举的比较运算

#### 枚举转换

```python
print(VIP(2))
# VIP.GREEN
```

#### 枚举小结

```python
@unique
class VIP(IntEnum):
    YELLOW = 1
    GREEN = 2
    BLUE = 3
```

### 装饰器

作用：简化代码，无需修改代码也可以实现某种功能。

- 符合开闭原则

#### 无参函数

```python
def decorator(func):
    def wrapper():
        print(time.time())
        func()
 return wrapper

@decorator
def f1():
    print('this is a function.')
    
f1()
# 1621778760.9516318
# this is a function.
```

#### 带输入参数的函数

```python
def decorator(func):
    def wrapper(*args):
        print(time.time())
        func(*args)
    return wrapper

@decorator
def f1(a, b):
    print('this is a function.')
    print(a + b)
    
f(1,2)    
# 1621779209.5720232
# this is a function.
# 3
```

## 14 函数式编程

### 一切皆对象

- 函数可以作为另一个函数的参数
- 函数可以作为另一个函数的返回值

### 闭包

- 定义：函数 + 环境变量
- 保存现场
- 函数式编程思想的一种体现

#### 闭包的经典误区

### lambda表达式

```python
lambda：param_list: expression

def add(x, y):
    return x + y

f = lambda: x,y: x +y
print(f(1,2))
# 3
```

### 三元表达式

```python
条件为真时返回的结果 if 条件 else 条件为假时返回的结果

print(1 if 2 > 3 else 0)
# 0
```

### 高阶函数

```python
res = list(map(lambda x: x * x, [1, 2, 3]))
# [1, 4, 9]
res = list(filter(lambda x: x > 2, [1, 2, 3]))
# [3]
res = reduce(lambda x, y: x + y, [1, 2, 3], 10)
# 16
```

#### map

映射

#### reduce

聚合

#### filter

过滤

### 命令式编程VS函数式编程

## 15 原生爬虫

### 正则分析html

```python
import re
from urllib import request
class Spider():
    '''
     This is a class
    '''
    url = 'https://www.panda.tv/cate/kingglory'
    root_pattern = '<div class="video-info">([\s\S]*?)</div>'
    name_pattern = '</i>([\s\S]*?)</span>'
    number_pattern = '<span class="video-number">([\s\S]*?)</span>'

    def __fetch_content(self):
        '''
        ''' 
        r = request.urlopen(Spider.url) # This is a HTTP request 

        # This is ....
        htmls = r.read() 
        htmls = str(htmls, encoding='utf-8')
        return htmls

    def __analysis(self, htmls):
        root_html = re.findall(Spider.root_pattern, htmls)

        anchors = []
        for html in root_html:
            name = re.findall(Spider.name_pattern, html)
            number = re.findall(Spider.number_pattern, html)
            anchor = {'name':name, 'number':number}
            anchors.append(anchor)

        return anchors
    
    def __refine(self, anchors):
        l =  lambda anchor: {
            'name':anchor['name'][0].strip(),
            'number':anchor['number'][0]
            }
        return map(l, anchors)

    def __sort(self, anchors):
        #filter
        anchors = sorted(anchors, key=self.__sort_seed,reverse=True)
        return anchors

    def __sort_seed(self, anchor):
        r = re.findall('\d*', anchor['number'])
        number = float(r[0])
        if '万' in anchor['number']:
            number *= 10000
        return number

    def __show(self, anchors):
        for rank in range(0, len(anchors)):
            print('rank ' + str(rank + 1)
            + '  : '+ anchors[rank]['name']
            + '   ' + anchors[rank]['number'])

    def go(self):
        htmls = self.__fetch_content()
        anchors = self.__analysis(htmls)
        anchors = list(self.__refine(anchors))
        anchors = self.__sort(anchors)
        self.__show(anchors)

spider = Spider()
spider.go()
```

## 16 Pythonic与python杂记

### 用字典映射代替switch case语句

```python
week = {
    0: "Sunday",
    1: "Monday"
}
day = 0
dayName = week[day]
# dayName = week.get(day, 'unkown')
# Sunday
```

### 列表推导式

```python
a = [1, 2, 3]
# b = list(map(lambda x: x*x, list(filter(lambda x: x > 2, a))))
b = [i*i for i in a if i > 2]
#[9]
```

#### 字典如何编写列表推导式

```python
a = {
    'name': 'kris',
    'age': 30
}
b = [key for key, value in a.items()]
print(b)
```

### iterator与generator

#### iterator

```python
class BookCollection:

    def __init__(self):
        self.data = ['红楼梦', '水浒传', '西游记']
        self.cur = 0

    def __iter__(self):
        return self

    def __next__(self):
        if self.cur >= len(self.data):
            raise StopIteration()
        value = self.data[self.cur]
        self.cur += 1
        return value
    
    for book in BookCollection():
        print(book)
```

#### generator

```python
def generator(max_value):
    n = 0
    while n <= max_value:
        n += 1
        yield n
    g = generator(100)
    
    print(next(g))
    # 1
    print(next(g))        
    # 2
```

### None

空

```python
print(type(None))
# NoneType
```

#### 如何判空？

```python
if not a:
    print 'a is empty.'

if a is None:
    print 'a is None.'
```

#### 对象存在不一定是True

- None表示不存在
- False表示假

### 装饰器的副作用

### f关键字做字符串拼接

```python
name = 'kris'
print(f'name is {name}')
```