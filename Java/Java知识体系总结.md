[TOC]

# JAVA基础

## 异常体系



# 并发编程

**常见名词**

- JUC
- CAS
- AQS

## 1  线程基础

### 线程创建

1、创建线程有哪几种方式？哪种方式更好？

本质一样，run方法内容来源不同

```java
public void run() {
    if (target != null) {
        target.run();
    }
}
```

2、创建线程时，既传入了Runnable对象，又继承了Thread类重写run方法，此时执行哪一个run方法？

### 线程启动

1、Thread#start方法调用两次会怎样？

```java
/**
 * Causes this thread to begin execution; the Java Virtual Machine
 * calls the <code>run</code> method of this thread.
 * <p>
 * The result is that two threads are running concurrently: the
 * current thread (which returns from the call to the
 * <code>start</code> method) and the other thread (which executes its
 * <code>run</code> method).
 * <p>
 * It is never legal to start a thread more than once.
 * In particular, a thread may not be restarted once it has completed
 * execution.
 *
 * @exception  IllegalThreadStateException  if the thread was already
 *               started.
 * @see        #run()
 * @see        #stop()
 */
public synchronized void start() {
    /**
     * This method is not invoked for the main method thread or "system"
     * group threads created/set up by the VM. Any new functionality added
     * to this method in the future may have to also be added to the VM.
     *
     * A zero status value corresponds to state "NEW".
     */
    if (threadStatus != 0)
        throw new IllegalThreadStateException();

    /* Notify the group that this thread is about to be started
     * so that it can be added to the group's list of threads
     * and the group's unstarted count can be decremented. */
    group.add(this);

    boolean started = false;
    try {
        start0();
        started = true;
    } finally {
        try {
            if (!started) {
                group.threadStartFailed(this);
            }
        } catch (Throwable ignore) {
            /* do nothing. If start0 threw a Throwable then
              it will be passed up the call stack */
        }
    }
}
```

2、调用Thread#start方法启动线程和和直接调用Thread#run方法的区别。

### 线程停止

1、如何正确停止线程？

**使用通知方式而不是强制停止**

- 请求方：发出中断信号
- 被停止方：及时检测并响应中断
- 子方法调用方：恢复中断状态

### 生命周期

### Thread和Object类方法

#### sleep

暂停当前线程，线程进入阻塞状态，把CPU时间片让给其他线程执行

- 线程睡眠期间会释放占用的CPU时间片，但不会释放对象锁（如果当前线程获取到对象锁）
- 线程苏醒时，继续执行

#### yield

当前线程让出CPU时间片，线程状态转换为就绪状态，并重新竞争CPU的调度权

- yield设计初衷是为了防止过度使用CPU
- 让其他线程有机会运行，但不能保证某个特定的线程能够获得CPU资源；是否可以获取CPU时间片运行取决于调度器

##### 使用场景

当前线程让出CPU的使用权

#### join

主线程等待子线程执行完毕之后才会执行

```java
synchronized (thread) {
    thread.wait();
}
```

#### wait/notify

必须放在synchronized同步代码块中执行，为了避免使用者出现Lost Wake-Up Problem（通知丢失问题）

- wait会释放当前持有的对象锁

##### 生产者消费者模式演示

```java
public class ProducerConsumerModel {
    public static void main(String[] args) {
        EventStorage eventStorage = new EventStorage();
        Producer producer = new Producer(eventStorage);
        Consumer consumer = new Consumer(eventStorage);
        new Thread(producer).start();
        new Thread(consumer).start();
    }
}

class Producer implements Runnable {

    private EventStorage storage;

    public Producer(
            EventStorage storage) {
        this.storage = storage;
    }

    @Override
    public void run() {
        for (int i = 0; i < 100; i++) {
            storage.put();
        }
    }
}

class Consumer implements Runnable {

    private EventStorage storage;

    public Consumer(
            EventStorage storage) {
        this.storage = storage;
    }

    @Override
    public void run() {
        for (int i = 0; i < 100; i++) {
            storage.take();
        }
    }
}

class EventStorage {

    private int maxSize;
    private LinkedList<Date> storage;

    public EventStorage() {
        maxSize = 10;
        storage = new LinkedList<>();
    }

    public synchronized void put() {
        while (storage.size() == maxSize) {
            try {
                wait();
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
        }
        storage.add(new Date());
        System.out.println("仓库里有了" + storage.size() + "个产品。");
        notify();
    }

    public synchronized void take() {
        while (storage.size() == 0) {
            try {
                wait();
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
        }
        System.out.println("拿到了" + storage.poll() + "，现在仓库还剩下" + storage.size());
        notify();
    }
}
```

##### 两个线程交替打印0~100的奇偶数

###### 用synchronized关键字实现

```java
public class WaitNotifyPrintOddEvenSyn {

    private static int count;

    private static final Object lock = new Object();

    //新建2个线程
    //1个只处理偶数，第二个只处理奇数（用位运算）
    //用synchronized来通信
    public static void main(String[] args) {
        new Thread(new Runnable() {
            @Override
            public void run() {
                while (count < 100) {
                    synchronized (lock) {
                        if ((count & 1) == 0) {
                            System.out.println(Thread.currentThread().getName() + ":" + count++);
                        }
                    }
                }
            }
        }, "偶数").start();

        new Thread(new Runnable() {
            @Override
            public void run() {
                while (count < 100) {
                    synchronized (lock) {
                        if ((count & 1) == 1) {
                            System.out.println(Thread.currentThread().getName() + ":" + count++);
                        }
                    }
                }
            }
        }, "奇数").start();
    }
}
```

###### 用wait和notify方法实现

```java
public class WaitNotifyPrintOddEveWait {

    private static int count = 0;
    private static final Object lock = new Object();


    public static void main(String[] args) {
        new Thread(new TurningRunner(), "偶数").start();
        new Thread(new TurningRunner(), "奇数").start();
    }

    //1. 拿到锁，我们就打印
    //2. 打印完，唤醒其他线程，自己就休眠
    static class TurningRunner implements Runnable {

        @Override
        public void run() {
            while (count <= 100) {
                synchronized (lock) {
                    //拿到锁就打印
                    System.out.println(Thread.currentThread().getName() + ":" + count++);
                    lock.notify();
                    if (count <= 100) {
                        try {
                            //如果任务还没结束，就让出当前的锁，并休眠
                            lock.wait();
                        } catch (InterruptedException e) {
                            e.printStackTrace();
                        }
                    }
                }
            }
        }
    }
}
```

