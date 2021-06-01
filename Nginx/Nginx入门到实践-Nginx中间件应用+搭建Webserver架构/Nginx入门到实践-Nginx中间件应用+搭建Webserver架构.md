[TOC]



## 基础篇

### 什么是Nginx？

- http中间件服务

- http代理服务

### 常见的中间件服务

#### 常见的http服务

### Nginx的优势

#### 采用多路IO多路复用Epoll模型

##### 优势

###### 实现IO流非阻塞模式

```python
while true {
    for i in stream[] {
        if i has data {
            read until unavailable;
        }
    }
}
```

### Nginx的CPU亲和

#### 什么是CPU的亲和性？

是一种把CPU核心和Nginx工作进程绑定方式，把每一个worker进程绑定在cpu上执行，减少cpu的cache miss，获得更好的性能。

### Nginx的目录

| 路径                                                         | 类型           | 作用                                       |
| ------------------------------------------------------------ | -------------- | ------------------------------------------ |
| /etc/nginx<br/>/etc/nginx/nginx.conf<br/>/etc/nginx/conf.d<br/>/etc/nginx/conf.d/default.conf | 目录、配置文件 | Nginx的主配置文件                          |
| /etc/nginx/koi-utf<br/>/etc/nginx/koi-win<br/>/etc/nginx/win-utf | 配置文件       | 编码转换映射转化文件                       |
| /etc/nginx/mime.types                                        | 配置文件       | 设置http协议的Content-Type与扩展名对应关系 |
| /usr/lib64/nginx/modules<br/>/etc/nginx/modules              | 目录           | Nginx模块目录                              |
| /usr/sbin/nginx<br/>/usr/sbin/nginx-debug                    | 命令           | Nginx服务的启动管理终端命令                |
| /var/cache/nginx                                             | 目录           | Nginx的缓存目录                            |
| /var/log/nginx                                               | 目录           | Nginx的日志目录                            |

#### 查看nginx的安装目录

```bash
rpm -ql nginx
```

### Nginx的配置语法