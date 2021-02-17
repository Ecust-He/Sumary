[TOC]

## 七大设计原则

### 开闭原则

#### 编码演示

##### 课程接口

```java
public interface ICourse {
    Integer getId();
    String getName();
    Double getPrice();
}
```

##### Java课程

```java
public class JavaCourse implements ICourse{
    private Integer Id;
    private String name;
    private Double price;

    public JavaCourse(Integer id, String name, Double price) {
        this.Id = id;
        this.name = name;
        this.price = price;
    }

    public Double getPrice() {
        return this.price;
    }

}
```

##### Java打折课程

```java
public class JavaDiscountCourse extends JavaCourse {

    public JavaDiscountCourse(Integer id, String name, Double price) {
        super(id, name, price);
    }

    public Double getDiscountPrice(){
        return super.getPrice()*0.8;
    }
}
```

### 依赖倒置原则

#### 编码演示

##### 课程接口

```java
public interface ICourse {
    void studyCourse();
}
```

##### Java课程

```java
public class JavaCourse implements ICourse {

    @Override
    public void studyCourse() {
        System.out.println("Geely在学习Java课程");
    }
}
```

##### 学生类

```java
public class Geely {

    public void setiCourse(ICourse iCourse) {
        this.iCourse = iCourse;
    }

    private ICourse iCourse;

    public void studyImoocCourse(){
        iCourse.studyCourse();
    }
}
```

##### 测试类

```java
public class Test {

    //v1
//    public static void main(String[] args) {
//        Geely geely = new Geely();
//        geely.studyJavaCourse();
//        geely.studyFECourse();
//    }

    //v2
//    public static void main(String[] args) {
//        Geely geely = new Geely();
//        geely.studyImoocCourse(new JavaCourse());
//        geely.studyImoocCourse(new FECourse());
//        geely.studyImoocCourse(new PythonCourse());
//    }

    //v3
//    public static void main(String[] args) {
//        Geely geely = new Geely(new JavaCourse());
//        geely.studyImoocCourse();
//    }
    public static void main(String[] args) {
        Geely geely = new Geely();
        geely.setiCourse(new JavaCourse());
        geely.studyImoocCourse();

        geely.setiCourse(new FECourse());
        geely.studyImoocCourse();
    }
}
```

### 单一职责原则

#### 编码演示1

##### 鸟类

```java
public class Bird {
    public void mainMoveMode(String birdName){
        if("鸵鸟".equals(birdName)){
            System.out.println(birdName+"用脚走");
        }else{
            System.out.println(birdName+"用翅膀飞");
        }
    }
}
```

##### 会飞的鸟

```java
public class FlyBird {
    public void mainMoveMode(String birdName){
        System.out.println(birdName+"用翅膀飞");
    }
}
```

##### 会行走的鸟

```java
public class WalkBird {
    public void mainMoveMode(String birdName){
        System.out.println(birdName+"用脚走");
    }
}
```

##### 测试类

```java
public class Test {
    public static void main(String[] args) {
//        Bird bird = new Bird();
//        bird.mainMoveMode("大雁");
//        bird.mainMoveMode("鸵鸟");

        FlyBird flyBird = new FlyBird();
        flyBird.mainMoveMode("大雁");

        WalkBird walkBird = new WalkBird();
        walkBird.mainMoveMode("鸵鸟");

    }
}
```

#### 编码演示2

##### 课程接口类

```java
public interface ICourse {
    String getCourseName();
    byte[] getCourseVideo();

    void studyCourse();
    void refundCourse();

}
```

##### 课程内容接口类

```java
public interface ICourseContent {
    String getCourseName();
    byte[] getCourseVideo();
}
```

##### 课程管理接口类

```java
public interface ICourseManager {
    void studyCourse();
    void refundCourse();
}
```

##### 课程实现类

```java
public class CourseImpl implements ICourseManager,ICourseContent {
    @Override
    public void studyCourse() {

    }

    @Override
    public void refundCourse() {

    }

    @Override
    public String getCourseName() {
        return null;
    }

    @Override
    public byte[] getCourseVideo() {
        return new byte[0];
    }
}
```

### 接口隔离原则

#### 编码演示

```java
public interface IEatAnimalAction {
    void eat();
}
```

```java
public interface ISwimAnimalAction {
    void swim();
}
```

```java
public class Dog implements ISwimAnimalAction,IEatAnimalAction {

    @Override
    public void eat() {

    }

    @Override
    public void swim() {

    }
}
```

### 迪米特法则

#### 编码演示

##### 课程类

```java
public class Course {
}
```

##### 领导类

```java
public class TeamLeader {
    public void checkNumberOfCourses(){
        List<Course> courseList = new ArrayList<Course>();
        for(int i = 0 ;i < 20;i++){
            courseList.add(new Course());
        }
        System.out.println("在线课程的数量是："+courseList.size());
    }
}
```

##### 老板类

```java
public class Boss {

    public void commandCheckNumber(TeamLeader teamLeader){
        teamLeader.checkNumberOfCourses();
    }
}
```

##### 测试类

```java
public class Test {
    public static void main(String[] args) {
        Boss boss = new Boss();
        TeamLeader teamLeader = new TeamLeader();
        boss.commandCheckNumber(teamLeader);
    }
}
```

### 里式替换原则

#### 编码演示1

##### 父类

```java
public class Base {
    public void method(HashMap map){
        System.out.println("父类被执行");
    }
}
```

##### 子类

```java
public class Child extends Base {
//    @Override
//    public void method(HashMap map) {
//        System.out.println("子类HashMap入参方法被执行");
//    }

    public void method(Map map) {
        System.out.println("子类HashMap入参方法被执行");
    }
}
```

##### 测试类