#### 面试题

##### 1、sleep和yield方法的区别？

- 调用sleep方法，线程苏醒后继续之前的执行；调用yield方法，当前线程进入就绪状态，是否运行取决于CPU的调度

##### 2、sleep和wait方法的区别？

- wait是Object类的方法，而sleep是Thread类的方法
- wait方法用于线程间通信，sleep方法用于短时间暂停当前线程
- wait方法需要在synchronized同步代码块中执行，而sleep方法不需要
- wait方法会让出对象锁，而sleep方法不会释放对象锁
- sleep方法可以通过interrupt方法中断，而wait方法不可以
- sleep方法需要等到睡眠时间自动苏醒，而wait方法可以通过notify或notifyAll方法随时唤醒

### 异常处理

为线程自定义异常处理类

## 2  底层原理

### Java内存模型

线程对内存的访问（读写操作）没有统一的规范，会出现**缓存不一致**问题；为了解决这个问题，JVM屏蔽了底层硬件和操作系统的差异，抽象出**主内存**和**工作内存**的概念；Java内存模型主要围绕在并发过程中如何处理**原子性**、**可见性**、**有序性**三大特性展开的，可以解除缓存一致性问题。

#### 特性

- 原子性
- 可见性
- 有序性

### happens-before原则

#### 定义

A happens before B，即A先于B发生，则A的执行结果对B可见。

#### 规则

JMM定义了如下的happens-before原则，已此保证有序性。

- 程序顺序规则：一个线程（单线程）中的每一个操作，happens-before于该线程后续任意操作

- 锁规则：对一个线程的**加锁**happens-before线程的**解锁**

- volatile规则：volatile变量的**写**happens-before volatile变量的**读**

- 传递性：A happens-before B，且B happens-before C，则A happens-before C

### 指令重排序

定义：编译器和处理器为了优化程序性能而对指令进行重新排序。

#### as-if-serial语义

不管指令如何重排序，（单线程）最终的执行结果不能改变。

### volatile关键字

#### 性质

- 可见性
- 禁止指令重排序优化（底层实现是使用**内存屏障**）

#### 使用场景

- 纯赋值操作（不适用i++等复合操作）
- 作为触发器

### 单例模式

| 模式                        | 特点                     | 是否线程安全       |
| --------------------------- | ------------------------ | ------------------ |
| 懒汉模式                    | 懒加载，使用时进行实例化 | 否                 |
| 懒汉模式+synchronized关键字 |                          | 是（不推荐）       |
| 双重检查模式                |                          | 是（推荐，面试用） |
| 饿汉模式（静态常量）        | 类加载时进行初始化       | 是                 |
| 饿汉模式（静态代码块）      | 类加载时进行初始化       | 是                 |
| 内部静态类                  |                          | 是                 |
| 枚举类                      |                          | 是                 |

```java
// 双重检查
public class Singleton {

    // volatile禁止指令重排序
    private volatile static Singleton instance;

    private Singleton() {}

    public static Singleton getInstance() {
        if (instance == null) {
            synchronized (Singleton6.class) {
                if (instance == null) {
                    // 创建对象 1、在堆区分配内存 2、实例化对象 3、引用指向对象
                    instance = new Singleton();
                }
            }
        }
        return instance;
    }
}
```

## 3  线程安全

多线程竞争共享资源时，需要保证同一时刻有且只有一个线程在操作共享资源

### 并发原子类

### 并发容器

#### CopyOnWriteArrayList

##### 源码实现

底层数据结构是**数组**，当向数组中添加、删除元素时，拷贝出一个新的数组并在新数组上完成更新操作，当更新操作完成之后再将数组的引用指向更新后的新数组。

加锁时使用ReentrantLock锁，加锁的粒度是整个方法体。

#### ConcurrentHashMap

##### 源码实现

底层数据结构是**数组+链表+红黑树**，并通过cas + synchronized关键字保证线程安全

###### 构造方法

若不指定初始容量（无参数构造方法），默认是16；

```java
/**
 * The default initial table capacity.  Must be a power of 2
 * (i.e., at least 1) and at most MAXIMUM_CAPACITY.
 */	
private static final int DEFAULT_CAPACITY = 16;
```

若指定初始化容量，如果大于MAXIMUM_CAPACITY的一半则取MAXIMUM_CAPACITY；否则取大于输入参数且最近的**2的整数次幂**作为初始化容量；

```java
public ConcurrentHashMap(int initialCapacity) {
    if (initialCapacity < 0)
        throw new IllegalArgumentException();
    int cap = ((initialCapacity >= (MAXIMUM_CAPACITY >>> 1)) ?
               MAXIMUM_CAPACITY :
               tableSizeFor(initialCapacity + (initialCapacity >>> 1) + 1));
    this.sizeCtl = cap;
}
```

###### put方法

