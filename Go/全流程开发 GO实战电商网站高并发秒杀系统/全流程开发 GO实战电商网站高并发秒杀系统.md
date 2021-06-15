[TOC]



## 1 秒杀系统需求整理 & 系统设计

## 需求分析

## 系统架构设计

## 2  环境搭建之初识RabbitMQ

### RabbitMQ介绍

#### 定义和特征

1. RabbitMQ是面向消息的中间件，用于组件之间的解耦，只要体现在消息的发送者和消费者之间无强依赖关系；
2. 特点：高可用，扩展性，多语言客户端，管理界面等
3. 使用场景：流量削峰、异步处理、应用解耦等

### RabbitMQ常用命令及管理界面

#### 常用命令

##### 启动命令

```bash
systemctl start rabbitmq-server
```

##### 停止命令

```bash
rabbitmqctl stop
```

##### 插件管理命令

###### 查看插件列表

```bash
rabbitmq-plugins list
```

###### 启用管理界面

```bash
rabbitmq-plugins enable rabbitmq_management
```

### RabbitMQ核心概念

- VirtualHost
- Connection
- Exchange
- Channel
- Queue
- Binding

### RabbitMQ工作模式

#### Simple模式

##### 创建RabbitMQ实例

```go
//连接信息
const MQURL = "amqp://guest:guest@127.0.0.1:5672/"

//rabbitMQ结构体
type RabbitMQ struct {
   conn      *amqp.Connection
   channel   *amqp.Channel
   //队列名称
   QueueName string
   //交换机名称
   Exchange  string
   //bind Key 名称
   Key string
   //连接信息
   Mqurl     string
}

//创建结构体实例
func NewRabbitMQ(queueName string,exchange string ,key string) *RabbitMQ {
   return &RabbitMQ{QueueName:queueName,Exchange:exchange,Key:key,Mqurl:MQURL}
}

//断开channel 和 connection
func (r *RabbitMQ) Destory() {
   r.channel.Close()
   r.conn.Close()
}
//错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
   if err != nil {
      log.Fatalf("%s:%s", message, err)
      panic(fmt.Sprintf("%s:%s", message, err))
   }
}

//创建简单模式下RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
   //创建RabbitMQ实例
   rabbitmq := NewRabbitMQ(queueName,"","")
   var err error
    //获取connection
   rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
   rabbitmq.failOnErr(err, "failed to connect rabb"+
      "itmq!")
   //获取channel
   rabbitmq.channel, err = rabbitmq.conn.Channel()
   rabbitmq.failOnErr(err, "failed to open a channel")
   return rabbitmq
}
```

##### 生产者发送消息

```go
//直接模式队列生产
func (r *RabbitMQ) PublishSimple(message string) {
   //1.申请队列，如果队列不存在会自动创建，存在则跳过创建
   _, err := r.channel.QueueDeclare(
      r.QueueName,
      //是否持久化
      false,
      //是否自动删除
      false,
      //是否具有排他性
      false,
      //是否阻塞处理
      false,
      //额外的属性
      nil,
   )
   if err != nil {
      fmt.Println(err)
   }
   //调用channel 发送消息到队列中
   r.channel.Publish(
      r.Exchange,
      r.QueueName,
      //如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
      false,
      //如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
      false,
      amqp.Publishing{
         ContentType: "text/plain",
         Body:        []byte(message),
      })
}
```

##### 消费者接收消息

```go
//simple 模式下消费者
func (r *RabbitMQ) ConsumeSimple() {
   //1.申请队列，如果队列不存在会自动创建，存在则跳过创建
   q, err := r.channel.QueueDeclare(
      r.QueueName,
      //是否持久化
      false,
      //是否自动删除
      false,
      //是否具有排他性
      false,
      //是否阻塞处理
      false,
      //额外的属性
      nil,
   )
   if err != nil {
      fmt.Println(err)
   }

   //接收消息
   msgs, err :=r.channel.Consume(
      q.Name, // queue
      //用来区分多个消费者
      "",     // consumer
      //是否自动应答
      true,   // auto-ack
      //是否独有
      false,  // exclusive
      //设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
      false,  // no-local
      //列是否阻塞
      false,  // no-wait
      nil,    // args
   )
   if err != nil {
      fmt.Println(err)
   }

   forever := make(chan bool)
    //启用协程处理消息
   go func() {
      for d := range msgs {
         //消息逻辑处理，可以自行设计逻辑
         log.Printf("Received a message: %s", d.Body)
      }
   }()
   log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
   <-forever
}
```