```java
public class Test {
    public static void main(String[] args) {
        Base child = new Child();
        HashMap hashMap = new HashMap();
        child.method(hashMap);
    }
}
```

#### 编码演示2

##### 父类

```java
public abstract class Base {
    public abstract Map method();

}
```

##### 子类

```java
public class Child extends Base {
    @Override
    public HashMap method() {
        HashMap hashMap = new HashMap();
        System.out.println("子类method被执行");
        hashMap.put("message","子类method被执行");
        return hashMap;
    }
}
```

##### 测试类

```java
public class Test {
    public static void main(String[] args) {
        Child child = new Child();
        System.out.println(child.method());
    }
}
```

### 组合聚合原则

#### 编码演示

##### 数据库连接接口

```java
public abstract class DBConnection {
//    public String getConnection(){
//        return "MySQL数据库连接";
//    }
    public abstract String getConnection();
}
```

##### mysql数据库连接

```java
public class MySQLConnection extends DBConnection {
    @Override
    public String getConnection() {
        return "MySQL数据库连接";
    }
}
```

##### postgresql数据库连接

```java
public class PostgreSQLConnection extends DBConnection {
    @Override
    public String getConnection() {
        return "PostgreSQL数据库连接";
    }
}
```

##### 商品数据库操作

```java
public class ProductDao{
    private DBConnection dbConnection;

    public void setDbConnection(DBConnection dbConnection) {
        this.dbConnection = dbConnection;
    }

    public void addProduct(){
        String conn = dbConnection.getConnection();
        System.out.println("使用"+conn+"增加产品");
    }
}
```

##### 测试类

```java
public class Test {
    public static void main(String[] args) {
        ProductDao productDao = new ProductDao();
        productDao.setDbConnection(new PostgreSQLConnection());
        productDao.addProduct();
    }
}
```

## 设计模式

### 创建型

#### 简单工厂

##### 编码演示

###### 视频抽象类

```java
public abstract class Video {
    public abstract void produce();
}
```

###### Java视频类

```java
public class JavaVideo extends Video {
    @Override
    public void produce() {
        System.out.println("录制Java课程视频");
    }
}
```

###### 视频工厂类

```java
public class VideoFactory {
    public Video getVideo(Class c){
        Video video = null;
        try {
            video = (Video) Class.forName(c.getName()).newInstance();
        } catch (InstantiationException e) {
            e.printStackTrace();
        } catch (IllegalAccessException e) {
            e.printStackTrace();
        } catch (ClassNotFoundException e) {
            e.printStackTrace();
        }
        return video;
    }

    public Video getVideo(String type){
        if("java".equalsIgnoreCase(type)){
            return new JavaVideo();
        }else if("python".equalsIgnoreCase(type)){
            return new PythonVideo();
        }
        return null;
    }
}
```

###### 测试类

```java
public class Test {
    public static void main(String[] args) {
//        VideoFactory videoFactory = new VideoFactory();
//        Video video = videoFactory.getVideo("java");
//        if(video == null){
//            return;
//        }
//        video.produce();

        VideoFactory videoFactory = new VideoFactory();
        Video video = videoFactory.getVideo(JavaVideo.class);
        if(video == null){
            return;
        }
        video.produce();
	}
}
```

#### 工厂方法

##### 编码演示

###### 视频抽象类

```java
public abstract class Video {
    public abstract void produce();
}
```

###### Java视频类

```java
public class JavaVideo extends Video {
    @Override
    public void produce() {
        System.out.println("录制Java课程视频");
    }
}
```

###### 视频工厂

```java
public abstract class VideoFactory {

    public abstract Video getVideo();

//    public Video getVideo(Class c){
//        Video video = null;
//        try {
//            video = (Video) Class.forName(c.getName()).newInstance();
//        } catch (InstantiationException e) {
//            e.printStackTrace();
//        } catch (IllegalAccessException e) {
//            e.printStackTrace();
//        } catch (ClassNotFoundException e) {
//            e.printStackTrace();
//        }
//        return video;
//    }
//
//
//    public Video getVideo(String type){
//        if("java".equalsIgnoreCase(type)){
//            return new JavaVideo();
//        }else if("python".equalsIgnoreCase(type)){
//            return new PythonVideo();
//        }
//        return null;
//    }

}
```

###### Java视频工厂

```java
public class JavaVideoFactory extends VideoFactory {
    @Override
    public Video getVideo() {
        return new JavaVideo();
    }
}
```

###### 测试类

```java
public class Test {
    public static void main(String[] args) {
        VideoFactory videoFactory = new PythonVideoFactory();
        VideoFactory videoFactory2 = new JavaVideoFactory();
        VideoFactory videoFactory3 = new FEVideoFactory();
        Video video = videoFactory.getVideo();
        video.produce();
    }
}
```

#### 抽象工厂

##### 编码演示

###### 文章抽象类

```java
public abstract class Article {
    public abstract void produce();
}
```

###### Java文章类

```java
public class JavaArticle extends Article {
    @Override
    public void produce() {
        System.out.println("编写Java课程手记");
    }
}
```

###### 视频抽象类

```java
public abstract class Video {
    public abstract void produce();
}
```

###### Java视频类

```java
public class JavaVideo extends Video {
    @Override
    public void produce() {
        System.out.println("录制Java课程视频");
    }
}
```

###### 课程工厂抽象类

```java
public interface CourseFactory {
    Video getVideo();
    Article getArticle();

}
```

###### Java课程工厂

```java
public class JavaCourseFactory implements CourseFactory {
    @Override
    public Video getVideo() {
        return new JavaVideo();
    }

    @Override
    public Article getArticle() {
        return new JavaArticle();
    }
}
```

###### 测试类

```java
public class Test {
    public static void main(String[] args) {
        CourseFactory courseFactory = new JavaCourseFactory();
        Video video = courseFactory.getVideo();
        Article article = courseFactory.getArticle();
        video.produce();
        article.produce();
    }
}
```

