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

#### 复数

```go
# 欧拉公式
func euler() {
   fmt.Printf("%.3f\n",
      cmplx.Exp(1i*math.Pi)+1)
}
```

#### 强制类型转换

类型转换是强制的

```go
func triangle() {
	var a, b int = 3, 4
	fmt.Println(calcTriangle(a, b))
}

func calcTriangle(a, b int) int {
	var c int
	c = int(math.Sqrt(float64(a*a + b*b)))
	return c
}
```

### 常量与枚举

#### 常量定义

常量可以作为任意类型使用，不需要强制转换

```go
func consts() {
	const (
		filename = "abc.txt"
		a, b     = 3, 4
	)
	var c int
	c = int(math.Sqrt(a*a + b*b))
	fmt.Println(filename, c)
}
```

#### 枚举定义

```go
func enums() {
   const (
      cpp = iota # 自增
      _
      python
      golang
      javascript
   )

   const (
      b = 1 << (10 * iota)
      kb
      mb
      gb
      tb
      pb
   )

   fmt.Println(cpp, javascript, python, golang)
   fmt.Println(b, kb, mb, gb, tb, pb)
}
```

### 流程控制

#### 条件语句

if条件语句可以进行赋值操作，赋值变量作用域在if语句内

```go
const filename = "abc.txt"
if contents, err := ioutil.ReadFile(filename); err != nil {
    fmt.Println(err)
} else {
    fmt.Printf("%s\n", contents)
}
```

#### 分支语句

- case语句不需要加break
- case语句后面可以不加表达式

```go
func grade(score int) string {
	g := ""
	switch {
	case score < 0 || score > 100:
		panic(fmt.Sprintf(
			"Wrong score: %d", score))
	case score < 60:
		g = "F"
	case score < 80:
		g = "C"
	case score < 90:
		g = "B"
	case score <= 100:
		g = "A"
	}
	return g
}
```

#### 循环语句

- for循环语句可以省略初始条件、结束条件、递增条件

```go
# 十进制数转换成二进制数
func convertToBin(n int) string {
	result := ""
	for ; n > 0; n /= 2 {
		lsb := n % 2
		result = strconv.Itoa(lsb) + result
	}
	return result
}
```

```go
# 打印文件内容
func printFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	printFileContents(file)
}

func printFileContents(reader io.Reader) {
	scanner := bufio.NewScanner(reader)

    # for循环省略初始条件和递增条件，与while等价
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

# 死循环(省略初始条件、结束条件和递增条件)
func forever() {
	for {
		fmt.Println("abc")
	}
}
```

### 函数

- 函数可以有多个返回值
- 函数作为入参
- 函数作为返回值
- 可变参数列表
- 没有默认参数、可选参数

```go
func eval(a, b int, op string) (int, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		q, _ := div(a, b)
		return q, nil
	default:
		return 0, fmt.Errorf(
			"unsupported operation: %s", op)
	}
}

func div(a, b int) (q, r int) {
	return a / b, a % b
}
```

```go
# 函数式编程
func apply(op func(int, int) int, a, b int) int {
	return op(a, b)
}

fmt.Println("pow(3, 4) is:", apply(
    func(a int, b int) int {
        return int(math.Pow(
            float64(a), float64(b)))
    }, 3, 4))
```

```go
# 求和
func sum(numbers ...int) int {
	s := 0
	for i := range numbers {
		s += numbers[i]
	}
	return s
}
```

#### 建议

- 函数返回时，显示指定函数返回值名称，函数被调用时会自动提示
- 函数返回时，使用return关键字显示返回（多个）结果

### 指针

- 指针不能运算
- 参数传递时，Go语言只有值传递

## 2 内建容器

### 数组

#### 数组定义

- 数字定义在变量前

```go
# 声明数组 [0 0 0 0 0]
var arr1 [5]int

# 声明并赋值 [1 3 5]
arr2 := [3]int{1, 3, 5}

# 不指定数组大小，声明并赋值
arr3 := [...]int{2, 4, 6, 8, 10}

var grid [4][5]int
```

#### 数组遍历

- for与range结合使用
- 可以通过_省略变量

```go
# 打印数组
func printArray(arr [5]int) {
	arr[0] = 100
	for i, v := range arr {
		fmt.Println(i, v)
	}
}
```

#### 值类型

- 很少使用

### 切片

#### 概念

- slice本身没有数据，只是对底层array的一个view
- slice可以再slice
- slice支持向后扩展，但不可以向前扩展
- 切片操作时，qian不可以超越len(s)，向后扩展时hou不可超越底层数组的cap(s)

#### 实现

  <img src="Slice的实现.png" style="zoom:80%;" />