#### Work模式

**一个消息只能被一个消费者获取，即一个消息只能被消费一次**

#### Publish模式

消息被路投递到多个队列，一个消息被多个消费者获取。

##### 创建RabbitMQ实例

```go
//订阅模式创建RabbitMQ实例
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
   //创建RabbitMQ实例
   rabbitmq := NewRabbitMQ("",exchangeName,"")
   var err error
   //获取connection
   rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
   rabbitmq.failOnErr(err,"failed to connect rabbitmq!")
   //获取channel
   rabbitmq.channel, err = rabbitmq.conn.Channel()
   rabbitmq.failOnErr(err, "failed to open a channel")
   return rabbitmq
}
```

##### 生产者发送消息

```go
//订阅模式生产
func (r *RabbitMQ) PublishPub(message string) {
   //1.尝试创建交换机
   err := r.channel.ExchangeDeclare(
      r.Exchange,
      "fanout",
      true,
      false,
      //true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
      false,
      false,
      nil,
   )

   r.failOnErr(err, "Failed to declare an excha"+
      "nge")

   //2.发送消息
   err = r.channel.Publish(
      r.Exchange,
      "",
      false,
      false,
      amqp.Publishing{
         ContentType: "text/plain",
         Body:        []byte(message),
      })
}
```

##### 消费者接收消息

```go
//订阅模式消费端代码
func (r *RabbitMQ) RecieveSub() {
   //1.试探性创建交换机
   err := r.channel.ExchangeDeclare(
      r.Exchange,
      //交换机类型
      "fanout",
      true,
      false,
      //YES表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
      false,
      false,
      nil,
   )
   r.failOnErr(err, "Failed to declare an exch"+
      "ange")
    //2.试探性创建队列，这里注意队列名称不要写
   q, err := r.channel.QueueDeclare(
      "", //随机生产队列名称
      false,
      false,
      true,
      false,
      nil,
   )
   r.failOnErr(err, "Failed to declare a queue")

   //绑定队列到 exchange 中
   err = r.channel.QueueBind(
      q.Name,
      //在pub/sub模式下，这里的key要为空
      "",
      r.Exchange,
      false,
      nil)

   //消费消息
   messges, err := r.channel.Consume(
      q.Name,
      "",
      true,
      false,
      false,
      false,
      nil,
   )

   forever := make(chan bool)

   go func() {
      for d := range messges {
         log.Printf("Received a message: %s", d.Body)
      }
   }()

   fmt.Println("退出请按 CTRL+C\n")
   <-forever
}
```

#### Routing模型

一个消息被多个消费者获取。并且消息的目标队列可被**生产者指定**。

##### 创建RabbitMQ实例

```go
//路由模式
//创建RabbitMQ实例
func NewRabbitMQRouting(exchangeName string,routingKey string) *RabbitMQ {
   //创建RabbitMQ实例
   rabbitmq := NewRabbitMQ("",exchangeName,routingKey)
   var err error
   //获取connection
   rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
   rabbitmq.failOnErr(err,"failed to connect rabbitmq!")
   //获取channel
   rabbitmq.channel, err = rabbitmq.conn.Channel()
   rabbitmq.failOnErr(err, "failed to open a channel")
   return rabbitmq
}
```

##### 生产者发送消息

```go
//路由模式发送消息
func (r *RabbitMQ) PublishRouting(message string )  {
   //1.尝试创建交换机
   err := r.channel.ExchangeDeclare(
      r.Exchange,
      //要改成direct
      "direct",
      true,
      false,
      false,
      false,
      nil,
   )

   r.failOnErr(err, "Failed to declare an excha"+
      "nge")

   //2.发送消息
   err = r.channel.Publish(
      r.Exchange,
      //要设置
      r.Key,
      false,
      false,
      amqp.Publishing{
         ContentType: "text/plain",
         Body:        []byte(message),
      })
}
```

##### 消费者接收消息

