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

## 2 集合

### ArrayList 源码解析和设计思路

#### 整体架构

- DEFAULT_CAPACITY表示数组的初始大小，默认是10
- size表示当前数组的大小，int类型，没有被volatile关键字修饰，非线程安全
- modCount统计数组被修改版本的次数

##### 类注释

- 允许put null值，会自动扩容
- size、isEmpty、get、set等方法时间复杂度是O(1)
- add方法时间复杂度是O(1)
- 非线程安全的，多线程并发场景下，推荐使用线程安全类（List list = Collections.synchronizedList(new ArrayList(...));）
- 增强的for循环，迭代器迭代过程中，如果有修改操作（add、remove方法），会抛出ConcurrentModificationException异常，采用fail-fast机制

#### 源码分析

##### 初始化

###### 无参直接初始化

```java
private static final Object[] DEFAULTCAPACITY_EMPTY_ELEMENTDATA = {};

//无参数直接初始化，数组大小为空
public ArrayList() {
    this.elementData = DEFAULTCAPACITY_EMPTY_ELEMENTDATA;
}
```

###### 指定大小初始化

```java
public ArrayList(int initialCapacity) {
    if (initialCapacity > 0) {
        this.elementData = new Object[initialCapacity];
    } else if (initialCapacity == 0) {
        this.elementData = EMPTY_ELEMENTDATA;
    } else {
        throw new IllegalArgumentException("Illegal Capacity: "+
                                           initialCapacity);
    }
}
```

###### 指定初始数据初始化

```java
//指定初始数据初始化
public ArrayList(Collection<? extends E> c) {
    //elementData 是保存数组的容器，默认为 null
    elementData = c.toArray();
    //如果给定的集合（c）数据有值
    if ((size = elementData.length) != 0) {
        // c.toArray might (incorrectly) not return Object[] (see 6260652)
        //如果集合元素类型不是 Object 类型，我们会转成 Object
        if (elementData.getClass() != Object[].class) {
            elementData = Arrays.copyOf(elementData, size, Object[].class);
        }
    } else {
        // 给定集合（c）无值，则默认空数组
        this.elementData = EMPTY_ELEMENTDATA;
    }
}
```

##### 新增和扩容实现

新增就是往数组中添加元素，主要分两步：

- 判断是否需要扩容，如果需要进行扩容操作，否则直接赋值

###### 添加元素源码

```java
public boolean add(E e) {
  //确保数组大小是否足够，不够执行扩容，size 为当前数组的大小
  ensureCapacityInternal(size + 1);  // Increments modCount!!
  //直接赋值，线程不安全的
  elementData[size++] = e;
  return true;
}
```

###### 扩容源码

```java
private void ensureCapacityInternal(int minCapacity) {
  //如果初始化数组大小时，有给定初始值，以给定的大小为准，不走 if 逻辑
  if (elementData == DEFAULTCAPACITY_EMPTY_ELEMENTDATA) {
    minCapacity = Math.max(DEFAULT_CAPACITY, minCapacity);
  }
  //确保容积足够
  ensureExplicitCapacity(minCapacity);
}

private void ensureExplicitCapacity(int minCapacity) {
  //记录数组被修改
  modCount++;
  // 如果我们期望的最小容量大于目前数组的长度，那么就扩容
  if (minCapacity - elementData.length > 0)
    grow(minCapacity);
}

//扩容，并把现有数据拷贝到新的数组里面去
private void grow(int minCapacity) {
  int oldCapacity = elementData.length;
  // oldCapacity >> 1 是把 oldCapacity 除以 2 的意思
  int newCapacity = oldCapacity + (oldCapacity >> 1);

  // 如果扩容后的值 < 我们的期望值，扩容后的值就等于我们的期望值
  if (newCapacity - minCapacity < 0)
    newCapacity = minCapacity;

  // 如果扩容后的值 > jvm 所能分配的数组的最大值，那么就用 Integer 的最大值
  if (newCapacity - MAX_ARRAY_SIZE > 0)
    newCapacity = hugeCapacity(minCapacity);
 
  // 通过复制进行扩容
  elementData = Arrays.copyOf(elementData, newCapacity);
}
```

- 扩容的规则：扩容后的大小是原容量大小 + 容量大小的一半，直白的说，扩容后的大小是原来容量的1.5倍
- Arraylist中的数组最大值是Integer.MAX_VALUE，超过这个值，JVM就不会给数组分配内存了
- 新增时，并没有对数值进行严格校验，所以ArrayList是运行null值的

###### 扩容的本质

本质：数组之间的拷贝，会新建一个符合预期容量的新数组，然后将老数组中的数据拷贝过去

##### 删除

###### 根据索引删除

###### 根据值删除

