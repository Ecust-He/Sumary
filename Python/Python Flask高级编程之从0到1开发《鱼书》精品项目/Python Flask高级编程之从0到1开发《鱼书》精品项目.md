[TOC]

## 1  Python Flask高级编程之从0到1开发《鱼书》精品项目

### Flask的基本原理与核心知识

#### 使用官方推荐的pipenv创建虚拟环境

##### 安装pipenv

```bash
pip install pipenv
```

##### 为项目创建虚拟环境

```bash
pipenv install
```

##### 查看虚拟环境安装目录

```bash
pipenv --venv
```

##### 进入虚拟环境

```pip
pipenv shell
```

##### 查看安装列表

```python
pip list
```

##### 安装flask

```bash
pipenv install flask
```

### Flask最小原型与唯一URL原则

### 路由的另一种注册方法

```python
# @app.route('/hello')
def hello():
    # 基于类的视图（即插视图）
    return 'Hello, Kris'
app.add_url_rule('/hello', view_func=hello)
```

### app.run相关参数

```python
app = Flask(__name__)
#加载配置文件
app.config.from_object('app.config')

app.run(host='0.0.0.0', debug=app.config['DEBUG'], port=81, threaded=True)
```

### flask配置文件

**app/config.py**

```python
DEBUG = True
```

### 响应对象：Response

## 2  数据与flask路由

### 书籍搜索与查询

#### 搜索关键字

#### ISBN号关键字

#### requests vs urllib

#### 使用jsonify

## 3  蓝图、模型与CodeFirst

### 应用、蓝图与视图函数

### 使用蓝图注册视图函数

```python
# 放在__init__py文件中
web = Blueprint('web', __name__)

# 放在__init__py文件中
app.register_blueprint(web)

@web.route('/test')
def test1():
    return ''
```

### 单蓝图多模块拆分视图函数

```python
# __init__.py文件导入book.py模块，实现路由注册
from app.web import book
```

### request对象

### WTForms参数验证

### 拆分配置文件

```python
from flask import current_app
```

### 数据库表创建方式

1. Database First
2. Model First
3. Code First

### 模型映射到数据库

#### base.py

```python
from datetime import datetime

__author__ = '七月'

from flask_sqlalchemy import SQLAlchemy as _SQLAlchemy, BaseQuery
from sqlalchemy import Column, Integer, SmallInteger
from contextlib import contextmanager


class SQLAlchemy(_SQLAlchemy):
    @contextmanager
    def auto_commit(self):
        try:
            yield
            self.session.commit()
        except Exception as e:
            db.session.rollback()
            raise e


class Query(BaseQuery):
    def filter_by(self, **kwargs):
        if 'status' not in kwargs.keys():
            kwargs['status'] = 1
        return super(Query, self).filter_by(**kwargs)


db = SQLAlchemy(query_class=Query)


class Base(db.Model):
    __abstract__ = True
    create_time = Column('create_time', Integer)
    status = Column(SmallInteger, default=1)

    def __init__(self):
        self.create_time = int(datetime.now().timestamp())

    def set_attrs(self, attrs_dict):
        for key, value in attrs_dict.items():
            if hasattr(self, key) and key != 'id':
                setattr(self, key, value)

    @property
    def create_datetime(self):
        if self.create_time:
            return datetime.fromtimestamp(self.create_time)
        else:
            return None

    def delete(self):
        self.status = 0
```

#### book.py

```python
from sqlalchemy import Column, Integer, String

from app.models.base import db

__author__ = '七月'


class Book(db.Model):
    id = Column(Integer, primary_key=True, autoincrement=True)
    title = Column(String(50), nullable=False)
    author = Column(String(30), default='未名')
    binding = Column(String(20))
    publisher = Column(String(50))
    price = Column(String(20))
    pages = Column(Integer)
    pubdate = Column(String(20))
    isbn = Column(String(15), nullable=False, unique=True)
    summary = Column(String(1000))
    image = Column(String(50))

    # MVC M Model 只有数据 = 数据表
    # ORM 对象关系映射 Code First

    def sample(self):
        pass
```

### ORM与CodeFirst的区别

CodeFirst专注于业务模型的设计，而不是专注数据库的设计

## 4  flask核心机制

### flask中经典错误working outside application context

```python
from flask import Flask, current_app, request, Request

app = Flask(__name__)
# 应用上下文 对象 Flask
# 请求上下文 对象 Request
# Flask AppContext
# Request RequestContext
# 离线应用、单元测试
# ctx = app.app_context()
# ctx.push()
# a = current_app
# d = current_app.config['DEBUG']
# ctx.pop()

with app.app_context():
    a = current_app
    d = current_app.config['DEBUG']
```

### AppContext、RequestContext、Flask与Request之间的关系

### 详解flask上下文与出入栈

### flask上下文与with语句

```python
# 文件读写
# try:
#     f = open(r'D:\t.txt')
#     print(f.read())
# finally:
#     f.close()
#
# with open(r'') as f:
#     print(f.read())
# 1. 连接数据库
# 2. sql
# 3. 释放资源
#
# __enter__
# __exit__
# try
# except
# finally:

# 实现了上下文协议的对象使用with
# 上下文管理器
# __enter__ __exit__
# 上下文表达式必须要返回一个上下文管理器
# with
```

### 详解上下文管理器的\__exit\__方法

```python
class MyResource:
    def __enter__(self):
        print('connect to resource')
        return self

    def __exit__(self, exc_type, exc_value, tb):
        if tb:
            print('process exception')
        else:
            print('no exception')
        print('close resource connection')

    def query(self):
        print('query data')


try:
    with MyResource() as resource:
        # 1 / 0
        resource.query()
except Exception as ex:
    pass
```