```java
/** Implementation for put and putIfAbsent */
final V putVal(K key, V value, boolean onlyIfAbsent) {
    if (key == null || value == null) throw new NullPointerException();
    int hash = spread(key.hashCode());
    int binCount = 0;
    for (Node<K,V>[] tab = table;;) {
        Node<K,V> f; int n, i, fh;
        // Node为空，初始化Node
        if (tab == null || (n = tab.length) == 0)
            tab = initTable();
        // CAS对指定位置的节点进行原子操作
        else if ((f = tabAt(tab, i = (n - 1) & hash)) == null) {
            if (casTabAt(tab, i, null,
                         new Node<K,V>(hash, key, value, null)))
                break;                   // no lock when adding to empty bin
        }
        // 如果Node的hash值等于-1，Map进行扩容
        else if ((fh = f.hash) == MOVED)
            tab = helpTransfer(tab, f);
        else {
            V oldVal = null;
            synchronized (f) {
                if (tabAt(tab, i) == f) {
                    if (fh >= 0) {
                        binCount = 1;
                        for (Node<K,V> e = f;; ++binCount) {
                            K ek;
                            if (e.hash == hash &&
                                ((ek = e.key) == key ||
                                 (ek != null && key.equals(ek)))) {
                                // hash冲突时，替换原来的值
                                oldVal = e.val;
                                if (!onlyIfAbsent)
                                    e.val = value;
                                break;
                            }
                            Node<K,V> pred = e;
                            if ((e = e.next) == null) {
                                // 添加后继节点
                                pred.next = new Node<K,V>(hash, key,
                                                          value, null);
                                break;
                            }
                        }
                    }
                    else if (f instanceof TreeBin) {
                        Node<K,V> p;
                        binCount = 2;
                        if ((p = ((TreeBin<K,V>)f).putTreeVal(hash, key,
                                                       value)) != null) {
                            oldVal = p.val;
                            if (!onlyIfAbsent)
                                p.val = value;
                        }
                    }
                }
            }
            if (binCount != 0) {
                if (binCount >= TREEIFY_THRESHOLD)
                    treeifyBin(tab, i);
                if (oldVal != null)
                    return oldVal;
                break;
            }
        }
    }
    addCount(1L, binCount);
    return null;
}
```

### 不可变对象

#### final关键字

|          | 说明                     |
| -------- | ------------------------ |
| 修饰类   | 类不可被继承             |
| 修饰方法 | 方法不可在子类中重写     |
| 修饰对象 | 引用不可重新指向新的对象 |

### 锁

#### 分类

| 分类依据                           | 判断结果 | 类型                                       |
| ---------------------------------- | -------- | ------------------------------------------ |
| 是否对资源加锁                     | 是       | 悲观锁                                     |
|                                    | 否       | 乐观锁                                     |
| 不同线程是否可以共享同一把锁       | 可以     | 共享锁                                     |
|                                    | 不可以   | 独占锁                                     |
| 同一个线程是否可以重复获取同一把锁 | 可以     | 可重入锁                                   |
|                                    | 不可以   | 不可重入锁                                 |
| 并发冲突时是否自旋转               | 自旋     | 自旋锁                                     |
|                                    | 阻塞     | 非自旋锁                                   |
| 加锁时是否可中断                   | 可以     | 可中断锁                                   |
|                                    | 不可以   | 不可中断锁                                 |
| 获取锁时是否可以插队               | 可以     | 非公平锁（先尝试插队，如果插队失败再排队） |
|                                    | 不可以   | 公平锁                                     |

#### sychornized

##### 底层原理

###### 对象锁机制

```bash
javap -v xx.class
```

查看java的字节码文件，monitorenter和monitorexit两个指令；字面理解就是监视进入，监视退出。可以理解为代码块执行前的加锁，和退出同步时的解锁。

```bash
    4: monitorenter
    5: getstatic    #9    // Field java/lang/System.out:Ljava/io/PrintStream; 
    8: ldc           #15   // String hello
    10: invokevirtual #17  // Method java/io/PrintStream.println:(Ljava/lang/String;)V
    13: aload_1
    14: monitorexit
```

##### 特点

不需要显示释放锁

#### Lock锁

##### Sync源码实现

```java
abstract static class Sync extends AbstractQueuedSynchronizer {
    private static final long serialVersionUID = -5179523762034025860L;

    /**
     * Performs {@link Lock#lock}. The main reason for subclassing
     * is to allow fast path for nonfair version.
     */
    abstract void lock();

    /**
     * Performs non-fair tryLock.  tryAcquire is implemented in
     * subclasses, but both need nonfair try for trylock method.
     */
    final boolean nonfairTryAcquire(int acquires) {
        final Thread current = Thread.currentThread();
        int c = getState();
        if (c == 0) {
            if (compareAndSetState(0, acquires)) {
                // 将当前线程设置为锁的持有者
                setExclusiveOwnerThread(current);
                return true;
            }
        }
        // 可重入锁逻辑
        else if (current == getExclusiveOwnerThread()) {
            int nextc = c + acquires;
            if (nextc < 0) // overflow
                throw new Error("Maximum lock count exceeded");
            setState(nextc);
            return true;
        }
        return false;
    }

    protected final boolean tryRelease(int releases) {
        int c = getState() - releases;
        if (Thread.currentThread() != getExclusiveOwnerThread())
            throw new IllegalMonitorStateException();
        boolean free = false;
        if (c == 0) {
            free = true;
            setExclusiveOwnerThread(null);
        }
        setState(c);
        return free;
    }

    protected final boolean isHeldExclusively() {
        // While we must in general read state before owner,
        // we don't need to do so to check if current thread is owner
        return getExclusiveOwnerThread() == Thread.currentThread();
    }

    final ConditionObject newCondition() {
        return new ConditionObject();
    }

    // Methods relayed from outer class

    final Thread getOwner() {
        return getState() == 0 ? null : getExclusiveOwnerThread();
    }

    final int getHoldCount() {
        return isHeldExclusively() ? getState() : 0;
    }

    final boolean isLocked() {
        return getState() != 0;
    }

    /**
     * Reconstitutes the instance from a stream (that is, deserializes it).
     */
    private void readObject(java.io.ObjectInputStream s)
        throws java.io.IOException, ClassNotFoundException {
        s.defaultReadObject();
        setState(0); // reset to unlocked state
    }
}
```

##### ReentrantLock锁

###### 公平锁和非公平锁

- 先尝试插队，如果插队失败再排队
- 非公平锁性能好，但可能会导致饥饿情况；由于频繁的线程上下文切换，公平锁性能较差

**NonfairSync 非公平锁源码实现**

```java
static final class NonfairSync extends Sync {
    private static final long serialVersionUID = 7316153563782823691L;

    /**
     * Performs lock.  Try immediate barge, backing up to normal
     * acquire on failure.
     */
    final void lock() {
        if (compareAndSetState(0, 1))
           // 将当前线程设置为锁的持有者
            setExclusiveOwnerThread(Thread.currentThread());
        else
            acquire(1);
    }

    protected final boolean tryAcquire(int acquires) {
        return nonfairTryAcquire(acquires);
    }
}
```

**FairSync公平锁源码实现**

