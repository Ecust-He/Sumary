[TOC]



# 并发编程

**常见名词**

- JUC
- CAS
- AQS

## 线程基础

### 线程创建

1、创建线程有哪几种方式？哪种方式更好？

```java
public void run() {
    if (target != null) {
        target.run();
    }
}
```

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

### 线程生命周期

## 线程安全

### 并发原子类

### 并发容器

### 锁

#### 死锁问题

## ThreadLocal

### 使用场景

## 线程治理

### 线程池

#### 获取子线程执行结果



## 线程协作

### CountDownLatch门闩

#### 使用场景

- 一等多（所有运动员都到达了终点）
- 多等一（一声枪响）

# 高效编程

## 流式编程