```go
//路由模式接受消息
func (r *RabbitMQ) RecieveRouting() {
   //1.试探性创建交换机
   err := r.channel.ExchangeDeclare(
      r.Exchange,
      //交换机类型
      "direct",
      true,
      false,
      false,
      false,
      nil,
   )
   r.failOnErr(err, "Failed to declare an exch"+
      "ange")
   //2.试探性创建队列，这里注意队列名称不要写
   q, err := r.channel.QueueDeclare(
      "", //随机生产队列名称
      false,
      false,
      true,
      false,
      nil,
   )
   r.failOnErr(err, "Failed to declare a queue")

   //绑定队列到 exchange 中
   err = r.channel.QueueBind(
      q.Name,
      //需要绑定key
      r.Key,
      r.Exchange,
      false,
      nil)

   //消费消息
   messges, err := r.channel.Consume(
      q.Name,
      "",
      true,
      false,
      false,
      false,
      nil,
   )

   forever := make(chan bool)

   go func() {
      for d := range messges {
         log.Printf("Received a message: %s", d.Body)
      }
   }()

   fmt.Println("退出请按 CTRL+C\n")
   <-forever
}
```

#### Topic模式

一个消息被多个消费者获取。消息的目标Queue可用BindKey以通配符（#：一个或多个词，*:一个词）的方式指定。

##### 创建RabbitMQ实例

```go
//话题模式
//创建RabbitMQ实例
func NewRabbitMQTopic(exchangeName string,routingKey string) *RabbitMQ {
   //创建RabbitMQ实例
   rabbitmq := NewRabbitMQ("",exchangeName,routingKey)
   var err error
   //获取connection
   rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
   rabbitmq.failOnErr(err,"failed to connect rabbitmq!")
   //获取channel
   rabbitmq.channel, err = rabbitmq.conn.Channel()
   rabbitmq.failOnErr(err, "failed to open a channel")
   return rabbitmq
}
```

##### 生产者发送消息

```go
//话题模式发送消息
func (r *RabbitMQ) PublishTopic(message string )  {
   //1.尝试创建交换机
   err := r.channel.ExchangeDeclare(
      r.Exchange,
      //要改成topic
      "topic",
      true,
      false,
      false,
      false,
      nil,
   )

   r.failOnErr(err, "Failed to declare an excha"+
      "nge")

   //2.发送消息
   err = r.channel.Publish(
      r.Exchange,
      //要设置
      r.Key,
      false,
      false,
      amqp.Publishing{
         ContentType: "text/plain",
         Body:        []byte(message),
      })
}
```

##### 消费者接收消息

```go
//话题模式接受消息
//要注意key,规则
//其中“*”用于匹配一个单词，“#”用于匹配多个单词（可以是零个）
//匹配 imooc.* 表示匹配 imooc.hello, 但是imooc.hello.one需要用imooc.#才能匹配到
func (r *RabbitMQ) RecieveTopic() {
   //1.试探性创建交换机
   err := r.channel.ExchangeDeclare(
      r.Exchange,
      //交换机类型
      "topic",
      true,
      false,
      false,
      false,
      nil,
   )
   r.failOnErr(err, "Failed to declare an exch"+
      "ange")
   //2.试探性创建队列，这里注意队列名称不要写
   q, err := r.channel.QueueDeclare(
      "", //随机生产队列名称
      false,
      false,
      true,
      false,
      nil,
   )
   r.failOnErr(err, "Failed to declare a queue")

   //绑定队列到 exchange 中
   err = r.channel.QueueBind(
      q.Name,
      //在pub/sub模式下，这里的key要为空
      r.Key,
      r.Exchange,
      false,
      nil)

   //消费消息
   messges, err := r.channel.Consume(
      q.Name,
      "",
      true,
      false,
      false,
      false,
      nil,
   )

   forever := make(chan bool)

   go func() {
      for d := range messges {
         log.Printf("Received a message: %s", d.Body)
      }
   }()

   fmt.Println("退出请按 CTRL+C\n")
   <-forever
}
```

## 2  环境搭建之Iris 框架入门

### MVC目录结构

#### main.go

```go
import (
   "github.com/kataras/iris/v12"
   "github.com/kataras/iris/v12/mvc"
   "imooc-iris/web/controllers"
)

func main()  {
   //1.创建实例
   app:=iris.New()
   app.Logger().SetLevel("debug")
   app.RegisterView(iris.HTML("./web/views", ".html"))
   //2.注册控制器
   mvc.New(app.Party("/hello")).Handle(new(controllers.MovieController))
   //3.启动服务
   app.Run(
      iris.Addr("localhost:8080"),
   )
}
```

## 3  后台管理功能开发之商品管理功能开发