#### 建造者模式

##### 编码演示

```java
public class Course {

    private String courseName;
    private String coursePPT;
    private String courseVideo;
    private String courseArticle;

    //question & answer
    private String courseQA;

    public Course(CourseBuilder courseBuilder) {
        this.courseName = courseBuilder.courseName;
        this.coursePPT = courseBuilder.coursePPT;
        this.courseVideo = courseBuilder.courseVideo;
        this.courseArticle = courseBuilder.courseArticle;
        this.courseQA = courseBuilder.courseQA;
    }


    @Override
    public String toString() {
        return "Course{" +
                "courseName='" + courseName + '\'' +
                ", coursePPT='" + coursePPT + '\'' +
                ", courseVideo='" + courseVideo + '\'' +
                ", courseArticle='" + courseArticle + '\'' +
                ", courseQA='" + courseQA + '\'' +
                '}';
    }

    public static class CourseBuilder{
        private String courseName;
        private String coursePPT;
        private String courseVideo;
        private String courseArticle;

        //question & answer
        private String courseQA;

        public CourseBuilder buildCourseName(String courseName){
            this.courseName = courseName;
            return this;
        }


        public CourseBuilder buildCoursePPT(String coursePPT) {
            this.coursePPT = coursePPT;
            return this;
        }

        public CourseBuilder buildCourseVideo(String courseVideo) {
            this.courseVideo = courseVideo;
            return this;
        }

        public CourseBuilder buildCourseArticle(String courseArticle) {
            this.courseArticle = courseArticle;
            return this;
        }

        public CourseBuilder buildCourseQA(String courseQA) {
            this.courseQA = courseQA;
            return this;
        }

        public Course build(){
            return new Course(this);
        }
    }
}
```

###### 测试类

```java
public class Test {
    public static void main(String[] args) {
        Course course = new Course.CourseBuilder().buildCourseName("Java设计模式精讲").buildCoursePPT("Java设计模式精讲PPT").buildCourseVideo("Java设计模式精讲视频").build();
        System.out.println(course);

        Set<String> set = ImmutableSet.<String>builder().add("a").add("b").build();

        System.out.println(set);
    }
}
```

#### 单例模式

##### 编码演示

###### 1、懒汉模式

```java
public class LazySingleton {
    private static LazySingleton lazySingleton = null;
    private LazySingleton(){
        if(lazySingleton != null){
            throw new RuntimeException("单例构造器禁止反射调用");
        }
    }
    public synchronized static LazySingleton getInstance(){
        if(lazySingleton == null){
            lazySingleton = new LazySingleton();
        }
        return lazySingleton;
    }
}
```

###### 2、DoubleCheck双重检查

```java
public class LazyDoubleCheckSingleton {
    private volatile static LazyDoubleCheckSingleton lazyDoubleCheckSingleton = null;
    private LazyDoubleCheckSingleton(){

    }
    public static LazyDoubleCheckSingleton getInstance(){
        if(lazyDoubleCheckSingleton == null){
            synchronized (LazyDoubleCheckSingleton.class){
                if(lazyDoubleCheckSingleton == null){
                    lazyDoubleCheckSingleton = new LazyDoubleCheckSingleton();
                    //1.分配内存给这个对象
//                  //3.设置lazyDoubleCheckSingleton 指向刚分配的内存地址
                    //2.初始化对象
//                    intra-thread semantics
//                    ---------------//3.设置lazyDoubleCheckSingleton 指向刚分配的内存地址
                }
            }
        }
        return lazyDoubleCheckSingleton;
    }
}
```

###### 3、静态内部类

```java
public class StaticInnerClassSingleton {
    private static class InnerClass{
        private static StaticInnerClassSingleton staticInnerClassSingleton = new StaticInnerClassSingleton();
    }
    public static StaticInnerClassSingleton getInstance(){
        return InnerClass.staticInnerClassSingleton;
    }
    private StaticInnerClassSingleton(){
        if(InnerClass.staticInnerClassSingleton != null){
            throw new RuntimeException("单例构造器禁止反射调用");
        }
    }
}
```

###### 4、饿汉模式

```java
public class HungrySingleton implements Serializable,Cloneable{

    private final static HungrySingleton hungrySingleton;

    static{
        hungrySingleton = new HungrySingleton();
    }
    private HungrySingleton(){
        if(hungrySingleton != null){
            throw new RuntimeException("单例构造器禁止反射调用");
        }
    }
    public static HungrySingleton getInstance(){
        return hungrySingleton;
    }

    private Object readResolve(){
        return hungrySingleton;
    }

    @Override
    protected Object clone() throws CloneNotSupportedException {
        return getInstance();
    }
}
```

###### 5、枚举单例

```java
public enum EnumInstance {
    INSTANCE{
        protected  void printTest(){
            System.out.println("Geely Print Test");
        }
    };
    protected abstract void printTest();
    private Object data;

    public Object getData() {
        return data;
    }

    public void setData(Object data) {
        this.data = data;
    }
    public static EnumInstance getInstance(){
        return INSTANCE;
    }

}
```

#### 原型模式

##### 编码演示

###### 1、邮件类

```java
public class Mail implements Cloneable{
    private String name;
    private String emailAddress;
    private String content;
    public Mail(){
        System.out.println("Mail Class Constructor");
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getEmailAddress() {
        return emailAddress;
    }

    public void setEmailAddress(String emailAddress) {
        this.emailAddress = emailAddress;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

    @Override
    public String toString() {
        return "Mail{" +
                "name='" + name + '\'' +
                ", emailAddress='" + emailAddress + '\'' +
                ", content='" + content + '\'' +
                '}'+super.toString();
    }

    @Override
    protected Object clone() throws CloneNotSupportedException {
        System.out.println("clone mail object");
        return super.clone();
    }
}
```

