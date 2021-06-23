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

## 3  环境搭建之Iris 框架入门

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

## 4  后台管理功能开发之商品管理功能开发

### 商品模型开发

```go
package datamodels

type Product struct {
   ID           int64  `json:"id" sql:"ID" imooc:"ID"`
   ProductName  string `json:"ProductName" sql:"productName" imooc:"ProductName"`
   ProductNum   int64  `json:"ProductNum" sql:"productNum" imooc:"ProductNum"`
   ProductImage string `json:"ProductImage" sql:"productImage" imooc:"ProductImage"`
   ProductUrl   string `json:"ProductUrl" sql:"productUrl" imooc:"ProductUrl"`
}
```

### 商品repository开发

#### 数据库连接

```go
import (
   "database/sql"
)

//创建mysql 连接
func NewMysqlConn() (db *sql.DB, err error) {
   db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/imooc?charset=utf8")
   return
}
```

#### 插入

```go
func (p *ProductManager) Insert(product *datamodels.Product) (productId int64,err error) {
   //1.判断连接是否存在
   if err=p.Conn();err != nil{
      return
   }
   //2.准备sql
   sql :="INSERT product SET productName=?,productNum=?,productImage=?,productUrl=?"
   stmt,errSql := p.mysqlConn.Prepare(sql)
   defer stmt.Close()
   if errSql !=nil {
      return 0,errSql
   }
   //3.传入参数
   result,errStmt:=stmt.Exec(product.ProductName,product.ProductNum,product.ProductImage,product.ProductUrl)
   if errStmt !=nil {
      return 0,errStmt
   }
   return result.LastInsertId()
}
```

#### 删除

```go
//商品的删除
func (p *ProductManager)Delete(productID int64) bool  {
   //1.判断连接是否存在
   if err:=p.Conn();err != nil{
      return false
   }
   sql := "delete from product where ID=?"
   stmt,err := p.mysqlConn.Prepare(sql)
   defer stmt.Close()
   if err!= nil {
      return false
   }
   _,err = stmt.Exec(strconv.FormatInt(productID,10))
   if err !=nil {
      return false
   }
   return true
}
```

#### 更新

```go
//商品的更新
func (p *ProductManager)Update(product *datamodels.Product) error {
   //1.判断连接是否存在
   if err:=p.Conn();err != nil{
      return err
   }

   sql := "Update product set productName=?,productNum=?,productImage=?,productUrl=? where ID="+strconv.FormatInt(product.ID,10)

   stmt,err := p.mysqlConn.Prepare(sql)
   defer stmt.Close()
   if err !=nil {
      return err
   }

   _,err = stmt.Exec(product.ProductName,product.ProductNum,product.ProductImage,product.ProductUrl)
   if err !=nil {
      return err
   }
   return nil
}
```

#### 根据商品ID查询商品

```go
//根据商品ID查询商品
func (p *ProductManager) SelectByKey(productID int64) (productResult *datamodels.Product,err error) {
   //1.判断连接是否存在
   if err=p.Conn();err != nil{
      return &datamodels.Product{},err
   }
   sql := "Select * from "+p.table+" where ID ="+strconv.FormatInt(productID,10)
   row,errRow :=p.mysqlConn.Query(sql)
   defer row.Close()
   if errRow !=nil {
      return &datamodels.Product{},errRow
   }
   result := common.GetResultRow(row)
   if len(result)==0{
      return &datamodels.Product{},nil
   }
   productResult = &datamodels.Product{}
   common.DataToStructByTagSql(result,productResult)
   return
}
```

##### 将rows转换为map

```go
//获取返回值， 获取一条
func GetResultRow(rows *sql.Rows) map[string]string{
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([][]byte, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}	
	record := make(map[string]string)
	for rows.Next() {
		//将行数据保存到record字典 
		rows.Scan(scanArgs...)
		for i, v := range values {
			if v != nil {
				//fmt.Println(reflect.TypeOf(col))
				record[columns[i]] = string(v)
			}
		}
	}
	return record
}
```

##### 根据结构体中sql标签映射数据到结构体中并且转换类型

