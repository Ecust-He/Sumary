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

##### 替换、删除