###### 2、邮件工具类

```java
public class MailUtil {
    public static void sendMail(Mail mail){
        String outputContent = "向{0}同学,邮件地址:{1},邮件内容:{2}发送邮件成功";
        System.out.println(MessageFormat.format(outputContent,mail.getName(),mail.getEmailAddress(),mail.getContent()));
    }
    public static void saveOriginMailRecord(Mail mail){
        System.out.println("存储originMail记录,originMail:"+mail.getContent());
    }
}
```

###### 3、测试类

```java
public class Test {
    public static void main(String[] args) throws CloneNotSupportedException {
        Mail mail = new Mail();
        mail.setContent("初始化模板");
        System.out.println("初始化mail:"+mail);
        for(int i = 0;i < 10;i++){
            Mail mailTemp = (Mail) mail.clone();
            mailTemp.setName("姓名"+i);
            mailTemp.setEmailAddress("姓名"+i+"@imooc.com");
            mailTemp.setContent("恭喜您，此次慕课网活动中奖了");
            MailUtil.sendMail(mailTemp);
            System.out.println("克隆的mailTemp:"+mailTemp);
        }
        MailUtil.saveOriginMailRecord(mail);
    }
}
```

### 结构型

#### 外观模式

##### 编码演示

###### 积分礼物类

```java
public class PointsGift {
    private String name;

    public PointsGift(String name) {
        this.name = name;
    }

    public String getName() {
        return name;
    }
}
```

###### 资格校验服务类

```java
public class QualifyService {
    public boolean isAvailable(PointsGift pointsGift){
        System.out.println("校验"+pointsGift.getName()+" 积分资格通过,库存通过");
        return true;
    }
}
```

###### 积分支付服务类

```java
public class PointsPaymentService {
    public boolean pay(PointsGift pointsGift){
        //扣减积分
        System.out.println("支付"+pointsGift.getName()+" 积分成功");
        return true;
    }

}
```

###### 物流系统服务类

```java
public class ShippingService {
    public String shipGift(PointsGift pointsGift){
        //物流系统的对接逻辑
        System.out.println(pointsGift.getName()+"进入物流系统");
        String shippingOrderNo = "666";
        return shippingOrderNo;
    }
}
```

###### 积分兑换服务类

```java
public class GiftExchangeService {
    private QualifyService qualifyService = new QualifyService();
    private PointsPaymentService pointsPaymentService = new PointsPaymentService();
    private ShippingService shippingService = new ShippingService();

    public void giftExchange(PointsGift pointsGift){
        if(qualifyService.isAvailable(pointsGift)){
            //资格校验通过
            if(pointsPaymentService.pay(pointsGift)){
                //如果支付积分成功
                String shippingOrderNo = shippingService.shipGift(pointsGift);
                System.out.println("物流系统下单成功,订单号是:"+shippingOrderNo);
            }
        }
    }

}
```

###### 测试类

```java
public class Test {
    public static void main(String[] args) {
        PointsGift pointsGift = new PointsGift("T恤");
        GiftExchangeService giftExchangeService = new GiftExchangeService();
        giftExchangeService.giftExchange(pointsGift);
    }
}
```

#### 装饰器模式

##### 编码演示

###### 煎饼抽象类

```java
public abstract class ABattercake {
    protected abstract String getDesc();
    protected abstract int cost();
}
```

###### 抽象装饰器类

```java
public abstract class AbstractDecorator extends ABattercake {
    private ABattercake aBattercake;

    public AbstractDecorator(ABattercake aBattercake) {
        this.aBattercake = aBattercake;
    }

    protected abstract void doSomething();

    @Override
    protected String getDesc() {
        return this.aBattercake.getDesc();
    }

    @Override
    protected int cost() {
        return this.aBattercake.cost();
    }
}
```

###### 煎饼类

```java
public class Battercake extends ABattercake {
    @Override
    protected String getDesc() {
        return "煎饼";
    }

    @Override
    protected int cost() {
        return 8;
    }
}
```

###### 鸡蛋装饰器类

```java
public class EggDecorator extends AbstractDecorator {
    public EggDecorator(ABattercake aBattercake) {
        super(aBattercake);
    }

    @Override
    protected void doSomething() {

    }

    @Override
    protected String getDesc() {
        return super.getDesc()+" 加一个鸡蛋";
    }

    @Override
    protected int cost() {
        return super.cost()+1;
    }
}
```

###### 香肠装饰器类

```java
public class SausageDecorator extends AbstractDecorator{
    public SausageDecorator(ABattercake aBattercake) {
        super(aBattercake);
    }

    @Override
    protected void doSomething() {

    }

    @Override
    protected String getDesc() {
        return super.getDesc()+" 加一根香肠";
    }

    @Override
    protected int cost() {
        return super.cost()+2;
    }
}
```

###### 测试类

```java
public class Test {
    public static void main(String[] args) {
        ABattercake aBattercake;
        aBattercake = new Battercake();
        aBattercake = new EggDecorator(aBattercake);
        aBattercake = new EggDecorator(aBattercake);
        aBattercake = new SausageDecorator(aBattercake);
        System.out.println(aBattercake.getDesc()+" 销售价格:"+aBattercake.cost());
    }
}
```

#### 适配器模式

##### 编码演示1

###### 交流220V电源类

```java
public class AC220 {
    public int outputAC220V(){
        int output = 220;
        System.out.println("输出交流电"+output+"V");
        return output;
    }
}
```

###### 直流5V电源类

```java
public interface DC5 {
    int outputDC5V();
}
```

###### 电源适配器