```go
//根据结构体中sql标签映射数据到结构体中并且转换类型
func DataToStructByTagSql(data map[string]string, obj interface{}) {
	objValue := reflect.ValueOf(obj).Elem()
	for i := 0; i < objValue.NumField(); i++ {
		//获取sql对应的值
		value := data[objValue.Type().Field(i).Tag.Get("sql")]
		//获取对应字段的名称
		name := objValue.Type().Field(i).Name
		//获取对应字段类型
		structFieldType := objValue.Field(i).Type()
		//获取变量类型，也可以直接写"string类型"
		val := reflect.ValueOf(value)
		var err error
		if structFieldType != val.Type() {
			//类型转换
			val, err = TypeConversion(value, structFieldType.Name()) //类型转换
			if err != nil {

			}
		}
		//设置类型值
		objValue.FieldByName(name).Set(val)
	}
}

//类型转换
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	//else if .......增加其他一些类型的转换

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}
```

#### 查询所有商品

```go
//获取所有商品
func (p *ProductManager)SelectAll()(productArray []*datamodels.Product,errProduct error){
   //1.判断连接是否存在
   if err:=p.Conn();err!= nil{
      return nil,err
   }
   sql := "Select * from "+p.table
   rows,err := p.mysqlConn.Query(sql)
   defer  rows.Close()
   if err !=nil {
      return nil ,err
   }

   result:= common.GetResultRows(rows)
   if len(result)==0{
      return nil,nil
   }

   for _,v :=range result{
      product := &datamodels.Product{}
      common.DataToStructByTagSql(v,product)
      productArray=append(productArray, product)
   }
   return
}
```

### 商品管理功能service实现

```go
type IProductService interface {
   GetProductByID(int64) (*datamodels.Product,error)
   GetAllProduct()([]*datamodels.Product,error)
   DeleteProductByID(int64)bool
   InsertProduct(product *datamodels.Product)(int64,error)
   UpdateProduct(product *datamodels.Product)error
}

type ProductService struct {
   productRepository repositories.IProduct
}

//初始化函数
func NewProductService(repository repositories.IProduct)IProductService   {
   return &ProductService{repository}
}

func (p *ProductService)GetProductByID(productID int64) (*datamodels.Product,error)  {
   return p.productRepository.SelectByKey(productID)
}

func (p *ProductService) GetAllProduct() ([]*datamodels.Product, error) {
   return p.productRepository.SelectAll()
}

func (p *ProductService)DeleteProductByID(productID int64)bool {
   return p.productRepository.Delete(productID)
}

func (p *ProductService) InsertProduct(product *datamodels.Product) (int64, error) {
   return p.productRepository.Insert(product)
}

func (p *ProductService) UpdateProduct(product *datamodels.Product)error  {
   return p.productRepository.Update(product)
}
```

### 商品管理功能Contoller&View开发

```go
type ProductController struct {
   Ctx iris.Context
   ProductService services.IProductService
}

func (p *ProductController) GetAll() mvc.View {
   productArray ,_:=p.ProductService.GetAllProduct()
   return mvc.View{
      Name:"product/view.html",
      Data:iris.Map{
         "productArray":productArray,
      },
   }
}

//修改商品
func (p *ProductController) PostUpdate ()  {
   product :=&datamodels.Product{}
   p.Ctx.Request().ParseForm()
   dec := common.NewDecoder(&common.DecoderOptions{TagName:"imooc"})
   if err:= dec.Decode(p.Ctx.Request().Form,product);err!=nil {
      p.Ctx.Application().Logger().Debug(err)
   }
   err:=p.ProductService.UpdateProduct(product)
   if err !=nil {
      p.Ctx.Application().Logger().Debug(err)
   }
   p.Ctx.Redirect("/product/all")
}
```

### Golang 模板(template)的基本语法

## 5  后台管理功能开发之订单功能开发

### 订单model开发

```go
type Order struct {
   ID int64 `sql:"ID"`
   UserId int64 `sql:"userID"`
   ProductId int64 `sql:"productID"`
   OrderStatus int64 `sql:"orderStatus"`
}

const (
   OrderWait = iota
   OrderSuccess  //1
   OrderFailed //2
)
```