```java
static final class FairSync extends Sync {
    private static final long serialVersionUID = -3000897897090466540L;

    final void lock() {
        acquire(1);
    }

    /**
     * Fair version of tryAcquire.  Don't grant access unless
     * recursive call or no waiters or is first.
     */
    protected final boolean tryAcquire(int acquires) {
        final Thread current = Thread.currentThread();
        int c = getState();
        if (c == 0) {
            if (!hasQueuedPredecessors() &&
                compareAndSetState(0, acquires)) {
                setExclusiveOwnerThread(current);
                return true;
            }
        }
        else if (current == getExclusiveOwnerThread()) {
            int nextc = c + acquires;
            if (nextc < 0)
                throw new Error("Maximum lock count exceeded");
            setState(nextc);
            return true;
        }
        return false;
    }
}
```

```java
// 主要用来判断线程是否需要排队
public final boolean hasQueuedPredecessors() {
    // The correctness of this depends on head being initialized
    // before tail and on head.next being accurate if the current
    // thread is first in queue.
    Node t = tail; // Read fields in reverse initialization order
    Node h = head;
    Node s;
    return h != t &&
        ((s = h.next) == null || s.thread != Thread.currentThread());
    // 结果一：返回false，不需要排队
    	// 情况一：队列为空
    	// 情况二：头节点有后继节点，且后继节点和头节点是相同线程
}
```

**代码演示**

```java
public class FairLock {

    public static void main(String[] args) {
        PrintQueue printQueue = new PrintQueue();
        Thread thread[] = new Thread[10];
        for (int i = 0; i < 10; i++) {
            thread[i] = new Thread(new Job(printQueue));
        }
        for (int i = 0; i < 10; i++) {
            thread[i].start();
            try {
                Thread.sleep(100);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
        }
    }
}

class Job implements Runnable {

    PrintQueue printQueue;

    public Job(PrintQueue printQueue) {
        this.printQueue = printQueue;
    }

    @Override
    public void run() {
        System.out.println(Thread.currentThread().getName() + "开始打印");
        printQueue.printJob(new Object());
        System.out.println(Thread.currentThread().getName() + "打印完毕");
    }
}

class PrintQueue {

    private Lock queueLock = new ReentrantLock(false);

    public void printJob(Object document) {
        queueLock.lock();
        try {
            int duration = new Random().nextInt(10) + 1;
            System.out.println(Thread.currentThread().getName() + "正在打印，需要" + duration);
            Thread.sleep(duration * 1000);
        } catch (InterruptedException e) {
            e.printStackTrace();
        } finally {
            queueLock.unlock();
        }

        queueLock.lock();
        try {
            int duration = new Random().nextInt(10) + 1;
            System.out.println(Thread.currentThread().getName() + "正在打印，需要" + duration+"秒");
            Thread.sleep(duration * 1000);
        } catch (InterruptedException e) {
            e.printStackTrace();
        } finally {
            queueLock.unlock();
        }
    }
}
```

###### 可重入锁

###### 读写锁

**特点**

- 多读一写
- 写锁可降级，读锁不可升级
- 若当前持有的锁是读锁，后面是一个写线程，那么写线程后面的读线程无法获取读锁，这是jdk为了避免写线程过分饥渴采取的策略。

**读锁源码实现**

```java
public final void acquireShared(int arg) {
    if (tryAcquireShared(arg) < 0)
        doAcquireShared(arg);
}

// 返回值小于 0 代表没有获取到共享锁(读锁)，大于 0 代表获取到
protected final int tryAcquireShared(int unused) {
	/*
	 * Walkthrough:
	 * 1. If write lock held by another thread, fail.
	 * 2. Otherwise, this thread is eligible for
	 *    lock wrt state, so ask if it should block
	 *    because of queue policy. If not, try
	 *    to grant by CASing state and updating count.
	 *    Note that step does not check for reentrant
	 *    acquires, which is postponed to full version
	 *    to avoid having to check hold count in
	 *    the more typical non-reentrant case.
	 * 3. If step 2 fails either because thread
	 *    apparently not eligible or CAS fails or count
	 *    saturated, chain to version with full retry loop.
	 */
	Thread current = Thread.currentThread();
	int c = getState();
    // 有其他线程持有写锁，当前线程不能获取到读锁
	if (exclusiveCount(c) != 0 &&
		getExclusiveOwnerThread() != current)
		return -1;
	int r = sharedCount(c);
	if (!readerShouldBlock() &&
		r < MAX_COUNT &&
		compareAndSetState(c, c + SHARED_UNIT)) {
        // 第一个获取读锁的线程
		if (r == 0) {
			firstReader = current;
			firstReaderHoldCount = 1;
		} else if (firstReader == current) {// 重入获取读锁
			firstReaderHoldCount++;
		} else {
			HoldCounter rh = cachedHoldCounter;
			if (rh == null || rh.tid != getThreadId(current))
				cachedHoldCounter = rh = readHolds.get();
			else if (rh.count == 0)
				readHolds.set(rh);
			rh.count++;
		}
		return 1;
	}
	return fullTryAcquireShared(current);
}
```

**写锁源码实现**

```java
public void lock() {
    sync.acquire(1);
}

public final void acquire(int arg) {
    // 如果获取锁失败，则进入到阻塞队列等待
	if (!tryAcquire(arg) &&
		acquireQueued(addWaiter(Node.EXCLUSIVE), arg))
		selfInterrupt();
}

protected final boolean tryAcquire(int acquires) {
	/*
	 * Walkthrough:
	 * 1. If read count nonzero or write count nonzero
	 *    and owner is a different thread, fail.
	 * 2. If count would saturate, fail. (This can only
	 *    happen if count is already nonzero.)
	 * 3. Otherwise, this thread is eligible for lock if
	 *    it is either a reentrant acquire or
	 *    queue policy allows it. If so, update state
	 *    and set owner.
	 */
	Thread current = Thread.currentThread();
	int c = getState();
	int w = exclusiveCount(c);
	if (c != 0) {
		// (Note: if c != 0 and w == 0 then shared count != 0)
        // 有其他线程获取写锁，则不能获取写锁
		if (w == 0 || current != getExclusiveOwnerThread())
			return false;
		if (w + exclusiveCount(acquires) > MAX_COUNT)
			throw new Error("Maximum lock count exceeded");
		// Reentrant acquire
		setState(c + acquires);
		return true;
	}
	if (writerShouldBlock() ||
		!compareAndSetState(c, c + acquires))
		return false;
	setExclusiveOwnerThread(current);
	return true;
}
```