```java
public class PowerAdapter implements DC5{
    private AC220 ac220 = new AC220();

    @Override
    public int outputDC5V() {
        int adapterInput = ac220.outputAC220V();
        //变压器...
        int adapterOutput = adapterInput/44;

        System.out.println("使用PowerAdapter输入AC:"+adapterInput+"V"+"输出DC:"+adapterOutput+"V");
        return adapterOutput;
    }
}
```

###### 测试类

```java
public class Test {
    public static void main(String[] args) {
        DC5 dc5 = new PowerAdapter();
        dc5.outputDC5V();
    }
}
```

##### 编码演示2（对象适配器）

###### 被适配者类

```java
public class Adaptee {
    public void adapteeRequest(){
        System.out.println("被适配者的方法");
    }
}
```

###### 目标类

```java
public interface Target {
    void request();
}
```

###### 适配器类

```java
public class Adapter implements Target{
    private Adaptee adaptee = new Adaptee();

    @Override
    public void request() {
        //...
        adaptee.adapteeRequest();
        //...
    }
}
```

###### 测试类

```java
public class Test {
    public static void main(String[] args) {
        Target adapterTarget = new Adapter();
        adapterTarget.request();
    }
}
```

##### 编码演示3（类适配器）

###### 被适配者类

```java
public class Adaptee {
    public void adapteeRequest(){
        System.out.println("被适配者的方法");
    }

}
```

###### 目标类

```java
public interface Target {
    void request();
}
```

###### 适配器类

```java
public class Adapter extends Adaptee implements Target{
    @Override
    public void request() {
        //...
        super.adapteeRequest();
        //...
    }
}
```

###### 测试类

```java
public class Test {
    public static void main(String[] args) {
        Target adapterTarget = new Adapter();
        adapterTarget.request();
    }
}
```

#### 享元模式

##### 编码演示

###### 部门员工接口

```java
public interface Employee {
    void report();
}
```

###### 员工工厂类

```java
public class EmployeeFactory {
    private static final Map<String,Employee> EMPLOYEE_MAP = new HashMap<String,Employee>();

    public static Employee getManager(String department){
        Manager manager = (Manager) EMPLOYEE_MAP.get(department);

        if(manager == null){
            manager = new Manager(department);
            System.out.print("创建部门经理:"+department);
            String reportContent = department+"部门汇报:此次报告的主要内容是......";
            manager.setReportContent(reportContent);
            System.out.println(" 创建报告:"+reportContent);
            EMPLOYEE_MAP.put(department,manager);
        }
        return manager;
    }
}
```

###### 部门经理类

```java
public class Manager implements Employee {
    @Override
    public void report() {
        System.out.println(reportContent);
    }
    private String title = "部门经理";
    private String department;
    private String reportContent;

    public void setReportContent(String reportContent) {
        this.reportContent = reportContent;
    }

    public Manager(String department) {
        this.department = department;
    }
}
```

###### 测试类

```java
public class Test {
    private static final String departments[] = {"RD","QA","PM","BD"};

    public static void main(String[] args) {
        for(int i=0; i<10; i++){
            String department = departments[(int)(Math.random() * departments.length)];
            Manager manager = (Manager) EmployeeFactory.getManager(department);
            manager.report();

        }
    }
}
```

#### 组合模式

##### 编码演示

###### 目录组件抽象类

```java
public abstract class CatalogComponent {
    public void add(CatalogComponent catalogComponent){
        throw new UnsupportedOperationException("不支持添加操作");
    }
    
    public void remove(CatalogComponent catalogComponent){
        throw new UnsupportedOperationException("不支持删除操作");
    }

    public String getName(CatalogComponent catalogComponent){
        throw new UnsupportedOperationException("不支持获取名称操作");
    }

    public double getPrice(CatalogComponent catalogComponent){
        throw new UnsupportedOperationException("不支持获取价格操作");
    }

    public void print(){
        throw new UnsupportedOperationException("不支持打印操作");
    }
}
```

###### 课程类

```java
public class Course extends CatalogComponent {
    private String name;
    private double price;

    public Course(String name, double price) {
        this.name = name;
        this.price = price;
    }

    @Override
    public String getName(CatalogComponent catalogComponent) {
        return this.name;
    }

    @Override
    public double getPrice(CatalogComponent catalogComponent) {
        return this.price;
    }

    @Override
    public void print() {
        System.out.println("Course Name:"+name+" Price:"+price);
    }

}
```

###### 课程目录类

```java
public class CourseCatalog extends CatalogComponent {
    private List<CatalogComponent> items = new ArrayList<CatalogComponent>();
    private String name;
    private Integer level;


    public CourseCatalog(String name,Integer level) {
        this.name = name;
        this.level = level;
    }

    @Override
    public void add(CatalogComponent catalogComponent) {
        items.add(catalogComponent);
    }

    @Override
    public String getName(CatalogComponent catalogComponent) {
        return this.name;
    }

    @Override
    public void remove(CatalogComponent catalogComponent) {
        items.remove(catalogComponent);
    }

    @Override
    public void print() {
        System.out.println(this.name);
        for(CatalogComponent catalogComponent : items){
            if(this.level != null){
                for(int  i = 0; i < this.level; i++){
                    System.out.print("  ");
                }
            }
            catalogComponent.print();
        }
    }

}
```

###### 测试类

```java
public class Test {
    public static void main(String[] args) {
        CatalogComponent linuxCourse = new Course("Linux课程",11);
        CatalogComponent windowsCourse = new Course("Windows课程",11);

        CatalogComponent javaCourseCatalog = new CourseCatalog("Java课程目录",2);

        CatalogComponent mmallCourse1 = new Course("Java电商一期",55);
        CatalogComponent mmallCourse2 = new Course("Java电商二期",66);
        CatalogComponent designPattern = new Course("Java设计模式",77);

        javaCourseCatalog.add(mmallCourse1);
        javaCourseCatalog.add(mmallCourse2);
        javaCourseCatalog.add(designPattern);

        CatalogComponent imoocMainCourseCatalog = new CourseCatalog("慕课网课程主目录",1);
        imoocMainCourseCatalog.add(linuxCourse);
        imoocMainCourseCatalog.add(windowsCourse);
        imoocMainCourseCatalog.add(javaCourseCatalog);
        
        imoocMainCourseCatalog.print();
    }
}
```