```java
public boolean remove(Object o) {
  // 如果要删除的值是 null，找到第一个值是 null 的删除
  if (o == null) {
    for (int index = 0; index < size; index++)
      if (elementData[index] == null) {
        fastRemove(index);
        return true;
      }
  } else {
    // 如果要删除的值不为 null，找到第一个和要删除的值相等的删除
    for (int index = 0; index < size; index++)
      // 这里是根据  equals 来判断值相等的，相等后再根据索引位置进行删除
      if (o.equals(elementData[index])) {
        fastRemove(index);
        return true;
      }
  }
  return false;
}

private void fastRemove(int index) {
  // 记录数组的结构要发生变动了
  modCount++;
  // numMoved 表示删除 index 位置的元素后，需要从 index 后移动多少个元素到前面去
  // 减 1 的原因，是因为 size 从 1 开始算起，index 从 0开始算起
  int numMoved = size - index - 1;
  if (numMoved > 0)
    // 从 index +1 位置开始被拷贝，拷贝的起始位置是 index，长度是 numMoved
    System.arraycopy(elementData, index+1, elementData, index, numMoved);
  //数组最后一个位置赋值 null，帮助 GC
  elementData[--size] = null;
}
```

###### 批量删除

- 新增元素时没有对null值校验，所以删除元素时，允许删除null值
- 找到值在数组中的索引位置是通过equals方法判断的，如果数组中元素不是基本数据类型，需要关注equals的具体实现

##### 迭代器

```java
int cursor;// 迭代过程中，下一个元素的位置，默认从 0 开始。
int lastRet = -1; // 新增场景：表示上一次迭代过程中，索引的位置；删除场景：为 -1。
int expectedModCount = modCount;// expectedModCount 表示迭代过程中，期望的版本号；modCount 表示数组实际的版本号。

public boolean hasNext() {
  return cursor != size;//cursor 表示下一个元素的位置，size 表示实际大小，如果两者相等，说明已经没有元素可以迭代了，如果不等，说明还可以迭代
}

public E next() {
  //迭代过程中，判断版本号有无被修改，有被修改，抛 ConcurrentModificationException 异常
  checkForComodification();
  //本次迭代过程中，元素的索引位置
  int i = cursor;
  if (i >= size)
    throw new NoSuchElementException();
  Object[] elementData = ArrayList.this.elementData;
  if (i >= elementData.length)
    throw new ConcurrentModificationException();
  // 下一次迭代时，元素的位置，为下一次迭代做准备
  cursor = i + 1;
  // 返回元素值
  return (E) elementData[lastRet = i];
}
// 版本号比较
final void checkForComodification() {
  if (modCount != expectedModCount)
    throw new ConcurrentModificationException();
}

public void remove() {
  // 如果上一次操作时，数组的位置已经小于 0 了，说明数组已经被删除完了
  if (lastRet < 0)
    throw new IllegalStateException();
  //迭代过程中，判断版本号有无被修改，有被修改，抛 ConcurrentModificationException 异常
  checkForComodification();

  try {
    ArrayList.this.remove(lastRet);
    cursor = lastRet;
    // -1 表示元素已经被删除，这里也防止重复删除
    lastRet = -1;
    // 删除元素时 modCount 的值已经发生变化，在此赋值给 expectedModCount
    // 这样下次迭代时，两者的值是一致的了
    expectedModCount = modCount;
  } catch (IndexOutOfBoundsException ex) {
    throw new ConcurrentModificationException();
  }
}
```

### LinkedList 源码解析

#### 整体架构

LinkedList底层数据结构是一个双向链表

##### 源码分析

###### 节点增加

**从尾部追加（add）**

```java
// 从尾部开始追加节点
void linkLast(E e) {
    // 把尾节点数据暂存
    final Node<E> l = last;
    // 新建新的节点，初始化入参含义：
    // l 是新节点的前一个节点，当前值是尾节点值
    // e 表示当前新增节点，当前新增节点后一个节点是 null
    final Node<E> newNode = new Node<>(l, e, null);
    // 新建节点追加到尾部
    last = newNode;
    //如果链表为空（l 是尾节点，尾节点为空，链表即空），头部和尾部是同一个节点，都是新建的节点
    if (l == null)
        first = newNode;
    //否则把前尾节点的下一个节点，指向当前尾节点。
    else
        l.next = newNode;
    //大小和版本更改
    size++;
    modCount++;
}
```

**从头部追加（addFirst）**

```java
// 从头部追加
private void linkFirst(E e) {
    // 头节点赋值给临时变量
    final Node<E> f = first;
    // 新建节点，前一个节点指向null，e 是新建节点，f 是新建节点的下一个节点，目前值是头节点的值
    final Node<E> newNode = new Node<>(null, e, f);
    // 新建节点成为头节点
    first = newNode;
    // 头节点为空，就是链表为空，头尾节点是一个节点
    if (f == null)
        last = newNode;
    //上一个头节点的前一个节点指向当前节点
    else
        f.prev = newNode;
    size++;
    modCount++;
}
```

###### 节点删除

**从头部删除**

```java
//从头删除节点 f 是链表头节点
private E unlinkFirst(Node<E> f) {
    // 拿出头节点的值，作为方法的返回值
    final E element = f.item;
    // 拿出头节点的下一个节点
    final Node<E> next = f.next;
    //帮助 GC 回收头节点
    f.item = null;
    f.next = null;
    // 头节点的下一个节点成为头节点
    first = next;
    //如果 next 为空，表明链表为空
    if (next == null)
        last = null;
    //链表不为空，头节点的前一个节点指向 null
    else
        next.prev = null;
    //修改链表大小和版本
    size--;
    modCount++;
    return element;
}
```

###### 节点查询