**使用场景**

- 读写分离场景，降低加锁粒度

**代码示例**

```java
// 这是一个关于缓存操作的故事
class CachedData {
    Object data;
    volatile boolean cacheValid;
    // 读写锁实例
    final ReentrantReadWriteLock rwl = new ReentrantReadWriteLock();

    void processCachedData() {
        // 获取读锁
        rwl.readLock().lock();
        if (!cacheValid) { // 如果缓存过期了，或者为 null
            // 释放掉读锁，然后获取写锁 (后面会看到，没释放掉读锁就获取写锁，会发生死锁情况)
            rwl.readLock().unlock();
            rwl.writeLock().lock();

            try {
                if (!cacheValid) { // 重新判断，因为在等待写锁的过程中，可能前面有其他写线程执行过了
                    data = ...
                    cacheValid = true;
                }
                // 获取读锁 (持有写锁的情况下，是允许获取读锁的，称为 “锁降级”，反之不行。)
                rwl.readLock().lock();
            } finally {
                // 释放写锁，此时还剩一个读锁
                rwl.writeLock().unlock(); // Unlock write, still hold read
            }
        }

        try {
            use(data);
        } finally {
            // 释放读锁
            rwl.readLock().unlock();
        }
    }
}
```

```java
演示非公平和公平的ReentrantReadWriteLock的策略
    public class NonfairBargeDemo {

    private static ReentrantReadWriteLock reentrantReadWriteLock = new ReentrantReadWriteLock(
            true);

    private static ReentrantReadWriteLock.ReadLock readLock = reentrantReadWriteLock.readLock();
    private static ReentrantReadWriteLock.WriteLock writeLock = reentrantReadWriteLock.writeLock();

    private static void read() {
        System.out.println(Thread.currentThread().getName() + "开始尝试获取读锁");
        readLock.lock();
        try {
            System.out.println(Thread.currentThread().getName() + "得到读锁，正在读取");
            try {
                Thread.sleep(20);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
        } finally {
            System.out.println(Thread.currentThread().getName() + "释放读锁");
            readLock.unlock();
        }
    }

    private static void write() {
        System.out.println(Thread.currentThread().getName() + "开始尝试获取写锁");
        writeLock.lock();
        try {
            System.out.println(Thread.currentThread().getName() + "得到写锁，正在写入");
            try {
                Thread.sleep(40);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
        } finally {
            System.out.println(Thread.currentThread().getName() + "释放写锁");
            writeLock.unlock();
        }
    }

    public static void main(String[] args) {
        new Thread(()->write(),"Thread1").start();
        new Thread(()->read(),"Thread2").start();
        new Thread(()->read(),"Thread3").start();
        new Thread(()->write(),"Thread4").start();
        new Thread(()->read(),"Thread5").start();
        new Thread(new Runnable() {
            @Override
            public void run() {
                Thread thread[] = new Thread[1000];
                for (int i = 0; i < 1000; i++) {
                    thread[i] = new Thread(() -> read(), "子线程创建的Thread" + i);
                }
                for (int i = 0; i < 1000; i++) {
                    thread[i].start();
                }
            }
        }).start();
    }
}
```

#### 基于volatile + CAS实现同步锁

**CAS是乐观锁的一种典型实现**

乐观锁的实现主要两个步骤：

- 冲突检测
- 数据更新

##### CAS的缺陷

- ABA问题
- 循环时间长，开销大
- 只能保证一个共享变量的原子操作

#### 死锁问题

##### 死锁的四个必要条件

- 互斥条件（同一时刻，某个资源只能被一个线程占用）
- 请求与保持（**吃着碗里的看着锅里的**，即线程占用至少一个资源，又提出了新的资源请求，而该资源被其他线程占用，所以当前线程处于等待状态）
- 不可剥夺（已持有的资源在线程未使用完之前，不可被强行剥夺，只能等待获取资源的线程主动释放）
- 循环等待（**你等我，我等你**，线程间形成了头尾相接的循环等待资源关系）

##### 典型案例

###### 银行转账问题

```java
/**
 * 描述：     多人同时转账，依然很危险
 */
public class MultiTransferMoney {

    private static final int NUM_ACCOUNTS = 500;
    private static final int NUM_MONEY = 1000;
    private static final int NUM_ITERATIONS = 1000000;
    private static final int NUM_THREADS = 20;

    public static void main(String[] args) {

        Random rnd = new Random();
        Account[] accounts = new Account[NUM_ACCOUNTS];
        for (int i = 0; i < accounts.length; i++) {
            accounts[i] = new Account(NUM_MONEY);
        }
        class TransferThread extends Thread {

            @Override
            public void run() {
                for (int i = 0; i < NUM_ITERATIONS; i++) {
                    int fromAcct = rnd.nextInt(NUM_ACCOUNTS);
                    int toAcct = rnd.nextInt(NUM_ACCOUNTS);
                    int amount = rnd.nextInt(NUM_MONEY);
                    TransferMoney.transferMoney(accounts[fromAcct], accounts[toAcct], amount);
                }
                System.out.println("运行结束");
            }
        }
        for (int i = 0; i < NUM_THREADS; i++) {
            new TransferThread().start();
        }
    }
}
```

```java
public class TransferMoney implements Runnable {
    
    static Object lock = new Object();

    public static void transferMoney(Account from, Account to, int amount) {
        class Helper {
            public void transfer() {
                if (from.balance - amount < 0) {
                    System.out.println("余额不足，转账失败。");
                    return;
                }
                from.balance -= amount;
                to.balance = to.balance + amount;
                System.out.println("成功转账" + amount + "元");
            }
        }
        int fromHash = System.identityHashCode(from);
        int toHash = System.identityHashCode(to);
        if (fromHash < toHash) {
            synchronized (from) {
                synchronized (to) {
                    new Helper().transfer();
                }
            }
        }
        else if (fromHash > toHash) {
            synchronized (to) {
                synchronized (from) {
                    new Helper().transfer();
                }
            }
        }else  {
            synchronized (lock) {
                synchronized (to) {
                    synchronized (from) {
                        new Helper().transfer();
                    }
                }
            }
        }

    }

    static class Account {

        public Account(int balance) {
            this.balance = balance;
        }

        int balance;

    }
}
```