#### 桥接模式

##### 编码演示

###### 银行账户接口

```java
public interface Account {
    Account openAccount();
    void showAccountType();
}
```

###### 银行抽象类

```java
public abstract class Bank {
    protected Account account;
    public Bank(Account account){
        this.account = account;
    }
    abstract Account openAccount();
}
```

###### 农业银行类

```java
public class ABCBank extends Bank {
    public ABCBank(Account account) {
        super(account);
    }

    @Override
    Account openAccount() {
        System.out.println("打开中国农业银行账号");
        account.openAccount();
        return account;
    }
}
```

###### 定期账户类

```java
public class DepositAccount implements Account {
    @Override
    public Account openAccount() {
        System.out.println("打开定期账号");
        return new DepositAccount();
    }

    @Override
    public void showAccountType() {
        System.out.println("这是一个定期账号");
    }
}
```

###### 测试类

```java
public class Test {
    public static void main(String[] args) {
        Bank icbcBank = new ICBCBank(new DepositAccount());
        Account icbcAccount = icbcBank.openAccount();
        icbcAccount.showAccountType();

        Bank icbcBank2 = new ICBCBank(new SavingAccount());
        Account icbcAccount2 = icbcBank2.openAccount();
        icbcAccount2.showAccountType();

        Bank abcBank = new ABCBank(new SavingAccount());
        Account abcAccount = abcBank.openAccount();
        abcAccount.showAccountType();
    }
}
```

#### 代理模式

##### 静态代理

```java
public class OrderServiceStaticProxy {
    private IOrderService iOrderService;

    public int saveOrder(Order order){
        beforeMethod(order);
        iOrderService = new OrderServiceImpl();
        int result = iOrderService.saveOrder(order);
        afterMethod();
        return result;
    }

    private void beforeMethod(Order order){
        int userId = order.getUserId();
        int dbRouter = userId % 2;
        System.out.println("静态代理分配到【db"+dbRouter+"】处理数据");

        //todo 设置dataSource;
        DataSourceContextHolder.setDBType("db"+String.valueOf(dbRouter));
        System.out.println("静态代理 before code");
    }
    private void afterMethod(){
        System.out.println("静态代理 after code");
    }
}
```

##### 动态代理

```java
public class OrderServiceDynamicProxy implements InvocationHandler {

    private Object target;

    public OrderServiceDynamicProxy(Object target) {
        this.target = target;
    }

    public Object bind(){
        Class cls = target.getClass();
        return Proxy.newProxyInstance(cls.getClassLoader(),cls.getInterfaces(),this);
    }




    @Override
    public Object invoke(Object proxy, Method method, Object[] args) throws Throwable {
        Object argObject = args[0];
        beforeMethod(argObject);
        Object object = method.invoke(target,args);
        afterMethod();
        return object;
    }

    private void beforeMethod(Object obj){
        int userId = 0;
        System.out.println("动态代理 before code");
        if(obj instanceof Order){
            Order order = (Order)obj;
            userId = order.getUserId();
        }
        int dbRouter = userId % 2;
        System.out.println("动态代理分配到【db"+dbRouter+"】处理数据");

        //todo 设置dataSource;
        DataSourceContextHolder.setDBType("db"+String.valueOf(dbRouter));
    }

    private void afterMethod(){
        System.out.println("动态代理 after code");
    }
}
```

###### 测试类

```java
public class Test {
    public static void main(String[] args) {
        Order order = new Order();
//        order.setUserId(2);
        order.setUserId(1);
        IOrderService orderServiceDynamicProxy = (IOrderService) new OrderServiceDynamicProxy(new OrderServiceImpl()).bind();

        orderServiceDynamicProxy.saveOrder(order);
    }
}
```

### 行为型

#### 模板方法模式

##### 编码演示

###### 课程模板抽象类

```java
public abstract class ACourse {

    protected final void makeCourse(){
        this.makePPT();
        this.makeVideo();
        if(needWriteArticle()){
            this.writeArticle();
        }
        this.packageCourse();
    }

    final void makePPT(){
        System.out.println("制作PPT");
    }
    final void makeVideo(){
        System.out.println("制作视频");
    }
    final void writeArticle(){
        System.out.println("编写手记");
    }
    //钩子方法
    protected boolean needWriteArticle(){
        return false;
    }
    abstract void packageCourse();
}
```

###### 设计模式课程类

```java
public class DesignPatternCourse extends ACourse {
    @Override
    void packageCourse() {
        System.out.println("提供课程Java源代码");
    }

    @Override
    protected boolean needWriteArticle() {
        return true;
    }
}
```

###### 测试类

```java
public class Test {
    public static void main(String[] args) {
//        System.out.println("后端设计模式课程start---");
//        ACourse designPatternCourse = new DesignPatternCourse();
//        designPatternCourse.makeCourse();
//        System.out.println("后端设计模式课程end---");

        System.out.println("前端课程start---");
        ACourse feCourse = new FECourse(false);
        feCourse.makeCourse();
        System.out.println("前端课程end---");
    }
}
```

#### 迭代器模式

##### 编码演示

###### 课程迭代器接口

```java
public interface CourseIterator {
    Course nextCourse();
    boolean isLastCourse();
}
```

###### 课程迭代器实现

