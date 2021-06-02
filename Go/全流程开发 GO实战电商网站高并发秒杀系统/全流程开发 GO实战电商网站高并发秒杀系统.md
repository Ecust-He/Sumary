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

##### rabbitMQ实例的创建

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