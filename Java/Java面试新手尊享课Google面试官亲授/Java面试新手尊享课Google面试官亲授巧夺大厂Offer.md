[TOC]



## 1 操作系统

### 进程和线程

- 了解面试者侧重点
- 进一步展开

进程之间不能共享内存，进程间通信方法有

<img src="进程.png" style="zoom:50%;" />

<img src="线程.png" style="zoom:50%;" />

### 存储和寻址

#### 存储

<img src="存储.png" style="zoom:50%;" />

#### 寻址空间

32位      4G

64位      10^19Bytes

#### 寻址过程

<img src="寻址过程.png" style="zoom: 50%;" />

### 进程间通信

#### 进程间通信的方法

##### 文件

##### signal

###### 示例1

Crtl + c

###### 示例2

```bash
#kill -l
kill -2 pid
# 效果同Ctrl + c
kill -9 pid
# 强制杀进程
```

##### 消息队列

##### 管道/命令管道

```bash
cat oki.log | grep "ERROR" --color
```

##### 共享内存

##### 同步机制

- 信号量

##### Socket

## 2 网络

### 网络基础

<img src="网络基础.png" style="zoom:50%;" />

#### 网络传输

##### 不可靠

- 丢包，重复包
- 出错
- 乱序

##### 不安全

- 中间人攻击
- 窃取
- 篡改

### 滑动窗口问题

- TCP协议中使用，解决网络传输不可靠性
- 维持发送方/接收方缓冲区

### 网络抓包演示

### TCP链接建立与断开

#### 建立连接（三次握手），断开连接（四次挥手）

- 网络是不可靠的，任何包都有可能丢
- 遇见问题，解决问题，不断迭代

##### 建立连接（三次握手）

<img src="三次握手.png" style="zoom:50%;" />



##### 断开连接（四次挥手）

<img src="四次挥手.png" style="zoom:50%;" />

## 3 数据库

### 关系型数据库

- 基于关系代数理论
- 缺点：表结构不直观，实现复杂，速度慢
- 优点：健壮性高，社区庞大

### JOIN和GROUP BY

- 笛卡尔积
- group by 字段需要写在select语句中，select语句多配合聚合函数使用
- 子查询

```mysql
SELECT p.*, cat_min.categoryName FROM `product` p join (SELECT p.`categoryId`, p.`categoryName`, MIN(p.`price`)  AS min_price FROM `product` p JOIN `category` c on p.`categoryId` = c.`categoryId` GROUP BY p.`categoryId`, p.`categoryName`) AS cat_min on p.`categoryId` = cat_min.categoryId WHERE p.`price` = cat_min.min_price
```

### 事务和乐观锁

#### 事务

##### ACID属性

- Atomic
- Consistency
- Isolution
- Durability

##### 隔离级别

- Read Uncommitted
- Read Committed
- Repetable Reads
- Serializable

```mysql
# SET SESSION TRANSACTION ISOLATION LEVEL REPEATABLE READ;
BEGIN;
SET AUTOCOMMIT = 0;
SELECT @@tx_isolation;
SELECT count FROM `product` WHERE `productId` = 2;

SELECT count FROM `product` WHERE `productId` = 2 for update;
# 加行级锁

UPDATE `product` SET `count` = 49 WHERE `productId` = 2;
COMMIT;
# ROLLBACK;
```

#### 乐观锁

- 版本控制
- timestamp

```mysql
BEGIN;

SELECT count FROM `product` WHERE `productId` = 2;

UPDATE `product` SET `count` = 49 WHERE `productId` = 2 AND `count` = 50;

COMMIT;
```