```java
// 根据链表索引位置查询节点
Node<E> node(int index) {
    // 如果 index 处于队列的前半部分，从头开始找，size >> 1 是 size 除以 2 的意思。
    if (index < (size >> 1)) {
        Node<E> x = first;
        // 直到 for 循环到 index 的前一个 node 停止
        for (int i = 0; i < index; i++)
            x = x.next;
        return x;
    } else {// 如果 index 处于队列的后半部分，从尾开始找
        Node<E> x = last;
        // 直到 for 循环到 index 的后一个 node 停止
        for (int i = size - 1; i > index; i--)
            x = x.prev;
        return x;
    }
}
```

LinkedList采用了简单的二分查找法，若查找的索引在链表的前半部分，从头部开始查找；反之，从尾部开始查找。通过这种方式，使得循环次数减少了一半，大大的提高了查找速率，值得借鉴。

###### 迭代器

### List源码面试问题

**1、说说你对ArrayList的理解？**

ArrayList 底层数据结构是个数组，其 API 都做了一层对数组底层访问的封装，比如说 add 方法的过程是……（这里可以引用我们在 ArrayList 源码解析中 add 的过程）。

**2、ArrayList 无参数构造器构造，现在 add 一个值进去，此时数组的大小是多少，下一次扩容前最大可用大小是多少？**

此处数组的大小是 1，下一次扩容前最大可用大小是 10

**3、 数组初始化，被加入一个值后，如果我使用 addAll 方法，一下子加入 15 个值，那么最终数组的大小是多少？**

数组在加入一个值后，实际大小是 1，最大可用大小是 10 ，现在需要一下子加入 15 个值，那我们期望数组的大小值就是 16，此时数组最大可用大小只有 10，明显不够，需要扩容，扩容后的大小是：10 + 10 ／2 = 15，这时候发现扩容后的大小仍然不到我们期望的值 16，这时候源码中有一种策略如下：

```java
// newCapacity 本次扩容的大小，minCapacity 我们期望的数组最小大小
// 如果扩容后的值 < 我们的期望值，我们的期望值就等于本次扩容的大小
if (newCapacity - minCapacity < 0)
    newCapacity = minCapacity;
```

所以最终数组扩容后的大小为 16。

**4、 现在我有一个很大的数组需要拷贝，原数组大小是 5k，请问如何快速拷贝？**

因为原数组比较大，如果新建新数组的时候，不指定数组大小的话，就会频繁扩容，频繁扩容就会有大量拷贝的工作，造成拷贝的性能低下，所以回答说新建数组时，指定新数组的大小为 5k 即可。

**5、为什么说扩容会消耗性能？**

扩容底层使用的是 System.arraycopy 方法，会把原数组的数据全部拷贝到新数组上，所以性能消耗比较严重。

**6、ArrayList 和 LinkedList 有何不同？**

可以先从底层数据结构开始说起，然后以某一个方法为突破口深入，比如：最大的不同是两者底层的数据结构不同，ArrayList 底层是数组，LinkedList 底层是双向链表，两者的数据结构不同也导致了操作的 API 实现有所差异，拿新增实现来说，ArrayList 会先计算并决定是否扩容，然后把新增的数据直接赋值到数组上，而 LinkedList 仅仅只需要改变插入节点和其前后节点的指向位置关系即可。

**7、ArrayList 和 LinkedList 应用场景有何不同？**

ArrayList 更适合于快速的查找匹配，不适合频繁新增删除，像工作中经常会对元素进行匹配查询的场景比较合适，LinkedList 更适合于经常新增和删除，对查询反而很少的场景。

**8、ArrayList 和 LinkedList 两者有没有最大容量？**

ArrayList 有最大容量的，为 Integer 的最大值，大于这个值 JVM 是不会为数组分配内存空间的，LinkedList 底层是双向链表，理论上可以无限大。但源码中，LinkedList 实际大小用的是 int 类型，这也说明了 LinkedList 不能超过 Integer 的最大值，不然会溢出。

**9、ArrayList 和 LinkedList 是如何对 null 值进行处理的？**

ArrayList 允许 null 值新增，也允许 null 值删除。删除 null 值时，是从头开始，找到第一值是 null 的元素删除；LinkedList 新增删除时对 null 值没有特殊校验，是允许新增和删除的。

**10、 ArrayList 和 LinedList 是线程安全的么，为什么？**

当两者作为非共享变量时，比如说仅仅是在方法里面的局部变量时，是没有线程安全问题的，只有当两者是共享变量时，才会有线程安全问题。主要的问题点在于多线程环境下，所有线程任何时刻都可对数组和链表进行操作，这会导致值被覆盖，甚至混乱的情况。

如果有线程安全问题，在迭代的过程中，会频繁报 ConcurrentModificationException 的错误，意思是在我当前循环的过程中，数组或链表的结构被其它线程修改了。

**11、你能描述下双向链表么？**

如果和面试官面对面沟通的话，你可以去画一下，可以把 《LinkedList 源码解析》中的 LinkedList 的结构画出来，如果是电话面试，可以这么描述：双向链表中双向的意思是说前后节点之间互相有引用，链表的节点我们称为 Node。Node 有三个属性组成：其前一个节点，本身节点的值，其下一个节点，假设 A、B 节点相邻，A 节点的下一个节点就是 B，B 节点的上一个节点就是 A，两者互相引用，在链表的头部节点，我们称为头节点。头节点的前一个节点是 null，尾部称为尾节点，尾节点的后一个节点是 null，如果链表数据为空的话，头尾节点是同一个节点，本身是 null，指向前后节点的值也是 null。