```java
public class CourseIteratorImpl implements CourseIterator {

    private List courseList;
    private int position;
    private Course course;
    public CourseIteratorImpl(List courseList){
        this.courseList=courseList;
    }

    @Override
    public Course nextCourse() {
        System.out.println("返回课程,位置是: "+position);
        course=(Course)courseList.get(position);
        position++;
        return course;
    }

    @Override
    public boolean isLastCourse(){
        if(position< courseList.size()){
            return false;
        }
        return true;
    }
}
```

###### 测试类

```java
public class Test {
    public static void main(String[] args) {
        Course course1 = new Course("Java电商一期");
        Course course2 = new Course("Java电商二期");
        Course course3 = new Course("Java设计模式精讲");
        Course course4 = new Course("Python课程");
        Course course5 = new Course("算法课程");
        Course course6 = new Course("前端课程");

        CourseAggregate courseAggregate = new CourseAggregateImpl();

        courseAggregate.addCourse(course1);
        courseAggregate.addCourse(course2);
        courseAggregate.addCourse(course3);
        courseAggregate.addCourse(course4);
        courseAggregate.addCourse(course5);
        courseAggregate.addCourse(course6);

        System.out.println("-----课程列表-----");
        printCourses(courseAggregate);

        courseAggregate.removeCourse(course4);
        courseAggregate.removeCourse(course5);

        System.out.println("-----删除操作之后的课程列表-----");
        printCourses(courseAggregate);
    }

    public static void printCourses(CourseAggregate courseAggregate){
        CourseIterator courseIterator= courseAggregate.getCourseIterator();
        while(!courseIterator.isLastCourse()){
            Course course=courseIterator.nextCourse();
            System.out.println(course.getName());
        }
    }
}
```

#### 策略模式

##### 编码演示

###### 促销策略接口

```java
public interface PromotionStrategy {
    void doPromotion();
}
```

###### 促销活动类

```java
public class PromotionActivity {
    private PromotionStrategy promotionStrategy;

    public PromotionActivity(PromotionStrategy promotionStrategy) {
        this.promotionStrategy = promotionStrategy;
    }

    public void executePromotionStrategy(){
        promotionStrategy.doPromotion();
    }
}
```

###### 促销策略工厂

```java
public class PromotionStrategyFactory {
    private static Map<String,PromotionStrategy> PROMOTION_STRATEGY_MAP = new HashMap<String, PromotionStrategy>();
    static {
        PROMOTION_STRATEGY_MAP.put(PromotionKey.LIJIAN,new LiJianPromotionStrategy());
        PROMOTION_STRATEGY_MAP.put(PromotionKey.FANXIAN,new FanXianPromotionStrategy());
        PROMOTION_STRATEGY_MAP.put(PromotionKey.MANJIAN,new ManJianPromotionStrategy());
    }

    private static final PromotionStrategy NON_PROMOTION = new EmptyPromotionStrategy();

    private PromotionStrategyFactory(){

    }

    public static PromotionStrategy getPromotionStrategy(String promotionKey){
        PromotionStrategy promotionStrategy = PROMOTION_STRATEGY_MAP.get(promotionKey);
        return promotionStrategy == null ? NON_PROMOTION : promotionStrategy;
    }

    private interface PromotionKey{
        String LIJIAN = "LIJIAN";
        String FANXIAN = "FANXIAN";
        String MANJIAN = "MANJIAN";
    }
}
```

###### 测试类

```java
public class Test {
    public static void main(String[] args) {
        PromotionActivity promotionActivity618 = new PromotionActivity(new LiJianPromotionStrategy());
        PromotionActivity promotionActivity1111 = new PromotionActivity(new FanXianPromotionStrategy());

        promotionActivity618.executePromotionStrategy();
        promotionActivity1111.executePromotionStrategy();
    }
}
```

#### 解释器模式

##### 编码演示

###### 解释器接口

```java
public interface Interpreter {
    int interpret();
}
```

###### 数字解释器

```java
public class NumberInterpreter implements Interpreter {
    private int number;
    public NumberInterpreter(int number){
        this.number=number;
    }
    public NumberInterpreter(String number){
        this.number=Integer.parseInt(number);
    }
    @Override
    public int interpret(){
        return this.number;
    }
}
```

###### 加法解释器

```java
public class AddInterpreter implements Interpreter {
    private Interpreter firstExpression,secondExpression;
    public AddInterpreter(Interpreter firstExpression, Interpreter secondExpression){
        this.firstExpression=firstExpression;
        this.secondExpression=secondExpression;
    }
    @Override
    public int interpret(){
        return this.firstExpression.interpret()+this.secondExpression.interpret();
    }
}
```

###### 操作工具类

```java
public class OperatorUtil {
    public static boolean isOperator(String symbol) {
        return (symbol.equals("+") || symbol.equals("*"));
    }

    public static Interpreter getExpressionObject(Interpreter firstExpression, Interpreter secondExpression, String symbol) {
        if (symbol.equals("+")) {
            return new AddInterpreter(firstExpression, secondExpression);
        } else if (symbol.equals("*")) {
            return new MultiInterpreter(firstExpression, secondExpression);
        }
        return null;
    }
}
```

###### 表达式解释器

```java
public class GeelyExpressionParser {
    private Stack<Interpreter> stack = new Stack<Interpreter>();

    public int parse(String str) {
        String[] strItemArray = str.split(" ");
        for (String symbol : strItemArray) {
            if (!OperatorUtil.isOperator(symbol)) {
                Interpreter numberExpression = new NumberInterpreter(symbol);
                stack.push(numberExpression);
                System.out.println(String.format("入栈: %d", numberExpression.interpret()));
            } else {
                //是运算符号，可以计算
                Interpreter firstExpression = stack.pop();
                Interpreter secondExpression = stack.pop();
                System.out.println(String.format("出栈: %d 和 %d",
                        firstExpression.interpret(), secondExpression.interpret()));
                Interpreter operator = OperatorUtil.getExpressionObject(firstExpression, secondExpression, symbol);
                System.out.println(String.format("应用运算符: %s", operator));
                int result = operator.interpret();
                NumberInterpreter resultExpression = new NumberInterpreter(result);
                stack.push(resultExpression);
                System.out.println(String.format("阶段结果入栈: %d", resultExpression.interpret()));
            }
        }
        int result = stack.pop().interpret();
        return result;
    }
}
```

