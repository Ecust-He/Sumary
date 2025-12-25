# 基础篇

## init函数是什么时候执行的？
1、程序执行前包初始化时（只会初始化一次）
2、package依赖顺序为：main -> A -> B -> C
   执行顺序为：C -> B -> A -> main

## go语言中new和make关键字的区别？
1、make不仅会分配内存空间而且还会对内存进行初始化（不为nil）
2、make返回的是类型本身（T）;new返回的是指向类型的指针（*T）
3、make只能为slice、map、channel类型进行分配及初始化内存空间，new 可以为任意类型

## 数组和切片有什么区别？

## float或切片可以作为map类型的key吗？
只有可比较的类型才能做map的key

## interface可以比较吗？
空接口     eface（_type和data）
非空接口   iface（tab和data）

类型和数据都相等，interface才相等
未显示初始化的接口的类型和数据都为nil

## 空struct{}有什么用？
不占内存空间

## 如何优雅的处理go语言中的错误？


## 如何判断两个对象是否完全相等？


# 高级特性
## context最佳实践
<- ctx.Done()表示任务取消或超时
### 使用场景
#### 取消任务
取消请求时，通知goroutine任务退出

#### 超时控制

#### 通过函数签名传递context参数









