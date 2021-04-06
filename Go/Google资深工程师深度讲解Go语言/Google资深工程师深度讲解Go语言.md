[TOC]



## 1 基础语法

### 变量定义

#### 使用var关键字

##### 使用场景

- 可放在函数内或直接放在包内
- 使用var（）集中定义变量

```go
# 声明变量，不赋初值（默认ZeroValue）
func variableZeroValue() {
	var a int
	var s string
	fmt.Printf("%d %q\n", a, s)
}

# 声明变量，同时赋初值
func variableInitialValue() {
	var a, b int = 3, 4
	var s string = "abc"
	fmt.Println(a, b, s)
}

# 不显示指明变量类型，由编译器自动决定类型（类型推断）
func variableTypeDeduction() {
	var a, b, c, s = 3, 4, true, "def"
	fmt.Println(a, b, c, s)
}
```

#### 短语句赋值

##### 使用场景

- 只能在函数内使用

```go
func variableShorter() {
   a, b, c, s := 3, 4, true, "def"
   b = 5
   fmt.Println(a, b, c, s)
}
```

### 内建变量类型

#### 