**12、描述下双向链表的新增和删除**

新增：我们可以选择从链表头新增，也可以选择从链表尾新增，如果是从链表尾新增的话，直接把当前节点追加到尾节点之后，本身节点自动变为尾节点。

删除：把删除节点的后一个节点的 prev 指向其前一个节点，把删除节点的前一个节点的 next 指向其后一个节点，最后把删除的节点置为 null 即可。

### HashMap源码分析

HashMap 底层的数据结构主要是：数组 + 链表 + 红黑树。其中当链表的长度大于等于 8 时，链表会转化成红黑树，当红黑树的大小小于等于 6 时，红黑树会转化成链表

#### 类注释

- 允许 null 值，不同于 HashTable ，是线程不安全的；
- load factor（影响因子） 默认值是 0.75， 是均衡了时间和空间损耗算出来的值，较高的值会减少空间开销（扩容减少，数组大小增长速度变慢），但增加了查找成本（hash 冲突增加，链表长度变长），不扩容的条件：数组容量 > 需要的数组大小 /load factor；
- 如果有很多数据需要储存到 HashMap 中，建议 HashMap 的容量一开始就设置成足够的大小，这样可以防止在其过程中不断的扩容，影响性能；
- HashMap 是非线程安全的，我们可以自己在外部加锁，或者通过 Collections#synchronizedMap 来实现线程安全，Collections#synchronizedMap 的实现是在每个方法上加上了 synchronized 锁；
- 在迭代过程中，如果 HashMap 的结构被修改，会快速失败。

#### 常见属性

```java

 //桶上的红黑树大小小于等于6时，红黑树转化成链表
 static final int UNTREEIFY_THRESHOLD = 6;

 //当数组容量大于 64 时，链表才会转化成红黑树
 static final int MIN_TREEIFY_CAPACITY = 64;

 //记录迭代过程中 HashMap 结构是否发生变化，如果有变化，迭代时会 fail-fast
 transient int modCount;

 //HashMap 的实际大小，可能不准(因为当你拿到这个值的时候，可能又发生了变化)
 transient int size;

 //存放数据的数组
 transient Node<K,V>[] table;

 // 扩容的门槛，有两种情况
 // 如果初始化时，给定数组大小的话，通过 tableSizeFor 方法计算，数组大小永远接近于 2 的幂次方，比如你给定初始化大小 19，实际上初始化大小为 32，为 2 的 5 次方。
 // 如果是通过 resize 方法进行扩容，大小 = 数组容量 * 0.75
 int threshold;

 //链表的节点
 static class Node<K,V> implements Map.Entry<K,V> {
 
 //红黑树的节点
 static final class TreeNode<K,V> extends LinkedHashMap.Entry<K,V> {
```

#### 新增

```java
// 入参 hash：通过 hash 算法计算出来的值。
// 入参 onlyIfAbsent：false 表示即使 key 已经存在了，仍然会用新值覆盖原来的值，默认为 false
final V putVal(int hash, K key, V value, boolean onlyIfAbsent,
               boolean evict) {
    // n 表示数组的长度，i 为数组索引下标，p 为 i 下标位置的 Node 值
    Node<K,V>[] tab; Node<K,V> p; int n, i;
    //如果数组为空，使用 resize 方法初始化
    if ((tab = table) == null || (n = tab.length) == 0)
        n = (tab = resize()).length;
    // 如果当前索引位置是空的，直接生成新的节点在当前索引位置上
    if ((p = tab[i = (n - 1) & hash]) == null)
        tab[i] = newNode(hash, key, value, null);
    // 如果当前索引位置有值的处理方法，即我们常说的如何解决 hash 冲突
    else {
        // e 当前节点的临时变量
        Node<K,V> e; K k;
        // 如果 key 的 hash 和值都相等，直接把当前下标位置的 Node 值赋值给临时变量
        if (p.hash == hash &&
            ((k = p.key) == key || (key != null && key.equals(k))))
            e = p;
        // 如果是红黑树，使用红黑树的方式新增
        else if (p instanceof TreeNode)
            e = ((TreeNode<K,V>)p).putTreeVal(this, tab, hash, key, value);
        // 是个链表，把新节点放到链表的尾端
        else {
            // 自旋
            for (int binCount = 0; ; ++binCount) {
                // e = p.next 表示从头开始，遍历链表
                // p.next == null 表明 p 是链表的尾节点
                if ((e = p.next) == null) {
                    // 把新节点放到链表的尾部 
                    p.next = newNode(hash, key, value, null);
                    // 当链表的长度大于等于 8 时，链表转红黑树
                    if (binCount >= TREEIFY_THRESHOLD - 1)
                        treeifyBin(tab, hash);
                    break;
                }
                // 链表遍历过程中，发现有元素和新增的元素相等，结束循环
                if (e.hash == hash &&
                    ((k = e.key) == key || (key != null && key.equals(k))))
                    break;
                //更改循环的当前元素，使 p 在遍历过程中，一直往后移动。
                p = e;
            }
        }
        // 说明新节点的新增位置已经找到了
        if (e != null) {
            V oldValue = e.value;
            // 当 onlyIfAbsent 为 false 时，才会覆盖值 
            if (!onlyIfAbsent || oldValue == null)
                e.value = value;
            afterNodeAccess(e);
            // 返回老值
            return oldValue;
         }
    }
    // 记录 HashMap 的数据结构发生了变化
    ++modCount;
    //如果 HashMap 的实际大小大于扩容的门槛，开始扩容
    if (++size > threshold)
        resize();
    afterNodeInsertion(evict);
    return null;
}
```

