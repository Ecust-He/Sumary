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