###### 哲学家就餐问题

```java
public class DiningPhilosophers {

    public static class Philosopher implements Runnable {

        private Object leftChopstick;

        private Object rightChopstick;

        public Philosopher(Object leftChopstick, Object rightChopstick) {
            this.leftChopstick = leftChopstick;
            this.rightChopstick = rightChopstick;
        }

        @Override
        public void run() {
            try {
                while (true) {
                    doAction("Thinking");
                    synchronized (leftChopstick) {
                        doAction("Picked up left chopstick");
                        synchronized (rightChopstick) {
                            doAction("Picked up right chopstick - eating");
                            doAction("Put down right chopstick");
                        }
                        doAction("Put down left chopstick");
                    }
                }
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
        }

        private void doAction(String action) throws InterruptedException {
            System.out.println(Thread.currentThread().getName() + " " + action);
            Thread.sleep((long) (Math.random() * 10));
        }
    }

    public static void main(String[] args) {
        Philosopher[] philosophers = new Philosopher[5];
        Object[] chopsticks = new Object[philosophers.length];
        for (int i = 0; i < chopsticks.length; i++) {
            chopsticks[i] = new Object();
        }
        for (int i = 0; i < philosophers.length; i++) {
            Object leftChopstick = chopsticks[i];
            Object rightChopstick = chopsticks[(i + 1) % chopsticks.length];
            if (i == philosophers.length - 1) {
                philosophers[i] = new Philosopher(rightChopstick, leftChopstick);
            } else {
                philosophers[i] = new Philosopher(leftChopstick, rightChopstick);
            }
            new Thread(philosophers[i], "哲学家" + (i + 1) + "号").start();
        }
    }
}
```

#### 活锁问题

```java
public class LiveLock {

    static class Spoon {

        private Diner owner;

        public Spoon(Diner owner) {
            this.owner = owner;
        }

        public Diner getOwner() {
            return owner;
        }

        public void setOwner(Diner owner) {
            this.owner = owner;
        }

        public synchronized void use() {
            System.out.printf("%s吃完了!", owner.name);
        }
    }

    static class Diner {

        private String name;
        private boolean isHungry;

        public Diner(String name) {
            this.name = name;
            isHungry = true;
        }

        public void eatWith(Spoon spoon, Diner spouse) {
            while (isHungry) {
                if (spoon.owner != this) {
                    try {
                        Thread.sleep(1);
                    } catch (InterruptedException e) {
                        e.printStackTrace();
                    }
                    continue;
                }
                Random random = new Random();
                if (spouse.isHungry && random.nextInt(10) < 9) {
                    System.out.println(name + ": 亲爱的" + spouse.name + "你先吃吧");
                    spoon.setOwner(spouse);
                    continue;
                }

                spoon.use();
                isHungry = false;
                System.out.println(name + ": 我吃完了");
                spoon.setOwner(spouse);

            }
        }
    }


    public static void main(String[] args) {
        Diner husband = new Diner("牛郎");
        Diner wife = new Diner("织女");

        Spoon spoon = new Spoon(husband);

        new Thread(new Runnable() {
            @Override
            public void run() {
                husband.eatWith(spoon, wife);
            }
        }).start();

        new Thread(new Runnable() {
            @Override
            public void run() {
                wife.eatWith(spoon, husband);
            }
        }).start();
    }
}
```

#### 注意事项

1. 锁的选取需要综合考虑是否超时获取锁、获取到锁时是否可被中断等条件
2. 加锁的粒度尽可能小
3. 防止死锁

### ThreadLocal

#### 源码解析

每个Thread维护一个ThreadLocalMap（自定义的Hash Map），key为ThreadLocal本身，value为ThreadLocal中保存的数据。

#### 使用场景

#### 优点

- 线程隔离，线程安全
- 避免参数层层传递，降低代码耦合度

#### 代码示例

**打印日期**

```java
// 使用ThreadLocal给每个线程分配一个SimpleDateFormat对象，保证线程安全
public class ThreadLocalNormalUsage {

    public static ExecutorService threadPool = Executors.newFixedThreadPool(10);

    public static void main(String[] args) throws InterruptedException {
        for (int i = 0; i < 1000; i++) {
            int finalI = i;
            threadPool.submit(() -> System.out.println( new ThreadLocalNormalUsage05().date(finalI)));
        }
        threadPool.shutdown();
    }

    public String date(int seconds) {
        //参数的单位是毫秒，从1970.1.1 00:00:00 GMT计时
        Date date = new Date(1000 * seconds);
        SimpleDateFormat dateFormat = ThreadSafeFormatter.dateFormatThreadLocal2.get();
        return dateFormat.format(date);
    }
}

class ThreadSafeFormatter {

    public static ThreadLocal<SimpleDateFormat> dateFormatThreadLocal = new ThreadLocal<SimpleDateFormat>() {
        @Override
        protected SimpleDateFormat initialValue() {
            return new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
        }
    };

    public static ThreadLocal<SimpleDateFormat> dateFormatThreadLocal2 = ThreadLocal
            .withInitial(() -> new SimpleDateFormat("yyyy-MM-dd HH:mm:ss"));
}
```

## 4  线程治理

### 线程池

- 用于管理线程，避免增加创建线程和销毁线程的资源消耗
- 提高响应速度
- 线程复用

#### 核心参数

| 参数            | 说明                                   |
| --------------- | -------------------------------------- |
| corePoolSize    | 核心线程数                             |
| maximumPoolSize | 最大线程数 = 核心线程数 + 非核心线程数 |
| keepAliveTime   | 非核心线程闲置超时时长                 |
| unit            | keepAliveTime的单位                    |
| workQueue       | 阻塞队列                               |
| threadFactory   | 线程工厂                               |
| handler         | 拒绝策略                               |

#### 工作流程

![](./线程池工作流程图.png)

#### 拒绝策略