##### 链表的新增

###### 红黑树新增节点过程

1. 首先判断新增的节点在红黑树上是不是已经存在，判断手段有如下两种：1）如果节点没有实现 Comparable 接口，使用 equals 进行判断；2） 如果节点自己实现了 Comparable 接口，使用 compareTo 进行判断。
2. 新增的节点如果已经在红黑树上，直接返回；不在的话，判断新增节点是在当前节点的左边还是右边，左边值小，右边值大；
3. 自旋递归 1 和 2 步，直到当前节点的左边或者右边的节点为空时，停止自旋，当前节点即为我们新增节点的父节点；
4. 把新增节点放到当前节点的左边或右边为空的地方，并于当前节点建立父子节点关系；
5. 进行着色和旋转，结束

```java
//入参 h：key 的hash值
final TreeNode<K,V> putTreeVal(HashMap<K,V> map, Node<K,V>[] tab,
                               int h, K k, V v) {
    Class<?> kc = null;
    boolean searched = false;
    //找到根节点
    TreeNode<K,V> root = (parent != null) ? root() : this;
    //自旋
    for (TreeNode<K,V> p = root;;) {
        int dir, ph; K pk;
        // p hash 值大于 h，说明 p 在 h 的右边
        if ((ph = p.hash) > h)
            dir = -1;
        // p hash 值小于 h，说明 p 在 h 的左边
        else if (ph < h)
            dir = 1;
        //要放进去key在当前树中已经存在了(equals来判断)
        else if ((pk = p.key) == k || (k != null && k.equals(pk)))
            return p;
        //自己实现的Comparable的话，不能用hashcode比较了，需要用compareTo
        else if ((kc == null &&
                  //得到key的Class类型，如果key没有实现Comparable就是null
                  (kc = comparableClassFor(k)) == null) ||
                  //当前节点pk和入参k不等
                 (dir = compareComparables(kc, k, pk)) == 0) {
            if (!searched) {
                TreeNode<K,V> q, ch;
                searched = true;
                if (((ch = p.left) != null &&
                     (q = ch.find(h, k, kc)) != null) ||
                    ((ch = p.right) != null &&
                     (q = ch.find(h, k, kc)) != null))
                    return q;
            }
            dir = tieBreakOrder(k, pk);
        }

        TreeNode<K,V> xp = p;
        //找到和当前hashcode值相近的节点(当前节点的左右子节点其中一个为空即可)
        if ((p = (dir <= 0) ? p.left : p.right) == null) {
            Node<K,V> xpn = xp.next;
            //生成新的节点
            TreeNode<K,V> x = map.newTreeNode(h, k, v, xpn);
            //把新节点放在当前子节点为空的位置上
            if (dir <= 0)
                xp.left = x;
            else
                xp.right = x;
            //当前节点和新节点建立父子，前后关系
            xp.next = x;
            x.parent = x.prev = xp;
            if (xpn != null)
                ((TreeNode<K,V>)xpn).prev = x;
            //balanceInsertion 对红黑树进行着色或旋转，以达到更多的查找效率，着色或旋转的几种场景如下
            //着色：新节点总是为红色；如果新节点的父亲是黑色，则不需要重新着色；如果父亲是红色，那么必须通过重新着色或者旋转的方法，再次达到红黑树的5个约束条件
            //旋转： 父亲是红色，叔叔是黑色时，进行旋转
            //如果当前节点是父亲的右节点，则进行左旋
            //如果当前节点是父亲的左节点，则进行右旋
          
            //moveRootToFront 方法是把算出来的root放到根节点上
            moveRootToFront(tab, balanceInsertion(root, x));
            return null;
        }
    }
}
```

红黑树的新增，要求大家对红黑树的数据结构有一定的了解。面试的时候，一般只会问到新增节点到红黑树上大概是什么样的一个过程，着色和旋转的细节不会问，因为很难说清楚，但我们要清楚着色指的是给红黑树的节点着上红色或黑色，旋转是为了让红黑树更加平衡，提高查询的效率，总的来说都是为了满足红黑树的 5 个原则：

1. 节点是红色或黑色
2. 根是黑色
3. 所有叶子都是黑色
4. 从任一节点到其每个叶子的所有简单路径都包含相同数目的黑色节点
5. 从每个叶子到根的所有路径上不能有两个连续的红色节点

#### 查找

HashMap查找主要分为以下三个步骤：

- 根据 hash 算法定位数组的索引位置，equals 判断当前节点是否是我们需要寻找的 key，是的话直接返回，不是的话往下
- 判断当前节点有无 next 节点，有的话判断是链表类型，还是红黑树类型
- 分别走链表和红黑树不同类型的查找方法