### 订单repository开发

```go
type IOrderRepository interface {
   Conn() error
   Insert(*datamodels.Order) (int64,error)
   Delete(int64) bool
   Update(*datamodels.Order) error
   SelectByKey (int64) (*datamodels.Order,error)
   SelectAll ()([]*datamodels.Order,error)
   SelectAllWithInfo()(map[int]map[string]string,error)
}

func NewOrderMangerRepository(table string,sql *sql.DB) IOrderRepository  {
   return &OrderMangerRepository{table:table,mysqlConn:sql}
}

type OrderMangerRepository struct {
   table string
   mysqlConn *sql.DB
}

func (o *OrderMangerRepository)Conn() error   {
   if o.mysqlConn ==nil {
      mysql,err := common.NewMysqlConn()
      if err !=nil {
         return err
      }
      o.mysqlConn=mysql
   }
   if o.table == "" {
      o.table = "`order`"
   }

   return nil
}

func (o *OrderMangerRepository) Insert(order *datamodels.Order) (productID int64,err error) {
   if err =o.Conn(); err !=nil {
      return
   }

   sql :="INSERT "+o.table+" set userID=?,productID=?,orderStatus=?"
   stmt ,errStmt := o.mysqlConn.Prepare(sql)
   defer stmt.Close()
   if errStmt !=nil {
      return productID,errStmt
   }
   result,errResult :=stmt.Exec(order.UserId,order.ProductId,order.OrderStatus)
   if errResult!=nil  {
      return productID,errResult
   }
   return result.LastInsertId()
}

func (o *OrderMangerRepository)Delete(orderID int64) (isOk bool)   {
   if err :=o.Conn();err !=nil {
      return
   }
   sql :="delete from "+o.table+" where ID =?"
   stmt,errStmt:=o.mysqlConn.Prepare(sql)
   defer stmt.Close()
   if errStmt !=nil {
      return
   }
   _,err := stmt.Exec(orderID)
   if err !=nil {
      return
   }
   return true
}

func (o *OrderMangerRepository)Update(order *datamodels.Order) (err error) {
   if errConn := o.Conn(); errConn != nil {
      return errConn
   }

   sql := "Update " + o.table + " set userID=?,productID=?,orderStatus=? Where ID=" + strconv.FormatInt(order.ID, 10)
   stmt, errStmt := o.mysqlConn.Prepare(sql)
   defer stmt.Close()
   if errStmt != nil {
      return errStmt
   }
   _, err = stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
   return
}

func (o *OrderMangerRepository)SelectByKey (orderID int64) (order *datamodels.Order,err error)  {
   if errConn := o.Conn(); errConn != nil {
      return &datamodels.Order{},errConn
   }

   sql :="Select * From "+o.table+" where ID="+ strconv.FormatInt(orderID, 10)
   row,errRow :=o.mysqlConn.Query(sql)
   defer row.Close()
   if errRow!=nil{
      return &datamodels.Order{},errRow
   }

   result := common.GetResultRow(row)
   if len(result) == 0 {
      return &datamodels.Order{},err
   }

   order = &datamodels.Order{}
   common.DataToStructByTagSql(result,order)
   return
}

func (o *OrderMangerRepository)SelectAll ()(orderArray []*datamodels.Order,err error)  {
   if errConn := o.Conn(); errConn != nil {
      return nil,errConn
   }
   sql:="Select * from "+o.table
   rows ,errRows:=o.mysqlConn.Query(sql)
   defer rows.Close()
   if errRows!=nil {
      return nil,errRows
   }
   result := common.GetResultRows(rows)
   if len(result) == 0 {
      return nil ,err
   }

   for _,v :=range result{
      order :=&datamodels.Order{}
      common.DataToStructByTagSql(v,order)
      orderArray=append(orderArray,order)
   }
   return 
}

func (o *OrderMangerRepository) SelectAllWithInfo() (OrderMap map[int]map[string]string, err error) {
   if errConn := o.Conn(); errConn != nil {
      return nil, errConn
   }
   sql := "Select o.ID,p.productName,o.orderStatus From imooc.order as o left join product as p on o.productID=p.ID"
   rows, errRows := o.mysqlConn.Query(sql)
   defer rows.Close()
   if errRows != nil {
      return nil, errRows
   }
   return common.GetResultRows(rows), err
}
```