当任务队列已满且线程池中线程数量达到maximumPoolSize时，执行任务的拒绝策略。

| 拒绝策略            | 说明                                           |
| ------------------- | ---------------------------------------------- |
| AbortPolicy（默认） | 丢弃任务并抛出RejectedExecutionException       |
| DiscardPolicy       | 丢弃任务，但不抛出异常                         |
| DiscardOldestPolicy | 丢弃队列中最前面的任务，并重新添加被拒绝的任务 |
| CallerRunsPolicy    | 又调用线程处理该任务                           |

#### 工作队列

#### 常见线程池

| 分类     | Executors的静态方法     | 阻塞队列            | 优劣 | 使用场景       |
| -------- | ----------------------- | ------------------- | ---- | -------------- |
| 固定线程 | newFixedThreadPool      | LinkedBlockingQueue | OOM  | 任务少、时间长 |
| 单线程   | newSingleThreadExecutor | LinkedBlockingQueue |      | 串行执行任务   |
| 带缓存   | newCachedThreadPool     | SynchronousQueue    |      | 任务多、时间短 |
| 可调度   | newScheduledThreadPool  | DelayedWorkQueue    |      | 周期性执行任务 |

##### newFixedThreadPool

##### newSingleThreadExecutor

##### newCachedThreadPool

##### newScheduledThreadPool

#### 操作子线程

##### 暂停/恢复子线程执行

```java
public class PauseableThreadPool extends ThreadPoolExecutor {

    private final ReentrantLock lock = new ReentrantLock();
    private Condition unpaused = lock.newCondition();
    private boolean isPaused;


    public PauseableThreadPool(int corePoolSize, int maximumPoolSize, long keepAliveTime,
            TimeUnit unit,
            BlockingQueue<Runnable> workQueue) {
        super(corePoolSize, maximumPoolSize, keepAliveTime, unit, workQueue);
    }

    public PauseableThreadPool(int corePoolSize, int maximumPoolSize, long keepAliveTime,
            TimeUnit unit, BlockingQueue<Runnable> workQueue,
            ThreadFactory threadFactory) {
        super(corePoolSize, maximumPoolSize, keepAliveTime, unit, workQueue, threadFactory);
    }

    public PauseableThreadPool(int corePoolSize, int maximumPoolSize, long keepAliveTime,
            TimeUnit unit, BlockingQueue<Runnable> workQueue,
            RejectedExecutionHandler handler) {
        super(corePoolSize, maximumPoolSize, keepAliveTime, unit, workQueue, handler);
    }

    public PauseableThreadPool(int corePoolSize, int maximumPoolSize, long keepAliveTime,
            TimeUnit unit, BlockingQueue<Runnable> workQueue,
            ThreadFactory threadFactory, RejectedExecutionHandler handler) {
        super(corePoolSize, maximumPoolSize, keepAliveTime, unit, workQueue, threadFactory,
                handler);
    }

    // 在每个任务执行之前添加钩子函数
    @Override
    protected void beforeExecute(Thread t, Runnable r) {
        super.beforeExecute(t, r);
        lock.lock();
        try {
            while (isPaused) {
                unpaused.await();
            }
        } catch (InterruptedException e) {
            e.printStackTrace();
        } finally {
            lock.unlock();
        }
    }

    private void pause() {
        lock.lock();
        try {
            isPaused = true;
        } finally {
            lock.unlock();
        }
    }

    public void resume() {
        lock.lock();
        try {
            isPaused = false;
            unpaused.signalAll();
        } finally {
            lock.unlock();
        }
    }
}
```

##### 取消子线程执行

##### 获取子线程执行结果

```java
// 批量提交任务，并获取子线程的执行结果
public class MultiFutures {

    public static void main(String[] args) throws InterruptedException {
        ExecutorService service = Executors.newFixedThreadPool(20);
        ArrayList<Future> futures = new ArrayList<>();
        for (int i = 0; i < 20; i++) {
            Future<Integer> future = service.submit(new CallableTask());
            futures.add(future);
        }
        Thread.sleep(5000);
        for (int i = 0; i < 20; i++) {
            Future<Integer> future = futures.get(i);
            try {
                Integer integer = future.get();
                System.out.println(integer);
            } catch (InterruptedException e) {
                e.printStackTrace();
            } catch (ExecutionException e) {
                e.printStackTrace();
            }
        }
    }

    static class CallableTask implements Callable<Integer> {

        @Override
        public Integer call() throws Exception {
            Thread.sleep(3000);
            return new Random().nextInt();
        }
    }
}
```

##### 获取子线程执行时抛出的异常

```java
// 抛出异常时机是在执行get方法
public class GetException {

    public static void main(String[] args) {
        ExecutorService service = Executors.newFixedThreadPool(20);
        Future<Integer> future = service.submit(new CallableTask());

        try {
            for (int i = 0; i < 5; i++) {
                System.out.println(i);
                Thread.sleep(500);
            }
            System.out.println(future.isDone());
            future.get();
            System.out.println(future.isDone());
        } catch (InterruptedException e) {
            e.printStackTrace();
            System.out.println("InterruptedException异常");
        } catch (ExecutionException e) {
            e.printStackTrace();
            System.out.println("ExecutionException异常");
        }
    }


    static class CallableTask implements Callable<Integer> {

        @Override
        public Integer call() throws Exception {
            throw new IllegalArgumentException("Callable抛出异常");
        }
    }
}
```

## 5  线程协作

同步工具类，用于控制并发流程

### CountDownLatch闭锁

#### 源码解析

#### 使用场景

**控制子流程的开始和结束**

- 一等多（当所有运动员都到达了终点，比赛结束）
- 多等一（一声枪响，比赛开始）

#### 代码示例