###### 测试类

```java
public class Test {
    public static void main(String[] args) {
        String geelyInputStr="6 100 11 + *";
        GeelyExpressionParser expressionParser=new GeelyExpressionParser();
        int result=expressionParser.parse(geelyInputStr);
        System.out.println("解释器计算结果: "+result);
    }
}
```

#### 观察者模式

##### 编码演示1

###### 课程类（被观察对象）

```java
public class Course extends Observable{
    private String courseName;

    public Course(String courseName) {
        this.courseName = courseName;
    }

    public String getCourseName() {
        return courseName;
    }

    public void produceQuestion(Course course, Question question){
        System.out.println(question.getUserName()+"在"+course.courseName+"提交了一个问题");
        setChanged();
        notifyObservers(question);
    }
}
```

###### 老师类（观察者）

```java
public class Teacher implements Observer{
    private String teacherName;

    public Teacher(String teacherName) {
        this.teacherName = teacherName;
    }
    
    @Override
    public void update(Observable o, Object arg) {
        Course course = (Course)o;
        Question question = (Question)arg;
        System.out.println(teacherName+"老师的"+course.getCourseName()+"课程接收到一个"+question.getUserName()+"提交的问答:"+question.getQuestionContent());
    }
}
```

###### 测试类

```java
public class Test {
    public static void main(String[] args) {
        Course course = new Course("Java设计模式精讲");
        Teacher teacher = new Teacher("Alpha");
        course.addObserver(teacher);
        //业务逻辑代码
        Question question = new Question();
        question.setUserName("Geely");
        question.setQuestionContent("Java的主函数如何编写");
        course.produceQuestion(course,question);
    }
}
```

##### 编码演示2（Guava）

```java
public class GuavaEvent {
    @Subscribe
    public void subscribe(String str){
        //业务逻辑
        System.out.println("执行subscribe方法,传入的参数是:" + str);
    }
}
```

```java
public class GuavaEventTest {
    public static void main(String[] args) {
        EventBus eventbus = new EventBus();
        GuavaEvent guavaEvent = new GuavaEvent();
        eventbus.register(guavaEvent);
        eventbus.post("post的内容");
    }
}
```

#### 备忘录模式

##### 编码演示

###### 文章类

```java
public class Article {

    private String title;
    private String content;
    private String imgs;

    public Article(String title, String content, String imgs) {
        this.title = title;
        this.content = content;
        this.imgs = imgs;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

    public String getImgs() {
        return imgs;
    }

    public void setImgs(String imgs) {
        this.imgs = imgs;
    }

    public ArticleMemento saveToMemento() {
        ArticleMemento articleMemento = new ArticleMemento(this.title,this.content,this.imgs);
        return articleMemento;
    }

    public void undoFromMemento(ArticleMemento articleMemento) {

        this.title = articleMemento.getTitle();
        this.content = articleMemento.getContent();
        this.imgs = articleMemento.getImgs();
    }
}
```

###### 文章备忘录类

```java
public class ArticleMemento {
    private String title;
    private String content;
    private String imgs;

    public ArticleMemento(String title, String content, String imgs) {
        this.title = title;
        this.content = content;
        this.imgs = imgs;
    }

    public String getTitle() {
        return title;
    }

    public String getContent() {
        return content;
    }

    public String getImgs() {
        return imgs;
    }
}
```

###### 文章备忘录管理类

```java
public class ArticleMementoManager {

    private final Stack<ArticleMemento> ARTICLE_MEMENTO_STACK = new Stack<ArticleMemento>();

    public ArticleMemento getMemento()
    {
        ArticleMemento articleMemento= ARTICLE_MEMENTO_STACK.pop();
        return articleMemento;
    }

    public void addMemento(ArticleMemento articleMemento)
    {
        ARTICLE_MEMENTO_STACK.push(articleMemento);
    }

}
```

###### 测试类

```java
public class Test {

    public static void main(String[] args) {
        ArticleMementoManager articleMementoManager = new ArticleMementoManager();

        Article article= new Article("如影随行的设计模式A","手记内容A","手记图片A");

        ArticleMemento articleMemento = article.saveToMemento();

        articleMementoManager.addMemento(articleMemento);
        System.out.println("标题:"+article.getTitle()+" 内容:"+article.getContent()+" 图片:"+article.getImgs()+" 暂存成功");

        System.out.println("手记完整信息:"+article);
        System.out.println("修改手记start");

        article.setTitle("如影随行的设计模式B");
        article.setContent("手记内容B");
        article.setImgs("手记图片B");

        System.out.println("修改手记end");
        System.out.println("手记完整信息:"+article);

        articleMemento = article.saveToMemento();
        articleMementoManager.addMemento(articleMemento);

        article.setTitle("如影随行的设计模式C");
        article.setContent("手记内容C");
        article.setImgs("手记图片C");

        System.out.println("暂存回退start");

        System.out.println("回退出栈1次");
        articleMemento = articleMementoManager.getMemento();
        article.undoFromMemento(articleMemento);

        System.out.println("回退出栈2次");
        articleMemento = articleMementoManager.getMemento();
        article.undoFromMemento(articleMemento);
        
        System.out.println("暂存回退end");
        System.out.println("手记完整信息:"+article);
    }
}
```