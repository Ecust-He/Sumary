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

###### 