```java
// 采用自旋方式从链表中查找 key，e 初始为为链表的头节点
do {
    // 如果当前节点 hash 等于 key 的 hash，并且 equals 相等，当前节点就是我们要找的节点
    // 当 hash 冲突时，同一个 hash 值上是一个链表的时候，我们是通过 equals 方法来比较 key 是否相等的
    if (e.hash == hash &&
        ((k = e.key) == key || (key != null && key.equals(k))))
        return e;
    // 否则，把当前节点的下一个节点拿出来继续寻找
} while ((e = e.next) != null);
```

### TreeMap 和 LinkedHashMap 核心源码解析

#### TreeMap

##### 整体架构

```java
//比较器，如果外部有传进来 Comparator 比较器，首先用外部的
//如果外部比较器为空，则使用 key 自己实现的 Comparable#compareTo 方法
//比较手段和上面日常工作中的比较 demo 是一致的
private final Comparator<? super K> comparator;

//红黑树的根节点
private transient Entry<K,V> root;

//红黑树的已有元素大小
private transient int size = 0;

//树结构变化的版本号，用于迭代过程中的快速失败场景
private transient int modCount = 0;

//红黑树的节点
static final class Entry<K,V> implements Map.Entry<K,V> {}
```

##### 新增节点

1、判断红黑树的节点是否为空，为空的话，新增的节点直接作为根节点

```java
Entry<K,V> t = root;
//红黑树根节点为空，直接新建
if (t == null) {
    // compare 方法限制了 key 不能为 null
    compare(key, key); // type (and possibly null) check
    // 成为根节点
    root = new Entry<>(key, value, null);
    size = 1;
    modCount++;
    return null;
}
```

2、根据红黑树左小右大的特性，进行判断，找到应该新增节点的父节点

```java
Comparator<? super K> cpr = comparator;
if (cpr != null) {
    //自旋找到 key 应该新增的位置，就是应该挂载那个节点的头上
    do {
        //一次循环结束时，parent 就是上次比过的对象
        parent = t;
        // 通过 compare 来比较 key 的大小
        cmp = cpr.compare(key, t.key);
        //key 小于 t，把 t 左边的值赋予 t，因为红黑树左边的值比较小，循环再比
        if (cmp < 0)
            t = t.left;
        //key 大于 t，把 t 右边的值赋予 t，因为红黑树右边的值比较大，循环再比
        else if (cmp > 0)
            t = t.right;
        //如果相等的话，直接覆盖原值
        else
            return t.setValue(value);
        // t 为空，说明已经到叶子节点了
    } while (t != null);
}
```

3、在父节点的左边或右边插入新增节点

```java
//cmp 代表最后一次对比的大小，小于 0 ，代表 e 在上一节点的左边
if (cmp < 0)
    parent.left = e;
//cmp 代表最后一次对比的大小，大于 0 ，代表 e 在上一节点的右边，相等的情况第二步已经处理了。
else
    parent.right = e;
```

4、着色旋转，达到平衡，结束

#### LinkedHashMap

##### 整体架构

- 按照插入顺序进行访问

- 实现了访问最少最先删除功能，其目的是把很久都没有访问的 key 自动删除

##### 按照插入顺序访问

```java
// 链表头
transient LinkedHashMap.Entry<K,V> head;

// 链表尾
transient LinkedHashMap.Entry<K,V> tail;

// 继承 Node，为数组的每个元素增加了 before 和 after 属性
static class Entry<K,V> extends HashMap.Node<K,V> {
    Entry<K,V> before, after;
    Entry(int hash, K key, V value, Node<K,V> next) {
        super(hash, key, value, next);
    }
}

// 控制两种访问模式的字段，默认 false
// true 按照访问顺序，会把经常访问的 key 放到队尾
// false 按照插入顺序提供访问
final boolean accessOrder;

// 新增节点，并追加到链表的尾部
Node<K,V> newNode(int hash, K key, V value, Node<K,V> e) {
    // 新增节点
    LinkedHashMap.Entry<K,V> p =
        new LinkedHashMap.Entry<K,V>(hash, key, value, e);
    // 追加到链表的尾部
    linkNodeLast(p);
    return p;
}
// link at the end of list
private void linkNodeLast(LinkedHashMap.Entry<K,V> p) {
    LinkedHashMap.Entry<K,V> last = tail;
    // 新增节点等于位节点
    tail = p;
    // last 为空，说明链表为空，首尾节点相等
    if (last == null)
        head = p;
    // 链表有数据，直接建立新增节点和上个尾节点之间的前后关系即可
    else {
        p.before = last;
        last.after = p;
    }
}

// 初始化时，默认从头节点开始访问
LinkedHashIterator() {
    // 头节点作为第一个访问的节点
    next = head;
    expectedModCount = modCount;
    current = null;
}

final LinkedHashMap.Entry<K,V> nextNode() {
    LinkedHashMap.Entry<K,V> e = next;
    if (modCount != expectedModCount)// 校验
        throw new ConcurrentModificationException();
    if (e == null)
        throw new NoSuchElementException();
    current = e;
    next = e.after; // 通过链表的 after 结构，找到下一个迭代的节点
    return e;
}
```

##### 访问最少删除策略

###### 示例