### 订单管理service实现

```go
type IOrderService interface {
   GetOrderByID(int64) (*datamodels.Order,error)
   DeleteOrderByID(int64) bool
   UpdateOrder(*datamodels.Order) error
   InsertOrder(*datamodels.Order) (int64 ,error)
   GetAllOrder()([]*datamodels.Order,error)
   GetAllOrderInfo()(map[int]map[string]string,error)
}

func NewOrderService(repository repositories.IOrderRepository ) IOrderService  {
   return &OrderService{OrderRepository:repository}
}

type OrderService struct {
   OrderRepository repositories.IOrderRepository
}

func (o *OrderService)GetOrderByID(orderID int64) (order *datamodels.Order,err error)  {
   return o.OrderRepository.SelectByKey(orderID)
}

func (o *OrderService)DeleteOrderByID(orderID int64) (isOk bool)  {
   isOk = o.OrderRepository.Delete(orderID)
   return
}

func (o *OrderService)UpdateOrder(order *datamodels.Order) error{
   return o.OrderRepository.Update(order)
}

func (o *OrderService)InsertOrder(order *datamodels.Order) (orderID int64,err error)  {
   return o.OrderRepository.Insert(order)
}

func (o *OrderService) GetAllOrder()([]*datamodels.Order,error) {
   return o.OrderRepository.SelectAll()
}

func (o *OrderService) GetAllOrderInfo()(map[int]map[string]string,error) {
   return o.OrderRepository.SelectAllWithInfo()
}
```

### 订单管理Controller&View实现

```go
type OrderController struct {
   Ctx iris.Context
   OrderService services.IOrderService
}

func (o *OrderController) Get() mvc.View {
   orderArray,err:=o.OrderService.GetAllOrderInfo()
   if err !=nil {
      o.Ctx.Application().Logger().Debug("查询订单信息失败")
   }

   return mvc.View{
      Name:"order/view.html",
      Data:iris.Map{
         "order":orderArray,
      },
   }
   
}
```

### Go语言中的Tag语法

Tag在运行时可以通过reflection包来读取

```go
package main
import (
    "fmt"
    "reflect"
)
type T struct {
    f1     string "f one"
    f2     string
    f3     string `f three`
    f4, f5 int64  `f four and five`
}
func main() {
    t := reflect.TypeOf(T{})
    f1, _ := t.FieldByName("f1")
    fmt.Println(f1.Tag) // f one
    f4, _ := t.FieldByName("f4")
    fmt.Println(f4.Tag) // f four and five
    f5, _ := t.FieldByName("f5")
    fmt.Println(f5.Tag) // f four and five
}
```

#### 格式

Tags可以由键值对来组成，通过空格符来分割键值 —key1:“value1” key2:“value2” key3:“value3”。如果Tags格式没问题的话，我们可以通过Lookup或者Get来获取键值对的值。
Lookup回传两个值 —对应的值和是否找到

```go
type T struct {
    f string `one:"1" two:"2"blank:""`
}
func main() {
    t := reflect.TypeOf(T{})
    f, _ := t.FieldByName("f")
    fmt.Println(f.Tag) // one:"1" two:"2"blank:""
    v, ok := f.Tag.Lookup("one")
    fmt.Printf("%s, %t\n", v, ok) // 1, true
    v, ok = f.Tag.Lookup("blank")
    fmt.Printf("%s, %t\n", v, ok) // , true
    v, ok = f.Tag.Lookup("five")
    fmt.Printf("%s, %t\n", v, ok) // , false
}
```

#### 转化

```go
type T1 struct {
     f int `json:"foo"`
 }
 type T2 struct {
     f int `json:"bar"`
 }
 t1 := T1{10}
 var t2 T2
 t2 = T2(t1)
 fmt.Println(t2) // {10}
```

## 6  秒杀前台功能开发之用户注册登录

### model开发

