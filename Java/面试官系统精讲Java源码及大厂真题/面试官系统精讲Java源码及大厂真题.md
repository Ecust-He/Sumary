[TOC]



## 1 基础

### String、Long 源码解析和面试题

#### String

##### 不可变性

- HashMap的key建议使用不可变类，比如String
- 不可变指的是类一旦被初始化就不能被改变，如果被修改将会是新的类

###### 源码分析

```java
public final class String
    implements java.io.Serializable, Comparable<String>, CharSequence {
    /** The value is used for character storage. */
    private final char value[];
}
```

- String类被final关键字修饰，不可被继承；同时String类的方法也不可被覆写
- String类的底层数据结构是char类型数组，被final关键字修饰，说明value一旦被赋值则无法被修改；同时value被private访问修饰符修饰，外面无法访问，说明一旦赋值，value的内存地址无法被修改

**String的大多数操作方法，会返回新的String**

```java
String str ="hello world !!";
// 这种写法是替换不掉的，必须接受 replace 方法返回的参数才行，这样才行：
System.out.println(str.replace("l", "dd"));
```

##### 字符串乱码

原因：不同环境默认的文件编码方式不一致，需要强制规定文件编码格式

##### 相等判断

###### 源码分析

```java
public boolean equals(Object anObject) {
    // 判断内存地址是否相同
    if (this == anObject) {
        return true;
    }
    // 待比较的对象是否是 String，如果不是 String，直接返回不相等
    if (anObject instanceof String) {
        String anotherString = (String)anObject;
        int n = value.length;
        // 两个字符串的长度是否相等，不等则直接返回不相等
        if (n == anotherString.value.length) {
            char v1[] = value;
            char v2[] = anotherString.value;
            int i = 0;
            // 依次比较每个字符是否相等，若有一个不等，直接返回不相等
            while (n-- != 0) {
                if (v1[i] != v2[i])
                    return false;
                i++;
            }
            return true;
        }
    }
    return false;
}
```

##### 替换和删除

##### 拆分和合并

###### Guava工具类

```java
String a =",a, ,  b  c ,";
// Splitter 是 Guava 提供的 API 
List<String> list = Splitter.on(',')
    .trimResults()// 去掉空格
    .omitEmptyStrings()// 去掉空值
    .splitToList(a);
log.info("Guava 去掉空格的分割方法：{}",JSON.toJSONString(list));
// 打印出的结果为：
["a","b  c"]
```

```java
// 依次 join 多个字符串，Joiner 是 Guava 提供的 API
Joiner joiner = Joiner.on(",").skipNulls();
String result = joiner.join("hello",null,"china");
log.info("依次 join 多个字符串:{}",result);

List<String> list = Lists.newArrayList(new String[]{"hello","china",null});
log.info("自动删除 list 中空值:{}",joiner.join(list));
// 输出的结果为；
依次 join 多个字符串:hello,china
自动删除 list 中空值:hello,china
```

#### Long

##### 缓存机制

缓存了-128到127内的所有Long值，如果命中缓存就不会进行初始化。

```java
private static class LongCache {
    private LongCache(){}

    static final Long cache[] = new Long[-(-128) + 127 + 1];

    static {
        // 缓存 Long 值，注意这里是 i - 128 ，所以再拿的时候就需要 + 128
        for(int i = 0; i < cache.length; i++)
            cache[i] = new Long(i - 128);
    }
}
```

```java
public static Long valueOf(long l) {
    final int offset = 128;
    if (l >= -128 && l <= 127) { // will cache
        return LongCache.cache[(int)l + offset];
    }
    return new Long(l);
}
```

#### 面试题

**1、为什么使用Long时，大多使用valueof方法，少使用parseLong方法？**

因为Long本身有缓存机制，缓存了-128到127范围内的值，如果命中缓存，会从缓存中取值，减小资源的开销。

**2、如何解决String乱码问题？**

乱码问题的根源

- 字符集不支持中文字符
- 二进制转换时字符集不匹配

乱码解决方法

- 在可以指定字符集的地方强制指定字符编码，如果new String()和getBytes()
- 使用UTF-8编码完整支持中文字符集

### Java常用关键字理解

#### Static

静态的、全局的，一旦被修饰，数据可以在一定范围内共享，需要注意并发读写问题。

##### 修饰的对象

###### 类变量

- 如果该变量被public修饰，可直接使用类名.static变量访问。
- 并发读写共享变量时，存在线程安全问题

###### 方法

- 如果该变量被public修饰，代表该方法和当前类无关，任何类都可以访问。
- 该方法内部只能调用同样被static修饰的方法

###### 方法块

静态代码块，在类启动之前，进行一些初始化操作。

```java
public static List<String> list = new ArrayList();
// 进行一些初始化的工作
static {
    list.add("1");
}
```

##### 初始化时机

```tex
父类静态变量初始化
父类静态块初始化
子类静态变量初始化
子类静态块初始化
main 方法执行
父类构造器初始化
子类构造器初始化
```

- 父类的静态变量和静态代码块比子类优先初始化
- 静态变量和静态代码块比类构造器优先初始化
- 被static修饰的方法，在类初始化的时候并不会初始化，只有在调用时，才被执行

#### final

##### 三种使用场景

- 被final修饰的类，无法被继承
- 被final修饰的方法，无法被重写
- 被final修饰的变量，在声明时，必须完成初始化且以后不能修改内存地址（注意：无法修改地址和无法修改内容的区别，比如List和Map集合类的使用）

#### try、catch、finally

```java
public void testCatchFinally() {
  try {
    # print1
    log.info("try is run");
    if (true) {
      throw new RuntimeException("try exception");
    }
  } catch (Exception e) {
    # print2
    log.info("catch is run");
    if (true) {
      throw new RuntimeException("catch exception");
    }
  } finally {
    # print3
    log.info("finally is run");
  }
```

- finally先执行后，再抛出catch到的异常
- 最终捕获的异常是catch的异常，try抛出的异常已经catch捕获（可以使用日志打印）

#### volatile

- volatile关键字常用来修饰共享变量，当共享变量的值被修改后，可被其他线程及时的看到

- 多核CPU下，为提高效率，线程直接与CPU缓存打交道，而不是内存，因为CPU缓存执行速度快；CPU缓存中的值与内存中值可能不是实时同步，需要一种同步机制，即主内存更新后，需要及时的通知CPU缓存更新最新的值

#### transient

使用transient关键自修饰的变量，在序列化时，会忽略该变量

#### default

default关键字修饰的接口，必须要有默认实现，子类无需强制实现

#### 面试题

**1、如何证明static静态变量和类无关？**

- 无需初始化类，可直接使用静态变量
- 静态变量只会初始化一次
- 不仅仅是静态变量，静态方法快也和类无关

**2、常常看见变量和方法被static和final关键字修饰，为什么这么做？**

- 变量和方法与类无关，可直接使用比较方便
- 强调变量内存地址不可变，方法不可继承，强调方法内部的稳定性

**3、catch发生了异常，finally还会执行吗？**

会的。并且在finally执行完成之后，才会抛出catch的异常

### Arrays、Collections、Objects 常用方法源码解析

#### 工具类通用特征

- 类的构造方法是私有的
- 工具类的方法必须被static final关键字修饰。这样的话可以保证方法的不可变性
- 尽量不在工具方法中，对共享变量做修改操作访问；如果必须使用要加锁，不然会有线程安全问题

#### Arrays

#### Collections

#### Objects

##### 相等判断

##### 为空判断

#### 面试题