```java
public void testAccessOrder() {
  // 新建 LinkedHashMap
  LinkedHashMap<Integer, Integer> map = new LinkedHashMap<Integer, Integer>(4,0.75f,true) {
    {
      put(10, 10);
      put(9, 9);
      put(20, 20);
      put(1, 1);
    }

    @Override
    // 覆写了删除策略的方法，我们设定当节点个数大于 3 时，就开始删除头节点
    protected boolean removeEldestEntry(Map.Entry<Integer, Integer> eldest) {
      return size() > 3;
    }
  };

  log.info("初始化：{}",JSON.toJSONString(map));
  Assert.assertNotNull(map.get(9));
  log.info("map.get(9)：{}",JSON.toJSONString(map));
  Assert.assertNotNull(map.get(20));
  log.info("map.get(20)：{}",JSON.toJSONString(map));

}
```

######  元素被转移到队尾

```java
public V get(Object key) {
    Node<K,V> e;
    // 调用 HashMap  get 方法
    if ((e = getNode(hash(key), key)) == null)
        return null;
    // 如果设置了 LRU 策略
    if (accessOrder)
    // 这个方法把当前 key 移动到队尾
        afterNodeAccess(e);
    return e.value;
}
```

###### 删除策略

```java
// 删除很少被访问的元素，被 HashMap 的 put 方法所调用
void afterNodeInsertion(boolean evict) { 
    // 得到元素头节点
    LinkedHashMap.Entry<K,V> first;
    // removeEldestEntry 来控制删除策略，如果队列不为空，并且删除策略允许删除的情况下，删除头节点
    if (evict && (first = head) != null && removeEldestEntry(first)) {
        K key = first.key;
        // removeNode 删除头节点
        removeNode(hash(key), key, null, false, true);
    }
}
```

#### Map源码面试问题

**1、说一说HashMap底层数据结构**

HashMap 底层是数组 + 链表 + 红黑树的数据结构，数组的主要作用是方便快速查找，时间复杂度是 O(1)，默认大小是 16，数组的下标索引是通过 key 的 hashcode 计算出来的，数组元素叫做 Node，当多个 key 的 hashcode 一致，但 key 值不同时，单个 Node 就会转化成链表，链表的查询复杂度是 O(n)，当链表的长度大于等于 8 并且数组的大小超过 64 时，链表就会转化成红黑树，红黑树的查询复杂度是 O(log(n))，简单来说，最坏的查询次数相当于红黑树的最大深度。

**2、HashMap、TreeMap、LinkedHashMap 三者有啥相同点，有啥不同点？**

**相同点：**

1. 三者在特定的情况下，都会使用红黑树；
2. 底层的 hash 算法相同；
3. 在迭代的过程中，如果 Map 的数据结构被改动，都会报 ConcurrentModificationException 的错误。

**不同点：**

1. HashMap 数据结构以数组为主，查询非常快，TreeMap 数据结构以红黑树为主，利用了红黑树左小右大的特点，可以实现 key 的排序，LinkedHashMap 在 HashMap 的基础上增加了链表的结构，实现了插入顺序访问和最少访问删除两种策略;
2. 由于三种 Map 底层数据结构的差别，导致了三者的使用场景的不同，TreeMap 适合需要根据 key 进行排序的场景，LinkedHashMap 适合按照插入顺序访问，或需要删除最少访问元素的场景，剩余场景我们使用 HashMap 即可，我们工作中大部分场景基本都在使用 HashMap；
3. 由于三种 map 的底层数据结构的不同，导致上层包装的 api 略有差别。

**3、说一下Map的hash算法**

```java
static final int hash(Object key) {
    int h;
    return (key == null) ? 0 : (h = key.hashCode()) ^ (h >>> 16);
}
key 在数组中的位置公式：tab[(n - 1) & hash]
```

这其实是一个数学问题，源码中就是通过以上代码来计算 hash 的，首先计算出 key 的 hashcode，因为 key 是 Object，所以会根据 key 的不同类型进行 hashcode 的计算，接着计算 h ^ (h >>> 16) ，这么做的好处是使大多数场景下，算出来的 hash 值比较分散。

一般来说，hash 值算出来之后，要计算当前 key 在数组中的索引下标位置时，可以采用取模的方式，就是索引下标位置 = hash 值 % 数组大小，这样做的好处，就是可以保证计算出来的索引下标值可以均匀的分布在数组的各个索引位置上，但取模操作对于处理器的计算是比较慢的，数学上有个公式，当 b 是 2 的幂次方时，a % b = a &（b-1），所以此处索引位置的计算公式我们可以更换为： (n-1) & hash。

**此问题可以延伸出三个小问题：**

1、为什么不用 key % 数组大小，而是需要用 key 的 hash 值 % 数组大小。

如果 key 是数字，直接用 key % 数组大小是完全没有问题的，但我们的 key 还有可能是字符串，是复杂对象，这时候用 字符串或复杂对象 % 数组大小是不行的，所以需要先计算出 key 的 hash 值。

2、计算 hash 值时，为什么需要右移 16 位

hash 算法是 h ^ (h >>> 16)，为了使计算出的 hash 值更分散，所以选择先将 h 无符号右移 16 位，然后再于 h 异或时，就能达到 h 的高 16 位和低 16 位都能参与计算，减少了碰撞的可能性。

3、为什么把取模操作换成了 & 操作？

key.hashCode() 算出来的 hash 值还不是数组的索引下标，为了随机的计算出索引的下表位置，我们还会用 hash 值和数组大小进行取模，这样子计算出来的索引下标比较均匀分布。