```go
func updateSlice(s []int) {
	s[0] = 100
}

func main() {
	arr := [...]int{0, 1, 2, 3, 4, 5, 6, 7}
    
    # arr[2:6] = [2 3 4 5]
    # arr[:6] = [0 1 2 3 4 5]
  
	s1 := arr[2:]
	# s1 = [2 3 4 5 6 7]
	s2 := arr[:]
	# s2 = [0 1 2 3 4 5 6 7]

	fmt.Println("After updateSlice(s1)")
	updateSlice(s1)
	fmt.Println(s1)
    # [100 3 4 5 6 7]
	fmt.Println(arr)
    # [0 1 100 3 4 5 6 7]

	fmt.Println("After updateSlice(s2)")
	updateSlice(s2)
	fmt.Println(s2)
    # [100 1 100 3 4 5 6 7]
	fmt.Println(arr)
    # [100 1 100 3 4 5 6 7]

	fmt.Println("Reslice")
    # [100 1 100 3 4 5 6 7]
	s2 = s2[:5]
    # [100 1 100 3 4]
	s2 = s2[2:]
	# [100 3 4]

	fmt.Println("Extending slice")
	arr[0], arr[2] = 0, 2
	fmt.Println("arr =", arr)
    # arr = [0 1 2 3 4 5 6 7]
	s1 = arr[2:6]
    # s1=[2 3 4 5], len(s1)=4, cap(s1)=6
	s2 = s1[3:5] // [s1[3], s1[4]]
	# s2=[5 6], len(s2)=2, cap(s2)=3
    # slice可以向后扩展，但不可以向前扩展

	s3 := append(s2, 10)
    # [5 6 10]
	s4 := append(s3, 11)
    # [5 6 10 11]
	s5 := append(s4, 12)
    # [5 6 10 11 12]
	// s4 and s5 no longer view arr.
	# arr = [0 1 2 3 4 5 6 10]

	// Uncomment to run sliceOps demo.
	// If we see undefined: sliceOps
	// please try go run slices.go sliceops.go
	fmt.Println("Uncomment to see sliceOps demo")
	// sliceOps()
}
```

#### 操作

- 添加元素时如果超越cap，系统会重新分配更大的底层数组
- 由于值传递的关系，添加元素时必须接收append的返回值

```go
func printSlice(s []int) {
	fmt.Printf("%v, len=%d, cap=%d\n",
		s, len(s), cap(s))
}

func sliceOps() {
	fmt.Println("Creating slice")
	var s []int // Zero value for slice is nil

    # 添加元素
	for i := 0; i < 100; i++ {
		printSlice(s)
		s = append(s, 2*i+1)
	}
	fmt.Println(s)
	# cap *2自动扩容
    
	s1 := []int{2, 4, 6, 8}
    # [2 4 6 8], len=4, cap=4

    # 声明并赋值
	s2 := make([]int, 16)
    # [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0], len=16, cap=16
	s3 := make([]int, 10, 32)
	# [0 0 0 0 0 0 0 0 0 0], len=10, cap=32

    # 拷贝元素
	fmt.Println("Copying slice")
	copy(s2, s1)
	printSlice(s2)
    # [2 4 6 8 0 0 0 0 0 0 0 0 0 0 0 0], len=16, cap=16

    # 删除元素
	fmt.Println("Deleting elements from slice")
	s2 = append(s2[:3], s2[4:]...)
	printSlice(s2)
    # [2 4 6 0 0 0 0 0 0 0 0 0 0 0 0], len=15, cap=16

	fmt.Println("Popping from front")
	front := s2[0]
	s2 = s2[1:]

	fmt.Println(front)
	printSlice(s2)

	fmt.Println("Popping from back")
	tail := s2[len(s2)-1]
	s2 = s2[:len(s2)-1]

	fmt.Println(tail)
	printSlice(s2)
}
```

#### 引用类型

### Map

- Map的key使用hash表存储，必须可比较

#### 操作

```go
func main() {
    # 声明并赋值
	m := map[string]string{
		"name":    "ccmouse",
		"course":  "golang",
		"site":    "imooc",
		"quality": "notbad",
	}

	m2 := make(map[string]int) // m2 == empty map

	var m3 map[string]int // m3 == nil
    # Go语言中的nil参与运算
    
	fmt.Println(m, m2, m3)
	# map[course:golang name:ccmouse quality:notbad site:imooc] map[] map[]
    
	fmt.Println("Traversing map m")
    # map遍历
	for k, v := range m {
		fmt.Println(k, v)
	}

	fmt.Println("Getting values")
    # 查找元素
    # 若key不存在，获取value的值为Zero Value
	courseName := m["course"]
	fmt.Println(`m["course"] =`, courseName)
    
    # 判断key值是否存在
	if causeName, ok := m["cause"]; ok {
		fmt.Println(causeName)
	} else {
		fmt.Println("key 'cause' does not exist")
	}

	fmt.Println("Deleting values")
	name, ok := m["name"]
	fmt.Printf("m[%q] before delete: %q, %v\n",
		"name", name, ok)

    # 删除元素
	delete(m, "name")
	name, ok = m["name"]
	fmt.Printf("m[%q] after delete: %q, %v\n",
		"name", name, ok)
}
```

#### LeetCode真题

寻找最长不重复字符子串

#### 引用类型