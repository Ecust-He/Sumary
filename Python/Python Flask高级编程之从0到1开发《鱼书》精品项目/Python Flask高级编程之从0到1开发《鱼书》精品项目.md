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