```java
public static void main(String[] args) throws InterruptedException {
    CountDownLatch begin = new CountDownLatch(1);
    CountDownLatch end = new CountDownLatch(5);
    
    ExecutorService service = Executors.newFixedThreadPool(5);
    for (int i = 0; i < 5; i++) {
        final int no = i + 1;
        Runnable runnable = new Runnable() {
            @Override
            public void run() {
                System.out.println("No." + no + "准备完毕，等待发令枪");
                try {
                    begin.await();
                    System.out.println("No." + no + "开始跑步了");
                    Thread.sleep((long) (Math.random() * 10000));
                    System.out.println("No." + no + "跑到终点了");
                } catch (InterruptedException e) {
                    e.printStackTrace();
                } finally {
                    end.countDown();
                }
            }
        };
        service.submit(runnable);
    }
    //裁判员检查发令枪...
    Thread.sleep(5000);
    System.out.println("发令枪响，比赛开始！");
    begin.countDown();

    end.await();
    System.out.println("所有人到达终点，比赛结束");
}
```

### Semaphore信号量

用户可以根据业务需要，通过获取和释放多个许可证，控制并发数

#### 源码解析

#### 使用场景

- 接口限流

### CyclicBarrier循环栅栏

循环阻塞一组线程执行

#### 源码分析

#### 使用场景

#### 代码演示

```java
public class CyclicBarrierDemo {
    public static void main(String[] args) {
        CyclicBarrier cyclicBarrier = new CyclicBarrier(5, new Runnable() {
            @Override
            public void run() {
                System.out.println("所有人都到场了， 大家统一出发！");
            }
        });
        for (int i = 0; i < 10; i++) {
            new Thread(new Task(i, cyclicBarrier)).start();
        }
    }

    static class Task implements Runnable{
        private int id;
        private CyclicBarrier cyclicBarrier;

        public Task(int id, CyclicBarrier cyclicBarrier) {
            this.id = id;
            this.cyclicBarrier = cyclicBarrier;
        }

        @Override
        public void run() {
            System.out.println("线程" + id + "现在前往集合地点");
            try {
                Thread.sleep((long) (Math.random()*10000));
                System.out.println("线程"+id+"到了集合地点，开始等待其他人到达");
                cyclicBarrier.await();
                System.out.println("线程"+id+"出发了");
            } catch (InterruptedException e) {
                e.printStackTrace();
            } catch (BrokenBarrierException e) {
                e.printStackTrace();
            }
        }
    }
}
```

### Condition

wait/notify等待唤醒机制的**升级版**

#### 源码分析

#### 使用场景

#### 代码演示

**生产者消费者模式**

```java
public class ConditionDemo {

    private int queueSize = 10;
    private PriorityQueue<Integer> queue = new PriorityQueue<Integer>(queueSize);
    private Lock lock = new ReentrantLock();
    private Condition notFull = lock.newCondition();
    private Condition notEmpty = lock.newCondition();

    public static void main(String[] args) {
        ConditionDemo2 conditionDemo2 = new ConditionDemo2();
        Producer producer = conditionDemo2.new Producer();
        Consumer consumer = conditionDemo2.new Consumer();
        producer.start();
        consumer.start();
    }

    class Consumer extends Thread {

        @Override
        public void run() {
            consume();
        }

        private void consume() {
            while (true) {
                lock.lock();
                try {
                    while (queue.size() == 0) {
                        System.out.println("队列空，等待数据");
                        try {
                            notEmpty.await();
                        } catch (InterruptedException e) {
                            e.printStackTrace();
                        }
                    }
                    queue.poll();
                    notFull.signalAll();
                    System.out.println("从队列里取走了一个数据，队列剩余" + queue.size() + "个元素");
                } finally {
                    lock.unlock();
                }
            }
        }
    }

    class Producer extends Thread {

        @Override
        public void run() {
            produce();
        }

        private void produce() {
            while (true) {
                lock.lock();
                try {
                    while (queue.size() == queueSize) {
                        System.out.println("队列满，等待有空余");
                        try {
                            notFull.await();
                        } catch (InterruptedException e) {
                            e.printStackTrace();
                        }
                    }
                    queue.offer(1);
                    notEmpty.signalAll();
                    System.out.println("向队列插入了一个元素，队列剩余空间" + (queueSize - queue.size()));
                } finally {
                    lock.unlock();
                }
            }
        }
    }
}
```

# 高效编程

## Lambda表达式

函数式编程风格，将函数作为方法入参（处理逻辑参数化）

### 函数式接口

#### 定义

- FunctionalInterface注解修饰
- 有且只有一个抽象方法

#### 常用函数式接口

| 接口           |                    | 说明     |
| -------------- | ------------------ | -------- |
| Predicate<T>   | boolean test(T t); | 断言     |
| Supplier<T>    | T get();           | 生成对象 |
| Consumer<T>    | void accept(T t);  | 消费数据 |
| Function<T, R> | R apply(T t);      | 数据转换 |

### 方法引用

通过方法引用替代Lambda表达式

|      | 语法             | Lambda表达式写法          | 双冒号::写法     |
| ---- | ---------------- | ------------------------- | ---------------- |
|      | 对象名::成员方法 | str -> str.toUpperCase()  | str::toUpperCase |
|      | 类名::静态方法   | n -> Math.abs(n)          | Math::abs        |
|      | 类名::构造方法   | name -> new Person(name)  | Person::new      |
|      | 数组::构造方法   | length -> new int[length] | int[]::new       |

## 流式编程

集合是面向存储，流是面向对象

### 流的组成

支持链式调用

#### 数据源

#### 中间操作

| 分类   | 方法                                                  | 描述 |
| ------ | ----------------------------------------------------- | ---- |
| 无状态 | peek                                                  | 遍历 |
|        | filter                                                | 过滤 |
|        | map、mapToInt、mapToLong、mapToDouble                 | 映射 |
|        | flatMap、flatMapToInt、flatMapToLong、flatMapToDouble | 平铺 |
| 有状态 | distinct                                              | 去重 |
|        | sorted                                                | 排序 |
|        | limit                                                 | 截取 |
|        | skip                                                  | 跳过 |

#### 终端操作

| 分类       | 方法                          | 描述 |
| ---------- | ----------------------------- | ---- |
| 非短路操作 | foreach                       | 遍历 |
|            | count                         | 计数 |
|            | min、max                      | 最值 |
|            | reduce                        | 归约 |
|            | collect                       | 收集 |
| 短路操作   | findFirst、findAny            | 查找 |
|            | anyMatch、allMatch、noneMatch | 匹配 |

### 流的创建

## Optional的使用

解决引用存在和引用缺失问题

## TWR关闭资源