### 阅读源码解决db.create_all的问题

```python
with app.app_context():
    db.create_all()
```

## 5  Flask中的多线程与线程隔离技术

### 什么是进程

```python
# 资源是稀缺的
# 计算机资源 竞争计算机的资源
# 进程
# 至少有1个进程
# 进程是竞争计算机资源的基本单位

# 单核CPU 永远只能够执行一个应用程序？
# 在不同的应用程序进程之间切换
# Pycharm、风暴英雄、QQ
# 进程调度
# 算法 挂起 切换到另外一个进程 操作系统原理
# 进程/线程 开销是非常大 上下文

# 4核
```

### 线程概念

```python
# 线程 线程是进程的一部分 1个线程 多个线程
# CPU 粒度太大了 更小的单元 CPU的资源
# 线程

# 进程 分配资源 内存资源
# 线程 利用CPU执行代码

# 代码 指令 CPU来执行 资源
# 访问资源
# 线程属于进程 访问进程资源
```

### 多线程

```python
def worker():
    print('i am thread')
    t = threading.current_thread()
    time.sleep(8)
    print(t.getName())

new_t = threading.Thread(target=worker, name='qiyue_thread')
new_t.start()
new_t = threading.Thread(target=worker, name='qiyue_thread')
new_t.start()
new_t = threading.Thread(target=worker, name='qiyue_thread')
new_t.start()

# 更加充分的利用CPU的性能优势
# 异步编程
# 单核CPU
# 4核 A核 B核 并行的执行程序
# python不能充分利用多核CPU优势
# python的多线程是鸡肋
# 6 6 = 6
# GIL 全局解释器锁 global interpreter lock
# 锁 线程安全

# 内存资源 一个进程 有多个线程 共享
# 线程不安全

# a = 3

# A a+=1
# print(a)

# B a+=1
# print(a)

# 锁
# 细粒度锁 程序员 主动加锁
# 粗粒度锁 解释器 GIL 多核cpu 1个线程执行 一定程度上保证线程安全
# a+=1
# bytecode
# python cpython jpython
# 多进程 进程通信技术
t = threading.current_thread()
print(t.getName())

# 主线程
```

### 全局解释器锁GIL

### 对于IO密集型程序，多线程是有意义的

```python
# python多线程到底是不是鸡肋
# GIL node.js 单进程 单线程
# 10 线程 非常严重的依赖CPU计算 CPU密集型程序
# IO密集型的程序 查询数据库、请求网络资源、读写文件

# IO密集型 等待
```

### 开启flask多线程所带来的问题

```python
# flask web框架
# 请求 线程
# 10请求 flask开启多少个线程来处理请求
# webserver
# Java、PHP nginx Apache tomcat IIS
```

```python
if __name__ == '__main__':
    app.run(host='0.0.0.0', debug=app.config['DEBUG'], port=81, threaded=True)
    #单进程、单线程
    # processes = 1
    # 10 个请求
```

### 线程隔离

### Flask中的线程隔离对象Local

```python
# 原理 字典 保存数据
# 操作数据
# werkzeug local Local 字典
# LocalStack Local 字典
# Local使用字典的方式实现的线程隔离
# 线程隔离的栈结构
# 封装 如果一次封装解决不了问题，那么就再来一次
# 编程也是一种艺术 含蓄

# {thread_id1:value1,thread_id2:value2...}
# L 线程隔离的对象
# t1 L.a t2 L.a
# Local

# 线程隔离的对象 被线程隔离的对象
# 操作线程隔离的对象（Local、LocalStack）来实现其他对象的线程隔离


from werkzeug.local import Local

class A:
    b = 1


my_obj = Local()
my_obj.b = 1


def worker():
    # 新线程
    my_obj.b = 2
    print('in new thread b is:' + str(my_obj.b))


new_t = threading.Thread(target=worker, name='qiyue_thread')
new_t.start()
time.sleep(1)
# 主线程
print('in main thread b is:' + str(my_obj.b))
```

### Flask中的线程隔离栈：LocalStack

### LocalStack的基本用法

```python
from werkzeug.local import LocalStack

# push、pop、top

s = LocalStack()
s.push(1)
print(s.top)
print(s.top)
print(s.pop())
print(s.top)

s.push(1)
s.push(2)
# 栈 后进先出
print(s.top)
print(s.top)
print(s.pop())
print(s.top)

# s[7] 栈、序列（列表）
# 数据结构 限制了某些能力
```

### LocalStack作为线程隔离对象的意义

```python
from werkzeug.local import LocalStack


my_stack = LocalStack()
my_stack.push(1)
print('in main thread after push , value is:' + str(my_stack.top))


def worker():
    # 新线程
    print('in new thread before push, value is:' + str(my_stack.top))
    my_stack.push(2)
    print('in new thread after push,value is:' + str(my_stack.top))


new_t = threading.Thread(target=worker, name='qiyue_thread')
new_t.start()
time.sleep(1)
# 主线程
print('finally, in main thread value is:' + str(my_stack.top))
```

### Flask中被线程隔离的对象

### 梳理串接Flask的一些名词

```python
# 以线程ID号作为key的字典->Local->LocalStack

# AppContext RequestContext -> LocalStack

# Flask -> AppContext   Request -> RequestContext

# current_app ->(LocalStack.top = AppContext top.app=Flask)

# request ->(LocalStack.top = RequestContext top.request= Request)
```