取模操作处理器计算比较慢，处理器对 & 操作就比较擅长，换成了 & 操作，是有数学上证明的支撑，为了提高了处理器处理的速度。

4、为什么提倡数组大小是 2 的幂次方？

因为只有大小是 2 的幂次方时，才能使 hash 值 % n(数组大小) == (n-1) & hash 公式成立。

**4、为解决 hash 冲突，大概有哪些办法？**

1. 好的 hash 算法，细问的话复述一下上题的 hash 算法;
2. 自动扩容，当数组大小快满的时候，采取自动扩容，可以减少 hash 冲突;
3. hash 冲突发生时，采用链表来解决;
4. hash 冲突严重时，链表会自动转化成红黑树，提高遍历速度。

**5、HashMap 是如何扩容的？**

扩容的时机：

1. put 时，发现数组为空，进行初始化扩容，默认扩容大小为 16;
2. put 成功后，发现现有数组大小大于扩容的门阀值时，进行扩容，扩容为老数组大小的 2 倍;

扩容的门阀是 threshold，每次扩容时 threshold 都会被重新计算，门阀值等于数组的大小 * 影响因子（0.75）。

新数组初始化之后，需要将老数组的值拷贝到新数组上，链表和红黑树都有自己拷贝的方法。

**6、hash 冲突时怎么办？**

hash 冲突指的是 key 值的 hashcode 计算相同，但 key 值不同的情况。

如果桶中元素原本只有一个或已经是链表了，新增元素直接追加到链表尾部；

如果桶中元素已经是链表，并且链表个数大于等于 8 时，此时有两种情况：

1. 如果此时数组大小小于 64，数组再次扩容，链表不会转化成红黑树;
2. 如果数组大小大于 64 时，链表就会转化成红黑树。

这里不仅仅判断链表个数大于等于 8，还判断了数组大小，数组容量小于 64 没有立即转化的原因，猜测主要是因为红黑树占用的空间比链表大很多，转化也比较耗时，所以数组容量小的情况下冲突严重，我们可以先尝试扩容，看看能否通过扩容来解决冲突的问题。

**7、为什么链表个数大于等于 8 时，链表要转化成红黑树了？**

当链表个数太多了，遍历可能比较耗时，转化成红黑树，可以使遍历的时间复杂度降低，但转化成红黑树，有空间和转化耗时的成本，我们通过泊松分布公式计算，正常情况下，链表个数出现 8 的概念不到千万分之一，所以说正常情况下，链表都不会转化成红黑树，这样设计的目的，是为了防止非正常情况下，比如 hash 算法出了问题时，导致链表个数轻易大于等于 8 时，仍然能够快速遍历。

**延伸问题：红黑树什么时候转变成链表？**

当节点的个数小于等于 6 时，红黑树会自动转化成链表，主要还是考虑红黑树的空间成本问题，当节点个数小于等于 6 时，遍历链表也很快，所以红黑树会重新变成链表。

**8、HashMap 在 put 时，如果数组中已经有了这个 key，我不想把 value 覆盖怎么办？取值时，如果得到的 value 是空时，想返回默认值怎么办？**

如果数组有了 key，但不想覆盖 value ，可以选择 putIfAbsent 方法，这个方法有个内置变量 onlyIfAbsent，内置是 true ，就不会覆盖，我们平时使用的 put 方法，内置 onlyIfAbsent 为 false，是允许覆盖的。

取值时，如果为空，想返回默认值，可以使用 getOrDefault 方法，方法第一参数为 key，第二个参数为你想返回的默认值，如 map.getOrDefault(“2”,“0”)，当 map 中没有 key 为 2 的值时，会默认返回 0，而不是空。

**9、DTO 作为 Map 的 key 时，有无需要注意的点？**

看是什么类型的 Map，如果是 HashMap 的话，一定需要覆写 equals 和 hashCode 方法，因为在 get 和 put 的时候，需要通过 equals 方法进行相等的判断；如果是 TreeMap 的话，DTO 需要实现 Comparable 接口，因为 TreeMap 会使用 Comparable 接口进行判断 key 的大小；如果是 LinkedHashMap 的话，和 HashMap 一样的。

**10、LinkedHashMap 中的 LRU 是什么意思，是如何实现的？**

LRU ，英文全称：Least recently used，中文叫做最近最少访问，在 LinkedHashMap 中，也叫做最少访问删除策略，我们可以通过 removeEldestEntry 方法设定一定的策略，使最少被访问的元素，在适当的时机被删除，原理是在 put 方法执行的最后，LinkedHashMap 会去检查这种策略，如果满足策略，就删除头节点。

保证头节点就是最少访问的元素的原理是：LinkedHashMap 在 get 的时候，都会把当前访问的节点，移动到链表的尾部，慢慢的，就会使头部的节点都是最少被访问的元素。

**11、为什么推荐 TreeMap 的元素最好都实现 Comparable 接口？但 key 是 String 的时候，我们却没有额外的工作呢？**

因为 TreeMap 的底层就是通过排序来比较两个 key 的大小的，所以推荐 key 实现 Comparable 接口，是为了往你希望的排序顺序上发展， 而 String 本身已经实现了 Comparable 接口，所以使用 String 时，我们不需要额外的工作，不仅仅是 String ，其他包装类型也都实现了 Comparable 接口，如 Long、Double、Short 等等。