```go
type User struct {
   ID           int64  `json:"id" form:"ID" sql:"ID"`
   NickName     string `json:"nickName" form:"nickName" sql:"nickName"`
   UserName     string `json:"userName" form:"userName" sql:"userName"`
   HashPassword string `json:"-" form:"passWord" sql:"passWord"`
}
```

### repository开发

```go
type IUserRepository interface {
   Conn() error
   Select(userName string) (user *datamodels.User, err error)
   Insert(user *datamodels.User) (userId int64, err error)
}

func NewUserRepository(table string, db *sql.DB) IUserRepository {
   return &UserManagerRepository{table, db}
}

type UserManagerRepository struct {
   table     string
   mysqlConn *sql.DB
}

func (u *UserManagerRepository) Conn() (err error) {
   if u.mysqlConn == nil {
      mysql, errMysql := common.NewMysqlConn()
      if errMysql != nil {
         return errMysql
      }
      u.mysqlConn = mysql
   }
   if u.table == "" {
      u.table = "user"
   }
   return
}

func (u *UserManagerRepository) Select(userName string) (user *datamodels.User, err error) {
   if userName == "" {
      return &datamodels.User{}, errors.New("条件不能为空！")
   }
   if err = u.Conn(); err != nil {
      return &datamodels.User{}, err
   }

   sql := "Select * from " + u.table + " where userName=?"
   rows, errRows := u.mysqlConn.Query(sql, userName)
   defer rows.Close()
   if errRows != nil {
      return &datamodels.User{}, errRows
   }

   result := common.GetResultRow(rows)
   if len(result) == 0 {
      return &datamodels.User{}, errors.New("用户不存在！")
   }

   user = &datamodels.User{}
   common.DataToStructByTagSql(result, user)
   return
}

func (u *UserManagerRepository) Insert(user *datamodels.User) (userId int64, err error) {
   if err = u.Conn(); err != nil {
      return
   }

   sql := "INSERT " + u.table + " SET nickName=?,userName=?,passWord=?"
   stmt, errStmt := u.mysqlConn.Prepare(sql)
   defer stmt.Close()
   if errStmt != nil {
      return userId, errStmt
   }
   result, errResult := stmt.Exec(user.NickName, user.UserName, user.HashPassword)
   if errResult != nil {
      return userId, errResult
   }
   return result.LastInsertId()
}

func (u *UserManagerRepository) SelectByID(userId int64) (user *datamodels.User, err error) {
   if err = u.Conn(); err != nil {
      return &datamodels.User{}, err
   }
   sql := "select * from " + u.table + " where ID=" + strconv.FormatInt(userId, 10)
   row, errRow := u.mysqlConn.Query(sql)
   defer row.Close()
   if errRow != nil {
      return &datamodels.User{}, errRow
   }
   result := common.GetResultRow(row)
   if len(result) == 0 {
      return &datamodels.User{}, errors.New("用户不存在！")
   }
   user = &datamodels.User{}
   common.DataToStructByTagSql(result, user)
   return
}
```

### service开发

```go
type IUserService interface {
   IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool)
   AddUser(user *datamodels.User) (userId int64, err error)
}

func NewService(repository repositories.IUserRepository) IUserService {
   return &UserService{repository}
}

type UserService struct {
   UserRepository repositories.IUserRepository
}

func (u *UserService) IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool) {

   user, err := u.UserRepository.Select(userName)

   if err != nil {
      return
   }
   isOk, _ = ValidatePassword(pwd, user.HashPassword)

   if !isOk {
      return &datamodels.User{}, false
   }

   return
}

func (u *UserService) AddUser(user *datamodels.User) (userId int64, err error) {
   pwdByte, errPwd := GeneratePassword(user.HashPassword)
   if errPwd != nil {
      return userId, errPwd
   }
   user.HashPassword = string(pwdByte)
   return u.UserRepository.Insert(user)
}

func GeneratePassword(userPassword string) ([]byte, error) {
   return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

func ValidatePassword(userPassword string, hashed string) (isOK bool, err error) {
   if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
      return false, errors.New("密码比对错误！")
   }
   return true, nil
}
```

### Controller&View开发

